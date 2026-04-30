package mcp

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/gigol/irecall/mcp/irecallapi"
)

const (
	DefaultBaseURL     = "http://127.0.0.1:9527"
	EnvBaseURL         = "IRECALL_BASE_URL"
	EnvAPITokenFile    = "IRECALL_API_TOKEN_FILE"
	EnvAPIToken        = "IRECALL_API_TOKEN"
	maxSecretFileBytes = 4096
)

type Config struct {
	BaseURL     string
	APIToken    string
	HTTPTimeout time.Duration
}

func LoadConfig(baseURLOverride string, tokenFileOverride string, timeout time.Duration) (Config, error) {
	baseURL := strings.TrimSpace(baseURLOverride)
	if baseURL == "" {
		baseURL = strings.TrimSpace(os.Getenv(EnvBaseURL))
	}
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return Config{}, fmt.Errorf("parse base URL: %w", err)
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return Config{}, errors.New("iRecall base URL must include scheme and host")
	}

	token, err := resolveAPIToken(tokenFileOverride)
	if err != nil {
		return Config{}, err
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

func resolveAPIToken(tokenFileOverride string) (string, error) {
	tokenFile := strings.TrimSpace(tokenFileOverride)
	if tokenFile == "" {
		tokenFile = strings.TrimSpace(os.Getenv(EnvAPITokenFile))
	}
	if tokenFile != "" {
		return readSecretFile(tokenFile, "API token")
	}

	token := strings.TrimSpace(os.Getenv(EnvAPIToken))
	if token == "" {
		return "", errors.New("API token is required via --token-file, IRECALL_API_TOKEN_FILE, or IRECALL_API_TOKEN")
	}
	return token, nil
}

func readSecretFile(path string, label string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("open %s file %q: %w", label, path, err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("stat %s file %q: %w", label, path, err)
	}
	if !info.Mode().IsRegular() {
		return "", fmt.Errorf("%s file %q must be a regular file", label, path)
	}
	if runtime.GOOS != "windows" {
		if perm := info.Mode().Perm(); perm&0o037 != 0 {
			return "", fmt.Errorf("%s file %q must be owner-only or owner+group-read only (mode %03o)", label, path, perm)
		}
	}

	data, err := io.ReadAll(io.LimitReader(file, maxSecretFileBytes+1))
	if err != nil {
		return "", fmt.Errorf("read %s file %q: %w", label, path, err)
	}
	if len(data) > maxSecretFileBytes {
		return "", fmt.Errorf("%s file %q exceeds %d bytes", label, path, maxSecretFileBytes)
	}

	value := strings.TrimSpace(string(data))
	if value == "" {
		return "", fmt.Errorf("%s file %q is empty", label, path)
	}
	return value, nil
}

func (c Config) APIConfig() irecallapi.Config {
	return irecallapi.Config{
		BaseURL:     c.BaseURL,
		APIToken:    c.APIToken,
		HTTPTimeout: c.HTTPTimeout,
	}
}
