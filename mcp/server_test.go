package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	mcpclient "github.com/mark3labs/mcp-go/client"
	mcpproto "github.com/mark3labs/mcp-go/mcp"
)

func TestMCPToolsCallAuthenticatedRESTAPI(t *testing.T) {
	const token = "irc_test_token"

	var seen []restCall
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer "+token {
			http.Error(w, `{"error":"bad auth"}`, http.StatusUnauthorized)
			return
		}

		call := restCall{Method: r.Method, Path: r.URL.Path, Query: r.URL.RawQuery}
		defer func() { seen = append(seen, call) }()
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/api/app/bootstrap-state":
			if r.Method != http.MethodGet {
				http.NotFound(w, r)
				return
			}
			_, _ = w.Write([]byte(`{"productName":"iRecall","greeting":"Hi! Test","pages":["Recall"],"paths":{"rootDir":"/tmp/irecall"},"profile":{"displayName":"Tester"},"docs":{"mcp":"docs/MCP_OPENCLAW.md"}}`))
		case "/api/app/list-quotes":
			if r.Method != http.MethodGet {
				http.NotFound(w, r)
				return
			}
			_, _ = w.Write([]byte(`[{"ID":7,"Content":"stored quote","Tags":["test"]}]`))
		case "/api/app/add-quote":
			if r.Method != http.MethodPost {
				http.NotFound(w, r)
				return
			}
			var req struct {
				Content string `json:"content"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatalf("decode add quote request: %v", err)
			}
			call.Body = req.Content
			_, _ = w.Write([]byte(`{"ID":8,"Content":"` + req.Content + `"}`))
		case "/api/app/run-recall":
			if r.Method != http.MethodPost {
				http.NotFound(w, r)
				return
			}
			var req struct {
				Question string `json:"question"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatalf("decode recall request: %v", err)
			}
			call.Body = req.Question
			_, _ = w.Write([]byte(`{"question":"` + req.Question + `","keywords":["memory"],"quotes":[],"response":"grounded answer"}`))
		case "/api/app/save-recall-as-quote":
			if r.Method != http.MethodPost {
				http.NotFound(w, r)
				return
			}
			var req struct {
				Question string   `json:"question"`
				Response string   `json:"response"`
				Keywords []string `json:"keywords"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatalf("decode save recall request: %v", err)
			}
			call.Body = req.Question + "|" + req.Response
			_, _ = w.Write([]byte(`{"ID":9,"Content":"saved recall"}`))
		case "/api/app/update-quote":
			if r.Method != http.MethodPost {
				http.NotFound(w, r)
				return
			}
			var req struct {
				ID      int64  `json:"id"`
				Content string `json:"content"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatalf("decode update quote request: %v", err)
			}
			call.Body = req.Content
			_, _ = w.Write([]byte(`{"ID":10,"Content":"` + req.Content + `"}`))
		case "/api/app/delete-quotes":
			if r.Method != http.MethodPost {
				http.NotFound(w, r)
				return
			}
			var req struct {
				IDs []int64 `json:"ids"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatalf("decode delete quotes request: %v", err)
			}
			call.Body = encodeIDs(req.IDs)
			_, _ = w.Write([]byte(`{"ok":true}`))
		case "/api/app/list-recall-history":
			if r.Method != http.MethodGet {
				http.NotFound(w, r)
				return
			}
			_, _ = w.Write([]byte(`[{"ID":11,"Question":"old question","Response":"old response"}]`))
		case "/api/app/get-recall-history":
			if r.Method != http.MethodGet {
				http.NotFound(w, r)
				return
			}
			call.Body = r.URL.Query().Get("id")
			_, _ = w.Write([]byte(`{"ID":11,"Question":"old question","Response":"old response","Quotes":[{"ID":7,"Content":"stored quote"}]}`))
		case "/api/app/delete-recall-history":
			if r.Method != http.MethodPost {
				http.NotFound(w, r)
				return
			}
			var req struct {
				IDs []int64 `json:"ids"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatalf("decode delete history request: %v", err)
			}
			call.Body = encodeIDs(req.IDs)
			_, _ = w.Write([]byte(`{"ok":true}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := newInProcessTestClient(t, Config{BaseURL: server.URL, APIToken: token, HTTPTimeout: 5 * time.Second})
	defer client.Close()

	tools, err := client.ListTools(context.Background(), mcpproto.ListToolsRequest{})
	if err != nil {
		t.Fatalf("ListTools() error = %v", err)
	}
	assertToolNames(t, tools, []string{
		"irecall_health",
		"irecall_recall",
		"irecall_list_quotes",
		"irecall_add_quote",
		"irecall_save_recall_as_quote",
		"irecall_update_quote",
		"irecall_delete_quotes",
		"irecall_list_history",
		"irecall_get_history",
		"irecall_delete_history",
	})

	healthResult := callTool(t, client, "irecall_health", nil)
	assertToolTextContains(t, healthResult, `"ok": true`)
	assertToolTextNotContains(t, healthResult, `"pages"`)
	assertToolTextNotContains(t, healthResult, `"paths"`)
	assertToolTextContains(t, callTool(t, client, "irecall_list_quotes", map[string]any{"limit": 10, "offset": 20}), "stored quote")
	assertToolTextContains(t, callTool(t, client, "irecall_add_quote", map[string]any{"content": "new note"}), "new note")
	assertToolTextContains(t, callTool(t, client, "irecall_recall", map[string]any{"question": "what did I save?"}), "grounded answer")
	assertToolTextContains(t, callTool(t, client, "irecall_save_recall_as_quote", map[string]any{"question": "q", "response": "r", "keywords": []string{"k"}}), "saved recall")
	assertToolTextContains(t, callTool(t, client, "irecall_update_quote", map[string]any{"id": 10, "content": "updated note"}), "updated note")
	assertToolTextContains(t, callTool(t, client, "irecall_delete_quotes", map[string]any{"ids": []int64{10, 11}}), `"ok": true`)
	assertToolTextContains(t, callTool(t, client, "irecall_list_history", nil), "old question")
	assertToolTextContains(t, callTool(t, client, "irecall_get_history", map[string]any{"id": 11}), "stored quote")
	assertToolTextContains(t, callTool(t, client, "irecall_delete_history", map[string]any{"ids": []int64{11}}), `"ok": true`)

	want := []restCall{
		{Method: http.MethodGet, Path: "/api/app/bootstrap-state"},
		{Method: http.MethodGet, Path: "/api/app/list-quotes", Query: "limit=10&offset=20"},
		{Method: http.MethodPost, Path: "/api/app/add-quote", Body: "new note"},
		{Method: http.MethodPost, Path: "/api/app/run-recall", Body: "what did I save?"},
		{Method: http.MethodPost, Path: "/api/app/save-recall-as-quote", Body: "q|r"},
		{Method: http.MethodPost, Path: "/api/app/update-quote", Body: "updated note"},
		{Method: http.MethodPost, Path: "/api/app/delete-quotes", Body: "10,11"},
		{Method: http.MethodGet, Path: "/api/app/list-recall-history"},
		{Method: http.MethodGet, Path: "/api/app/get-recall-history", Query: "id=11", Body: "11"},
		{Method: http.MethodPost, Path: "/api/app/delete-recall-history", Body: "11"},
	}
	if !equalRESTCalls(seen, want) {
		t.Fatalf("REST calls = %#v, want %#v", seen, want)
	}
}

func TestMCPToolReturnsRESTErrorAsToolError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"invalid token"}`))
	}))
	defer server.Close()

	client := newInProcessTestClient(t, Config{BaseURL: server.URL, APIToken: "bad-token", HTTPTimeout: 5 * time.Second})
	defer client.Close()

	result := callTool(t, client, "irecall_list_quotes", nil)
	if !result.IsError {
		t.Fatalf("CallTool().IsError = false, want true")
	}
	assertToolTextContains(t, result, "invalid token")
}

type restCall struct {
	Method string
	Path   string
	Query  string
	Body   string
}

func newInProcessTestClient(t *testing.T, cfg Config) *mcpclient.Client {
	t.Helper()
	srv, err := NewServer(cfg, "test")
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}
	client, err := mcpclient.NewInProcessClient(srv)
	if err != nil {
		t.Fatalf("NewInProcessClient() error = %v", err)
	}
	if err := client.Start(context.Background()); err != nil {
		t.Fatalf("client.Start() error = %v", err)
	}
	initReq := mcpproto.InitializeRequest{}
	initReq.Params.ProtocolVersion = mcpproto.LATEST_PROTOCOL_VERSION
	initReq.Params.ClientInfo = mcpproto.Implementation{Name: "irecall-test-client", Version: "test"}
	if _, err := client.Initialize(context.Background(), initReq); err != nil {
		t.Fatalf("client.Initialize() error = %v", err)
	}
	return client
}

func callTool(t *testing.T, client *mcpclient.Client, name string, args any) *mcpproto.CallToolResult {
	t.Helper()
	req := mcpproto.CallToolRequest{}
	req.Params.Name = name
	req.Params.Arguments = args
	result, err := client.CallTool(context.Background(), req)
	if err != nil {
		t.Fatalf("CallTool(%s) protocol error = %v", name, err)
	}
	return result
}

func assertToolNames(t *testing.T, result *mcpproto.ListToolsResult, names []string) {
	t.Helper()
	found := make(map[string]bool)
	for _, tool := range result.Tools {
		found[tool.Name] = true
	}
	for _, name := range names {
		if !found[name] {
			t.Fatalf("tool %q missing from ListTools result %#v", name, result.Tools)
		}
	}
}

func assertToolTextContains(t *testing.T, result *mcpproto.CallToolResult, want string) {
	t.Helper()
	if !strings.Contains(toolText(result), want) {
		t.Fatalf("tool result %#v does not contain %q", result, want)
	}
}

func assertToolTextNotContains(t *testing.T, result *mcpproto.CallToolResult, unwanted string) {
	t.Helper()
	if strings.Contains(toolText(result), unwanted) {
		t.Fatalf("tool result %#v unexpectedly contains %q", result, unwanted)
	}
}

func toolText(result *mcpproto.CallToolResult) string {
	var out strings.Builder
	for _, content := range result.Content {
		if text, ok := content.(mcpproto.TextContent); ok {
			out.WriteString(text.Text)
		}
	}
	return out.String()
}

func equalRESTCalls(a, b []restCall) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func encodeIDs(ids []int64) string {
	parts := make([]string, 0, len(ids))
	for _, id := range ids {
		parts = append(parts, fmt.Sprintf("%d", id))
	}
	return strings.Join(parts, ",")
}
