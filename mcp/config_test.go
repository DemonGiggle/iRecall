package mcp

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestLoadConfigResolvesTokenSourcesInPriorityOrder(t *testing.T) {
	flagTokenPath := writeTokenFileForTest(t, "flag-token\n", 0o600)
	envTokenFilePath := writeTokenFileForTest(t, " env-file-token \n", 0o600)

	t.Setenv(EnvBaseURL, "http://env.example:9999/")
	t.Setenv(EnvAPITokenFile, envTokenFilePath)
	t.Setenv(EnvAPIToken, "env-token")

	cfg, err := LoadConfig("http://override.example:9527/", flagTokenPath, 0)
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}
	if cfg.BaseURL != "http://override.example:9527" {
		t.Fatalf("BaseURL = %q, want override without trailing slash", cfg.BaseURL)
	}
	if cfg.APIToken != "flag-token" {
		t.Fatalf("APIToken = %q, want flag token", cfg.APIToken)
	}
	if cfg.HTTPTimeout != 15*time.Second {
		t.Fatalf("HTTPTimeout = %v, want default", cfg.HTTPTimeout)
	}

	cfg, err = LoadConfig("", "", 5*time.Second)
	if err != nil {
		t.Fatalf("LoadConfig() with env file error = %v", err)
	}
	if cfg.BaseURL != "http://env.example:9999" {
		t.Fatalf("BaseURL = %q, want env base URL without trailing slash", cfg.BaseURL)
	}
	if cfg.APIToken != "env-file-token" {
		t.Fatalf("APIToken = %q, want env file token", cfg.APIToken)
	}
	if cfg.HTTPTimeout != 5*time.Second {
		t.Fatalf("HTTPTimeout = %v, want explicit timeout", cfg.HTTPTimeout)
	}

	t.Setenv(EnvAPITokenFile, "")
	cfg, err = LoadConfig("", "", 5*time.Second)
	if err != nil {
		t.Fatalf("LoadConfig() with env token fallback error = %v", err)
	}
	if cfg.APIToken != "env-token" {
		t.Fatalf("APIToken = %q, want env token fallback", cfg.APIToken)
	}
}

func TestLoadConfigRejectsEmptyOrMissingTokenSources(t *testing.T) {
	t.Setenv(EnvAPITokenFile, "")
	t.Setenv(EnvAPIToken, "")

	_, err := LoadConfig("", "", 5*time.Second)
	if err == nil {
		t.Fatal("LoadConfig() error = nil, want missing token failure")
	}
	if !strings.Contains(err.Error(), "--token-file") {
		t.Fatalf("LoadConfig() error = %q, want token source guidance", err)
	}

	emptyPath := writeTokenFileForTest(t, " \n", 0o600)
	_, err = LoadConfig("", emptyPath, 5*time.Second)
	if err == nil {
		t.Fatal("LoadConfig() error = nil, want empty token file failure")
	}
	if !strings.Contains(err.Error(), "is empty") {
		t.Fatalf("LoadConfig() error = %q, want empty token file guidance", err)
	}
}

func TestLoadConfigRejectsInsecureTokenFilePermissions(t *testing.T) {
	t.Parallel()

	if runtime.GOOS == "windows" {
		t.Skip("permission bits are not enforced on Windows")
	}

	allowedPath := writeTokenFileForTest(t, "flag-token\n", 0o640)
	cfg, err := LoadConfig("", allowedPath, 5*time.Second)
	if err != nil {
		t.Fatalf("LoadConfig() with 0640 token file error = %v, want success", err)
	}
	if cfg.APIToken != "flag-token" {
		t.Fatalf("APIToken = %q, want trimmed token from 0640 file", cfg.APIToken)
	}

	insecurePath := writeTokenFileForTest(t, "flag-token\n", 0o660)
	_, err = LoadConfig("", insecurePath, 5*time.Second)
	if err == nil {
		t.Fatal("LoadConfig() error = nil, want insecure permissions failure")
	}
	if !strings.Contains(err.Error(), "owner+group-read only") {
		t.Fatalf("LoadConfig() error = %q, want permission guidance", err)
	}
}

func writeTokenFileForTest(t *testing.T, contents string, mode os.FileMode) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "api-token")
	if err := os.WriteFile(path, []byte(contents), mode); err != nil {
		t.Fatalf("WriteFile(%q) error = %v", path, err)
	}
	if err := os.Chmod(path, mode); err != nil {
		t.Fatalf("Chmod(%q) error = %v", path, err)
	}
	return path
}
