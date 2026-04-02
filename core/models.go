package core

import "time"

// Quote is a single user-captured note with associated tags.
type Quote struct {
	ID           int64
	GlobalID     string
	AuthorUserID string
	AuthorName   string
	SourceUserID string
	SourceName   string
	Content      string
	Tags         []string
	Version      int64
	IsOwnedByMe  bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserProfile struct {
	UserID      string
	DisplayName string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

const ShareSchemaVersion = 1

type SharedQuoteEnvelope struct {
	SchemaVersion int                `json:"schema_version"`
	ExportedAt    time.Time          `json:"exported_at"`
	Quotes        []SharedQuoteEntry `json:"quotes"`
}

type SharedQuoteEntry struct {
	GlobalID     string    `json:"global_id"`
	AuthorUserID string    `json:"author_user_id"`
	AuthorName   string    `json:"author_name"`
	SourceUserID string    `json:"source_user_id"`
	SourceName   string    `json:"source_name"`
	Version      int64     `json:"version"`
	Content      string    `json:"content"`
	Tags         []string  `json:"tags"`
	CreatedAtUTC time.Time `json:"created_at_utc"`
	UpdatedAtUTC time.Time `json:"updated_at_utc"`
}

type ImportResult struct {
	Inserted   int
	Updated    int
	Duplicates int
	Stale      int
}

// Settings holds all persisted user preferences.
type Settings struct {
	Provider ProviderConfig
	Search   SearchConfig
}

// SearchConfig controls how candidate quotes are retrieved.
type SearchConfig struct {
	MaxResults   int     // max quotes returned per query (default: 5)
	MinRelevance float64 // FTS rank threshold; 0 = no filter (default: 0.0)
}

func DefaultSettings() *Settings {
	return &Settings{
		Provider: ProviderConfig{
			Host:  "localhost",
			Port:  11434,
			HTTPS: false,
			Model: "",
		},
		Search: SearchConfig{
			MaxResults:   5,
			MinRelevance: 0.0,
		},
	}
}
