package mcp

import (
	"context"
	"encoding/json"
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

		call := restCall{Method: r.Method, Path: r.URL.Path}
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
	})

	assertToolTextContains(t, callTool(t, client, "irecall_health", nil), `"ok": true`)
	assertToolTextContains(t, callTool(t, client, "irecall_list_quotes", nil), "stored quote")
	assertToolTextContains(t, callTool(t, client, "irecall_add_quote", map[string]any{"content": "new note"}), "new note")
	assertToolTextContains(t, callTool(t, client, "irecall_recall", map[string]any{"question": "what did I save?"}), "grounded answer")
	assertToolTextContains(t, callTool(t, client, "irecall_save_recall_as_quote", map[string]any{"question": "q", "response": "r", "keywords": []string{"k"}}), "saved recall")

	want := []restCall{
		{Method: http.MethodGet, Path: "/api/app/bootstrap-state"},
		{Method: http.MethodGet, Path: "/api/app/list-quotes"},
		{Method: http.MethodPost, Path: "/api/app/add-quote", Body: "new note"},
		{Method: http.MethodPost, Path: "/api/app/run-recall", Body: "what did I save?"},
		{Method: http.MethodPost, Path: "/api/app/save-recall-as-quote", Body: "q|r"},
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
	for _, content := range result.Content {
		if text, ok := content.(mcpproto.TextContent); ok && strings.Contains(text.Text, want) {
			return
		}
	}
	t.Fatalf("tool result %#v does not contain %q", result, want)
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
