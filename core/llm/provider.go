package llm

import "fmt"

// ProviderConfig describes a single OpenAI-compatible API endpoint.
type ProviderConfig struct {
	Host   string // hostname or IP, no scheme
	Port   int
	HTTPS  bool
	APIKey string // empty = no Authorization header
	Model  string
}

func (p ProviderConfig) BaseURL() string {
	scheme := "http"
	if p.HTTPS {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s:%d/v1", scheme, p.Host, p.Port)
}
