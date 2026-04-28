//go:build !wails

package main

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	irecallapp "github.com/gigol/irecall/app"
	"github.com/gigol/irecall/config"
)

func TestAuthCommandIssueTokenWithoutWebPasswordAndWritesTokenFile(t *testing.T) {
	root := t.TempDir()

	tokenPath := filepath.Join(t.TempDir(), "secrets", "irecall-api-token")
	var stdout bytes.Buffer
	err := runAuthCommand([]string{
		"issue-token",
		"--data-path", root,
		"--write-token-file", tokenPath,
	}, strings.NewReader(""), &stdout)
	if err != nil {
		t.Fatalf("runAuthCommand(issue-token) error = %v", err)
	}
	data, err := os.ReadFile(tokenPath)
	if err != nil {
		t.Fatalf("read token file: %v", err)
	}
	token := strings.TrimSpace(string(data))
	if strings.Contains(stdout.String(), token) {
		t.Fatalf("stdout leaked full token: %q", stdout.String())
	}
	if !strings.HasPrefix(token, "irc_") {
		t.Fatalf("token file content = %q, want iRecall token", token)
	}
	info, err := os.Stat(tokenPath)
	if err != nil {
		t.Fatalf("stat token file: %v", err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Fatalf("token file mode = %o, want 600", got)
	}

	runtimeApp, err := irecallapp.NewApp(root)
	if err != nil {
		t.Fatalf("NewApp() error = %v", err)
	}
	defer runtimeApp.Shutdown(context.Background())
	ok, err := runtimeApp.VerifyAPIToken(token)
	if err != nil {
		t.Fatalf("VerifyAPIToken() error = %v", err)
	}
	if !ok {
		t.Fatalf("VerifyAPIToken() = false, want true")
	}
}

func TestAuthCommandRotateAndRevokeTokenWithoutWebPassword(t *testing.T) {
	root := t.TempDir()

	firstToken, err := issueTokenForTest(root)
	if err != nil {
		t.Fatalf("issue first token: %v", err)
	}
	secondToken, err := issueTokenForTest(root)
	if err != nil {
		t.Fatalf("rotate token: %v", err)
	}
	if firstToken == secondToken {
		t.Fatalf("rotated token matched first token")
	}

	runtimeApp, err := irecallapp.NewApp(root)
	if err != nil {
		t.Fatalf("NewApp() error = %v", err)
	}
	ok, err := runtimeApp.VerifyAPIToken(firstToken)
	if err != nil {
		t.Fatalf("VerifyAPIToken(first) error = %v", err)
	}
	if ok {
		t.Fatalf("first token still valid after rotate")
	}
	ok, err = runtimeApp.VerifyAPIToken(secondToken)
	if err != nil {
		t.Fatalf("VerifyAPIToken(second) error = %v", err)
	}
	if !ok {
		t.Fatalf("second token not valid after rotate")
	}
	runtimeApp.Shutdown(context.Background())

	var stdout bytes.Buffer
	if err := runAuthCommand([]string{"revoke-token", "--data-path", root}, strings.NewReader(""), &stdout); err != nil {
		t.Fatalf("runAuthCommand(revoke-token) error = %v", err)
	}
	runtimeApp, err = irecallapp.NewApp(root)
	if err != nil {
		t.Fatalf("NewApp(after revoke) error = %v", err)
	}
	defer runtimeApp.Shutdown(context.Background())
	ok, err = runtimeApp.VerifyAPIToken(secondToken)
	if err != nil {
		t.Fatalf("VerifyAPIToken(after revoke) error = %v", err)
	}
	if ok {
		t.Fatalf("second token still valid after revoke")
	}
}

func TestAuthCommandAcceptsDeprecatedPasswordStdinWithoutPasswordConfigured(t *testing.T) {
	root := t.TempDir()
	err := runAuthCommand([]string{"issue-token", "--data-path", root, "--password-stdin"}, strings.NewReader("legacy-password\n"), &bytes.Buffer{})
	if err != nil {
		t.Fatalf("runAuthCommand() with deprecated --password-stdin error = %v", err)
	}
}

func TestAuthCommandDataPathDoesNotPersistPreferredRoot(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(t.TempDir(), "xdg-config"))
	originalRoot := config.RootPath()
	config.SetRootPath("")
	t.Cleanup(func() { config.SetRootPath(originalRoot) })

	root := t.TempDir()
	if err := runAuthCommand([]string{"issue-token", "--data-path", root}, strings.NewReader(""), &bytes.Buffer{}); err != nil {
		t.Fatalf("runAuthCommand(issue-token) error = %v", err)
	}
	preferredRoot, err := config.LoadPreferredRootPath()
	if err != nil {
		t.Fatalf("LoadPreferredRootPath() error = %v", err)
	}
	if preferredRoot != "" {
		t.Fatalf("preferred root = %q, want empty", preferredRoot)
	}
}

func issueTokenForTest(root string) (string, error) {
	tokenPath := filepath.Join(os.TempDir(), "irecall-token-test-"+filepath.Base(root))
	defer os.Remove(tokenPath)
	var stdout bytes.Buffer
	if err := runAuthCommand([]string{
		"rotate-token",
		"--data-path", root,
		"--write-token-file", tokenPath,
	}, strings.NewReader(""), &stdout); err != nil {
		return "", err
	}
	data, err := os.ReadFile(tokenPath)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}
