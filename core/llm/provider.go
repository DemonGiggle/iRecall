package llm

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
)

// ProviderConfig describes a single OpenAI-compatible API endpoint.
type ProviderConfig struct {
	Host         string // hostname or IP, no scheme
	Port         int
	HTTPS        bool
	APIKey       string // empty = no Authorization header
	Model        string // response/default model
	KeywordModel string // optional override for recall keyword extraction
}

func (p ProviderConfig) BaseURL() string {
	scheme := "http"
	if p.HTTPS {
		scheme = "https"
	}

	raw := strings.TrimSpace(p.Host)
	if raw == "" {
		return fmt.Sprintf("%s://:%d/v1", scheme, p.Port)
	}
	if !strings.Contains(raw, "://") {
		raw = scheme + "://" + raw
	}

	u, err := url.Parse(raw)
	if err != nil || u.Host == "" {
		return fmt.Sprintf("%s://%s:%d/v1", scheme, p.Host, p.Port)
	}

	u.Scheme = scheme
	host := u.Hostname()
	if host == "" {
		host = u.Host
	}
	if p.Port > 0 {
		u.Host = net.JoinHostPort(host, strconv.Itoa(p.Port))
	} else if u.Port() == "" {
		u.Host = host
	}

	basePath := strings.TrimRight(u.Path, "/")
	switch {
	case basePath == "":
		u.Path = "/v1"
	case strings.HasSuffix(basePath, "/v1"):
		u.Path = basePath
	default:
		u.Path = basePath + "/v1"
	}
	u.RawPath = ""
	return u.String()
}

func (p ProviderConfig) ResponseModel() string {
	return strings.TrimSpace(p.Model)
}

func (p ProviderConfig) RecallKeywordModel() string {
	if model := strings.TrimSpace(p.KeywordModel); model != "" {
		return model
	}
	return p.ResponseModel()
}

func (p ProviderConfig) ResponseProviderConfig() ProviderConfig {
	p.Model = p.ResponseModel()
	return p
}

func (p ProviderConfig) RecallKeywordProviderConfig() ProviderConfig {
	p.Model = p.RecallKeywordModel()
	return p
}
