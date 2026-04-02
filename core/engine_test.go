package core

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/gigol/irecall/core/db"
)

func TestParseJSONStringArray(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    []string
		wantErr bool
	}{
		{
			name:  "plain json array",
			input: `["emmc","flash memory","partition"]`,
			want:  []string{"emmc", "flash memory", "partition"},
		},
		{
			name:  "markdown fenced json array",
			input: "```json\n[\"emmc\", \"flash memory\"]\n```",
			want:  []string{"emmc", "flash memory"},
		},
		{
			name:  "extra prose before array",
			input: "Here you go: [\"alpha\", \"beta\"]",
			want:  []string{"alpha", "beta"},
		},
		{
			name:  "comma fallback",
			input: `"Alpha", beta, gamma`,
			want:  []string{"alpha", "beta", "gamma"},
		},
		{
			name:    "empty response",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := parseJSONStringArray(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("parseJSONStringArray() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestRefineQuoteDraft(t *testing.T) {
	t.Parallel()

	var gotRequest struct {
		Model    string `json:"model"`
		Messages []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"messages"`
		Stream bool `json:"stream"`
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/chat/completions" {
			t.Fatalf("path = %q, want /v1/chat/completions", r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&gotRequest); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"choices": []map[string]any{
				{
					"message": map[string]string{
						"content": "A clearer version of the original note.",
					},
				},
			},
		})
	}))
	defer srv.Close()

	engine := newTestEngine(t, srv.Listener.Addr().String())

	refined, err := engine.RefineQuoteDraft(context.Background(), "messy draft note")
	if err != nil {
		t.Fatalf("RefineQuoteDraft() error = %v", err)
	}
	if refined != "A clearer version of the original note." {
		t.Fatalf("RefineQuoteDraft() = %q", refined)
	}
	if gotRequest.Model != "test-model" {
		t.Fatalf("model = %q, want test-model", gotRequest.Model)
	}
	if gotRequest.Stream {
		t.Fatal("stream = true, want false")
	}
	if len(gotRequest.Messages) != 2 {
		t.Fatalf("message count = %d, want 2", len(gotRequest.Messages))
	}
	if gotRequest.Messages[0].Role != "system" {
		t.Fatalf("system role = %q, want system", gotRequest.Messages[0].Role)
	}
	if gotRequest.Messages[1].Role != "user" {
		t.Fatalf("user role = %q, want user", gotRequest.Messages[1].Role)
	}
	if gotRequest.Messages[1].Content == "" || gotRequest.Messages[1].Content == "messy draft note" {
		t.Fatalf("unexpected user prompt content = %q", gotRequest.Messages[1].Content)
	}
}

func newTestEngine(t *testing.T, host string) *Engine {
	t.Helper()

	store, err := db.Open(filepath.Join(t.TempDir(), "engine.db"))
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() {
		_ = store.Close()
	})

	settings := DefaultSettings()
	settings.Provider = ProviderConfig{
		Host:  host,
		Port:  0,
		HTTPS: false,
		Model: "test-model",
	}
	return New(store, settings)
}
