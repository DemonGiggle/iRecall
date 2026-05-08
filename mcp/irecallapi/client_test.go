package irecallapi

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestNewClientValidatesAndNormalizesConfig(t *testing.T) {
	t.Parallel()

	_, err := NewClient(Config{APIToken: "token"})
	if err == nil || !strings.Contains(err.Error(), "base URL is required") {
		t.Fatalf("NewClient() error = %v, want base URL required", err)
	}

	_, err = NewClient(Config{BaseURL: "http://example.com"})
	if err == nil || !strings.Contains(err.Error(), "API token is required") {
		t.Fatalf("NewClient() error = %v, want API token required", err)
	}

	_, err = NewClient(Config{BaseURL: "http://[::1", APIToken: "token"})
	if err == nil || !strings.Contains(err.Error(), "parse base URL") {
		t.Fatalf("NewClient() error = %v, want parse base URL failure", err)
	}

	client, err := NewClient(Config{
		BaseURL:  "http://example.com/",
		APIToken: "token",
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if client.baseURL != "http://example.com" {
		t.Fatalf("baseURL = %q, want trimmed trailing slash", client.baseURL)
	}
	if client.httpClient.Timeout != 15*time.Second {
		t.Fatalf("timeout = %v, want default 15s", client.httpClient.Timeout)
	}
}

func TestClientMethodsUseExpectedPathsHeadersAndPayloads(t *testing.T) {
	t.Parallel()

	const (
		quoteCreatedAt   = "2026-05-07T04:00:00Z"
		quoteUpdatedAt   = "2026-05-07T04:05:00Z"
		historyCreatedAt = "2026-05-07T04:10:00Z"
	)

	type requestSnapshot struct {
		Method      string
		Path        string
		Query       string
		Auth        string
		UserAgent   string
		ContentType string
		Body        string
	}

	var seen []requestSnapshot
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		snapshot := requestSnapshot{
			Method:      r.Method,
			Path:        r.URL.Path,
			Query:       r.URL.RawQuery,
			Auth:        r.Header.Get("Authorization"),
			UserAgent:   r.Header.Get("User-Agent"),
			ContentType: r.Header.Get("Content-Type"),
		}
		if r.Body != nil {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("read request body: %v", err)
			}
			snapshot.Body = strings.TrimSpace(string(body))
		}
		seen = append(seen, snapshot)

		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/app/bootstrap-state":
			_, _ = w.Write([]byte(`{"productName":"iRecall","greeting":"Hi! Tester","paths":{"rootDir":"/tmp/irecall"}}`))
		case "/api/app/count-quotes":
			_, _ = w.Write([]byte(`{"count":1}`))
		case "/api/app/count-recall-history":
			_, _ = w.Write([]byte(`{"count":1}`))
		case "/api/app/list-quotes":
			_, _ = w.Write([]byte(`[{"ID":7,"GlobalID":"quote-7","Content":"stored quote","Tags":["test"],"Version":2,"IsOwnedByMe":true,"CreatedAt":"` + quoteCreatedAt + `","UpdatedAt":"` + quoteUpdatedAt + `"}]`))
		case "/api/app/add-quote":
			_, _ = w.Write([]byte(`{"ID":8,"GlobalID":"quote-8","Content":"new note","Tags":["captured"],"CreatedAt":"` + quoteCreatedAt + `","UpdatedAt":"` + quoteUpdatedAt + `"}`))
		case "/api/app/run-recall":
			_, _ = w.Write([]byte(`{"question":"what did I save?","keywords":["memory"],"quotes":[{"ID":7,"GlobalID":"quote-7","Content":"stored quote","Tags":["test"],"CreatedAt":"` + quoteCreatedAt + `","UpdatedAt":"` + quoteUpdatedAt + `"}],"response":"grounded answer"}`))
		case "/api/app/save-recall-as-quote":
			_, _ = w.Write([]byte(`{"ID":9,"GlobalID":"quote-9","Content":"saved recall","Tags":["memory"],"CreatedAt":"` + quoteCreatedAt + `","UpdatedAt":"` + quoteUpdatedAt + `"}`))
		case "/api/app/update-quote":
			_, _ = w.Write([]byte(`{"ID":10,"GlobalID":"quote-10","Content":"updated note","CreatedAt":"` + quoteCreatedAt + `","UpdatedAt":"` + quoteUpdatedAt + `"}`))
		case "/api/app/delete-quotes":
			_, _ = w.Write([]byte(`{"ok":true}`))
		case "/api/app/list-recall-history":
			_, _ = w.Write([]byte(`[{"ID":11,"Question":"old question","Response":"old response","CreatedAt":"` + historyCreatedAt + `"}]`))
		case "/api/app/get-recall-history":
			_, _ = w.Write([]byte(`{"ID":11,"Question":"old question","Response":"old response","CreatedAt":"` + historyCreatedAt + `","Quotes":[{"ID":7,"GlobalID":"quote-7","Content":"stored quote","Tags":["test"],"CreatedAt":"` + quoteCreatedAt + `","UpdatedAt":"` + quoteUpdatedAt + `"}]}`))
		case "/api/app/delete-recall-history":
			_, _ = w.Write([]byte(`{"ok":true}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)

	bootstrap, err := client.BootstrapState(context.Background())
	if err != nil {
		t.Fatalf("BootstrapState() error = %v", err)
	}
	if bootstrap.ProductName != "iRecall" {
		t.Fatalf("product name = %q, want iRecall", bootstrap.ProductName)
	}

	count, err := client.CountQuotes(context.Background())
	if err != nil {
		t.Fatalf("CountQuotes() error = %v", err)
	}
	if count.Count != 1 {
		t.Fatalf("count = %d, want 1", count.Count)
	}

	quotes, err := client.ListQuotes(context.Background(), 10, 20)
	if err != nil {
		t.Fatalf("ListQuotes() error = %v", err)
	}
	if len(quotes) != 1 || quotes[0].GlobalID != "quote-7" {
		t.Fatalf("quotes = %+v, want one quote with global ID quote-7", quotes)
	}

	added, err := client.AddQuote(context.Background(), "new note")
	if err != nil {
		t.Fatalf("AddQuote() error = %v", err)
	}
	if added.Content != "new note" {
		t.Fatalf("added content = %q, want new note", added.Content)
	}

	recall, err := client.RunRecall(context.Background(), "what did I save?")
	if err != nil {
		t.Fatalf("RunRecall() error = %v", err)
	}
	if recall.Response != "grounded answer" {
		t.Fatalf("recall response = %q, want grounded answer", recall.Response)
	}

	saved, err := client.SaveRecallAsQuote(context.Background(), "q", "r", []string{"k"})
	if err != nil {
		t.Fatalf("SaveRecallAsQuote() error = %v", err)
	}
	if saved.GlobalID != "quote-9" {
		t.Fatalf("saved global ID = %q, want quote-9", saved.GlobalID)
	}

	updated, err := client.UpdateQuote(context.Background(), 10, "updated note")
	if err != nil {
		t.Fatalf("UpdateQuote() error = %v", err)
	}
	if updated.Content != "updated note" {
		t.Fatalf("updated content = %q, want updated note", updated.Content)
	}

	deleted, err := client.DeleteQuotes(context.Background(), []int64{10, 11})
	if err != nil {
		t.Fatalf("DeleteQuotes() error = %v", err)
	}
	if !deleted.OK {
		t.Fatal("DeleteQuotes() OK = false, want true")
	}

	historyCount, err := client.CountRecallHistory(context.Background())
	if err != nil {
		t.Fatalf("CountRecallHistory() error = %v", err)
	}
	if historyCount.Count != 1 {
		t.Fatalf("history count = %d, want 1", historyCount.Count)
	}

	history, err := client.ListRecallHistoryPage(context.Background(), 5, 10)
	if err != nil {
		t.Fatalf("ListRecallHistoryPage() error = %v", err)
	}
	if len(history) != 1 || history[0].ID != 11 {
		t.Fatalf("history = %+v, want one entry with ID 11", history)
	}

	entry, err := client.GetRecallHistory(context.Background(), 11)
	if err != nil {
		t.Fatalf("GetRecallHistory() error = %v", err)
	}
	if entry.ID != 11 || len(entry.Quotes) != 1 || entry.Quotes[0].GlobalID != "quote-7" {
		t.Fatalf("history entry = %+v, want linked quote quote-7", entry)
	}

	historyDeleted, err := client.DeleteRecallHistory(context.Background(), []int64{11})
	if err != nil {
		t.Fatalf("DeleteRecallHistory() error = %v", err)
	}
	if !historyDeleted.OK {
		t.Fatal("DeleteRecallHistory() OK = false, want true")
	}

	wantSeen := []requestSnapshot{
		{Method: http.MethodGet, Path: "/api/app/bootstrap-state", Auth: "Bearer test-token", UserAgent: "irecall-mcp/0"},
		{Method: http.MethodGet, Path: "/api/app/count-quotes", Auth: "Bearer test-token", UserAgent: "irecall-mcp/0"},
		{Method: http.MethodGet, Path: "/api/app/list-quotes", Query: "limit=10&offset=20", Auth: "Bearer test-token", UserAgent: "irecall-mcp/0"},
		{Method: http.MethodPost, Path: "/api/app/add-quote", Auth: "Bearer test-token", UserAgent: "irecall-mcp/0", ContentType: "application/json", Body: `{"content":"new note"}`},
		{Method: http.MethodPost, Path: "/api/app/run-recall", Auth: "Bearer test-token", UserAgent: "irecall-mcp/0", ContentType: "application/json", Body: `{"question":"what did I save?"}`},
		{Method: http.MethodPost, Path: "/api/app/save-recall-as-quote", Auth: "Bearer test-token", UserAgent: "irecall-mcp/0", ContentType: "application/json", Body: `{"question":"q","response":"r","keywords":["k"]}`},
		{Method: http.MethodPost, Path: "/api/app/update-quote", Auth: "Bearer test-token", UserAgent: "irecall-mcp/0", ContentType: "application/json", Body: `{"id":10,"content":"updated note"}`},
		{Method: http.MethodPost, Path: "/api/app/delete-quotes", Auth: "Bearer test-token", UserAgent: "irecall-mcp/0", ContentType: "application/json", Body: `{"ids":[10,11]}`},
		{Method: http.MethodGet, Path: "/api/app/count-recall-history", Auth: "Bearer test-token", UserAgent: "irecall-mcp/0"},
		{Method: http.MethodGet, Path: "/api/app/list-recall-history", Query: "limit=5&offset=10", Auth: "Bearer test-token", UserAgent: "irecall-mcp/0"},
		{Method: http.MethodGet, Path: "/api/app/get-recall-history", Query: "id=11", Auth: "Bearer test-token", UserAgent: "irecall-mcp/0"},
		{Method: http.MethodPost, Path: "/api/app/delete-recall-history", Auth: "Bearer test-token", UserAgent: "irecall-mcp/0", ContentType: "application/json", Body: `{"ids":[11]}`},
	}
	if !reflect.DeepEqual(seen, wantSeen) {
		t.Fatalf("seen requests = %#v, want %#v", seen, wantSeen)
	}
}

func TestClientDoJSONRejectsOversizedResponse(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{\"payload\":\"" + strings.Repeat("a", maxResponseBodySize) + "\"}"))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)

	var dst map[string]any
	err := client.doJSON(context.Background(), http.MethodGet, "/", nil, &dst)
	if !errors.Is(err, errResponseTooLarge) {
		t.Fatalf("doJSON() error = %v, want %v", err, errResponseTooLarge)
	}
}

func TestClientDoJSONReturnsStructuredAndPlainAPIErrors(t *testing.T) {
	t.Parallel()

	t.Run("structured JSON error", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error":"bad auth"}`))
		}))
		defer server.Close()

		client := newTestClient(t, server.URL)
		err := client.doJSON(context.Background(), http.MethodGet, "/", nil, &map[string]any{})
		assertAPIError(t, err, http.StatusUnauthorized, "bad auth", `{"error":"bad auth"}`)
	})

	t.Run("plain text error", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadGateway)
			_, _ = w.Write([]byte("upstream failed"))
		}))
		defer server.Close()

		client := newTestClient(t, server.URL)
		err := client.doJSON(context.Background(), http.MethodGet, "/", nil, &map[string]any{})
		assertAPIError(t, err, http.StatusBadGateway, "upstream failed", "upstream failed")
	})
}

func TestClientDoJSONHandlesEmptyAndInvalidResponses(t *testing.T) {
	t.Parallel()

	t.Run("empty success body", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		client := newTestClient(t, server.URL)
		var dst map[string]any
		if err := client.doJSON(context.Background(), http.MethodGet, "/", nil, &dst); err != nil {
			t.Fatalf("doJSON() error = %v, want nil", err)
		}
		if dst != nil {
			t.Fatalf("dst = %#v, want nil for empty body", dst)
		}
	})

	t.Run("invalid JSON body", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte("{"))
		}))
		defer server.Close()

		client := newTestClient(t, server.URL)
		err := client.doJSON(context.Background(), http.MethodGet, "/", nil, &map[string]any{})
		if err == nil || !strings.Contains(err.Error(), "decode response") {
			t.Fatalf("doJSON() error = %v, want decode response failure", err)
		}
	})
}

func newTestClient(t *testing.T, baseURL string) *Client {
	t.Helper()

	client, err := NewClient(Config{
		BaseURL:  baseURL,
		APIToken: "test-token",
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	return client
}

func assertAPIError(t *testing.T, err error, statusCode int, message, body string) {
	t.Helper()

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("error type = %T, want *APIError", err)
	}
	if apiErr.StatusCode != statusCode || apiErr.Message != message || apiErr.Body != body {
		t.Fatalf("APIError = %+v, want status=%d message=%q body=%q", apiErr, statusCode, message, body)
	}
}
