package core

import "time"

// Quote is a single user-captured note with associated tags.
type Quote struct {
	ID               int64
	GlobalID         string
	AuthorUserID     string
	AuthorName       string
	SourceUserID     string
	SourceName       string
	SourceBackend    string
	SourceNamespace  string
	SourceEntityType string
	SourceEntityID   string
	SourceLabel      string
	SourceURL        string
	Content          string
	Tags             []string
	Attachments      []QuoteAttachment
	Version          int64
	IsOwnedByMe      bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// QuoteAttachment describes an image managed by iRecall. StoragePath is kept
// internal to the persistence layer and is never exposed through this type.
type QuoteAttachment struct {
	ID        string
	Filename  string
	MediaType string
	Size      int64
	Width     int
	Height    int
	CreatedAt time.Time
}

// ImageInput is a new image supplied by a UI or API caller.
type ImageInput struct {
	Filename  string
	MediaType string
	Data      []byte
}

// QuoteMutationResult includes non-fatal warnings such as vision fallback.
type QuoteMutationResult struct {
	Quote    Quote
	Warnings []string
}

// QuoteKeywordRegeneration captures one explicit keyword refresh for a stored quote.
type QuoteKeywordRegeneration struct {
	QuoteID     int64
	GlobalID    string
	OldKeywords []string
	NewKeywords []string
	Changed     bool
	Quote       Quote
}

type UserProfile struct {
	UserID      string
	DisplayName string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

const ShareSchemaVersion = 2
const BundleSchemaVersion = 3

type SharedQuoteEnvelope struct {
	SchemaVersion int                `json:"schema_version"`
	ExportedAt    time.Time          `json:"exported_at"`
	Quotes        []SharedQuoteEntry `json:"quotes"`
}

type SharedQuoteEntry struct {
	GlobalID         string                  `json:"global_id"`
	AuthorUserID     string                  `json:"author_user_id"`
	AuthorName       string                  `json:"author_name"`
	SourceUserID     string                  `json:"source_user_id"`
	SourceName       string                  `json:"source_name"`
	SourceBackend    string                  `json:"source_backend,omitempty"`
	SourceNamespace  string                  `json:"source_namespace,omitempty"`
	SourceEntityType string                  `json:"source_entity_type,omitempty"`
	SourceEntityID   string                  `json:"source_entity_id,omitempty"`
	SourceLabel      string                  `json:"source_label,omitempty"`
	SourceURL        string                  `json:"source_url,omitempty"`
	Version          int64                   `json:"version"`
	Content          string                  `json:"content"`
	Tags             []string                `json:"tags"`
	Attachments      []SharedAttachmentEntry `json:"attachments,omitempty"`
	CreatedAtUTC     time.Time               `json:"created_at_utc"`
	UpdatedAtUTC     time.Time               `json:"updated_at_utc"`
}

type SharedAttachmentEntry struct {
	ID          string `json:"id"`
	Filename    string `json:"filename"`
	MediaType   string `json:"media_type"`
	Size        int64  `json:"size_bytes"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	SHA256      string `json:"sha256"`
	ArchivePath string `json:"archive_path"`
}

type ImportResult struct {
	Inserted   int
	Updated    int
	Duplicates int
	Stale      int
}

type RecallHistorySummary struct {
	ID        int64
	Question  string
	Response  string
	CreatedAt time.Time
}

type RecallHistoryEntry struct {
	ID        int64
	Question  string
	Response  string
	CreatedAt time.Time
	Quotes    []Quote
}

// Settings holds all persisted user preferences.
type Settings struct {
	Provider ProviderConfig
	Search   SearchConfig
	Debug    DebugConfig
	Theme    string
	Web      WebConfig
	RootDir  string
}

type DebugConfig struct {
	MockLLM bool
}

type WebConfig struct {
	Port int
}

// SearchConfig controls how candidate quotes are retrieved.
type SearchConfig struct {
	MaxResults   int     // max quotes returned per query (default: 5)
	MinRelevance float64 // normalized keyword-match threshold in [0,1]; 0 = no filter
}

func DefaultSettings() *Settings {
	return &Settings{
		Provider: ProviderConfig{
			Host:         "localhost",
			Port:         11434,
			HTTPS:        false,
			Model:        "",
			KeywordModel: "",
		},
		Search: SearchConfig{
			MaxResults:   5,
			MinRelevance: 0.0,
		},
		Debug: DebugConfig{
			MockLLM: false,
		},
		Theme: "violet",
		Web: WebConfig{
			Port: 9527,
		},
	}
}
