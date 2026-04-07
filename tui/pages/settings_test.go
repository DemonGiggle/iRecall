package pages

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gigol/irecall/core"
)

func TestSettingsPageFilterNarrowsModelSelection(t *testing.T) {
	page := NewSettingsPage(nil, 120, 40, core.DefaultSettings())
	page.models = []string{"gpt-4o", "gpt-4.1-mini", "llama3.2"}
	page.modelIdx = 0
	page.focused = fieldModelFilter
	page.inputs[fieldModelFilter].Focus()

	model, _ := page.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("mini")})
	page = model

	if got := page.SelectedModel(); got != "gpt-4.1-mini" {
		t.Fatalf("SelectedModel() = %q, want gpt-4.1-mini", got)
	}
	if got := page.filteredModels(); len(got) != 1 || got[0] != "gpt-4.1-mini" {
		t.Fatalf("filteredModels() = %v, want [gpt-4.1-mini]", got)
	}
}

func TestSettingsPageFetchPreservesMatchingSelection(t *testing.T) {
	settings := core.DefaultSettings()
	settings.Provider.Model = "gpt-4.1-mini"
	page := NewSettingsPage(nil, 120, 40, settings)

	model, _ := page.Update(ModelsFetchedMsg{Models: []string{"gpt-4o", "gpt-4.1-mini", "llama3.2"}})
	page = model

	if got := page.SelectedModel(); got != "gpt-4.1-mini" {
		t.Fatalf("SelectedModel() after fetch = %q, want gpt-4.1-mini", got)
	}
}

func TestSettingsPageFilterNoMatchesKeepsExistingSelection(t *testing.T) {
	settings := core.DefaultSettings()
	settings.Provider.Model = "gpt-4o"
	page := NewSettingsPage(nil, 120, 40, settings)
	page.models = []string{"gpt-4o", "gpt-4.1-mini"}
	page.modelIdx = 0
	page.focused = fieldModelFilter
	page.inputs[fieldModelFilter].Focus()

	model, _ := page.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("zzz")})
	page = model

	if got := page.SelectedModel(); got != "gpt-4o" {
		t.Fatalf("SelectedModel() with no filter matches = %q, want gpt-4o", got)
	}
	if got := len(page.filteredModels()); got != 0 {
		t.Fatalf("len(filteredModels()) = %d, want 0", got)
	}
}
