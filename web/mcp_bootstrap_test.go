//go:build !wails

package main

import (
	"bytes"
	"context"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	irecallapp "github.com/gigol/irecall/app"
	irecallmcp "github.com/gigol/irecall/mcp"
	mcpclient "github.com/mark3labs/mcp-go/client"
	mcpproto "github.com/mark3labs/mcp-go/mcp"
)

func TestOperatorBootstrapIssuesTokenAndMCPHealthChecksRealWebServer(t *testing.T) {
	root := t.TempDir()
	password := "Secret-pass-123!"
	setupAuthCommandPassword(t, root, password)

	tokenPath := filepath.Join(t.TempDir(), "secrets", "irecall-api-token")
	var stdout bytes.Buffer
	if err := runAuthCommand([]string{
		"issue-token",
		"--data-path", root,
		"--password-stdin",
		"--write-token-file", tokenPath,
	}, strings.NewReader(password+"\n"), &stdout); err != nil {
		t.Fatalf("runAuthCommand(issue-token) error = %v", err)
	}
	tokenData, err := os.ReadFile(tokenPath)
	if err != nil {
		t.Fatalf("read token file: %v", err)
	}
	token := strings.TrimSpace(string(tokenData))
	if token == "" {
		t.Fatalf("token file is empty")
	}
	if strings.Contains(stdout.String(), token) {
		t.Fatalf("stdout leaked full token: %q", stdout.String())
	}

	runtimeApp, err := irecallapp.NewApp(root)
	if err != nil {
		t.Fatalf("NewApp() error = %v", err)
	}
	t.Cleanup(func() { runtimeApp.Shutdown(context.Background()) })
	settings := *runtimeApp.GetSettings()
	settings.Debug.MockLLM = true
	if _, err := runtimeApp.SaveSettings(settings); err != nil {
		t.Fatalf("SaveSettings(MockLLM) error = %v", err)
	}
	if _, err := runtimeApp.AddQuote("operator bootstrap quote"); err != nil {
		t.Fatalf("AddQuote() error = %v", err)
	}

	handler := newTestServer(t, runtimeApp)
	httpServer := httptest.NewServer(handler)
	defer httpServer.Close()

	mcpClient := newOperatorBootstrapMCPClient(t, irecallmcp.Config{
		BaseURL:     httpServer.URL,
		APIToken:    token,
		HTTPTimeout: 5 * time.Second,
	})
	defer mcpClient.Close()

	health := operatorBootstrapCallTool(t, mcpClient, "irecall_health", nil)
	assertOperatorBootstrapToolTextContains(t, health, `"ok": true`)
	assertOperatorBootstrapToolTextNotContains(t, health, `"paths"`)
	assertOperatorBootstrapToolTextNotContains(t, health, `"pages"`)

	quotes := operatorBootstrapCallTool(t, mcpClient, "irecall_list_quotes", map[string]any{"limit": 10})
	assertOperatorBootstrapToolTextContains(t, quotes, "operator bootstrap quote")
}

func newOperatorBootstrapMCPClient(t *testing.T, cfg irecallmcp.Config) *mcpclient.Client {
	t.Helper()
	srv, err := irecallmcp.NewServer(cfg, "operator-bootstrap-test")
	if err != nil {
		t.Fatalf("mcp.NewServer() error = %v", err)
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
	initReq.Params.ClientInfo = mcpproto.Implementation{Name: "irecall-operator-bootstrap-test", Version: "test"}
	if _, err := client.Initialize(context.Background(), initReq); err != nil {
		t.Fatalf("client.Initialize() error = %v", err)
	}
	return client
}

func operatorBootstrapCallTool(t *testing.T, client *mcpclient.Client, name string, args any) *mcpproto.CallToolResult {
	t.Helper()
	req := mcpproto.CallToolRequest{}
	req.Params.Name = name
	req.Params.Arguments = args
	result, err := client.CallTool(context.Background(), req)
	if err != nil {
		t.Fatalf("CallTool(%s) protocol error = %v", name, err)
	}
	if result.IsError {
		t.Fatalf("CallTool(%s) returned tool error: %s", name, operatorBootstrapToolText(result))
	}
	return result
}

func assertOperatorBootstrapToolTextContains(t *testing.T, result *mcpproto.CallToolResult, want string) {
	t.Helper()
	if !strings.Contains(operatorBootstrapToolText(result), want) {
		t.Fatalf("tool result %#v does not contain %q", result, want)
	}
}

func assertOperatorBootstrapToolTextNotContains(t *testing.T, result *mcpproto.CallToolResult, unwanted string) {
	t.Helper()
	if strings.Contains(operatorBootstrapToolText(result), unwanted) {
		t.Fatalf("tool result %#v unexpectedly contains %q", result, unwanted)
	}
}

func operatorBootstrapToolText(result *mcpproto.CallToolResult) string {
	var out strings.Builder
	for _, content := range result.Content {
		if text, ok := content.(mcpproto.TextContent); ok {
			out.WriteString(text.Text)
		}
	}
	return out.String()
}
