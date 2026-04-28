//go:build !wails

package main

import (
	"context"
	"os"
	"testing"

	irecallapp "github.com/gigol/irecall/app"
)

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
