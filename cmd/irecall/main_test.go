package main

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUsageTextIncludesCurrentFlagsAndExamples(t *testing.T) {
	t.Parallel()

	fs := flag.NewFlagSet("irecall", flag.ContinueOnError)
	fs.Bool("debug", false, "enable debug logging")
	fs.String("data-path", "", "store database, config, and logs under this root path")
	fs.Bool("version", false, "print version and exit")

	text := usageText(fs, "irecall")

	for _, want := range []string{
		"Usage:",
		"-debug",
		"-data-path",
		"-version",
		"--version",
		"manual quote sharing via exported JSON",
		"asks for your display name",
		"/tmp/irecall-alice",
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
	code := run("irecall", []string{"--version"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("run(--version) exit code = %d, want 0", code)
	}
	if got := strings.TrimSpace(stdout.String()); got != "iRecall v9.9.9" {
		t.Fatalf("stdout = %q, want version output", got)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestRunReportsDirectorySetupErrors(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "not-a-directory")
	if err := os.WriteFile(filePath, []byte("block"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := run("irecall", []string{"-data-path", filePath}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("run(-data-path file) exit code = %d, want 1", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
	errText := stderr.String()
	if !strings.Contains(errText, "irecall: cannot create data directories:") {
		t.Fatalf("stderr = %q, want directory setup failure", errText)
	}
}
