//go:build e2e

package e2e

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	mcpclient "github.com/mark3labs/mcp-go/client"
	mcpproto "github.com/mark3labs/mcp-go/mcp"
)

func TestMCPProcessEndToEndAgainstLiveWebServer(t *testing.T) {
	repoRoot := repoRoot(t)
	dataRoot := filepath.Join(t.TempDir(), "data")
	tokenPath := filepath.Join(t.TempDir(), "secrets", "irecall-api-token")
	port := freePort(t)
	baseURL := fmt.Sprintf("http://127.0.0.1:%d", port)

	tokenOutput := runWebCommand(t, repoRoot, "issue-token",
		"auth",
		"issue-token",
		"--data-path", dataRoot,
		"--write-token-file", tokenPath,
	)

	tokenData, err := os.ReadFile(tokenPath)
	if err != nil {
		t.Fatalf("read token file: %v", err)
	}
	token := strings.TrimSpace(string(tokenData))
	if token == "" {
		t.Fatal("token file is empty")
	}
	if strings.Contains(tokenOutput, token) {
		t.Fatalf("auth output leaked the full token: %q", tokenOutput)
	}

	startWebServer(t, repoRoot, dataRoot, port)
	waitForHTTP200(t, baseURL+"/api/auth/status", 60*time.Second)

	client := newLoggedMCPClient(t, repoRoot, baseURL, tokenPath)
	initializeClient(t, client)

	health := callToolText(t, client, "irecall_health", nil)
	assertContains(t, health, `"ok": true`)
	assertNotContains(t, health, `"pages"`)
	assertNotContains(t, health, `"paths"`)

	countBefore := callToolText(t, client, "irecall_count_quotes", nil)
	assertContains(t, countBefore, `"count": 0`)

	quoteContent := "process e2e quote"
	added := callToolText(t, client, "irecall_add_quote", map[string]any{"content": quoteContent})
	assertContains(t, added, `"id":`)
	assertContains(t, added, `"content": "`+quoteContent+`"`)
	assertNotContainsAny(t, added, []string{`"ID":`, `"GlobalID":`, `"Content":`, `"CreatedAt":`, `"UpdatedAt":`})

	listed := callToolText(t, client, "irecall_list_quotes", map[string]any{"limit": 10})
	assertContains(t, listed, quoteContent)
	assertContains(t, listed, `"globalId":`)
	assertNotContainsAny(t, listed, []string{`"ID":`, `"GlobalID":`, `"Content":`, `"CreatedAt":`, `"UpdatedAt":`})

	countAfter := callToolText(t, client, "irecall_count_quotes", nil)
	assertContains(t, countAfter, `"count": 1`)
}

func runWebCommand(t *testing.T, repoRoot, name string, args ...string) string {
	t.Helper()

	command, commandArgs := webRunnerArgs(repoRoot)
	commandArgs = append(commandArgs, args...)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, command, commandArgs...)
	cmd.Dir = repoRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s failed: %v\n%s", name, err, string(output))
	}
	return string(output)
}

func startWebServer(t *testing.T, repoRoot, dataRoot string, port int) {
	t.Helper()

	command, args := webRunnerArgs(repoRoot)
	args = append(args,
		"--api-only",
		"--host", "127.0.0.1",
		"--port", fmt.Sprintf("%d", port),
		"--data-path", dataRoot,
	)
	startLoggedProcess(t, "irecall-web", repoRoot, command, args)
}

func newLoggedMCPClient(t *testing.T, repoRoot, baseURL, tokenPath string) *mcpclient.Client {
	t.Helper()

	command, args := mcpRunnerArgs(repoRoot)
	args = append(args,
		"--base-url", baseURL,
		"--token-file", tokenPath,
		"--timeout", "5s",
	)
	client, err := mcpclient.NewStdioMCPClient(command, nil, args...)
	if err != nil {
		t.Fatalf("NewStdioMCPClient() error = %v", err)
	}

	logPath := filepath.Join(t.TempDir(), "irecall-mcp.log")
	logFile, err := os.Create(logPath)
	if err != nil {
		_ = client.Close()
		t.Fatalf("create irecall-mcp log file: %v", err)
	}

	done := make(chan struct{})
	if stderr, ok := mcpclient.GetStderr(client); ok {
		go func() {
			_, _ = io.Copy(logFile, stderr)
			_ = logFile.Close()
			close(done)
		}()
	} else {
		_ = logFile.Close()
		close(done)
	}

	t.Cleanup(func() {
		_ = client.Close()
		<-done
		if t.Failed() {
			t.Logf("irecall-mcp stderr (%s):\n%s", logPath, readLogTail(logPath))
		}
	})

	return client
}

func initializeClient(t *testing.T, client *mcpclient.Client) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	initReq := mcpproto.InitializeRequest{}
	initReq.Params.ProtocolVersion = mcpproto.LATEST_PROTOCOL_VERSION
	initReq.Params.ClientInfo = mcpproto.Implementation{Name: "irecall-process-e2e", Version: "test"}
	if _, err := client.Initialize(ctx, initReq); err != nil {
		t.Fatalf("client.Initialize() error = %v", err)
	}
}

func callToolText(t *testing.T, client *mcpclient.Client, name string, args any) string {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	req := mcpproto.CallToolRequest{}
	req.Params.Name = name
	req.Params.Arguments = args
	result, err := client.CallTool(ctx, req)
	if err != nil {
		t.Fatalf("CallTool(%s) protocol error = %v", name, err)
	}
	if result.IsError {
		t.Fatalf("CallTool(%s) returned tool error: %s", name, toolText(result))
	}
	return toolText(result)
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

func webRunnerArgs(repoRoot string) (string, []string) {
	if path := strings.TrimSpace(os.Getenv("IRECALL_E2E_WEB_BIN")); path != "" {
		return path, nil
	}
	return "go", []string{"run", filepath.Join(repoRoot, "web")}
}

func mcpRunnerArgs(repoRoot string) (string, []string) {
	if path := strings.TrimSpace(os.Getenv("IRECALL_E2E_MCP_BIN")); path != "" {
		return path, nil
	}
	return "go", []string{"run", filepath.Join(repoRoot, "cmd", "irecall-mcp")}
}

func repoRoot(t *testing.T) string {
	t.Helper()

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), ".."))
}

func freePort(t *testing.T) int {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen on ephemeral port: %v", err)
	}
	defer listener.Close()

	addr, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		t.Fatalf("unexpected listener addr type %T", listener.Addr())
	}
	return addr.Port
}

func waitForHTTP200(t *testing.T, url string, timeout time.Duration) {
	t.Helper()

	deadline := time.Now().Add(timeout)
	client := &http.Client{Timeout: 2 * time.Second}
	var lastErr error
	for time.Now().Before(deadline) {
		resp, err := client.Get(url)
		if err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return
			}
			lastErr = fmt.Errorf("status %d", resp.StatusCode)
		} else {
			lastErr = err
		}
		time.Sleep(200 * time.Millisecond)
	}
	t.Fatalf("timeout waiting for %s: %v", url, lastErr)
}

type loggedProcess struct {
	name    string
	logPath string
	cancel  context.CancelFunc
	done    chan struct{}
}

func startLoggedProcess(t *testing.T, name, repoRoot, command string, args []string) *loggedProcess {
	t.Helper()

	logPath := filepath.Join(t.TempDir(), name+".log")
	logFile, err := os.Create(logPath)
	if err != nil {
		t.Fatalf("create %s log file: %v", name, err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Dir = repoRoot
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	if err := cmd.Start(); err != nil {
		_ = logFile.Close()
		t.Fatalf("start %s: %v", name, err)
	}

	done := make(chan struct{})
	go func() {
		_ = cmd.Wait()
		_ = logFile.Close()
		close(done)
	}()

	process := &loggedProcess{
		name:    name,
		logPath: logPath,
		cancel:  cancel,
		done:    done,
	}
	t.Cleanup(func() {
		process.cancel()
		<-process.done
		if t.Failed() {
			t.Logf("%s log (%s):\n%s", process.name, process.logPath, readLogTail(process.logPath))
		}
	})
	return process
}

func readLogTail(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Sprintf("read log failed: %v", err)
	}
	const maxBytes = 8192
	if len(data) <= maxBytes {
		return string(data)
	}
	return "...(truncated)...\n" + string(data[len(data)-maxBytes:])
}

func assertContains(t *testing.T, haystack, needle string) {
	t.Helper()
	if !strings.Contains(haystack, needle) {
		t.Fatalf("expected %q to contain %q", haystack, needle)
	}
}

func assertNotContains(t *testing.T, haystack, needle string) {
	t.Helper()
	if strings.Contains(haystack, needle) {
		t.Fatalf("expected %q not to contain %q", haystack, needle)
	}
}

func assertNotContainsAny(t *testing.T, haystack string, needles []string) {
	t.Helper()
	for _, needle := range needles {
		assertNotContains(t, haystack, needle)
	}
}
