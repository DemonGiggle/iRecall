package app

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/gigol/irecall/core"
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
	if len(state.Pages) != 4 {
		t.Fatalf("page count = %d, want 4", len(state.Pages))
	}
	if state.Pages[1] != "History" {
		t.Fatalf("pages = %v, want History tab in bootstrap state", state.Pages)
	}
	if state.Docs["uiDesign"] != "docs/UI_DESIGN.md" {
		t.Fatalf("ui design doc = %q, want docs/UI_DESIGN.md", state.Docs["uiDesign"])
	}
}

func TestDesktopBackendRecallHistoryLifecycle(t *testing.T) {
	t.Parallel()

	app, err := NewApp(filepath.Join(t.TempDir(), "desktop-history"))
	if err != nil {
		t.Fatalf("NewApp() error = %v", err)
	}
	t.Cleanup(func() { app.Shutdown(context.Background()) })

	if _, err := app.SaveUserProfile("Alice"); err != nil {
		t.Fatalf("SaveUserProfile() error = %v", err)
	}

	quote, err := app.AddQuote("History-enabled desktop quote")
	if err != nil {
		t.Fatalf("AddQuote() error = %v", err)
	}

	if _, err := app.engine.SaveRecallHistory(context.Background(),
		"How do I check history?",
		"Open the History tab and inspect the saved session.",
		[]core.Quote{*quote},
	); err != nil {
		t.Fatalf("SaveRecallHistory() error = %v", err)
	}

	history, err := app.ListRecallHistory()
	if err != nil {
		t.Fatalf("ListRecallHistory() error = %v", err)
	}
	if len(history) != 1 {
		t.Fatalf("history count = %d, want 1", len(history))
	}

	entry, err := app.GetRecallHistory(history[0].ID)
	if err != nil {
		t.Fatalf("GetRecallHistory() error = %v", err)
	}
	if entry.Question != "How do I check history?" {
		t.Fatalf("history question = %q, want exact saved question", entry.Question)
	}
	if len(entry.Quotes) != 1 || entry.Quotes[0].ID != quote.ID {
		t.Fatalf("history quotes = %+v, want original quote", entry.Quotes)
	}

	if err := app.DeleteRecallHistory([]int64{entry.ID}); err != nil {
		t.Fatalf("DeleteRecallHistory() error = %v", err)
	}

	history, err = app.ListRecallHistory()
	if err != nil {
		t.Fatalf("ListRecallHistory() after delete error = %v", err)
	}
	if len(history) != 0 {
		t.Fatalf("history count after delete = %d, want 0", len(history))
	}
}

func TestDesktopBackendSaveRecallAsQuote(t *testing.T) {
	t.Parallel()

	app, err := NewApp(filepath.Join(t.TempDir(), "desktop-recall-quote"))
	if err != nil {
		t.Fatalf("NewApp() error = %v", err)
	}
	t.Cleanup(func() { app.Shutdown(context.Background()) })

	if _, err := app.SaveUserProfile("Alice"); err != nil {
		t.Fatalf("SaveUserProfile() error = %v", err)
	}

	quote, err := app.SaveRecallAsQuote(
		"How do I export quotes?",
		"Use the share flow and select the quotes you want to export.",
		[]string{"export", "sharing"},
	)
	if err != nil {
		t.Fatalf("SaveRecallAsQuote() error = %v", err)
	}
	if quote.ID == 0 {
		t.Fatalf("quote id = %d, want persisted quote", quote.ID)
	}
	if quote.Content == "" || quote.Content[:9] != "Question:" {
		t.Fatalf("quote content = %q, want formatted recall quote", quote.Content)
	}
	if len(quote.Tags) == 0 {
		t.Fatalf("quote tags = %#v, want saved tags", quote.Tags)
	}
}

func TestDesktopBackendImportQuotesPayload(t *testing.T) {
	t.Parallel()

	app, err := NewApp(filepath.Join(t.TempDir(), "desktop-import-payload"))
	if err != nil {
		t.Fatalf("NewApp() error = %v", err)
	}
	t.Cleanup(func() { app.Shutdown(context.Background()) })

	if _, err := app.SaveUserProfile("Alice"); err != nil {
		t.Fatalf("SaveUserProfile() error = %v", err)
	}

	quote, err := app.AddQuote("payload import roundtrip")
	if err != nil {
		t.Fatalf("AddQuote() error = %v", err)
	}

	payload, err := app.PreviewQuoteExport([]int64{quote.ID})
	if err != nil {
		t.Fatalf("PreviewQuoteExport() error = %v", err)
	}

	target, err := NewApp(filepath.Join(t.TempDir(), "desktop-import-target"))
	if err != nil {
		t.Fatalf("NewApp(target) error = %v", err)
	}
	t.Cleanup(func() { target.Shutdown(context.Background()) })

	if _, err := target.SaveUserProfile("Bob"); err != nil {
		t.Fatalf("SaveUserProfile(target) error = %v", err)
	}

	result, err := target.ImportQuotesPayload(payload)
	if err != nil {
		t.Fatalf("ImportQuotesPayload() error = %v", err)
	}
	if result.Inserted != 1 {
		t.Fatalf("ImportQuotesPayload() result = %+v, want inserted=1", result)
	}
}
