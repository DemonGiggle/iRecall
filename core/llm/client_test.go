package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestChatUsesMaxTokensByDefault(t *testing.T) {
	t.Parallel()

	var got struct {
		MaxTokens           *int     `json:"max_tokens"`
		MaxCompletionTokens *int     `json:"max_completion_tokens"`
		Temperature         *float64 `json:"temperature"`
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"choices": []map[string]any{{
				"message": map[string]string{"content": "ok"},
			}},
		})
	}))
	defer srv.Close()

	maxTokens := 42
	client := NewClient(ProviderConfig{Host: srv.Listener.Addr().String(), Model: "test-model"})
	if _, err := client.Chat(context.Background(), []Message{{Role: "user", Content: "hi"}}, nil, ChatOptions{MaxTokens: &maxTokens}); err != nil {
		t.Fatalf("Chat() error = %v", err)
	}
	if got.MaxTokens == nil || *got.MaxTokens != maxTokens {
		t.Fatalf("max_tokens = %v, want %d", got.MaxTokens, maxTokens)
	}
	if got.MaxCompletionTokens != nil {
		t.Fatalf("max_completion_tokens = %v, want nil", *got.MaxCompletionTokens)
	}
}

func TestChatUsesMaxCompletionTokensForGPT5(t *testing.T) {
	t.Parallel()

	var got struct {
		MaxTokens           *int     `json:"max_tokens"`
		MaxCompletionTokens *int     `json:"max_completion_tokens"`
		Temperature         *float64 `json:"temperature"`
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"choices": []map[string]any{{
				"message": map[string]string{"content": "ok"},
			}},
		})
	}))
	defer srv.Close()

	maxTokens := 42
	temperature := 0.0
	client := NewClient(ProviderConfig{Host: srv.Listener.Addr().String(), Model: "gpt-5-nano"})
	if _, err := client.Chat(context.Background(), []Message{{Role: "user", Content: "hi"}}, nil, ChatOptions{Temperature: &temperature, MaxTokens: &maxTokens}); err != nil {
		t.Fatalf("Chat() error = %v", err)
	}
	if got.MaxTokens != nil {
		t.Fatalf("max_tokens = %v, want nil", *got.MaxTokens)
	}
	if got.MaxCompletionTokens == nil || *got.MaxCompletionTokens != 1024 {
		t.Fatalf("max_completion_tokens = %v, want %d", got.MaxCompletionTokens, 1024)
	}
	if got.Temperature != nil {
		t.Fatalf("temperature = %v, want nil", *got.Temperature)
	}
}
