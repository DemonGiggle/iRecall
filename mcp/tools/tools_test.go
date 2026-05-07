package tools

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gigol/irecall/mcp/irecallapi"
	mcpclient "github.com/mark3labs/mcp-go/client"
	mcpproto "github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
)

func TestRegisteredToolsHappyPath(t *testing.T) {
	t.Parallel()

	const (
		quoteCreatedAt   = "2026-05-07T04:00:00Z"
		quoteUpdatedAt   = "2026-05-07T04:05:00Z"
		historyCreatedAt = "2026-05-07T04:10:00Z"
	)

	type requestSnapshot struct {
		Method string
		Path   string
		Query  string
		Body   string
	}

	var seen []requestSnapshot
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			http.Error(w, `{"error":"bad auth"}`, http.StatusUnauthorized)
			return
		}

		snapshot := requestSnapshot{Method: r.Method, Path: r.URL.Path, Query: r.URL.RawQuery}
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
			_, _ = w.Write([]byte(`{"productName":"iRecall","profile":{"displayName":"Tester"},"paths":{"rootDir":"/tmp/irecall"},"pages":["Recall"]}`))
		case "/api/app/count-quotes":
			_, _ = w.Write([]byte(`{"count":1}`))
		case "/api/app/list-quotes":
			_, _ = w.Write([]byte(`[{"ID":7,"GlobalID":"quote-7","Content":"stored quote","Tags":["test"],"Version":2,"IsOwnedByMe":true,"CreatedAt":"` + quoteCreatedAt + `","UpdatedAt":"` + quoteUpdatedAt + `"}]`))
		case "/api/app/add-quote":
			_, _ = w.Write([]byte(`{"ID":8,"GlobalID":"quote-8","Content":"new note","Tags":["captured"],"CreatedAt":"` + quoteCreatedAt + `","UpdatedAt":"` + quoteUpdatedAt + `"}`))
		case "/api/app/update-quote":
			_, _ = w.Write([]byte(`{"ID":10,"GlobalID":"quote-10","Content":"updated note","CreatedAt":"` + quoteCreatedAt + `","UpdatedAt":"` + quoteUpdatedAt + `"}`))
		case "/api/app/delete-quotes":
			_, _ = w.Write([]byte(`{"ok":true}`))
		case "/api/app/save-recall-as-quote":
			_, _ = w.Write([]byte(`{"ID":9,"GlobalID":"quote-9","Content":"saved recall","Tags":["memory"],"CreatedAt":"` + quoteCreatedAt + `","UpdatedAt":"` + quoteUpdatedAt + `"}`))
		case "/api/app/run-recall":
			_, _ = w.Write([]byte(`{"question":"what did I save?","keywords":["memory"],"quotes":[{"ID":7,"GlobalID":"quote-7","Content":"stored quote","Tags":["test"],"CreatedAt":"` + quoteCreatedAt + `","UpdatedAt":"` + quoteUpdatedAt + `"}],"response":"grounded answer"}`))
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
	defer api.Close()

	client := newToolAPIClient(t, api.URL)
	mcp := newToolTestClient(t, client, RegisterHealthTool, RegisterQuoteTools, RegisterRecallTool, RegisterHistoryTools)
	defer mcp.Close()

	assertToolTextContains(t, callTool(t, mcp, "irecall_health", nil), `"ok": true`)
	assertToolTextContains(t, callTool(t, mcp, "irecall_health", nil), `"profilePresent": true`)
	assertToolTextContains(t, callTool(t, mcp, "irecall_count_quotes", nil), `"count": 1`)

	listQuotes := callTool(t, mcp, "irecall_list_quotes", nil)
	assertToolTextContains(t, listQuotes, `"limit": 20`)
	assertToolTextContains(t, listQuotes, `"globalId": "quote-7"`)
	assertToolTextNotContains(t, listQuotes, `"GlobalID"`)

	addQuote := callTool(t, mcp, "irecall_add_quote", map[string]any{"content": "  new note  "})
	assertToolTextContains(t, addQuote, `"content": "new note"`)

	updateQuote := callTool(t, mcp, "irecall_update_quote", map[string]any{"id": 10, "content": " updated note "})
	assertToolTextContains(t, updateQuote, `"content": "updated note"`)

	deleteQuotes := callTool(t, mcp, "irecall_delete_quotes", map[string]any{"ids": []int64{10, 11}})
	assertToolTextContains(t, deleteQuotes, `"ok": true`)

	saveRecall := callTool(t, mcp, "irecall_save_recall_as_quote", map[string]any{"question": "q", "response": "r", "keywords": []string{"k"}})
	assertToolTextContains(t, saveRecall, `"globalId": "quote-9"`)

	recall := callTool(t, mcp, "irecall_recall", map[string]any{"question": "  what did I save?  "})
	assertToolTextContains(t, recall, `"response": "grounded answer"`)

	listHistory := callTool(t, mcp, "irecall_list_history", nil)
	assertToolTextContains(t, listHistory, `"id": 11`)

	getHistory := callTool(t, mcp, "irecall_get_history", map[string]any{"id": 11})
	assertToolTextContains(t, getHistory, `"globalId": "quote-7"`)

	deleteHistory := callTool(t, mcp, "irecall_delete_history", map[string]any{"ids": []int64{11}})
	assertToolTextContains(t, deleteHistory, `"ok": true`)

	wantSeen := []requestSnapshot{
		{Method: http.MethodGet, Path: "/api/app/bootstrap-state"},
		{Method: http.MethodGet, Path: "/api/app/bootstrap-state"},
		{Method: http.MethodGet, Path: "/api/app/count-quotes"},
		{Method: http.MethodGet, Path: "/api/app/list-quotes", Query: "limit=20"},
		{Method: http.MethodPost, Path: "/api/app/add-quote", Body: `{"content":"new note"}`},
		{Method: http.MethodPost, Path: "/api/app/update-quote", Body: `{"id":10,"content":"updated note"}`},
		{Method: http.MethodPost, Path: "/api/app/delete-quotes", Body: `{"ids":[10,11]}`},
		{Method: http.MethodPost, Path: "/api/app/save-recall-as-quote", Body: `{"question":"q","response":"r","keywords":["k"]}`},
		{Method: http.MethodPost, Path: "/api/app/run-recall", Body: `{"question":"what did I save?"}`},
		{Method: http.MethodGet, Path: "/api/app/list-recall-history"},
		{Method: http.MethodGet, Path: "/api/app/get-recall-history", Query: "id=11"},
		{Method: http.MethodPost, Path: "/api/app/delete-recall-history", Body: `{"ids":[11]}`},
	}
	if !reflect.DeepEqual(seen, wantSeen) {
		t.Fatalf("seen requests = %#v, want %#v", seen, wantSeen)
	}
}

func TestRegisteredToolsValidationAndDownstreamErrors(t *testing.T) {
	t.Parallel()

	api := newErrorAPIServer(t)
	defer api.Close()
	client := newToolAPIClient(t, api.URL)
	mcp := newToolTestClient(t, client, RegisterHealthTool, RegisterQuoteTools, RegisterRecallTool, RegisterHistoryTools)
	defer mcp.Close()

	assertToolErrorContains(t, callTool(t, mcp, "irecall_list_quotes", map[string]any{"limit": -1}), "limit must be non-negative")
	assertToolErrorContains(t, callTool(t, mcp, "irecall_list_quotes", map[string]any{"offset": -1}), "offset must be non-negative")
	assertToolErrorContains(t, callTool(t, mcp, "irecall_add_quote", map[string]any{"content": "   "}), "content is required")
	assertToolErrorContains(t, callTool(t, mcp, "irecall_update_quote", map[string]any{"id": 0, "content": "x"}), "id must be positive")
	assertToolErrorContains(t, callTool(t, mcp, "irecall_update_quote", map[string]any{"id": 1, "content": "   "}), "content is required")
	assertToolErrorContains(t, callTool(t, mcp, "irecall_delete_quotes", map[string]any{"ids": []int64{}}), "ids is required")
	assertToolErrorContains(t, callTool(t, mcp, "irecall_delete_quotes", map[string]any{"ids": []int64{1, -2}}), "ids must contain only positive IDs")
	assertToolErrorContains(t, callTool(t, mcp, "irecall_save_recall_as_quote", map[string]any{"question": " ", "response": "r"}), "question is required")
	assertToolErrorContains(t, callTool(t, mcp, "irecall_save_recall_as_quote", map[string]any{"question": "q", "response": " "}), "response is required")
	assertToolErrorContains(t, callTool(t, mcp, "irecall_recall", map[string]any{"question": "   "}), "question is required")
	assertToolErrorContains(t, callTool(t, mcp, "irecall_get_history", map[string]any{"id": 0}), "id must be positive")
	assertToolErrorContains(t, callTool(t, mcp, "irecall_delete_history", map[string]any{"ids": []int64{}}), "ids is required")
	assertToolErrorContains(t, callTool(t, mcp, "irecall_delete_history", map[string]any{"ids": []int64{1, -2}}), "ids must contain only positive IDs")

	assertToolErrorContains(t, callTool(t, mcp, "irecall_health", nil), "Failed to reach the iRecall web API.")
	assertToolErrorContains(t, callTool(t, mcp, "irecall_count_quotes", nil), "Failed to count quotes in iRecall.")
	assertToolErrorContains(t, callTool(t, mcp, "irecall_recall", map[string]any{"question": "what happened?"}), "Failed to run recall in iRecall.")
	assertToolErrorContains(t, callTool(t, mcp, "irecall_list_history", nil), "Failed to list iRecall history.")
}

func TestContractHelpersAndJSONResult(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 5, 7, 4, 0, 0, 0, time.UTC)
	quote := irecallapi.Quote{
		ID:          7,
		GlobalID:    "quote-7",
		Content:     "stored quote",
		Tags:        []string{"test"},
		IsOwnedByMe: true,
		CreatedAt:   now,
		UpdatedAt:   now.Add(time.Minute),
	}

	if got := newQuoteResponse(nil); got != nil {
		t.Fatalf("newQuoteResponse(nil) = %#v, want nil", got)
	}
	if got := newQuoteResponses(nil); got != nil {
		t.Fatalf("newQuoteResponses(nil) = %#v, want nil", got)
	}
	if got := newRecallResultResponse(nil); got != nil {
		t.Fatalf("newRecallResultResponse(nil) = %#v, want nil", got)
	}
	if got := newRecallHistoryEntryResponse(nil); got != nil {
		t.Fatalf("newRecallHistoryEntryResponse(nil) = %#v, want nil", got)
	}
	if got := newRecallHistorySummaryResponses(nil); got != nil {
		t.Fatalf("newRecallHistorySummaryResponses(nil) = %#v, want nil", got)
	}

	gotQuote := newQuoteResponse(&quote)
	if gotQuote.GlobalID != "quote-7" || gotQuote.Content != "stored quote" || !gotQuote.IsOwnedByMe {
		t.Fatalf("newQuoteResponse() = %+v, want mapped quote fields", gotQuote)
	}

	gotList := newListQuotesResponse(20, 5, []irecallapi.Quote{quote})
	if gotList.Limit != 20 || gotList.Offset != 5 || len(gotList.Quotes) != 1 || gotList.Quotes[0].GlobalID != "quote-7" {
		t.Fatalf("newListQuotesResponse() = %+v, want mapped list response", gotList)
	}

	history := newRecallHistorySummaryResponses([]irecallapi.RecallHistorySummary{{ID: 11, Question: "q", Response: "r", CreatedAt: now}})
	if len(history) != 1 || history[0].ID != 11 {
		t.Fatalf("newRecallHistorySummaryResponses() = %+v, want mapped summary", history)
	}

	entry := newRecallHistoryEntryResponse(&irecallapi.RecallHistoryEntry{ID: 11, Question: "q", Response: "r", CreatedAt: now, Quotes: []irecallapi.Quote{quote}})
	if entry.ID != 11 || len(entry.Quotes) != 1 || entry.Quotes[0].GlobalID != "quote-7" {
		t.Fatalf("newRecallHistoryEntryResponse() = %+v, want mapped entry", entry)
	}

	recall := newRecallResultResponse(&irecallapi.RecallResult{Question: "q", Keywords: []string{"k"}, Quotes: []irecallapi.Quote{quote}, Response: "r"})
	if recall.Question != "q" || len(recall.Quotes) != 1 || recall.Quotes[0].GlobalID != "quote-7" {
		t.Fatalf("newRecallResultResponse() = %+v, want mapped recall response", recall)
	}

	result, err := jsonResult(map[string]any{"ok": true})
	if err != nil {
		t.Fatalf("jsonResult() error = %v", err)
	}
	assertToolTextContains(t, result, `"ok": true`)

	if _, err := jsonResult(make(chan int)); err == nil {
		t.Fatal("jsonResult(chan) error = nil, want marshal failure")
	}
}

func newErrorAPIServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(`{"error":"upstream failed"}`))
	}))
}

type toolRegistrar func(*mcpserver.MCPServer, *irecallapi.Client)

func newToolTestClient(t *testing.T, apiClient *irecallapi.Client, registrars ...toolRegistrar) *mcpclient.Client {
	t.Helper()

	srv := mcpserver.NewMCPServer("irecall-tools-test", "test")
	for _, register := range registrars {
		register(srv, apiClient)
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
	initReq.Params.ClientInfo = mcpproto.Implementation{Name: "irecall-tools-test", Version: "test"}
	if _, err := client.Initialize(context.Background(), initReq); err != nil {
		t.Fatalf("client.Initialize() error = %v", err)
	}
	return client
}

func newToolAPIClient(t *testing.T, baseURL string) *irecallapi.Client {
	t.Helper()

	client, err := irecallapi.NewClient(irecallapi.Config{
		BaseURL:     baseURL,
		APIToken:    "test-token",
		HTTPTimeout: 5 * time.Second,
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
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
		t.Fatalf("CallTool(%s) error = %v", name, err)
	}
	return result
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

func assertToolTextContains(t *testing.T, result *mcpproto.CallToolResult, want string) {
	t.Helper()
	if result.IsError {
		t.Fatalf("tool result unexpectedly errored: %s", toolText(result))
	}
	if !strings.Contains(toolText(result), want) {
		t.Fatalf("tool text %q does not contain %q", toolText(result), want)
	}
}

func assertToolTextNotContains(t *testing.T, result *mcpproto.CallToolResult, unwanted string) {
	t.Helper()
	if strings.Contains(toolText(result), unwanted) {
		t.Fatalf("tool text %q unexpectedly contains %q", toolText(result), unwanted)
	}
}

func assertToolErrorContains(t *testing.T, result *mcpproto.CallToolResult, want string) {
	t.Helper()
	if !result.IsError {
		t.Fatalf("tool result was not an error: %s", toolText(result))
	}
	if !strings.Contains(toolText(result), want) {
		t.Fatalf("tool error %q does not contain %q", toolText(result), want)
	}
}
