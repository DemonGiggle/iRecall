package pages

import (
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gigol/irecall/core"
	"github.com/gigol/irecall/tui/styles"
)

func TestSettingsPageFilterNarrowsModelSelection(t *testing.T) {
	page := NewSettingsPage(nil, 120, 40, core.DefaultSettings())
	page.models = []string{"gpt-4o", "gpt-4.1-mini", "llama3.2"}
	page.responseModel = "gpt-4o"
	page.focused = fieldModelFilter
	page.inputs[fieldModelFilter].Focus()

	model, _ := page.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("mini")})
	page = model

	if got := page.SelectedResponseModel(); got != "gpt-4.1-mini" {
		t.Fatalf("SelectedResponseModel() = %q, want gpt-4.1-mini", got)
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

	if got := page.SelectedResponseModel(); got != "gpt-4.1-mini" {
		t.Fatalf("SelectedResponseModel() after fetch = %q, want gpt-4.1-mini", got)
	}
}

func TestSettingsPageFilterNoMatchesKeepsExistingSelection(t *testing.T) {
	settings := core.DefaultSettings()
	settings.Provider.Model = "gpt-4o"
	page := NewSettingsPage(nil, 120, 40, settings)
	page.models = []string{"gpt-4o", "gpt-4.1-mini"}
	page.responseModel = "gpt-4o"
	page.focused = fieldModelFilter
	page.inputs[fieldModelFilter].Focus()

	model, _ := page.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("zzz")})
	page = model

	if got := page.SelectedResponseModel(); got != "gpt-4o" {
		t.Fatalf("SelectedResponseModel() with no filter matches = %q, want gpt-4o", got)
	}
	if got := len(page.filteredModels()); got != 0 {
		t.Fatalf("len(filteredModels()) = %d, want 0", got)
	}
}

func TestSettingsPageKeywordModelPreservesResponseFallback(t *testing.T) {
	settings := core.DefaultSettings()
	settings.Provider.Model = "gpt-4o"
	settings.Provider.KeywordModel = ""
	page := NewSettingsPage(nil, 120, 40, settings)
	page.models = []string{"gpt-4o", "gpt-4.1-mini"}
	page.focused = fieldModelFilter
	page.inputs[fieldModelFilter].Focus()

	model, _ := page.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("mini")})
	page = model

	if got := page.SelectedKeywordModel(); got != "" {
		t.Fatalf("SelectedKeywordModel() = %q, want response-model fallback", got)
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

func TestSettingsPageCurrentSettingsIncludesKeywordModel(t *testing.T) {
	settings := core.DefaultSettings()
	settings.Provider.Model = "gpt-4.1"
	settings.Provider.KeywordModel = "gpt-4.1-mini"
	page := NewSettingsPage(nil, 120, 40, settings)

	current, err := page.CurrentSettings()
	if err != nil {
		t.Fatalf("CurrentSettings() error = %v", err)
	}
	if current.Provider.Model != "gpt-4.1" {
		t.Fatalf("CurrentSettings().Provider.Model = %q, want gpt-4.1", current.Provider.Model)
	}
	if current.Provider.KeywordModel != "gpt-4.1-mini" {
		t.Fatalf("CurrentSettings().Provider.KeywordModel = %q, want gpt-4.1-mini", current.Provider.KeywordModel)
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

func TestSettingsPageShowsScrollbarWhenContentOverflows(t *testing.T) {
	page := NewSettingsPage(nil, 80, 12, core.DefaultSettings())

	if page.viewport.TotalLineCount() <= page.viewport.VisibleLineCount() {
		t.Fatalf("settings page should overflow at small height: total=%d visible=%d", page.viewport.TotalLineCount(), page.viewport.VisibleLineCount())
	}

	view := page.View()
	if !strings.Contains(view, "█") || !strings.Contains(view, "│") {
		t.Fatalf("settings view missing scrollbar markers:\n%s", view)
	}
}

func TestSettingsPageKeyNavigationMatchesRenderedOrder(t *testing.T) {
	page := NewSettingsPage(nil, 120, 40, core.DefaultSettings())
	page.focused = fieldTheme

	for _, want := range []settingsField{fieldMaxResults, fieldMinRelevance, fieldMockLLM, fieldRootDir} {
		model, _ := page.Update(tea.KeyMsg{Type: tea.KeyDown})
		page = model
		if page.focused != want {
			t.Fatalf("focused field after down = %v, want %v", page.focused, want)
		}
	}

	model, _ := page.Update(tea.KeyMsg{Type: tea.KeyUp})
	page = model
	if page.focused != fieldMockLLM {
		t.Fatalf("focused field after up = %v, want %v", page.focused, fieldMockLLM)
	}
}

func TestSettingsPageKeepsFocusedFieldVisibleInViewport(t *testing.T) {
	page := NewSettingsPage(nil, 80, 12, core.DefaultSettings())

	for steps := 0; steps < int(fieldCount) && page.focused != fieldRootDir; steps++ {
		model, _ := page.Update(tea.KeyMsg{Type: tea.KeyDown})
		page = model
	}

	if page.focused != fieldRootDir {
		t.Fatalf("focused field = %v, want %v", page.focused, fieldRootDir)
	}
	if page.viewport.YOffset == 0 {
		t.Fatalf("viewport y offset = %d, want > 0 after moving focus down", page.viewport.YOffset)
	}
	if !strings.Contains(page.View(), "Config root") {
		t.Fatalf("settings view missing focused lower field:\n%s", page.View())
	}
}

func TestSettingsPageViewportHeightUsesRenderedFooter(t *testing.T) {
	page := NewSettingsPage(nil, 80, 12, core.DefaultSettings())

	panelHeight := styles.Panel.GetVerticalFrameSize()
	footerHeight := lipgloss.Height(page.footerView())
	if got, want := page.viewport.Height, max(5, page.height-panelHeight-footerHeight); got != want {
		t.Fatalf("viewport height = %d, want %d", got, want)
	}

	page.statusMsg = "this is a long status message that should wrap onto multiple lines in a narrow viewport"
	page.isErr = true
	page.recalcViewport()

	footerHeight = lipgloss.Height(page.footerView())
	if got, want := page.viewport.Height, max(5, page.height-panelHeight-footerHeight); got != want {
		t.Fatalf("viewport height after wrapped status = %d, want %d", got, want)
	}
}

func TestSettingsPageTallFocusedBlockAlignsToTop(t *testing.T) {
	page := NewSettingsPage(nil, 48, 12, core.DefaultSettings())
	page.viewport.SetYOffset(20)

	focusStart := 8
	focusEnd := focusStart + page.viewport.Height + 2
	page.ensureVisible(focusStart, focusEnd)

	if page.viewport.YOffset > focusStart {
		t.Fatalf("viewport y offset = %d, want <= focus start %d for tall focused block", page.viewport.YOffset, focusStart)
	}
}
