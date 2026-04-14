package pages

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestQuoteImportPageImportsSharedFile(t *testing.T) {
	t.Parallel()

	exporter := newShareTestEngine(t)
	quote, err := exporter.AddQuote(context.Background(), "shared from alice")
	if err != nil {
		t.Fatalf("AddQuote() error = %v", err)
	}
	payload, err := exporter.ExportQuotes(context.Background(), []int64{quote.ID})
	if err != nil {
		t.Fatalf("ExportQuotes() error = %v", err)
	}
	path := filepath.Join(t.TempDir(), "incoming.json")
	if err := os.WriteFile(path, payload, 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	importer := newShareTestEngine(t)
	page := NewQuoteImportPage(importer, 120, 40)
	page.Reset()
	page.pathInput.SetValue(path)

	model, cmd := page.Update(tea.KeyMsg{Type: tea.KeyEnter})
	page = model
	if cmd == nil {
		t.Fatal("import command = nil, want command")
	}

	model, _ = page.Update(cmd())
	page = model
	if page.result == nil {
		t.Fatal("import result = nil, want result")
	}
	if page.result.Inserted != 1 || page.result.Updated != 0 {
		t.Fatalf("import result = %+v, want inserted=1 updated=0", *page.result)
	}

	quotes, err := importer.ListQuotes(context.Background())
	if err != nil {
		t.Fatalf("ListQuotes() error = %v", err)
	}
	if len(quotes) != 1 || quotes[0].Content != "shared from alice" {
		t.Fatalf("imported quotes = %+v", quotes)
	}
	if quotes[0].SourceBackend != "local" || quotes[0].SourceEntityType != "quote" {
		t.Fatalf("imported quote source provenance = %+v", quotes[0])
	}
}
