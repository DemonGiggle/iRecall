package core

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/gigol/irecall/core/db"
	"github.com/gigol/irecall/core/llm"
)

// Engine is the central orchestrator. It owns the DB store and LLM client
// and exposes the full recall workflow. No UI types are referenced here.
type Engine struct {
	store *db.Store
	llm   *llm.Client
	cfg   *Settings
}

// New creates an Engine from an open DB store and current settings.
func New(store *db.Store, cfg *Settings) *Engine {
	slog.Info("engine: creating engine", "provider_host", cfg.Provider.Host, "provider_port", cfg.Provider.Port, "model", cfg.Provider.Model)
	return &Engine{
		store: store,
		llm:   llm.NewClient(cfg.Provider),
		cfg:   cfg,
	}
}

func (e *Engine) Close() error {
	slog.Info("engine: closing")
	return e.store.Close()
}

// UpdateProvider rebuilds the LLM client after settings change.
func (e *Engine) UpdateProvider(cfg ProviderConfig) {
	slog.Info("engine: updating provider", "host", cfg.Host, "port", cfg.Port, "model", cfg.Model)
	e.cfg.Provider = cfg
	e.llm = llm.NewClient(cfg)
}

// UpdateSettings replaces the engine's in-memory settings.
func (e *Engine) UpdateSettings(s *Settings) {
	slog.Info("engine: updating settings", "host", s.Provider.Host, "port", s.Provider.Port, "model", s.Provider.Model, "max_results", s.Search.MaxResults)
	e.cfg = s
	e.llm = llm.NewClient(s.Provider)
}

// --- Quote management ---

// AddQuote stores a new quote and extracts tags via LLM.
// If tag extraction fails, the quote is still saved with no tags.
func (e *Engine) AddQuote(ctx context.Context, content string) (*Quote, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, fmt.Errorf("quote content is empty")
	}

	slog.Info("engine: adding quote", "content_len", len(content), "content_preview", truncate(content, 100))

	id, err := e.store.InsertQuote(content)
	if err != nil {
		slog.Error("engine: failed to insert quote", "error", err)
		return nil, fmt.Errorf("store quote: %w", err)
	}
	slog.Info("engine: quote inserted", "id", id)

	slog.Info("engine: extracting tags via LLM", "quote_id", id)
	tags, err := e.ExtractTags(ctx, content)
	if err != nil {
		slog.Error("engine: tag extraction failed, saving without tags", "quote_id", id, "error", err)
		tags = []string{}
	}
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
	return e.loadQuote(id)
}

// ListQuotes returns all quotes, newest first.
func (e *Engine) ListQuotes(ctx context.Context) ([]Quote, error) {
	slog.Debug("engine: listing quotes")
	rows, err := e.store.ListQuotes()
	if err != nil {
		slog.Error("engine: list quotes failed", "error", err)
		return nil, err
	}
	quotes := rowsToQuotes(rows)
	slog.Debug("engine: listed quotes", "count", len(quotes))
	return quotes, nil
}

// DeleteQuote removes a quote by ID.
func (e *Engine) DeleteQuote(ctx context.Context, id int64) error {
	slog.Info("engine: deleting quote", "id", id)
	return e.store.DeleteQuote(id)
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
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, fmt.Errorf("quote content is empty")
	}

	slog.Info("engine: updating quote", "id", id, "content_len", len(content), "content_preview", truncate(content, 100))
	if err := e.store.UpdateQuoteContent(id, content); err != nil {
		return nil, fmt.Errorf("update quote content: %w", err)
	}

	tags, err := e.ExtractTags(ctx, content)
	if err != nil {
		slog.Error("engine: tag extraction failed during update, saving without tags", "quote_id", id, "error", err)
		tags = []string{}
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

	return e.loadQuote(id)
}

// --- Recall workflow ---

// ExtractTags asks the LLM to produce keyword tags for a piece of text.
func (e *Engine) ExtractTags(ctx context.Context, text string) ([]string, error) {
	slog.Debug("engine: extracting tags", "text_len", len(text))
	msgs := []llm.Message{
		{
			Role: "system",
			Content: `You are a JSON keyword extractor. ` +
				`Output ONLY a valid JSON array of 3 to 8 short lowercase keyword strings. ` +
				`No explanation, no markdown, no code fences, no extra text — just the JSON array. ` +
				`Example output: ["emmc", "flash memory", "partition", "offset"]`,
		},
		{
			Role:    "user",
			Content: "Extract keyword tags for this text:\n" + text,
		},
	}
	zero := 0.0
	maxTok := 150
	raw, err := e.llm.Chat(ctx, msgs, nil, llm.ChatOptions{Temperature: &zero, MaxTokens: &maxTok})
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
	slog.Info("engine: extracted tags", "tags", tags)
	return tags, nil
}

// ExtractKeywords asks the LLM to produce search keywords for a question.
func (e *Engine) ExtractKeywords(ctx context.Context, question string) ([]string, error) {
	slog.Info("engine: extracting keywords for question", "question", question)
	msgs := []llm.Message{
		{
			Role: "system",
			Content: `You are a JSON search keyword extractor. ` +
				`Output ONLY a valid JSON array of 3 to 6 short lowercase keyword strings useful for searching a knowledge base. ` +
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
	raw, err := e.llm.Chat(ctx, msgs, nil, llm.ChatOptions{Temperature: &zero, MaxTokens: &maxTok})
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
	slog.Info("engine: extracted keywords", "keywords", keywords)
	return keywords, nil
}

// SearchQuotes runs a ranked FTS5 search using the given keywords.
func (e *Engine) SearchQuotes(ctx context.Context, keywords []string) ([]Quote, error) {
	slog.Info("engine: searching quotes", "keywords", keywords, "max_results", e.cfg.Search.MaxResults)
	rows, err := e.store.SearchQuotes(keywords, e.cfg.Search.MaxResults)
	if err != nil {
		slog.Error("engine: search failed", "error", err)
		return nil, err
	}
	quotes := rowsToQuotes(rows)
	slog.Info("engine: search complete", "result_count", len(quotes))
	for i, q := range quotes {
		slog.Debug("engine: search result", "index", i, "id", q.ID, "content_preview", truncate(q.Content, 80), "tags", q.Tags)
	}
	return quotes, nil
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
	val, err := e.store.GetSetting("settings")
	if err != nil {
		slog.Warn("engine: load settings failed, using defaults", "error", err)
		return DefaultSettings(), nil
	}
	if val == "" {
		slog.Debug("engine: no saved settings, using defaults")
		return DefaultSettings(), nil
	}
	var s Settings
	if err := json.Unmarshal([]byte(val), &s); err != nil {
		slog.Error("engine: unmarshal settings failed, using defaults", "error", err, "raw", val)
		return DefaultSettings(), nil
	}
	slog.Info("engine: settings loaded", "host", s.Provider.Host, "port", s.Provider.Port, "model", s.Provider.Model)
	return &s, nil
}

// SaveSettings persists settings to the DB and updates the engine.
func (e *Engine) SaveSettings(ctx context.Context, s *Settings) error {
	slog.Info("engine: saving settings", "host", s.Provider.Host, "port", s.Provider.Port, "model", s.Provider.Model)
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

// --- Helpers ---

func rowsToQuotes(rows []db.QuoteRow) []Quote {
	out := make([]Quote, len(rows))
	for i, r := range rows {
		tags := []string{}
		if r.Tags != "" {
			tags = strings.Split(r.Tags, ",")
		}
		out[i] = Quote{
			ID:        r.ID,
			Content:   r.Content,
			Tags:      tags,
			CreatedAt: time.Unix(r.CreatedAt, 0),
			UpdatedAt: time.Unix(r.UpdatedAt, 0),
		}
	}
	return out
}

func (e *Engine) loadQuote(id int64) (*Quote, error) {
	row, err := e.store.GetQuote(id)
	if err != nil {
		return nil, err
	}
	quotes := rowsToQuotes([]db.QuoteRow{row})
	return &quotes[0], nil
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
