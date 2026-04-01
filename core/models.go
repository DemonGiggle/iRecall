package core

import "time"

// Quote is a single user-captured note with associated tags.
type Quote struct {
	ID        int64
	Content   string
	Tags      []string
	CreatedAt time.Time
	UpdatedAt time.Time
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
