package mcpbridge

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gigol/irecall/mcp/irecallapi"
)

const (
	DefaultBaseURL = "http://127.0.0.1:9527"
	EnvBaseURL     = "IRECALL_BASE_URL"
	EnvAPIToken    = "IRECALL_API_TOKEN"
)

type Config struct {
	BaseURL     string
	APIToken    string
	HTTPTimeout time.Duration
}

func LoadConfig(baseURLOverride string, timeout time.Duration) (Config, error) {
	baseURL := strings.TrimSpace(baseURLOverride)
	if baseURL == "" {
		baseURL = strings.TrimSpace(os.Getenv(EnvBaseURL))
	}
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return Config{}, fmt.Errorf("parse %s: %w", EnvBaseURL, err)
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return Config{}, errors.New("iRecall base URL must include scheme and host")
	}

	token := strings.TrimSpace(os.Getenv(EnvAPIToken))
	if token == "" {
		return Config{}, fmt.Errorf("%s is required", EnvAPIToken)
	}
	if timeout <= 0 {
		timeout = 15 * time.Second
	}

	return Config{
		BaseURL:     strings.TrimRight(parsed.String(), "/"),
		APIToken:    token,
		HTTPTimeout: timeout,
	}, nil
}

func (c Config) APIConfig() irecallapi.Config {
	return irecallapi.Config{
		BaseURL:     c.BaseURL,
		APIToken:    c.APIToken,
		HTTPTimeout: c.HTTPTimeout,
	}
}
