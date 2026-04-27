package irecallapi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClientDoJSONRejectsOversizedResponse(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{\"payload\":\"" + strings.Repeat("a", maxResponseBodySize) + "\"}"))
	}))
	defer server.Close()

	client, err := NewClient(Config{
		BaseURL:  server.URL,
		APIToken: "test-token",
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	var dst map[string]any
	err = client.doJSON(context.Background(), http.MethodGet, "/", nil, &dst)
	if !errors.Is(err, errResponseTooLarge) {
		t.Fatalf("doJSON() error = %v, want %v", err, errResponseTooLarge)
	}
}
