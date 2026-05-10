//go:build !wails

package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"

	irecallapp "github.com/gigol/irecall/app"
	"github.com/gigol/irecall/core"
)

func TestOptionalBoolFlagBehaviors(t *testing.T) {
	t.Parallel()

	var flagValue optionalBoolFlag
	if flagValue.String() != "" {
		t.Fatalf("unset String() = %q, want empty", flagValue.String())
	}
	if !flagValue.IsBoolFlag() {
		t.Fatal("IsBoolFlag() = false, want true")
	}
	if err := flagValue.Set("true"); err != nil {
		t.Fatalf("Set(true) error = %v", err)
	}
	if !flagValue.set || !flagValue.value {
		t.Fatalf("flag state = %+v, want set=true value=true", flagValue)
	}
	if got := flagValue.String(); got != "true" {
		t.Fatalf("String() = %q, want true", got)
	}
	if err := flagValue.Set("not-bool"); err == nil {
		t.Fatal("Set(invalid) error = nil, want parse failure")
	}
}

func TestProviderStartupOptionsHasOverrides(t *testing.T) {
	t.Parallel()

	if (ProviderStartupOptions{}).HasOverrides() {
		t.Fatal("HasOverrides() = true for empty options, want false")
	}

	opts := ProviderStartupOptions{Model: "runtime-model"}
	if !opts.HasOverrides() {
		t.Fatal("HasOverrides() = false with model override, want true")
	}
}

func TestProviderStartupOptionsValidateRequiresAPIOnly(t *testing.T) {
	t.Parallel()

	opts := ProviderStartupOptions{Host: "provider.example"}
	err := opts.Validate(ServerOptions{})
	if err == nil {
		t.Fatal("Validate() error = nil, want api-only requirement")
	}
	if !strings.Contains(err.Error(), "--api-only") {
		t.Fatalf("Validate() error = %q, want --api-only guidance", err)
	}
}

func TestProviderStartupOptionsValidateRejectsInvalidPort(t *testing.T) {
	t.Parallel()

	err := (ProviderStartupOptions{Port: 70000}).Validate(ServerOptions{APIOnly: true})
	if err == nil {
		t.Fatal("Validate() error = nil, want invalid port error")
	}
	if !strings.Contains(err.Error(), "provider port") {
		t.Fatalf("Validate() error = %q, want provider port guidance", err)
	}
}

func TestProviderStartupOptionsResolveReadsAPIKeyFile(t *testing.T) {
	t.Parallel()

	apiKeyPath := filepath.Join(t.TempDir(), "provider-api-key")
	if err := os.WriteFile(apiKeyPath, []byte("  runtime-secret\n"), 0o600); err != nil {
		t.Fatalf("WriteFile(apiKeyPath) error = %v", err)
	}

	opts := ProviderStartupOptions{
		Host:         "provider.example/api",
		Port:         443,
		APIKeyPath:   apiKeyPath,
		Model:        "runtime-model",
		KeywordModel: "runtime-keyword-model",
	}
	if err := opts.HTTPS.Set("true"); err != nil {
		t.Fatalf("HTTPS.Set(true) error = %v", err)
	}

	resolved, err := opts.Resolve(core.ProviderConfig{
		Host:         "saved-provider",
		Port:         11434,
		HTTPS:        false,
		APIKey:       "saved-secret",
		Model:        "saved-model",
		KeywordModel: "saved-keyword-model",
	})
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	if resolved.Host != "provider.example/api" {
		t.Fatalf("resolved host = %q, want override", resolved.Host)
	}
	if resolved.Port != 443 {
		t.Fatalf("resolved port = %d, want 443", resolved.Port)
	}
	if !resolved.HTTPS {
		t.Fatalf("resolved HTTPS = false, want true")
	}
	if resolved.APIKey != "runtime-secret" {
		t.Fatalf("resolved API key = %q, want trimmed secret", resolved.APIKey)
	}
	if resolved.Model != "runtime-model" {
		t.Fatalf("resolved model = %q, want override", resolved.Model)
	}
	if resolved.KeywordModel != "runtime-keyword-model" {
		t.Fatalf("resolved keyword model = %q, want override", resolved.KeywordModel)
	}
}

func TestProviderStartupOptionsResolveRejectsBadAPIKeyFiles(t *testing.T) {
	t.Parallel()

	emptyPath := filepath.Join(t.TempDir(), "empty-key")
	if err := os.WriteFile(emptyPath, []byte(" \n"), 0o600); err != nil {
		t.Fatalf("WriteFile(emptyPath) error = %v", err)
	}
	tooLargePath := filepath.Join(t.TempDir(), "too-large-key")
	if err := os.WriteFile(tooLargePath, []byte(strings.Repeat("a", maxSecretFileBytes+1)), 0o600); err != nil {
		t.Fatalf("WriteFile(tooLargePath) error = %v", err)
	}

	tests := []struct {
		name    string
		path    string
		wantErr string
	}{
		{
			name:    "missing file",
			path:    filepath.Join(t.TempDir(), "missing-key"),
			wantErr: "open provider API key file",
		},
		{
			name:    "empty file",
			path:    emptyPath,
			wantErr: "is empty",
		},
		{
			name:    "too large file",
			path:    tooLargePath,
			wantErr: "exceeds",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := (ProviderStartupOptions{APIKeyPath: tt.path}).Resolve(core.ProviderConfig{})
			if err == nil {
				t.Fatal("Resolve() error = nil, want failure")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("Resolve() error = %q, want substring %q", err, tt.wantErr)
			}
		})
	}
}

func TestApplyAPIOnlyProviderStartupConfigNoOpCases(t *testing.T) {
	t.Parallel()

	runtimeApp := newTestApp(t)
	original := *runtimeApp.GetSettings()

	if err := applyAPIOnlyProviderStartupConfig(runtimeApp, ServerOptions{}, ProviderStartupOptions{}); err != nil {
		t.Fatalf("applyAPIOnlyProviderStartupConfig(non api-only empty opts) error = %v", err)
	}
	if !reflect.DeepEqual(*runtimeApp.GetSettings(), original) {
		t.Fatalf("settings changed for non api-only no-op: got %+v want %+v", *runtimeApp.GetSettings(), original)
	}

	if err := applyAPIOnlyProviderStartupConfig(runtimeApp, ServerOptions{APIOnly: true}, ProviderStartupOptions{}); err != nil {
		t.Fatalf("applyAPIOnlyProviderStartupConfig(no overrides) error = %v", err)
	}
	if !reflect.DeepEqual(*runtimeApp.GetSettings(), original) {
		t.Fatalf("settings changed for empty overrides no-op: got %+v want %+v", *runtimeApp.GetSettings(), original)
	}
}

func TestApplyAPIOnlyProviderStartupConfigRequiresInitializedApp(t *testing.T) {
	t.Parallel()

	err := applyAPIOnlyProviderStartupConfig(nil, ServerOptions{APIOnly: true}, ProviderStartupOptions{Host: "provider.example"})
	if err == nil || !strings.Contains(err.Error(), "app is not initialized") {
		t.Fatalf("applyAPIOnlyProviderStartupConfig(nil) error = %v, want initialized app failure", err)
	}
}

func TestApplyAPIOnlyProviderStartupConfigOverridesRuntimeWithoutPersisting(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(t.TempDir(), "xdg-config"))

	root := filepath.Join(t.TempDir(), "runtime-provider")
	runtimeApp, err := irecallapp.NewAppWithOptions(root, irecallapp.AppOptions{})
	if err != nil {
		t.Fatalf("NewAppWithOptions() error = %v", err)
	}

	savedSettings := *runtimeApp.GetSettings()
	savedSettings.Provider = core.ProviderConfig{
		Host:         "saved-provider",
		Port:         11434,
		HTTPS:        false,
		APIKey:       "saved-secret",
		Model:        "saved-model",
		KeywordModel: "saved-keyword-model",
	}
	if _, err := runtimeApp.SaveSettings(savedSettings); err != nil {
		t.Fatalf("SaveSettings() error = %v", err)
	}

	apiKeyPath := filepath.Join(t.TempDir(), "provider-api-key")
	if err := os.WriteFile(apiKeyPath, []byte("runtime-secret\n"), 0o600); err != nil {
		t.Fatalf("WriteFile(apiKeyPath) error = %v", err)
	}
	opts := ProviderStartupOptions{
		Host:         "runtime-provider",
		Port:         8443,
		APIKeyPath:   apiKeyPath,
		Model:        "runtime-model",
		KeywordModel: "runtime-keyword-model",
	}
	if err := opts.HTTPS.Set("true"); err != nil {
		t.Fatalf("HTTPS.Set(true) error = %v", err)
	}

	if err := applyAPIOnlyProviderStartupConfig(runtimeApp, ServerOptions{APIOnly: true}, opts); err != nil {
		t.Fatalf("applyAPIOnlyProviderStartupConfig() error = %v", err)
	}
	current := runtimeApp.GetSettings()
	if current == nil {
		t.Fatal("GetSettings() after apply returned nil")
	}
	if current.Provider.Host != "runtime-provider" ||
		current.Provider.Port != 8443 ||
		!current.Provider.HTTPS ||
		current.Provider.APIKey != "runtime-secret" ||
		current.Provider.Model != "runtime-model" ||
		current.Provider.KeywordModel != "runtime-keyword-model" {
		t.Fatalf("runtime provider = %+v, want applied overrides", current.Provider)
	}

	runtimeApp.Shutdown(context.Background())

	reopened, err := irecallapp.NewAppWithOptions(root, irecallapp.AppOptions{})
	if err != nil {
		t.Fatalf("NewAppWithOptions(reopen) error = %v", err)
	}
	t.Cleanup(func() { reopened.Shutdown(context.Background()) })

	reloaded := reopened.GetSettings()
	if reloaded == nil {
		t.Fatal("GetSettings() after reopen returned nil")
	}
	if reloaded.Provider.Host != "saved-provider" ||
		reloaded.Provider.Port != 11434 ||
		reloaded.Provider.HTTPS ||
		reloaded.Provider.APIKey != "saved-secret" ||
		reloaded.Provider.Model != "saved-model" ||
		reloaded.Provider.KeywordModel != "saved-keyword-model" {
		t.Fatalf("reloaded provider = %+v, want persisted settings without runtime override", reloaded.Provider)
	}
}

func TestAPIOnlyProviderStartupConfigEnablesLLMRoutes(t *testing.T) {
	t.Parallel()

	type observedRequest struct {
		Path            string
		Auth            string
		Model           string
		ReasoningEffort string
	}
	observedCh := make(chan observedRequest, 1)
	provider := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Model           string `json:"model"`
			ReasoningEffort string `json:"reasoning_effort"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		observedCh <- observedRequest{
			Path:            r.URL.Path,
			Auth:            r.Header.Get("Authorization"),
			Model:           req.Model,
			ReasoningEffort: req.ReasoningEffort,
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"choices": []map[string]any{
				{
					"message": map[string]any{
						"content": "refined output",
					},
				},
			},
		})
	}))
	defer provider.Close()

	host, port := providerHostPort(t, provider.URL)
	runtimeApp := newTestApp(t)
	tokenResult, err := runtimeApp.CreateAPIToken()
	if err != nil {
		t.Fatalf("CreateAPIToken() error = %v", err)
	}

	apiKeyPath := filepath.Join(t.TempDir(), "provider-api-key")
	if err := os.WriteFile(apiKeyPath, []byte("headless-secret\n"), 0o600); err != nil {
		t.Fatalf("WriteFile(apiKeyPath) error = %v", err)
	}
	if err := applyAPIOnlyProviderStartupConfig(runtimeApp, ServerOptions{APIOnly: true}, ProviderStartupOptions{
		Host:       host,
		Port:       port,
		APIKeyPath: apiKeyPath,
		Model:      "runtime-model",
	}); err != nil {
		t.Fatalf("applyAPIOnlyProviderStartupConfig() error = %v", err)
	}

	server := newTestServerWithOptions(t, runtimeApp, ServerOptions{APIOnly: true})
	req := httptest.NewRequest(http.MethodPost, "/api/app/refine-quote-draft", jsonBody(t, map[string]string{
		"content": "draft note",
	}))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenResult.Token)
	res := httptest.NewRecorder()
	server.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("POST /api/app/refine-quote-draft = %d, body = %s", res.Code, res.Body.String())
	}
	var refined string
	if err := json.Unmarshal(res.Body.Bytes(), &refined); err != nil {
		t.Fatalf("decode refine response: %v", err)
	}
	if refined != "refined output" {
		t.Fatalf("refined response = %q, want provider output", refined)
	}

	observed := <-observedCh
	if observed.Path != "/v1/chat/completions" {
		t.Fatalf("provider path = %q, want /v1/chat/completions", observed.Path)
	}
	if observed.Auth != "Bearer headless-secret" {
		t.Fatalf("provider auth = %q, want bearer secret", observed.Auth)
	}
	if observed.Model != "runtime-model" {
		t.Fatalf("provider model = %q, want runtime-model", observed.Model)
	}
	if observed.ReasoningEffort != "none" {
		t.Fatalf("provider reasoning_effort = %q, want none", observed.ReasoningEffort)
	}
}

func TestReadSecretFileSuccess(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "provider-secret")
	if err := os.WriteFile(path, []byte("  top-secret \n"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	value, err := readSecretFile(path, "provider API key")
	if err != nil {
		t.Fatalf("readSecretFile() error = %v", err)
	}
	if value != "top-secret" {
		t.Fatalf("readSecretFile() = %q, want trimmed secret", value)
	}
}

func providerHostPort(t *testing.T, rawURL string) (string, int) {
	t.Helper()

	parsed, err := url.Parse(rawURL)
	if err != nil {
		t.Fatalf("url.Parse(%q) error = %v", rawURL, err)
	}
	port, err := strconv.Atoi(parsed.Port())
	if err != nil {
		t.Fatalf("Atoi(%q) error = %v", parsed.Port(), err)
	}
	return parsed.Hostname(), port
}
