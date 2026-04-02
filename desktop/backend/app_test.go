package backend

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestDesktopBackendQuoteShareRoundTrip(t *testing.T) {
	t.Parallel()

	alice, err := NewApp(filepath.Join(t.TempDir(), "alice"))
	if err != nil {
		t.Fatalf("NewApp(alice) error = %v", err)
	}
	t.Cleanup(func() { alice.Shutdown(context.Background()) })

	if _, err := alice.SaveUserProfile("Alice"); err != nil {
		t.Fatalf("SaveUserProfile(alice) error = %v", err)
	}
	q, err := alice.AddQuote("desktop share roundtrip")
	if err != nil {
		t.Fatalf("AddQuote() error = %v", err)
	}

	exportPath := filepath.Join(t.TempDir(), "quotes", "share.json")
	if err := alice.ExportQuotesToFile([]int64{q.ID}, exportPath); err != nil {
		t.Fatalf("ExportQuotesToFile() error = %v", err)
	}
	if _, err := os.Stat(exportPath); err != nil {
		t.Fatalf("Stat(exportPath) error = %v", err)
	}

	bob, err := NewApp(filepath.Join(t.TempDir(), "bob"))
	if err != nil {
		t.Fatalf("NewApp(bob) error = %v", err)
	}
	t.Cleanup(func() { bob.Shutdown(context.Background()) })

	if _, err := bob.SaveUserProfile("Bob"); err != nil {
		t.Fatalf("SaveUserProfile(bob) error = %v", err)
	}
	result, err := bob.ImportQuotesFromFile(exportPath)
	if err != nil {
		t.Fatalf("ImportQuotesFromFile() error = %v", err)
	}
	if result.Inserted != 1 {
		t.Fatalf("import result = %+v, want inserted=1", result)
	}

	quotes, err := bob.ListQuotes()
	if err != nil {
		t.Fatalf("ListQuotes() error = %v", err)
	}
	if len(quotes) != 1 {
		t.Fatalf("quote count = %d, want 1", len(quotes))
	}
	if quotes[0].Content != "desktop share roundtrip" {
		t.Fatalf("imported content = %q, want desktop share roundtrip", quotes[0].Content)
	}
	if quotes[0].SourceName != "Alice" {
		t.Fatalf("source name = %q, want Alice", quotes[0].SourceName)
	}
}

func TestDesktopBackendBootstrapState(t *testing.T) {
	t.Parallel()

	app, err := NewApp(filepath.Join(t.TempDir(), "desktop"))
	if err != nil {
		t.Fatalf("NewApp() error = %v", err)
	}
	t.Cleanup(func() { app.Shutdown(context.Background()) })

	profile, err := app.SaveUserProfile("Alice")
	if err != nil {
		t.Fatalf("SaveUserProfile() error = %v", err)
	}

	state := app.BootstrapState()
	if state.ProductName != "iRecall" {
		t.Fatalf("product name = %q, want iRecall", state.ProductName)
	}
	if state.Greeting != "Hi! Alice" {
		t.Fatalf("greeting = %q, want Hi! Alice", state.Greeting)
	}
	if state.Profile == nil || state.Profile.DisplayName != profile.DisplayName {
		t.Fatalf("profile = %+v, want display name %q", state.Profile, profile.DisplayName)
	}
	if len(state.Pages) != 3 {
		t.Fatalf("page count = %d, want 3", len(state.Pages))
	}
	if state.Docs["uiDesign"] != "docs/UI_DESIGN.md" {
		t.Fatalf("ui design doc = %q, want docs/UI_DESIGN.md", state.Docs["uiDesign"])
	}
}
