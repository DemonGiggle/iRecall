//go:build !wails

package main

import (
	"context"
	"flag"
	"os"
	"strings"
	"testing"

	irecallapp "github.com/gigol/irecall/app"
)

func TestUsageTextIncludesVersionFlagAndExamples(t *testing.T) {
	t.Parallel()

	fs := flag.NewFlagSet("irecall-web", flag.ContinueOnError)
	fs.Bool("api-only", false, "run a headless API server without serving the frontend UI or requiring a web password at startup")
	fs.String("host", "0.0.0.0", "host/interface to bind the web server to")
	fs.Int("port", 0, "port to listen on (overrides saved web port)")
	fs.String("provider-host", "", "override provider host/server for --api-only mode")
	fs.String("provider-model", "", "override response/default provider model for --api-only mode")
	fs.Bool("version", false, "print version and exit")

	text := usageText(fs, "irecall-web")

	for _, want := range []string{
		"Usage:",
		"irecall-web auth <subcommand> [flags]",
		"-api-only",
		"-host string",
		`default "0.0.0.0"`,
		"-port int",
		"-provider-host string",
		"-provider-model string",
		"-version",
		"--version",
		"--api-only --provider-host",
		"auth issue-token",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("usage text missing %q:\n%s", want, text)
		}
	}
}

func TestFlagDefaultsTextIncludesRegisteredDefaults(t *testing.T) {
	t.Parallel()

	fs := flag.NewFlagSet("irecall-web", flag.ContinueOnError)
	fs.String("host", "127.0.0.1", "bind address")
	fs.Int("port", 9527, "listen port")

	text := flagDefaultsText(fs)
	for _, want := range []string{
		"-host string",
		`default "127.0.0.1"`,
		"-port int",
		"default 9527",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("flag defaults text missing %q:\n%s", want, text)
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

func TestServerOptionsValidateRejectsConflictingFlags(t *testing.T) {
	t.Parallel()

	err := (ServerOptions{
		APIOnly:               true,
		UnsafeNoPasswordCheck: true,
	}).Validate()
	if err == nil {
		t.Fatal("Validate() error = nil, want conflicting flags error")
	}
}

func TestEnsureWebPasswordConfiguredRequiresPasswordInNormalMode(t *testing.T) {
	runtimeApp, err := irecallapp.NewApp(t.TempDir())
	if err != nil {
		t.Fatalf("NewApp() error = %v", err)
	}
	t.Cleanup(func() { runtimeApp.Shutdown(context.Background()) })

	originalStdin := os.Stdin
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() error = %v", err)
	}
	defer reader.Close()
	defer writer.Close()
	os.Stdin = reader
	t.Cleanup(func() { os.Stdin = originalStdin })

	if err := ensureWebPasswordConfigured(runtimeApp); err == nil {
		t.Fatal("ensureWebPasswordConfigured() error = nil, want missing password error")
	}
}

func TestEnsureWebPasswordConfiguredAllowsExistingPassword(t *testing.T) {
	runtimeApp, err := irecallapp.NewApp(t.TempDir())
	if err != nil {
		t.Fatalf("NewApp() error = %v", err)
	}
	t.Cleanup(func() { runtimeApp.Shutdown(context.Background()) })

	if err := runtimeApp.SetupPassword("Secret-pass-123!", "Secret-pass-123!"); err != nil {
		t.Fatalf("SetupPassword() error = %v", err)
	}

	if err := ensureWebPasswordConfigured(runtimeApp); err != nil {
		t.Fatalf("ensureWebPasswordConfigured() error = %v", err)
	}
}
