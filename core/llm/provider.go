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
