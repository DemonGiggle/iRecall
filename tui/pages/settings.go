package pages

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/list"
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
	fieldMaxResults
	fieldMinRelevance
	fieldCount // sentinel
)

// SettingsPage manages LLM provider and search configuration.
type SettingsPage struct {
	engine *core.Engine

	inputs    [fieldCount]textinput.Model
	httpsOn   bool
	modelList list.Model
	models    []string

	focused   settingsField
	spinner   spinner.Model
	busy      bool
	statusMsg string
	isErr     bool

	width  int
	height int
}

type modelItem string

func (m modelItem) FilterValue() string { return string(m) }
func (m modelItem) Title() string       { return string(m) }
func (m modelItem) Description() string { return "" }

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

	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	ml := list.New(nil, delegate, width-8, 5)
	ml.SetShowHelp(false)
	ml.SetShowStatusBar(false)
	ml.SetFilteringEnabled(false)
	ml.Title = ""
	if s.Provider.Model != "" {
		ml.SetItems([]list.Item{modelItem(s.Provider.Model)})
	}

	return SettingsPage{
		engine:    engine,
		inputs:    inputs,
		httpsOn:   s.Provider.HTTPS,
		modelList: ml,
		focused:   fieldHost,
		spinner:   sp,
		width:     width,
		height:    height,
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
		case "tab", "shift+tab":
			// Tab cycles through fields. 'tab' handled at app level for page
			// switching only when no input is focused — we skip that here.
			// Handled at the app level; settings page only needs inner field focus.
			// Fall through to field cycling:
			dir := 1
			if msg.String() == "shift+tab" {
				dir = -1
			}
			p.cycleFocus(dir)

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
			items := make([]list.Item, len(msg.Models))
			for i, m := range msg.Models {
				items[i] = modelItem(m)
			}
			p.modelList.SetItems(items)
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

	// Update focused input.
	if p.isInputField(p.focused) {
		var cmd tea.Cmd
		p.inputs[p.focused], cmd = p.inputs[p.focused].Update(msg)
		cmds = append(cmds, cmd)
	}
	if p.focused == fieldFetchModels+1 { // model list
		var cmd tea.Cmd
		p.modelList, cmd = p.modelList.Update(msg)
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

	providerSection := lipgloss.JoinVertical(lipgloss.Left,
		styles.SectionHeader.Render("LLM Provider"),
		row("Host / IP", p.inputView(fieldHost)),
		row("Port", p.inputView(fieldPort)),
		row("HTTPS", httpsLabel),
		row("API Key", p.inputView(fieldAPIKey)),
		"",
		"  "+fetchBtn,
		"",
		row("Model", p.modelList.View()),
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

	helpLine := styles.HelpBar.Render("Tab/Shift+Tab: Move   Space: Toggle   Enter: Select   Ctrl+S: Save   Esc: Back")

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

func (p *SettingsPage) SetSize(width, height int) {
	p.width = width
	p.height = height
}

func (p *SettingsPage) LoadFrom(s *core.Settings) {
	p.inputs[fieldHost].SetValue(s.Provider.Host)
	p.inputs[fieldPort].SetValue(strconv.Itoa(s.Provider.Port))
	p.inputs[fieldAPIKey].SetValue(s.Provider.APIKey)
	p.inputs[fieldMaxResults].SetValue(strconv.Itoa(s.Search.MaxResults))
	p.inputs[fieldMinRelevance].SetValue(fmt.Sprintf("%.1f", s.Search.MinRelevance))
	p.httpsOn = s.Provider.HTTPS
}

// SelectedModel returns the currently selected model name (if any).
func (p *SettingsPage) SelectedModel() string {
	if sel := p.modelList.SelectedItem(); sel != nil {
		return string(sel.(modelItem))
	}
	return ""
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
		f == fieldMaxResults || f == fieldMinRelevance
}

func (p *SettingsPage) cycleFocus(dir int) {
	// Unfocus current.
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
