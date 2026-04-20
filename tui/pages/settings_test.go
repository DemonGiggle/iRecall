package pages

import (
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gigol/irecall/core"
	"github.com/gigol/irecall/tui/styles"
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

func TestSettingsPageThemeSelectionUpdatesCurrentSettingsAndPreview(t *testing.T) {
	settings := core.DefaultSettings()
	settings.Theme = "violet"
	page := NewSettingsPage(nil, 120, 40, settings)
	page.focused = fieldTheme

	model, _ := page.Update(tea.KeyMsg{Type: tea.KeyRight})
	page = model

	if got := page.SelectedTheme(); got != "forest" {
		t.Fatalf("SelectedTheme() = %q, want forest", got)
	}
	if got := styles.CurrentThemeName(); got != "forest" {
		t.Fatalf("CurrentThemeName() = %q, want forest", got)
	}
	current, err := page.CurrentSettings()
	if err != nil {
		t.Fatalf("CurrentSettings() error = %v", err)
	}
	if current.Theme != "forest" {
		t.Fatalf("CurrentSettings().Theme = %q, want forest", current.Theme)
	}
}

func TestSettingsPageShowsStoragePaths(t *testing.T) {
	settings := core.DefaultSettings()
	settings.RootDir = "/tmp/irecall-test"
	page := NewSettingsPage(nil, 120, 40, settings)
	view := page.View()

	for _, want := range []string{
		"Local Storage",
		filepath.FromSlash("/tmp/irecall-test/data"),
		filepath.FromSlash("/tmp/irecall-test/config"),
		filepath.FromSlash("/tmp/irecall-test/state"),
	} {
		if !strings.Contains(view, want) {
			t.Fatalf("settings view missing %q:\n%s", want, view)
		}
	}
}

func TestSettingsPageCurrentSettingsIncludesRootDir(t *testing.T) {
	page := NewSettingsPage(nil, 120, 40, core.DefaultSettings())
	page.inputs[fieldRootDir].SetValue("/tmp/irecall-alt")

	current, err := page.CurrentSettings()
	if err != nil {
		t.Fatalf("CurrentSettings() error = %v", err)
	}
	if current.RootDir != "/tmp/irecall-alt" {
		t.Fatalf("CurrentSettings().RootDir = %q, want /tmp/irecall-alt", current.RootDir)
	}
}

func TestSettingsPageCurrentSettingsRejectsMinRelevanceOutsideRange(t *testing.T) {
	page := NewSettingsPage(nil, 120, 40, core.DefaultSettings())
	page.inputs[fieldMinRelevance].SetValue("1.2")

	_, err := page.CurrentSettings()
	if err == nil || !strings.Contains(err.Error(), "between 0.0 and 1.0") {
		t.Fatalf("CurrentSettings() error = %v, want range validation", err)
	}
}

func TestSettingsPageMockLLMToggleUpdatesCurrentSettings(t *testing.T) {
	page := NewSettingsPage(nil, 120, 40, core.DefaultSettings())
	page.focused = fieldMockLLM

	model, _ := page.Update(tea.KeyMsg{Type: tea.KeySpace})
	page = model

	current, err := page.CurrentSettings()
	if err != nil {
		t.Fatalf("CurrentSettings() error = %v", err)
	}
	if !current.Debug.MockLLM {
		t.Fatal("CurrentSettings().Debug.MockLLM = false, want true")
	}
	if !strings.Contains(page.View(), "Mock LLM") {
		t.Fatalf("settings view missing debug control:\n%s", page.View())
	}
}
