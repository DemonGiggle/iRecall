package core

import (
	"context"
	"encoding/json"
	"fmt"
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
	return &Engine{
		store: store,
		llm:   llm.NewClient(cfg.Provider),
		cfg:   cfg,
	}
}

func (e *Engine) Close() error {
	return e.store.Close()
}

// UpdateProvider rebuilds the LLM client after settings change.
func (e *Engine) UpdateProvider(cfg ProviderConfig) {
	e.cfg.Provider = cfg
	e.llm = llm.NewClient(cfg)
}

// UpdateSettings replaces the engine's in-memory settings.
func (e *Engine) UpdateSettings(s *Settings) {
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

	id, err := e.store.InsertQuote(content)
	if err != nil {
		return nil, fmt.Errorf("store quote: %w", err)
	}

	tags, err := e.ExtractTags(ctx, content)
	if err != nil {
		// Best-effort: log and continue without tags.
		tags = []string{}
	}

	if len(tags) > 0 {
		tagIDs, err := e.store.UpsertTags(tags)
		if err == nil {
			_ = e.store.InsertQuoteTags(id, tagIDs)
			_ = e.store.UpdateQuoteFTS(id, tags)
		}
	}

	return &Quote{
		ID:        id,
		Content:   content,
		Tags:      tags,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// ListQuotes returns all quotes, newest first.
func (e *Engine) ListQuotes(ctx context.Context) ([]Quote, error) {
	rows, err := e.store.ListQuotes()
	if err != nil {
		return nil, err
	}
	return rowsToQuotes(rows), nil
}

// DeleteQuote removes a quote by ID.
func (e *Engine) DeleteQuote(ctx context.Context, id int64) error {
	return e.store.DeleteQuote(id)
}

// --- Recall workflow ---

// ExtractTags asks the LLM to produce keyword tags for a piece of text.
func (e *Engine) ExtractTags(ctx context.Context, text string) ([]string, error) {
	msgs := []llm.Message{
		{
			Role: "system",
			Content: `You are a keyword extractor. Given a piece of text, return a JSON array of ` +
				`3 to 8 short, lowercase keyword tags that best represent the core concepts. ` +
				`Return ONLY the JSON array with no explanation. ` +
				`Example: ["machine learning", "neural networks", "backpropagation"]`,
		},
		{Role: "user", Content: text},
	}
	raw, err := e.llm.Chat(ctx, msgs, nil)
	if err != nil {
		return nil, err
	}
	return parseJSONStringArray(raw)
}

// ExtractKeywords asks the LLM to produce search keywords for a question.
func (e *Engine) ExtractKeywords(ctx context.Context, question string) ([]string, error) {
	msgs := []llm.Message{
		{
			Role: "system",
			Content: `You are a search keyword extractor. Given a question, return a JSON array of ` +
				`3 to 6 lowercase keywords or short phrases most useful for searching a personal ` +
				`knowledge base. Return ONLY the JSON array.`,
		},
		{Role: "user", Content: question},
	}
	raw, err := e.llm.Chat(ctx, msgs, nil)
	if err != nil {
		return nil, err
	}
	return parseJSONStringArray(raw)
}

// SearchQuotes runs a ranked FTS5 search using the given keywords.
func (e *Engine) SearchQuotes(ctx context.Context, keywords []string) ([]Quote, error) {
	rows, err := e.store.SearchQuotes(keywords, e.cfg.Search.MaxResults)
	if err != nil {
		return nil, err
	}
	return rowsToQuotes(rows), nil
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
	var sb strings.Builder
	for i, q := range candidates {
		sb.WriteString(fmt.Sprintf("[%d] %s\n", i+1, q.Content))
	}

	systemPrompt := `You are a personal knowledge assistant. Answer the user's question using ONLY ` +
		`the reference notes provided below. Cite the notes by their number, e.g. [1]. ` +
		`If the notes do not contain enough information, say so clearly. Be concise and direct.` +
		"\n\nReference notes:\n" + sb.String()

	msgs := []llm.Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: question},
	}
	_, err := e.llm.Chat(ctx, msgs, tokenCh)
	return err
}

// --- Provider management ---

// FetchModels retrieves available model IDs from the provider.
func (e *Engine) FetchModels(ctx context.Context, cfg ProviderConfig) ([]string, error) {
	c := llm.NewClient(cfg)
	return c.FetchModels(ctx)
}

// TestProvider checks connectivity to the given provider config.
func (e *Engine) TestProvider(ctx context.Context, cfg ProviderConfig) error {
	c := llm.NewClient(cfg)
	_, err := c.FetchModels(ctx)
	return err
}

// --- Settings ---

// LoadSettings reads settings from the DB, falling back to defaults.
func (e *Engine) LoadSettings(ctx context.Context) (*Settings, error) {
	val, err := e.store.GetSetting("settings")
	if err != nil {
		return DefaultSettings(), nil
	}
	if val == "" {
		return DefaultSettings(), nil
	}
	var s Settings
	if err := json.Unmarshal([]byte(val), &s); err != nil {
		return DefaultSettings(), nil
	}
	return &s, nil
}

// SaveSettings persists settings to the DB and updates the engine.
func (e *Engine) SaveSettings(ctx context.Context, s *Settings) error {
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}
	if err := e.store.SetSetting("settings", string(data)); err != nil {
		return err
	}
	e.UpdateSettings(s)
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

// parseJSONStringArray parses a JSON string array, with comma-split fallback.
func parseJSONStringArray(s string) ([]string, error) {
	s = strings.TrimSpace(s)
	// Find the first '[' in case the model prefixes with text.
	if i := strings.Index(s, "["); i >= 0 {
		s = s[i:]
	}
	if j := strings.LastIndex(s, "]"); j >= 0 {
		s = s[:j+1]
	}
	var tags []string
	if err := json.Unmarshal([]byte(s), &tags); err != nil {
		// Fallback: split on commas and strip quotes/brackets.
		s = strings.Trim(s, "[]")
		for _, part := range strings.Split(s, ",") {
			part = strings.Trim(strings.TrimSpace(part), `"'`)
			if part != "" {
				tags = append(tags, strings.ToLower(part))
			}
		}
	}
	return tags, nil
}

// ProviderConfig is re-exported from the llm package for engine consumers.
type ProviderConfig = llm.ProviderConfig
