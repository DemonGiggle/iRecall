package main

import (
	"bytes"
	"flag"
	"strings"
	"testing"
)

func TestUsageTextIncludesFlagsAndExamples(t *testing.T) {
	t.Parallel()

	fs := flag.NewFlagSet("irecall-mcp", flag.ContinueOnError)
	fs.String("base-url", "", "override the iRecall web API base URL")
	fs.String("token-file", "", "read the API token from this file")
	fs.Duration("timeout", 0, "HTTP timeout for calls to the iRecall web API")
	fs.Bool("version", false, "print version and exit")

	text := usageText(fs, "irecall-mcp")

	for _, want := range []string{
		"Usage:",
		"-base-url",
		"-token-file",
		"-timeout",
		"-version",
		"--token-file ~/.config/irecall/mcp-api-token",
		"IRECALL_BASE_URL=http://127.0.0.1:9527",
		"stdio",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("usage text missing %q:\n%s", want, text)
		}
	}
}

func TestBinaryVersionPrefersInjectedValue(t *testing.T) {
	original := version
	version = "v1.2.3"
	t.Cleanup(func() { version = original })

	if got := binaryVersion(); got != "v1.2.3" {
		t.Fatalf("binaryVersion() = %q, want %q", got, "v1.2.3")
	}
}

func TestRunVersionWritesVersionAndExitsSuccess(t *testing.T) {
	original := version
	version = "v9.9.9"
	t.Cleanup(func() { version = original })

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := run("irecall-mcp", []string{"--version"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("run(--version) exit code = %d, want 0", code)
	}
	if got := strings.TrimSpace(stdout.String()); got != "iRecall MCP v9.9.9" {
		t.Fatalf("stdout = %q, want version output", got)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestRunReportsConfigErrors(t *testing.T) {
	t.Setenv("IRECALL_API_TOKEN_FILE", "")
	t.Setenv("IRECALL_API_TOKEN", "")
	t.Setenv("IRECALL_BASE_URL", "")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := run("irecall-mcp", nil, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("run() exit code = %d, want 1", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
	errText := stderr.String()
	if !strings.Contains(errText, "irecall-mcp:") {
		t.Fatalf("stderr = %q, want irecall-mcp prefix", errText)
	}
	if !strings.Contains(errText, "--token-file") {
		t.Fatalf("stderr = %q, want token source guidance", errText)
	}
}
