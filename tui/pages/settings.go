package pages

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gigol/irecall/core"
	"github.com/gigol/irecall/tui/styles"
)

// --- Messages ---

// SettingsSavedMsg signals a successful settings save.
type SettingsSavedMsg struct{}

// ModelsFetchedMsg carries the list of fetched model IDs.
type ModelsFetchedMsg struct {
	Models []string
	Err    error
}

// --- SettingsPage ---

type settingsField int

const (
	fieldHost settingsField = iota
	fieldPort
	fieldHTTPS
	fieldAPIKey
	fieldFetchModels
	fieldModelFilter
	fieldModel
	fieldTheme
	fieldMaxResults
	fieldMinRelevance
	fieldCount // sentinel
)

// SettingsPage manages LLM provider and search configuration.
type SettingsPage struct {
	engine *core.Engine

	inputs  [fieldCount]textinput.Model
	httpsOn bool

	models       []string // available model IDs
	modelIdx     int      // currently selected index (-1 = none)
	initialModel string   // model name from settings (before fetch)
	themes       []string
	themeIdx     int

	focused   settingsField
	spinner   spinner.Model
	busy      bool
	statusMsg string
	isErr     bool

	width  int
	height int
}

func NewSettingsPage(engine *core.Engine, width, height int, s *core.Settings) SettingsPage {
	makeInput := func(placeholder string, masked bool) textinput.Model {
		ti := textinput.New()
		ti.Placeholder = placeholder
		ti.CharLimit = 256
		if masked {
			ti.EchoMode = textinput.EchoPassword
			ti.EchoCharacter = '•'
		}
		return ti
	}

	var inputs [fieldCount]textinput.Model
	inputs[fieldHost] = makeInput("e.g. localhost", false)
	inputs[fieldPort] = makeInput("e.g. 11434", false)
	inputs[fieldAPIKey] = makeInput("optional", true)
	inputs[fieldModelFilter] = makeInput("type to filter", false)
	inputs[fieldMaxResults] = makeInput("1–20", false)
	inputs[fieldMinRelevance] = makeInput("0.0", false)

	// Populate from current settings.
	inputs[fieldHost].SetValue(s.Provider.Host)
	inputs[fieldPort].SetValue(strconv.Itoa(s.Provider.Port))
	inputs[fieldAPIKey].SetValue(s.Provider.APIKey)
	inputs[fieldMaxResults].SetValue(strconv.Itoa(s.Search.MaxResults))
	inputs[fieldMinRelevance].SetValue(fmt.Sprintf("%.1f", s.Search.MinRelevance))
	inputs[fieldHost].Focus()

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(styles.ColorAccent)

	return SettingsPage{
		engine:       engine,
		inputs:       inputs,
		httpsOn:      s.Provider.HTTPS,
		initialModel: s.Provider.Model,
		themes:       styles.ThemeNames(),
		themeIdx:     themeIndex(styles.ThemeNames(), s.Theme),
		modelIdx:     -1,
		focused:      fieldHost,
		spinner:      sp,
		width:        width,
		height:       height,
	}
}

func (p SettingsPage) Init() tea.Cmd {
	return textinput.Blink
}

func (p SettingsPage) Update(msg tea.Msg) (SettingsPage, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "down":
			p.cycleFocus(1)

		case "up":
			p.cycleFocus(-1)

		case " ":
			if p.focused == fieldHTTPS {
				p.httpsOn = !p.httpsOn
			}

		case "enter":
			if p.focused == fieldFetchModels {
				if p.busy {
					break
				}
				p.busy = true
				p.statusMsg = ""
				cmds = append(cmds, p.spinner.Tick, p.doFetchModels())
			}

		case "left":
			filtered := p.filteredModels()
			if p.focused == fieldModel && len(filtered) > 0 {
				current := p.SelectedModel()
				idx := p.filteredIndex(current)
				if idx < 0 {
					p.modelIdx = p.indexForModel(filtered[len(filtered)-1])
					break
				}
				idx--
				if idx < 0 {
					idx = len(filtered) - 1
				}
				p.modelIdx = p.indexForModel(filtered[idx])
			}
			if p.focused == fieldTheme && len(p.themes) > 0 {
				p.themeIdx--
				if p.themeIdx < 0 {
					p.themeIdx = len(p.themes) - 1
				}
				styles.ApplyTheme(p.SelectedTheme())
			}

		case "right":
			filtered := p.filteredModels()
			if p.focused == fieldModel && len(filtered) > 0 {
				current := p.SelectedModel()
				idx := p.filteredIndex(current)
				if idx < 0 {
					p.modelIdx = p.indexForModel(filtered[0])
					break
				}
				idx++
				if idx >= len(filtered) {
					idx = 0
				}
				p.modelIdx = p.indexForModel(filtered[idx])
			}
			if p.focused == fieldTheme && len(p.themes) > 0 {
				p.themeIdx++
				if p.themeIdx >= len(p.themes) {
					p.themeIdx = 0
				}
				styles.ApplyTheme(p.SelectedTheme())
			}

		case "ctrl+s":
			if p.busy {
				break
			}
			if err := p.save(); err != nil {
				p.statusMsg = "Error: " + err.Error()
				p.isErr = true
			} else {
				p.statusMsg = "Saved."
				p.isErr = false
			}
		}

	case ModelsFetchedMsg:
		p.busy = false
		if msg.Err != nil {
			p.statusMsg = "Fetch failed: " + msg.Err.Error()
			p.isErr = true
		} else if len(msg.Models) == 0 {
			p.statusMsg = "No models returned."
			p.isErr = false
		} else {
			p.models = msg.Models
			prev := p.SelectedModel()
			if prev == "" {
				prev = p.initialModel
			}
			p.syncModelSelection(prev)
			p.statusMsg = fmt.Sprintf("Fetched %d models.", len(msg.Models))
			p.isErr = false
		}

	case spinner.TickMsg:
		if p.busy {
			var cmd tea.Cmd
			p.spinner, cmd = p.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	// Update focused text input.
	if p.isInputField(p.focused) {
		var cmd tea.Cmd
		prevFilter := p.inputs[fieldModelFilter].Value()
		p.inputs[p.focused], cmd = p.inputs[p.focused].Update(msg)
		if p.focused == fieldModelFilter && p.inputs[fieldModelFilter].Value() != prevFilter {
			p.syncModelSelection(p.SelectedModel())
		}
		cmds = append(cmds, cmd)
	}

	return p, tea.Batch(cmds...)
}

func (p SettingsPage) View() string {
	row := func(label string, value string) string {
		return lipgloss.JoinHorizontal(lipgloss.Top,
			styles.FormLabel.Render(label),
			value,
		)
	}

	httpsLabel := "[ ] off"
	if p.httpsOn {
		httpsLabel = "[x] on"
	}
	if p.focused == fieldHTTPS {
		httpsLabel = styles.Accent.Render(httpsLabel) + styles.Muted.Render("  Space to toggle")
	}

	fetchBtn := styles.ButtonNormal.Render("Fetch Models")
	if p.focused == fieldFetchModels {
		fetchBtn = styles.ButtonFocused.Render("Fetch Models")
	}
	if p.busy {
		fetchBtn = p.spinner.View() + " Fetching..."
	}

	modelView := p.modelSelectorView()

	providerSection := lipgloss.JoinVertical(lipgloss.Left,
		styles.SectionHeader.Render("LLM Provider"),
		row("Host / IP", p.inputView(fieldHost)),
		row("Port", p.inputView(fieldPort)),
		row("HTTPS", httpsLabel),
		row("API Key", p.inputView(fieldAPIKey)),
		"",
		row("", fetchBtn),
		"",
		row("Filter", p.inputView(fieldModelFilter)),
		row("Model", modelView),
		row("Theme", p.themeSelectorView()),
	)

	searchSection := lipgloss.JoinVertical(lipgloss.Left,
		styles.SectionHeader.Render("Search"),
		row("Max ref quotes", p.inputView(fieldMaxResults)),
		row("Min relevance", p.inputView(fieldMinRelevance)),
	)

	var statusLine string
	if p.statusMsg != "" {
		if p.isErr {
			statusLine = styles.StatusErr.Render(p.statusMsg)
		} else {
			statusLine = styles.StatusOK.Render(p.statusMsg)
		}
	}

	helpLine := styles.HelpBar.Render("↑/↓: Move   type: Filter   ←/→: Cycle Model/Theme   space: Toggle   enter: Fetch   ctrl+s: Save   tab/shift+tab: Switch Page")

	return styles.Panel.Width(p.width - 4).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			providerSection,
			"",
			searchSection,
			"",
			statusLine,
			helpLine,
		),
	)
}

func (p *SettingsPage) modelSelectorView() string {
	if len(p.models) == 0 {
		name := p.initialModel
		if name == "" {
			name = "(none)"
		}
		if p.focused == fieldModel {
			return styles.Accent.Render(name) + styles.Muted.Render("  Fetch models first")
		}
		return styles.Muted.Render(name)
	}

	filtered := p.filteredModels()
	if len(filtered) == 0 {
		selected := p.SelectedModel()
		if selected == "" {
			selected = "(none)"
		}
		msg := "  No matches"
		if p.focused == fieldModel {
			return styles.Accent.Render(selected) + styles.Muted.Render(msg)
		}
		return selected + styles.Muted.Render(msg)
	}

	selected := p.SelectedModel()
	filteredIdx := p.filteredIndex(selected)
	if filteredIdx < 0 {
		filteredIdx = 0
		selected = filtered[0]
	}
	pos := fmt.Sprintf(" (%d/%d)", filteredIdx+1, len(filtered))

	if p.focused == fieldModel {
		return styles.Accent.Render("< "+selected+" >") +
			styles.Muted.Render(pos+"  ← / → to change")
	}
	return selected + styles.Muted.Render(pos)
}

func (p *SettingsPage) themeSelectorView() string {
	if len(p.themes) == 0 {
		return styles.Muted.Render("(none)")
	}
	name := p.SelectedTheme()
	pos := fmt.Sprintf(" (%d/%d)", p.themeIdx+1, len(p.themes))
	if p.focused == fieldTheme {
		return styles.Accent.Render("< "+name+" >") +
			styles.Muted.Render(pos+"  ← / → to change")
	}
	return name + styles.Muted.Render(pos)
}

func (p *SettingsPage) SetSize(width, height int) {
	p.width = width
	p.height = height
}

func (p *SettingsPage) LoadFrom(s *core.Settings) {
	p.inputs[fieldHost].SetValue(s.Provider.Host)
	p.inputs[fieldPort].SetValue(strconv.Itoa(s.Provider.Port))
	p.inputs[fieldAPIKey].SetValue(s.Provider.APIKey)
	p.inputs[fieldModelFilter].SetValue("")
	p.inputs[fieldMaxResults].SetValue(strconv.Itoa(s.Search.MaxResults))
	p.inputs[fieldMinRelevance].SetValue(fmt.Sprintf("%.1f", s.Search.MinRelevance))
	p.httpsOn = s.Provider.HTTPS
	p.initialModel = s.Provider.Model
	p.themeIdx = themeIndex(p.themes, s.Theme)
	styles.ApplyTheme(p.SelectedTheme())
	p.syncModelSelection(s.Provider.Model)
}

// SelectedModel returns the currently selected model name (if any).
func (p *SettingsPage) SelectedModel() string {
	if len(p.models) > 0 && p.modelIdx >= 0 && p.modelIdx < len(p.models) {
		return p.models[p.modelIdx]
	}
	return p.initialModel
}

func (p *SettingsPage) SelectedTheme() string {
	if len(p.themes) == 0 {
		return styles.CurrentThemeName()
	}
	if p.themeIdx < 0 || p.themeIdx >= len(p.themes) {
		return p.themes[0]
	}
	return p.themes[p.themeIdx]
}

// CurrentSettings builds a Settings from the form values.
func (p *SettingsPage) CurrentSettings() (*core.Settings, error) {
	port, err := strconv.Atoi(strings.TrimSpace(p.inputs[fieldPort].Value()))
	if err != nil || port < 1 || port > 65535 {
		return nil, fmt.Errorf("port must be a number between 1 and 65535")
	}
	maxResults, err := strconv.Atoi(strings.TrimSpace(p.inputs[fieldMaxResults].Value()))
	if err != nil || maxResults < 1 || maxResults > 20 {
		return nil, fmt.Errorf("max ref quotes must be between 1 and 20")
	}
	minRel, err := strconv.ParseFloat(strings.TrimSpace(p.inputs[fieldMinRelevance].Value()), 64)
	if err != nil {
		return nil, fmt.Errorf("min relevance must be a decimal number")
	}
	return &core.Settings{
		Provider: core.ProviderConfig{
			Host:   strings.TrimSpace(p.inputs[fieldHost].Value()),
			Port:   port,
			HTTPS:  p.httpsOn,
			APIKey: p.inputs[fieldAPIKey].Value(),
			Model:  p.SelectedModel(),
		},
		Search: core.SearchConfig{
			MaxResults:   maxResults,
			MinRelevance: minRel,
		},
		Theme: p.SelectedTheme(),
	}, nil
}

func (p *SettingsPage) save() error {
	s, err := p.CurrentSettings()
	if err != nil {
		return err
	}
	return p.engine.SaveSettings(context.Background(), s)
}

func (p *SettingsPage) doFetchModels() tea.Cmd {
	engine := p.engine
	host := p.inputs[fieldHost].Value()
	portStr := p.inputs[fieldPort].Value()
	apiKey := p.inputs[fieldAPIKey].Value()
	https := p.httpsOn
	return func() tea.Msg {
		port, _ := strconv.Atoi(portStr)
		cfg := core.ProviderConfig{
			Host:   host,
			Port:   port,
			HTTPS:  https,
			APIKey: apiKey,
		}
		models, err := engine.FetchModels(context.Background(), cfg)
		return ModelsFetchedMsg{Models: models, Err: err}
	}
}

func (p *SettingsPage) isInputField(f settingsField) bool {
	return f == fieldHost || f == fieldPort || f == fieldAPIKey ||
		f == fieldModelFilter || f == fieldMaxResults || f == fieldMinRelevance
}

func (p *SettingsPage) cycleFocus(dir int) {
	if p.isInputField(p.focused) {
		p.inputs[p.focused].Blur()
	}
	p.focused = settingsField((int(p.focused) + dir + int(fieldCount)) % int(fieldCount))
	if p.isInputField(p.focused) {
		p.inputs[p.focused].Focus()
	}
}

func (p *SettingsPage) inputView(f settingsField) string {
	return p.inputs[f].View()
}

func themeIndex(themes []string, name string) int {
	if len(themes) == 0 {
		return -1
	}
	for i, theme := range themes {
		if theme == name {
			return i
		}
	}
	return 0
}

func (p *SettingsPage) filteredModels() []string {
	filter := strings.ToLower(strings.TrimSpace(p.inputs[fieldModelFilter].Value()))
	if filter == "" {
		return p.models
	}
	filtered := make([]string, 0, len(p.models))
	for _, model := range p.models {
		if strings.Contains(strings.ToLower(model), filter) {
			filtered = append(filtered, model)
		}
	}
	return filtered
}

func (p *SettingsPage) filteredIndex(model string) int {
	for i, candidate := range p.filteredModels() {
		if candidate == model {
			return i
		}
	}
	return -1
}

func (p *SettingsPage) indexForModel(model string) int {
	for i, candidate := range p.models {
		if candidate == model {
			return i
		}
	}
	return -1
}

func (p *SettingsPage) syncModelSelection(preferred string) {
	if len(p.models) == 0 {
		p.modelIdx = -1
		return
	}
	if preferred != "" {
		if idx := p.indexForModel(preferred); idx >= 0 && p.filteredIndex(preferred) >= 0 {
			p.modelIdx = idx
			return
		}
	}
	filtered := p.filteredModels()
	if len(filtered) == 0 {
		return
	}
	p.modelIdx = p.indexForModel(filtered[0])
}
