package pages

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestQuoteEditorPreviewAcceptAndReject(t *testing.T) {
	t.Parallel()

	page := NewQuoteEditorPage(nil, 120, 40)
	page.Reset(QuoteEditorModeAdd, nil)
	page.textarea.SetValue("original draft")

	model, _ := page.Update(QuoteRefineDoneMsg{Refined: "refined draft"})
	page = model

	if !page.preview {
		t.Fatal("preview = false, want true")
	}
	if page.refined != "refined draft" {
		t.Fatalf("refined = %q, want refined draft", page.refined)
	}
	if page.textarea.Value() != "original draft" {
		t.Fatalf("textarea = %q, want original draft", page.textarea.Value())
	}
	view := page.View()
	if !containsAll(view, "Current Draft", "original draft", "Refined Draft", "refined draft") {
		t.Fatalf("preview comparison missing expected content:\n%s", view)
	}

	model, _ = page.Update(tea.KeyMsg{Type: tea.KeyEsc})
	page = model

	if page.preview {
		t.Fatal("preview = true after reject, want false")
	}
	if page.textarea.Value() != "original draft" {
		t.Fatalf("textarea after reject = %q, want original draft", page.textarea.Value())
	}

	model, _ = page.Update(QuoteRefineDoneMsg{Refined: "accepted draft"})
	page = model
	model, _ = page.Update(tea.KeyMsg{Type: tea.KeyEnter})
	page = model

	if page.preview {
		t.Fatal("preview = true after accept, want false")
	}
	if page.textarea.Value() != "accepted draft" {
		t.Fatalf("textarea after accept = %q, want accepted draft", page.textarea.Value())
	}
}

func TestQuoteEditorStagesAndRemovesImagePath(t *testing.T) {
	data, err := base64.StdEncoding.DecodeString("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNk+A8AAQUBAScY42YAAAAASUVORK5CYII=")
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "tiny.png")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}

	page := NewQuoteEditorPage(nil, 100, 35)
	page.Reset(QuoteEditorModeAdd, nil)
	model, _ := page.Update(tea.KeyMsg{Type: tea.KeyCtrlO})
	page = model
	if !page.attachmentMode {
		t.Fatal("attachment manager did not open")
	}
	page.attachmentInput.SetValue(path)
	model, _ = page.Update(tea.KeyMsg{Type: tea.KeyEnter})
	page = model
	if len(page.newImages) != 1 || page.newImages[0].Filename != "tiny.png" {
		t.Fatalf("staged images = %#v", page.newImages)
	}
	if !containsAll(page.View(), "tiny.png", "staged", "Attachments (1/5)") {
		t.Fatalf("attachment view missing metadata:\n%s", page.View())
	}
	model, _ = page.Update(tea.KeyMsg{Type: tea.KeyDelete})
	page = model
	if len(page.newImages) != 0 {
		t.Fatalf("images after delete = %#v", page.newImages)
	}
}

func containsAll(s string, parts ...string) bool {
	for _, part := range parts {
		if !strings.Contains(s, part) {
			return false
		}
	}
	return true
}
