package pages

import (
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
