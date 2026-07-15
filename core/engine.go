package core

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"
	"unicode"

	"github.com/gigol/irecall/core/db"
	"github.com/gigol/irecall/core/llm"
	"github.com/google/uuid"
)

// Engine is the central orchestrator. It owns the DB store and LLM client
// and exposes the full recall workflow. No UI types are referenced here.
type Engine struct {
	store          *db.Store
	llm            *llm.Client
	keywordLLM     *llm.Client
	cfg            *Settings
	profile        *UserProfile
	attachmentRoot string
}

const maxExtractedTags = 30

var genericTagBlacklist = map[string]struct{}{
	"content":     {},
	"idea":        {},
	"ideas":       {},
	"info":        {},
	"information": {},
	"item":        {},
	"items":       {},
	"note":        {},
	"notes":       {},
	"quote":       {},
	"quotes":      {},
	"summary":     {},
	"summaries":   {},
	"text":        {},
	"topic":       {},
	"topics":      {},
}

// New creates an Engine from an open DB store and current settings.
func New(store *db.Store, cfg *Settings) *Engine {
	slog.Info(
		"engine: creating engine",
		"provider_host", cfg.Provider.Host,
		"provider_port", cfg.Provider.Port,
		"response_model", cfg.Provider.ResponseModel(),
		"keyword_model", cfg.Provider.RecallKeywordModel(),
	)
	return &Engine{
		store:      store,
		llm:        llm.NewClient(cfg.Provider.ResponseProviderConfig()),
		keywordLLM: llm.NewClient(cfg.Provider.RecallKeywordProviderConfig()),
		cfg:        cfg,
	}
}

func (e *Engine) Close() error {
	slog.Info("engine: closing")
	return e.store.Close()
}

// UpdateProvider rebuilds the LLM client after settings change.
func (e *Engine) UpdateProvider(cfg ProviderConfig) {
	slog.Info(
		"engine: updating provider",
		"host", cfg.Host,
		"port", cfg.Port,
		"response_model", cfg.ResponseModel(),
		"keyword_model", cfg.RecallKeywordModel(),
	)
	e.cfg.Provider = cfg
	e.llm = llm.NewClient(cfg.ResponseProviderConfig())
	e.keywordLLM = llm.NewClient(cfg.RecallKeywordProviderConfig())
}

// UpdateSettings replaces the engine's in-memory settings.
func (e *Engine) UpdateSettings(s *Settings) {
	slog.Info(
		"engine: updating settings",
		"host", s.Provider.Host,
		"port", s.Provider.Port,
		"response_model", s.Provider.ResponseModel(),
		"keyword_model", s.Provider.RecallKeywordModel(),
		"max_results", s.Search.MaxResults,
	)
	e.cfg = s
	e.llm = llm.NewClient(s.Provider.ResponseProviderConfig())
	e.keywordLLM = llm.NewClient(s.Provider.RecallKeywordProviderConfig())
}

func (e *Engine) UpdateUserProfile(profile *UserProfile) {
	e.profile = profile
}

// --- Quote management ---

// AddQuote stores a new quote and extracts tags via LLM.
// If tag extraction fails, the quote is still saved with no tags.
func (e *Engine) AddQuote(ctx context.Context, content string) (*Quote, error) {
	return e.addQuoteWithSuggestedTags(ctx, content, nil)
}

func (e *Engine) AddQuoteWithImages(ctx context.Context, content string, images []ImageInput) (*QuoteMutationResult, error) {
	return e.addQuoteWithImagesAndSuggestedTags(ctx, content, nil, images)
}

func (e *Engine) SaveRecallAsQuote(ctx context.Context, question, response string, keywords []string) (*Quote, error) {
	question = strings.TrimSpace(question)
	response = strings.TrimSpace(response)
	if question == "" {
		return nil, fmt.Errorf("question is empty")
	}
	if response == "" {
		return nil, fmt.Errorf("response is empty")
	}

	if len(keywords) == 0 {
		var err error
		keywords, err = e.ExtractKeywords(ctx, question)
		if err != nil {
			return nil, fmt.Errorf("extract keywords: %w", err)
		}
	}

	content := formatRecallQuote(question, response)
	return e.addQuoteWithSuggestedTags(ctx, content, keywords)
}

func (e *Engine) addQuoteWithSuggestedTags(ctx context.Context, content string, suggestedTags []string) (*Quote, error) {
	result, err := e.addQuoteWithImagesAndSuggestedTags(ctx, content, suggestedTags, nil)
	if err != nil {
		return nil, err
	}
	return &result.Quote, nil
}

func (e *Engine) addQuoteWithImagesAndSuggestedTags(ctx context.Context, content string, suggestedTags []string, images []ImageInput) (*QuoteMutationResult, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, fmt.Errorf("quote content is empty")
	}
	prepared, err := prepareImages(images, 0)
	if err != nil {
		return nil, err
	}

	slog.Info("engine: adding quote", "content_len", len(content), "content_preview", truncate(content, 100))

	identity, err := e.quoteIdentityForNewQuote()
	if err != nil {
		return nil, err
	}

	id, err := e.store.InsertQuote(content, identity)
	if err != nil {
		slog.Error("engine: failed to insert quote", "error", err)
		return nil, fmt.Errorf("store quote: %w", err)
	}
	slog.Info("engine: quote inserted", "id", id)
	var storedRows []db.QuoteAttachmentRow
	for _, image := range prepared {
		row, writeErr := e.writePreparedImage(id, image)
		if writeErr != nil {
			for _, stored := range storedRows {
				_ = os.Remove(filepath.Join(e.attachmentRoot, stored.StoragePath))
			}
			_ = e.store.DeleteQuote(id)
			return nil, fmt.Errorf("store image: %w", writeErr)
		}
		if insertErr := e.store.InsertQuoteAttachment(row); insertErr != nil {
			_ = os.Remove(filepath.Join(e.attachmentRoot, row.StoragePath))
			for _, stored := range storedRows {
				_ = os.Remove(filepath.Join(e.attachmentRoot, stored.StoragePath))
			}
			_ = e.store.DeleteQuote(id)
			return nil, insertErr
		}
		storedRows = append(storedRows, row)
	}

	slog.Info("engine: extracting tags via LLM", "quote_id", id)
	extractedTags, warnings, err := e.extractTagsWithImageInputs(ctx, content, images)
	if err != nil {
		slog.Error("engine: tag extraction failed, saving without tags", "quote_id", id, "error", err)
		extractedTags = nil
	}
	tags, _ := normalizeTags(append(slices.Clone(suggestedTags), extractedTags...))
	slog.Info("engine: tags extracted", "quote_id", id, "tags", tags)

	if len(tags) > 0 {
		tagIDs, err := e.store.UpsertTags(tags)
		if err != nil {
			slog.Error("engine: upsert tags failed", "quote_id", id, "error", err)
		} else {
			slog.Debug("engine: tag IDs", "quote_id", id, "tag_ids", tagIDs)
			if err := e.store.InsertQuoteTags(id, tagIDs); err != nil {
				slog.Error("engine: insert quote-tags failed", "quote_id", id, "error", err)
			}
			if err := e.store.UpdateQuoteFTS(id, tags); err != nil {
				slog.Error("engine: update FTS failed", "quote_id", id, "error", err)
			}
		}
	} else {
		slog.Warn("engine: no tags for quote, FTS will only index content", "quote_id", id)
	}

	slog.Info("engine: add quote complete", "id", id, "tag_count", len(tags))
	quote, err := e.loadQuote(id)
	if err != nil {
		return nil, err
	}
	return &QuoteMutationResult{Quote: *quote, Warnings: warnings}, nil
}

// ListQuotes returns all quotes, newest first.
func (e *Engine) ListQuotes(ctx context.Context) ([]Quote, error) {
	return e.ListQuotesPage(ctx, 0, 0)
}

// CountQuotes returns the total number of stored quotes.
func (e *Engine) CountQuotes(ctx context.Context) (int64, error) {
	slog.Debug("engine: counting quotes")
	count, err := e.store.CountQuotes()
	if err != nil {
		slog.Error("engine: count quotes failed", "error", err)
		return 0, err
	}
	slog.Debug("engine: counted quotes", "count", count)
	return count, nil
}

// ListQuotesPage returns quotes newest first using limit/offset pagination.
// limit <= 0 returns all quotes. offset < 0 is treated as 0.
func (e *Engine) ListQuotesPage(ctx context.Context, limit, offset int) ([]Quote, error) {
	slog.Debug("engine: listing quotes", "limit", limit, "offset", offset)
	rows, err := e.store.ListQuotesPage(limit, offset)
	if err != nil {
		slog.Error("engine: list quotes failed", "error", err)
		return nil, err
	}
	quotes := rowsToQuotes(rows, e.localUserID())
	if err := e.enrichQuoteAttachments(quotes); err != nil {
		return nil, err
	}
	slog.Debug("engine: listed quotes", "count", len(quotes), "limit", limit, "offset", offset)
	return quotes, nil
}

// DeleteQuote removes a quote by ID.
func (e *Engine) DeleteQuote(ctx context.Context, id int64) error {
	slog.Info("engine: deleting quote", "id", id)
	attachments, err := e.store.ListQuoteAttachments(id)
	if err != nil {
		return err
	}
	if err := e.store.DeleteQuote(id); err != nil {
		return err
	}
	for _, a := range attachments {
		_ = os.Remove(filepath.Join(e.attachmentRoot, filepath.Base(a.StoragePath)))
	}
	return nil
}

// DeleteQuotes removes multiple quotes by ID.
func (e *Engine) DeleteQuotes(ctx context.Context, ids []int64) error {
	for _, id := range ids {
		if err := e.DeleteQuote(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

// UpdateQuote rewrites quote content, regenerates tags, and refreshes FTS.
func (e *Engine) UpdateQuote(ctx context.Context, id int64, content string) (*Quote, error) {
	rows, err := e.store.ListQuoteAttachments(id)
	if err != nil {
		return nil, err
	}
	retained := make([]string, len(rows))
	for i, row := range rows {
		retained[i] = row.ID
	}
	result, err := e.UpdateQuoteWithImages(ctx, id, content, retained, nil)
	if err != nil {
		return nil, err
	}
	return &result.Quote, nil
}

func (e *Engine) UpdateQuoteWithImages(ctx context.Context, id int64, content string, retainedIDs []string, images []ImageInput) (*QuoteMutationResult, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, fmt.Errorf("quote content is empty")
	}

	current, err := e.store.ListQuoteAttachments(id)
	if err != nil {
		return nil, err
	}
	byID := make(map[string]db.QuoteAttachmentRow, len(current))
	for _, row := range current {
		byID[row.ID] = row
	}
	seen := map[string]bool{}
	var retained []db.QuoteAttachmentRow
	for _, attachmentID := range retainedIDs {
		if seen[attachmentID] {
			continue
		}
		row, ok := byID[attachmentID]
		if !ok {
			return nil, fmt.Errorf("attachment %q does not belong to quote %d", attachmentID, id)
		}
		seen[attachmentID] = true
		retained = append(retained, row)
	}
	prepared, err := prepareImages(images, len(retained))
	if err != nil {
		return nil, err
	}
	totalBytes := int64(0)
	for _, row := range retained {
		totalBytes += row.Size
	}
	for _, image := range images {
		totalBytes += int64(len(image.Data))
	}
	if totalBytes > MaxQuoteImageBytes {
		return nil, fmt.Errorf("quote images exceed the 25 MiB total limit")
	}
	analysisImages := make([]ImageInput, 0, len(retained)+len(images))
	for _, row := range retained {
		data, readErr := e.attachmentData(row)
		if readErr != nil {
			return nil, readErr
		}
		analysisImages = append(analysisImages, ImageInput{Filename: row.Filename, MediaType: row.MediaType, Data: data})
	}
	analysisImages = append(analysisImages, images...)
	tags, warnings, tagErr := e.extractTagsWithImageInputs(ctx, content, analysisImages)
	if tagErr != nil {
		tags = []string{}
	}

	for _, image := range prepared {
		row, writeErr := e.writePreparedImage(id, image)
		if writeErr != nil {
			return nil, fmt.Errorf("store image: %w", writeErr)
		}
		if insertErr := e.store.InsertQuoteAttachment(row); insertErr != nil {
			_ = os.Remove(filepath.Join(e.attachmentRoot, row.StoragePath))
			return nil, insertErr
		}
	}

	slog.Info("engine: updating quote", "id", id, "content_len", len(content), "content_preview", truncate(content, 100))
	if err := e.store.UpdateQuoteContent(id, content); err != nil {
		return nil, fmt.Errorf("update quote content: %w", err)
	}

	var tagIDs []int64
	if len(tags) > 0 {
		tagIDs, err = e.store.UpsertTags(tags)
		if err != nil {
			slog.Error("engine: upsert tags failed during update", "quote_id", id, "error", err)
			tags = []string{}
			tagIDs = nil
		}
	}
	if err := e.store.ReplaceQuoteTags(id, tagIDs); err != nil {
		return nil, fmt.Errorf("replace quote tags: %w", err)
	}
	if err := e.store.UpdateQuoteFTS(id, tags); err != nil {
		return nil, fmt.Errorf("update quote fts: %w", err)
	}
	for _, row := range current {
		if seen[row.ID] {
			continue
		}
		if err := e.store.DeleteQuoteAttachment(row.ID); err != nil {
			return nil, err
		}
		_ = os.Remove(filepath.Join(e.attachmentRoot, filepath.Base(row.StoragePath)))
	}
	quote, err := e.loadQuote(id)
	if err != nil {
		return nil, err
	}
	return &QuoteMutationResult{Quote: *quote, Warnings: warnings}, nil
}

// RegenerateQuoteKeywords re-runs tag extraction for one stored quote and persists the result.
func (e *Engine) RegenerateQuoteKeywords(ctx context.Context, id int64, globalID string) (*QuoteKeywordRegeneration, error) {
	quote, err := e.loadQuoteByIdentifier(id, globalID)
	if err != nil {
		return nil, err
	}

	oldKeywords := slices.Clone(quote.Tags)
	slog.Info("engine: regenerating quote keywords", "id", quote.ID, "global_id", quote.GlobalID, "old_keywords", oldKeywords)

	var images []ImageInput
	rows, err := e.store.ListQuoteAttachments(quote.ID)
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		data, readErr := e.attachmentData(row)
		if readErr != nil {
			return nil, readErr
		}
		images = append(images, ImageInput{Filename: row.Filename, MediaType: row.MediaType, Data: data})
	}
	newKeywords, _, err := e.extractTagsWithImageInputs(ctx, quote.Content, images)
	if err != nil {
		return nil, fmt.Errorf("extract quote keywords: %w", err)
	}

	changed := !sameStringSet(oldKeywords, newKeywords)
	if changed {
		if err := e.store.ApplyQuoteTagUpdate(quote.ID, newKeywords, true); err != nil {
			return nil, err
		}
		quote, err = e.loadQuote(quote.ID)
		if err != nil {
			return nil, err
		}
	}

	return &QuoteKeywordRegeneration{
		QuoteID:     quote.ID,
		GlobalID:    quote.GlobalID,
		OldKeywords: oldKeywords,
		NewKeywords: slices.Clone(newKeywords),
		Changed:     changed,
		Quote:       *quote,
	}, nil
}

// RefineQuoteDraft asks the LLM to rewrite a draft quote for clarity while preserving intent.
func (e *Engine) RefineQuoteDraft(ctx context.Context, content string) (string, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return "", fmt.Errorf("quote content is empty")
	}
	if e.isMockLLMEnabled() {
		slog.Info("engine: mock llm enabled for quote refinement")
		return content, nil
	}

	slog.Info("engine: refining quote draft", "content_len", len(content), "content_preview", truncate(content, 100))

	msgs := []llm.Message{
		{
			Role: "system",
			Content: `You refine personal notes for clarity and readability. ` +
				`Keep the original meaning, facts, and intent. ` +
				`Do not add new information. ` +
				`Return only the rewritten note text with no explanation, no markdown, and no surrounding quotes.`,
		},
		{
			Role:    "user",
			Content: "Rewrite this note so it reads more clearly while preserving its meaning:\n\n" + content,
		},
	}

	zero := 0.0
	maxTok := 400
	noReasoning := "none"
	refined, err := e.llm.Chat(ctx, msgs, nil, llm.ChatOptions{
		Temperature:     &zero,
		MaxTokens:       &maxTok,
		ReasoningEffort: &noReasoning,
	})
	if err != nil {
		slog.Error("engine: refine quote draft LLM call failed", "error", err)
		return "", err
	}

	refined = strings.TrimSpace(refined)
	if refined == "" {
		return "", fmt.Errorf("provider returned empty refined note")
	}

	slog.Info("engine: refined quote draft", "response_len", len(refined))
	return refined, nil
}

// --- Recall workflow ---

// ExtractTags asks the LLM to produce keyword tags for a piece of text.
func (e *Engine) ExtractTags(ctx context.Context, text string) ([]string, error) {
	slog.Debug("engine: extracting tags", "text_len", len(text))
	msgs := []llm.Message{
		{
			Role: "system",
			Content: `You are a JSON keyword extractor. ` +
				`Output ONLY a valid JSON array of short lowercase English keyword strings. ` +
				`Translate non-English concepts into natural English search terms, while preserving proper nouns, product names, acronyms, and code identifiers. ` +
				`Prefer broad, relevant coverage. ` +
				`For dense or technical text, include a rich set of tags, usually up to 30. ` +
				`For short or simple text, return fewer tags when appropriate. ` +
				`No explanation, no markdown, no code fences, no extra text — just the JSON array. ` +
				`Example output: ["emmc", "flash memory", "partition", "offset"]`,
		},
		{
			Role:    "user",
			Content: "Extract keyword tags for this text:\n" + text,
		},
	}
	zero := 0.0
	maxTok := 384
	raw, err := e.keywordLLM.Chat(ctx, msgs, nil, llm.ChatOptions{Temperature: &zero, MaxTokens: &maxTok})
	if err != nil {
		slog.Error("engine: extract tags LLM call failed", "error", err)
		return nil, err
	}
	slog.Debug("engine: extract tags raw response", "raw", raw)
	tags, err := parseJSONStringArray(raw)
	if err != nil {
		slog.Error("engine: parse tags failed", "raw", raw, "error", err)
		return nil, err
	}
	normalized, stats := normalizeTags(tags)
	if !shouldRepairTags(normalized, stats, text) {
		slog.Info("engine: extracted tags", "tags", normalized)
		return normalized, nil
	}

	repaired, repairErr := e.repairTags(ctx, text, normalized)
	if repairErr != nil {
		slog.Warn("engine: tag repair failed, keeping initial tags", "error", repairErr, "tags", normalized)
		return normalized, nil
	}
	slog.Info("engine: repaired tags", "tags", repaired)
	return repaired, nil
}

func (e *Engine) extractTagsWithImageInputs(ctx context.Context, text string, images []ImageInput) ([]string, []string, error) {
	if e.isMockLLMEnabled() {
		tags, _ := normalizeTags(mockSplitKeywords(text))
		return tags, nil, nil
	}
	if len(images) == 0 {
		tags, err := e.ExtractTags(ctx, text)
		return tags, nil, err
	}
	parts := []llm.ContentPart{{Type: "text", Text: "Extract keyword tags for this text and its attached images:\n" + text}}
	for _, image := range images {
		mediaType := image.MediaType
		if mediaType == "" && len(image.Data) > 0 {
			mediaType = http.DetectContentType(image.Data[:min(len(image.Data), 512)])
		}
		parts = append(parts, llm.ContentPart{Type: "image_url", ImageURL: &llm.ImageURLValue{
			URL: "data:" + mediaType + ";base64," + base64.StdEncoding.EncodeToString(image.Data),
		}})
	}
	msgs := []llm.Message{
		{Role: "system", Content: `You are a JSON keyword extractor. Output ONLY a valid JSON array of short lowercase English keyword strings. Consider both the supplied note text and every attached image. Translate non-English concepts into natural English search terms, while preserving proper nouns, product names, acronyms, and code identifiers. Prefer high-signal concepts and return up to 30 tags. No explanation, markdown, or extra text.`},
		{Role: "user", Parts: parts},
	}
	zero := 0.0
	maxTok := 384
	raw, err := e.keywordLLM.Chat(ctx, msgs, nil, llm.ChatOptions{Temperature: &zero, MaxTokens: &maxTok})
	if err == nil {
		var parsed []string
		parsed, err = parseJSONStringArray(raw)
		if err == nil {
			normalized, stats := normalizeTags(parsed)
			if shouldRepairTags(normalized, stats, text) {
				if repaired, repairErr := e.repairTags(ctx, text, normalized); repairErr == nil {
					normalized = repaired
				}
			}
			return normalized, nil, nil
		}
	}
	slog.Warn("engine: image tag extraction failed; retrying text only", "error", err)
	tags, fallbackErr := e.ExtractTags(ctx, text)
	warning := "The configured model could not analyze the attached images; tags were generated from quote text only."
	if fallbackErr != nil {
		return nil, []string{warning}, fallbackErr
	}
	return tags, []string{warning}, nil
}

func (e *Engine) repairTags(ctx context.Context, text string, initial []string) ([]string, error) {
	msgs := []llm.Message{
		{
			Role: "system",
			Content: `You are repairing a JSON keyword extractor result. ` +
				`Return ONLY a valid JSON array of short lowercase English keyword strings. ` +
				`Translate non-English concepts into natural English search terms, while preserving proper nouns, product names, acronyms, and code identifiers. ` +
				`Prefer high-signal tags: specific technologies, entities, actions, domains, and concepts. ` +
				`Avoid generic labels such as "quote", "note", "text", "content", "topic", or "summary". ` +
				`Return up to 30 tags. For short or simple text, fewer tags are appropriate.`,
		},
		{
			Role: "user",
			Content: "Improve these draft tags for the text below.\n\n" +
				"Draft tags: " + mustMarshalJSONArray(initial) + "\n\n" +
				"Text:\n" + text,
		},
	}
	zero := 0.0
	maxTok := 384
	raw, err := e.keywordLLM.Chat(ctx, msgs, nil, llm.ChatOptions{Temperature: &zero, MaxTokens: &maxTok})
	if err != nil {
		return nil, err
	}
	tags, err := parseJSONStringArray(raw)
	if err != nil {
		return nil, err
	}
	normalized, _ := normalizeTags(tags)
	if len(normalized) == 0 {
		return nil, fmt.Errorf("repair produced no usable tags")
	}
	return normalized, nil
}

// ExtractKeywords asks the LLM to produce search keywords for a question.
func (e *Engine) ExtractKeywords(ctx context.Context, question string) ([]string, error) {
	slog.Info("engine: extracting keywords for question", "question", question)
	if e.isMockLLMEnabled() {
		slog.Info("engine: mock llm enabled for keyword extraction")
		return mockSplitKeywords(question), nil
	}
	msgs := []llm.Message{
		{
			Role: "system",
			Content: `You are a JSON search keyword extractor. ` +
				`Output ONLY a valid JSON array of 3 to 6 short lowercase English keyword strings useful for searching a knowledge base. ` +
				`Translate non-English questions into the same English keyword language used by stored quote tags. ` +
				`Preserve proper nouns, product names, acronyms, and code identifiers. ` +
				`No explanation, no markdown, no code fences, no extra text — just the JSON array. ` +
				`Example output: ["emmc", "flash", "partition"]`,
		},
		{
			Role:    "user",
			Content: "Extract search keywords for this question:\n" + question,
		},
	}
	zero := 0.0
	maxTok := 100
	raw, err := e.keywordLLM.Chat(ctx, msgs, nil, llm.ChatOptions{Temperature: &zero, MaxTokens: &maxTok})
	if err != nil {
		slog.Error("engine: extract keywords LLM call failed", "error", err)
		return nil, err
	}
	slog.Debug("engine: extract keywords raw response", "raw", raw)
	keywords, err := parseJSONStringArray(raw)
	if err != nil {
		slog.Error("engine: parse keywords failed", "raw", raw, "error", err)
		return nil, err
	}
	normalized, _ := normalizeTags(keywords)
	if len(normalized) == 0 {
		return nil, fmt.Errorf("keyword extraction produced no usable keywords")
	}
	if containsNonASCIIKeyword(normalized) {
		translated, translateErr := e.repairSearchKeywords(ctx, question, normalized)
		if translateErr != nil {
			slog.Warn("engine: search keyword translation repair failed, keeping initial keywords", "error", translateErr, "keywords", normalized)
		} else {
			normalized = translated
		}
	}
	keywords = appendCrossLanguageSearchFallbacks(normalized, question)
	slog.Info("engine: extracted keywords", "keywords", keywords)
	return keywords, nil
}

func (e *Engine) repairSearchKeywords(ctx context.Context, question string, initial []string) ([]string, error) {
	msgs := []llm.Message{
		{
			Role: "system",
			Content: `You are repairing a JSON search keyword extractor result. ` +
				`Return ONLY a valid JSON array of 3 to 6 short lowercase English keyword strings useful for searching a knowledge base. ` +
				`Translate every non-English keyword into natural English search terms. ` +
				`Preserve proper nouns, product names, acronyms, and code identifiers. ` +
				`No explanation, no markdown, no code fences, no extra text — just the JSON array.`,
		},
		{
			Role: "user",
			Content: "Improve these draft search keywords for the question below.\n\n" +
				"Draft keywords: " + mustMarshalJSONArray(initial) + "\n\n" +
				"Question:\n" + question,
		},
	}
	zero := 0.0
	maxTok := 100
	raw, err := e.keywordLLM.Chat(ctx, msgs, nil, llm.ChatOptions{Temperature: &zero, MaxTokens: &maxTok})
	if err != nil {
		return nil, err
	}
	keywords, err := parseJSONStringArray(raw)
	if err != nil {
		return nil, err
	}
	normalized, _ := normalizeTags(keywords)
	if len(normalized) == 0 {
		return nil, fmt.Errorf("repair produced no usable keywords")
	}
	return normalized, nil
}

func containsNonASCIIKeyword(keywords []string) bool {
	for _, kw := range keywords {
		if containsNonASCII(kw) {
			return true
		}
	}
	return false
}

func appendCrossLanguageSearchFallbacks(keywords []string, question string) []string {
	out := make([]string, 0, len(keywords)+8)
	seen := make(map[string]struct{}, len(keywords)+4)
	add := func(value string) {
		value = normalizeSearchToken(value)
		if value == "" {
			return
		}
		if _, ok := seen[value]; ok {
			return
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	for _, kw := range keywords {
		add(kw)
	}
	if containsNonASCII(question) {
		for _, token := range originalQuestionFallbackTokens(question) {
			add(token)
		}
	}
	return out
}

func originalQuestionFallbackTokens(text string) []string {
	matches := originalQuestionTokenPattern.FindAllString(text, -1)
	if len(matches) == 0 {
		return nil
	}
	out := make([]string, 0, len(matches))
	for _, match := range matches {
		match = normalizeSearchToken(match)
		if match != "" {
			out = append(out, match)
		}
	}
	return out
}

func normalizeSearchToken(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.Join(strings.Fields(value), " ")
	value = strings.Trim(value, `"'`)
	value = strings.TrimFunc(value, func(r rune) bool {
		return strings.ContainsRune(".,;:!?()[]{}", r)
	})
	return value
}

func containsNonASCII(value string) bool {
	for _, r := range value {
		if r > unicode.MaxASCII {
			return true
		}
	}
	return false
}

var originalQuestionTokenPattern = regexp.MustCompile(`[A-Za-z0-9_./#+-]{2,}|[^\x00-\x7F]+`)

// SearchQuotes runs a ranked FTS5 search using the given keywords.
func (e *Engine) SearchQuotes(ctx context.Context, keywords []string) ([]Quote, error) {
	minRelevance := e.cfg.Search.MinRelevance
	searchLimit := e.cfg.Search.MaxResults
	if minRelevance > 0 {
		searchLimit = max(searchLimit*5, 25)
	}
	slog.Info("engine: searching quotes", "keywords", keywords, "max_results", e.cfg.Search.MaxResults, "min_relevance", minRelevance, "search_limit", searchLimit)
	rows, err := e.store.SearchQuotes(keywords, searchLimit)
	if err != nil {
		slog.Error("engine: search failed", "error", err)
		return nil, err
	}
	quotes := rowsToQuotes(rows, e.localUserID())
	if err := e.enrichQuoteAttachments(quotes); err != nil {
		return nil, err
	}
	if minRelevance > 0 {
		filtered := make([]Quote, 0, len(quotes))
		for _, q := range quotes {
			score := quoteRelevanceScore(keywords, q)
			if score >= minRelevance {
				filtered = append(filtered, q)
			}
		}
		quotes = filtered
	}
	if len(quotes) > e.cfg.Search.MaxResults {
		quotes = quotes[:e.cfg.Search.MaxResults]
	}
	slog.Info("engine: search complete", "result_count", len(quotes))
	for i, q := range quotes {
		slog.Debug("engine: search result", "index", i, "id", q.ID, "content_preview", truncate(q.Content, 80), "tags", q.Tags)
	}
	return quotes, nil
}

func quoteRelevanceScore(keywords []string, quote Quote) float64 {
	normalized := normalizeSearchKeywords(keywords)
	if len(normalized) == 0 {
		return 0
	}

	haystack := strings.ToLower(quote.Content + "\n" + strings.Join(quote.Tags, "\n"))
	best := keywordCoverageScore(normalized, haystack)
	if containsNonASCIIKeyword(normalized) {
		ascii := make([]string, 0, len(normalized))
		nonASCII := make([]string, 0, len(normalized))
		for _, kw := range normalized {
			if containsNonASCII(kw) {
				nonASCII = append(nonASCII, kw)
			} else {
				ascii = append(ascii, kw)
			}
		}
		best = max(best, keywordCoverageScore(ascii, haystack), keywordCoverageScore(nonASCII, haystack))
	}
	return math.Round(best*100) / 100
}

func normalizeSearchKeywords(keywords []string) []string {
	normalized := make([]string, 0, len(keywords))
	seen := make(map[string]struct{}, len(keywords))
	for _, kw := range keywords {
		kw = normalizeSearchToken(kw)
		if kw == "" {
			continue
		}
		if _, ok := seen[kw]; ok {
			continue
		}
		seen[kw] = struct{}{}
		normalized = append(normalized, kw)
	}
	return normalized
}

func keywordCoverageScore(keywords []string, haystack string) float64 {
	if len(keywords) == 0 {
		return 0
	}
	matches := 0
	for _, kw := range keywords {
		if strings.Contains(haystack, kw) {
			matches++
		}
	}
	return float64(matches) / float64(len(keywords))
}

// GenerateResponse streams a synthesized answer grounded in candidate quotes.
// Tokens are sent to tokenCh; the channel is closed when streaming finishes.
// The caller must read tokenCh (or drain it) to avoid goroutine leaks.
func (e *Engine) GenerateResponse(
	ctx context.Context,
	question string,
	candidates []Quote,
	tokenCh chan<- string,
) error {
	slog.Info("engine: generating response", "question", question, "candidate_count", len(candidates))
	if e.isMockLLMEnabled() {
		slog.Info("engine: mock llm enabled for response generation")
		return e.generateMockResponse(candidates, tokenCh)
	}
	var sb strings.Builder
	for i, q := range candidates {
		sb.WriteString(fmt.Sprintf("[%d] %s\n", i+1, q.Content))
	}

	systemPrompt := `You are a retrieval assistant. Your only job is to answer the question using the reference notes below.
Rules:
- Use ONLY information from the reference notes. Do not add any knowledge, opinion, or context of your own.
- Cite each note you use by its number, e.g. [1].
- If the notes do not contain enough information to answer, say exactly: "The reference notes do not contain enough information to answer this question."
- Do not explain, elaborate, or offer suggestions beyond what the notes state.
- Be as brief as possible.

Reference notes:
` + sb.String()

	slog.Debug("engine: response system prompt", "prompt_len", len(systemPrompt))

	msgs := []llm.Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: question},
	}
	zero := 0.0
	_, err := e.llm.Chat(ctx, msgs, tokenCh, llm.ChatOptions{Temperature: &zero})
	if err != nil {
		slog.Error("engine: generate response failed", "error", err)
	}
	return err
}

// --- Provider management ---

// FetchModels retrieves available model IDs from the provider.
func (e *Engine) FetchModels(ctx context.Context, cfg ProviderConfig) ([]string, error) {
	slog.Info("engine: fetching models", "host", cfg.Host, "port", cfg.Port)
	c := llm.NewClient(cfg)
	models, err := c.FetchModels(ctx)
	if err != nil {
		slog.Error("engine: fetch models failed", "error", err)
		return nil, err
	}
	slog.Info("engine: models fetched", "count", len(models))
	return models, err
}

// TestProvider checks connectivity to the given provider config.
func (e *Engine) TestProvider(ctx context.Context, cfg ProviderConfig) error {
	slog.Info("engine: testing provider", "host", cfg.Host, "port", cfg.Port)
	c := llm.NewClient(cfg)
	_, err := c.FetchModels(ctx)
	if err != nil {
		slog.Error("engine: test provider failed", "error", err)
	}
	return err
}

// --- Settings ---

// LoadSettings reads settings from the DB, falling back to defaults.
func (e *Engine) LoadSettings(ctx context.Context) (*Settings, error) {
	slog.Debug("engine: loading settings")
	defaults := DefaultSettings()
	val, err := e.store.GetSetting("settings")
	if err != nil {
		slog.Warn("engine: load settings failed, using defaults", "error", err)
		return defaults, nil
	}
	if val == "" {
		slog.Debug("engine: no saved settings, using defaults")
		return defaults, nil
	}
	var s Settings
	if err := json.Unmarshal([]byte(val), &s); err != nil {
		slog.Error("engine: unmarshal settings failed, using defaults", "error", err, "raw", val)
		return defaults, nil
	}
	if strings.TrimSpace(s.Theme) == "" {
		s.Theme = defaults.Theme
	}
	if s.Web.Port < 1 || s.Web.Port > 65535 {
		s.Web.Port = defaults.Web.Port
	}
	slog.Info(
		"engine: settings loaded",
		"host", s.Provider.Host,
		"port", s.Provider.Port,
		"response_model", s.Provider.ResponseModel(),
		"keyword_model", s.Provider.RecallKeywordModel(),
	)
	return &s, nil
}

// SaveSettings persists settings to the DB and updates the engine.
func (e *Engine) SaveSettings(ctx context.Context, s *Settings) error {
	defaults := DefaultSettings()
	if strings.TrimSpace(s.Theme) == "" {
		s.Theme = defaults.Theme
	}
	if s.Web.Port < 1 || s.Web.Port > 65535 {
		return fmt.Errorf("web port must be a number between 1 and 65535")
	}
	slog.Info(
		"engine: saving settings",
		"host", s.Provider.Host,
		"port", s.Provider.Port,
		"response_model", s.Provider.ResponseModel(),
		"keyword_model", s.Provider.RecallKeywordModel(),
	)
	data, err := json.Marshal(s)
	if err != nil {
		slog.Error("engine: marshal settings failed", "error", err)
		return err
	}
	if err := e.store.SetSetting("settings", string(data)); err != nil {
		slog.Error("engine: persist settings failed", "error", err)
		return err
	}
	e.UpdateSettings(s)
	slog.Info("engine: settings saved and applied")
	return nil
}

func (e *Engine) isMockLLMEnabled() bool {
	return e != nil && e.cfg != nil && e.cfg.Debug.MockLLM
}

func mockSplitKeywords(text string) []string {
	parts := strings.Fields(strings.TrimSpace(text))
	if len(parts) == 0 {
		return nil
	}
	return parts
}

func (e *Engine) generateMockResponse(candidates []Quote, tokenCh chan<- string) error {
	response := strings.TrimSpace(joinQuoteContents(candidates))
	if tokenCh == nil {
		return nil
	}
	defer close(tokenCh)
	if response != "" {
		tokenCh <- response
	}
	return nil
}

func joinQuoteContents(candidates []Quote) string {
	if len(candidates) == 0 {
		return ""
	}
	parts := make([]string, 0, len(candidates))
	for _, q := range candidates {
		content := strings.TrimSpace(q.Content)
		if content != "" {
			parts = append(parts, content)
		}
	}
	return strings.Join(parts, "\n\n")
}

func (e *Engine) LoadUserProfile(ctx context.Context) (*UserProfile, error) {
	slog.Debug("engine: loading user profile")
	row, err := e.store.GetUserProfile()
	if err != nil {
		return nil, err
	}
	if row.UserID == "" {
		profile := &UserProfile{
			UserID:      uuid.NewString(),
			DisplayName: "",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		if err := e.SaveUserProfile(ctx, profile); err != nil {
			return nil, err
		}
		return profile, nil
	}
	profile := &UserProfile{
		UserID:      row.UserID,
		DisplayName: row.DisplayName,
		CreatedAt:   time.Unix(row.CreatedAt, 0),
		UpdatedAt:   time.Unix(row.UpdatedAt, 0),
	}
	e.profile = profile
	return profile, nil
}

func (e *Engine) SaveUserProfile(ctx context.Context, profile *UserProfile) error {
	_ = ctx
	now := time.Now()
	if profile.UserID == "" {
		profile.UserID = uuid.NewString()
	}
	if profile.CreatedAt.IsZero() {
		profile.CreatedAt = now
	}
	profile.DisplayName = strings.TrimSpace(profile.DisplayName)
	profile.UpdatedAt = now
	if err := e.store.SaveUserProfile(db.UserProfileRow{
		UserID:      profile.UserID,
		DisplayName: profile.DisplayName,
		CreatedAt:   profile.CreatedAt.Unix(),
		UpdatedAt:   profile.UpdatedAt.Unix(),
	}); err != nil {
		return err
	}
	if err := e.store.UpdateOwnedQuoteNames(profile.UserID, profile.DisplayName); err != nil {
		return err
	}
	e.profile = profile
	return nil
}

func (e *Engine) BootstrapQuoteIdentity(ctx context.Context) error {
	_ = ctx
	if e.profile == nil {
		return fmt.Errorf("user profile not loaded")
	}
	return e.store.BackfillQuoteIdentity(e.profile.UserID, e.profile.DisplayName, time.Now().Unix(), uuid.NewString)
}

func (e *Engine) SaveRecallHistory(ctx context.Context, question, response string, quotes []Quote) (*RecallHistoryEntry, error) {
	_ = ctx
	question = strings.TrimSpace(question)
	response = strings.TrimSpace(response)
	if question == "" {
		return nil, fmt.Errorf("history question is empty")
	}
	if response == "" {
		return nil, fmt.Errorf("history response is empty")
	}
	quoteIDs := make([]int64, 0, len(quotes))
	for _, q := range quotes {
		if q.ID > 0 {
			quoteIDs = append(quoteIDs, q.ID)
		}
	}
	id, err := e.store.InsertRecallHistory(question, response, quoteIDs)
	if err != nil {
		return nil, err
	}
	return e.GetRecallHistory(context.Background(), id)
}

func (e *Engine) ListRecallHistory(ctx context.Context) ([]RecallHistorySummary, error) {
	return e.ListRecallHistoryPage(ctx, 0, 0)
}

func (e *Engine) CountRecallHistory(ctx context.Context) (int64, error) {
	_ = ctx
	return e.store.CountRecallHistory()
}

func (e *Engine) ListRecallHistoryPage(ctx context.Context, limit, offset int) ([]RecallHistorySummary, error) {
	_ = ctx
	rows, err := e.store.ListRecallHistoryPage(limit, offset)
	if err != nil {
		return nil, err
	}
	out := make([]RecallHistorySummary, len(rows))
	for i, row := range rows {
		out[i] = RecallHistorySummary{
			ID:        row.ID,
			Question:  row.Question,
			Response:  row.Response,
			CreatedAt: time.Unix(row.CreatedAt, 0),
		}
	}
	return out, nil
}

func (e *Engine) GetRecallHistory(ctx context.Context, id int64) (*RecallHistoryEntry, error) {
	_ = ctx
	row, err := e.store.GetRecallHistory(id)
	if err != nil {
		return nil, err
	}
	quotes := rowsToQuotes(row.Quotes, e.localUserID())
	if err := e.enrichQuoteAttachments(quotes); err != nil {
		return nil, err
	}
	return &RecallHistoryEntry{
		ID:        row.ID,
		Question:  row.Question,
		Response:  row.Response,
		CreatedAt: time.Unix(row.CreatedAt, 0),
		Quotes:    quotes,
	}, nil
}

func (e *Engine) DeleteRecallHistory(ctx context.Context, ids []int64) error {
	_ = ctx
	return e.store.DeleteRecallHistory(ids)
}

// --- Helpers ---

func rowsToQuotes(rows []db.QuoteRow, localUserID string) []Quote {
	out := make([]Quote, len(rows))
	for i, r := range rows {
		tags := []string{}
		if r.Tags != "" {
			tags = strings.Split(r.Tags, ",")
		}
		out[i] = Quote{
			ID:               r.ID,
			GlobalID:         r.GlobalID,
			AuthorUserID:     r.AuthorUserID,
			AuthorName:       r.AuthorName,
			SourceUserID:     r.SourceUserID,
			SourceName:       r.SourceName,
			SourceBackend:    r.SourceBackend,
			SourceNamespace:  r.SourceNamespace,
			SourceEntityType: r.SourceEntityType,
			SourceEntityID:   r.SourceEntityID,
			SourceLabel:      r.SourceLabel,
			SourceURL:        r.SourceURL,
			Content:          r.Content,
			Tags:             tags,
			Version:          r.Version,
			IsOwnedByMe:      localUserID != "" && r.AuthorUserID == localUserID,
			CreatedAt:        time.Unix(r.CreatedAt, 0),
			UpdatedAt:        time.Unix(r.UpdatedAt, 0),
		}
	}
	return out
}

func (e *Engine) loadQuote(id int64) (*Quote, error) {
	row, err := e.store.GetQuote(id)
	if err != nil {
		return nil, err
	}
	quotes := rowsToQuotes([]db.QuoteRow{row}, e.localUserID())
	if err := e.enrichQuoteAttachments(quotes); err != nil {
		return nil, err
	}
	return &quotes[0], nil
}

func (e *Engine) loadQuoteByGlobalID(globalID string) (*Quote, error) {
	row, err := e.store.GetQuoteByGlobalID(globalID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("quote %q not found", globalID)
		}
		return nil, err
	}
	quotes := rowsToQuotes([]db.QuoteRow{row}, e.localUserID())
	return &quotes[0], nil
}

func (e *Engine) loadQuoteByIdentifier(id int64, globalID string) (*Quote, error) {
	globalID = strings.TrimSpace(globalID)
	switch {
	case id > 0:
		quote, err := e.loadQuote(id)
		if err != nil {
			return nil, err
		}
		if globalID != "" && quote.GlobalID != globalID {
			return nil, fmt.Errorf("quote id %d does not match global id %q", id, globalID)
		}
		return quote, nil
	case globalID != "":
		return e.loadQuoteByGlobalID(globalID)
	default:
		return nil, fmt.Errorf("quote id or global id is required")
	}
}

func (e *Engine) localUserID() string {
	if e.profile == nil {
		return ""
	}
	return e.profile.UserID
}

func formatRecallQuote(question, response string) string {
	return strings.TrimSpace("Question: " + strings.TrimSpace(question) + "\n\nResponse:\n" + strings.TrimSpace(response))
}

func (e *Engine) quoteIdentityForNewQuote() (db.QuoteIdentity, error) {
	if e.profile == nil || e.profile.UserID == "" {
		return db.QuoteIdentity{}, fmt.Errorf("user profile not loaded")
	}
	globalID := uuid.NewString()
	return db.QuoteIdentity{
		GlobalID:         globalID,
		AuthorUserID:     e.profile.UserID,
		AuthorName:       e.profile.DisplayName,
		SourceUserID:     e.profile.UserID,
		SourceName:       e.profile.DisplayName,
		SourceBackend:    "local",
		SourceNamespace:  "local:" + e.profile.UserID,
		SourceEntityType: "quote",
		SourceEntityID:   globalID,
		SourceLabel:      "Local quote",
		SourceURL:        "",
		Version:          1,
	}, nil
}

// parseJSONStringArray parses a JSON string array, with comma-split fallback.
func parseJSONStringArray(s string) ([]string, error) {
	s = strings.TrimSpace(stripMarkdownCodeFence(s))
	if s == "" {
		return nil, fmt.Errorf("empty model response")
	}
	// Find the first '[' in case the model prefixes with text.
	if i := strings.Index(s, "["); i >= 0 {
		s = s[i:]
	}
	if j := strings.LastIndex(s, "]"); j >= 0 {
		s = s[:j+1]
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, fmt.Errorf("empty model response")
	}
	var tags []string
	if err := json.Unmarshal([]byte(s), &tags); err != nil {
		slog.Debug("engine: JSON unmarshal failed, trying comma-split fallback", "input", s, "error", err)
		// Fallback: split on commas and strip quotes/brackets.
		s = strings.Trim(s, "[]")
		for _, part := range strings.Split(s, ",") {
			part = strings.Trim(strings.TrimSpace(part), `"'`)
			if part != "" {
				tags = append(tags, strings.ToLower(part))
			}
		}
		if len(tags) == 0 {
			return nil, fmt.Errorf("parse keyword array: %w", err)
		}
	}
	return tags, nil
}

type tagNormalizationStats struct {
	OriginalCount  int
	RemovedShort   int
	RemovedGeneric int
	RemovedDupes   int
}

func normalizeTags(tags []string) ([]string, tagNormalizationStats) {
	stats := tagNormalizationStats{OriginalCount: len(tags)}
	seen := make(map[string]struct{}, len(tags))
	normalized := make([]string, 0, len(tags))
	for _, tag := range tags {
		tag = strings.ToLower(strings.TrimSpace(tag))
		tag = strings.Join(strings.Fields(tag), " ")
		tag = strings.Trim(tag, `"'`)
		if len(tag) < 2 {
			stats.RemovedShort++
			continue
		}
		if _, blocked := genericTagBlacklist[tag]; blocked {
			stats.RemovedGeneric++
			continue
		}
		if _, ok := seen[tag]; ok {
			stats.RemovedDupes++
			continue
		}
		seen[tag] = struct{}{}
		normalized = append(normalized, tag)
	}
	if len(normalized) > maxExtractedTags {
		normalized = slices.Clone(normalized[:maxExtractedTags])
	}
	return normalized, stats
}

func shouldRepairTags(tags []string, stats tagNormalizationStats, text string) bool {
	if len(tags) == 0 {
		return true
	}
	textLen := len(strings.TrimSpace(text))
	if textLen >= 80 && len(tags) < 5 {
		return true
	}
	if textLen >= 200 && len(tags) < 8 {
		return true
	}
	if stats.RemovedGeneric >= 2 && len(tags) < 8 {
		return true
	}
	return false
}

func mustMarshalJSONArray(items []string) string {
	data, err := json.Marshal(items)
	if err != nil {
		return "[]"
	}
	return string(data)
}

func sameStringSet(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	left = slices.Clone(left)
	right = slices.Clone(right)
	slices.Sort(left)
	slices.Sort(right)
	return slices.Equal(left, right)
}

func stripMarkdownCodeFence(s string) string {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, "```") {
		return s
	}

	lines := strings.Split(s, "\n")
	if len(lines) == 0 {
		return s
	}
	if strings.HasPrefix(lines[0], "```") {
		lines = lines[1:]
	}
	if len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "```" {
		lines = lines[:len(lines)-1]
	}
	return strings.TrimSpace(strings.Join(lines, "\n"))
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

// ProviderConfig is re-exported from the llm package for engine consumers.
type ProviderConfig = llm.ProviderConfig
