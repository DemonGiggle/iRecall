package pages

import (
	"context"
	"fmt"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gigol/irecall/config"
	"github.com/gigol/irecall/core"
	"github.com/gigol/irecall/tui/styles"
)

// --- Messages ---

// SettingsSavedMsg reports the result of a settings save orchestrated by the app shell.
type SettingsSavedMsg struct {
	Settings     *core.Settings
	Err          error
	SwitchedRoot bool
}

// SaveSettingsRequestedMsg asks the app shell to persist and apply settings.
type SaveSettingsRequestedMsg struct {
	Settings *core.Settings
}

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
	fieldResponseModel
	fieldKeywordModel
	fieldTheme
	fieldRootDir
	fieldMaxResults
	fieldMinRelevance
	fieldMockLLM
	fieldCount // sentinel
)

// SettingsPage manages LLM provider and search configuration.
type SettingsPage struct {
	engine *core.Engine
	web    core.WebConfig
	debug  core.DebugConfig

	inputs  [fieldCount]textinput.Model
	httpsOn bool

	models        []string // available model IDs
	responseModel string
	keywordModel  string
	themes        []string
	themeIdx      int

	focused   settingsField
	spinner   spinner.Model
	busy      bool
	statusMsg string
	isErr     bool
	viewport  viewport.Model

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
	inputs[fieldRootDir] = makeInput("leave blank for default XDG/AppData folders", false)
	inputs[fieldMaxResults] = makeInput("1–20", false)
	inputs[fieldMinRelevance] = makeInput("0.0-1.0", false)

	// Populate from current settings.
	inputs[fieldHost].SetValue(s.Provider.Host)
	inputs[fieldPort].SetValue(strconv.Itoa(s.Provider.Port))
	inputs[fieldAPIKey].SetValue(s.Provider.APIKey)
	inputs[fieldRootDir].SetValue(s.RootDir)
	inputs[fieldMaxResults].SetValue(strconv.Itoa(s.Search.MaxResults))
	inputs[fieldMinRelevance].SetValue(fmt.Sprintf("%.1f", s.Search.MinRelevance))
	inputs[fieldHost].Focus()

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(styles.ColorAccent)

	page := SettingsPage{
		engine:        engine,
		inputs:        inputs,
		httpsOn:       s.Provider.HTTPS,
		responseModel: s.Provider.Model,
		keywordModel:  s.Provider.KeywordModel,
		web:           s.Web,
		debug:         s.Debug,
		themes:        styles.ThemeNames(),
		themeIdx:      themeIndex(styles.ThemeNames(), s.Theme),
		focused:       fieldHost,
		spinner:       sp,
		width:         width,
		height:        height,
	}
	page.recalcViewport()
	page.refreshViewport()
	return page
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

		case "pgdown":
			p.viewport.LineDown(max(1, p.viewport.Height-1))

		case "pgup":
			p.viewport.LineUp(max(1, p.viewport.Height-1))

		case "home":
			p.viewport.GotoTop()

		case "end":
			p.viewport.GotoBottom()

		case " ":
			if p.focused == fieldHTTPS {
				p.httpsOn = !p.httpsOn
			}
			if p.focused == fieldMockLLM {
				p.debug.MockLLM = !p.debug.MockLLM
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
			if p.focused == fieldResponseModel {
				p.cycleSelectedModel(-1, false)
			}
			if p.focused == fieldKeywordModel {
				p.cycleSelectedModel(-1, true)
			}
			if p.focused == fieldTheme && len(p.themes) > 0 {
				p.themeIdx--
				if p.themeIdx < 0 {
					p.themeIdx = len(p.themes) - 1
				}
				styles.ApplyTheme(p.SelectedTheme())
			}

		case "right":
			if p.focused == fieldResponseModel {
				p.cycleSelectedModel(1, false)
			}
			if p.focused == fieldKeywordModel {
				p.cycleSelectedModel(1, true)
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
			s, err := p.CurrentSettings()
			if err != nil {
				p.statusMsg = "Error: " + err.Error()
				p.isErr = true
			} else {
				return p, requestSaveSettingsCmd(s)
			}
		}

	case SettingsSavedMsg:
		if msg.Err != nil {
			p.statusMsg = "Error: " + msg.Err.Error()
			p.isErr = true
			break
		}
		if msg.Settings != nil {
			p.LoadFrom(msg.Settings)
		}
		p.statusMsg = "Saved."
		if msg.SwitchedRoot {
			p.statusMsg = "Saved. Switched storage root."
		}
		p.isErr = false

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
			p.syncSelectedModel(false)
			p.syncSelectedModel(true)
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
			p.syncSelectedModel(false)
			p.syncSelectedModel(true)
		}
		cmds = append(cmds, cmd)
	}

	p.refreshViewport()

	return p, tea.Batch(cmds...)
}

func (p SettingsPage) View() string {
	footer := p.footerView()
	body := lipgloss.JoinVertical(lipgloss.Left,
		p.settingsViewportView(),
		footer,
	)

	return styles.Panel.Width(max(24, p.width-4)).Render(body)
}

func (p SettingsPage) footerView() string {
	var statusLine string
	if p.statusMsg != "" {
		if p.isErr {
			statusLine = styles.StatusErr.Render(p.statusMsg)
		} else {
			statusLine = styles.StatusOK.Render(p.statusMsg)
		}
	}

	helpLinePrimary := styles.HelpBar.Render("↑/↓ move   type edit/filter   ←/→ cycle model/theme   space toggle")
	helpLineSecondary := styles.HelpBar.Render("enter fetch   ctrl+s save   pgup/pgdn/home/end scroll   tab/shift+tab page")
	return lipgloss.JoinVertical(lipgloss.Left,
		statusLine,
		helpLinePrimary,
		helpLineSecondary,
	)
}

func (p *SettingsPage) modelSelectorView(selected string, allowFallback bool, focused bool) string {
	emptyLabel := "(none)"
	if allowFallback {
		emptyLabel = "(use response model)"
	}
	if len(p.models) == 0 {
		name := selected
		if name == "" {
			name = emptyLabel
		}
		if focused {
			return styles.Accent.Render(name) + styles.Muted.Render("  Fetch models first")
		}
		return styles.Muted.Render(name)
	}

	filtered := p.filteredModels()
	if len(filtered) == 0 {
		if allowFallback && selected == "" {
			if focused {
				return styles.Accent.Render("< "+emptyLabel+" >") +
					styles.Muted.Render(" (1/1)")
			}
			return emptyLabel + styles.Muted.Render(" (1/1)")
		}
		if selected == "" {
			selected = emptyLabel
		}
		msg := "  No matches"
		if focused {
			return styles.Accent.Render(selected) + styles.Muted.Render(msg)
		}
		return selected + styles.Muted.Render(msg)
	}

	total := len(filtered)
	position := 0
	if allowFallback {
		total++
		if selected == "" {
			position = 1
			selected = emptyLabel
		}
	}
	if selected != emptyLabel {
		filteredIdx := p.filteredIndex(selected)
		if filteredIdx < 0 {
			filteredIdx = 0
			selected = filtered[0]
		}
		position = filteredIdx + 1
		if allowFallback {
			position++
		}
	}
	pos := fmt.Sprintf(" (%d/%d)", position, total)

	if focused {
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
	p.recalcViewport()
	p.refreshViewport()
}

func (p *SettingsPage) LoadFrom(s *core.Settings) {
	p.inputs[fieldHost].SetValue(s.Provider.Host)
	p.inputs[fieldPort].SetValue(strconv.Itoa(s.Provider.Port))
	p.inputs[fieldAPIKey].SetValue(s.Provider.APIKey)
	p.inputs[fieldModelFilter].SetValue("")
	p.inputs[fieldRootDir].SetValue(s.RootDir)
	p.inputs[fieldMaxResults].SetValue(strconv.Itoa(s.Search.MaxResults))
	p.inputs[fieldMinRelevance].SetValue(fmt.Sprintf("%.1f", s.Search.MinRelevance))
	p.httpsOn = s.Provider.HTTPS
	p.responseModel = s.Provider.Model
	p.keywordModel = s.Provider.KeywordModel
	p.web = s.Web
	p.debug = s.Debug
	p.themeIdx = themeIndex(p.themes, s.Theme)
	styles.ApplyTheme(p.SelectedTheme())
	p.syncSelectedModel(false)
	p.syncSelectedModel(true)
	p.refreshViewport()
}

func (p *SettingsPage) SelectedResponseModel() string {
	return p.responseModel
}

func (p *SettingsPage) SelectedKeywordModel() string {
	return p.keywordModel
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
	if minRel < 0 || minRel > 1 {
		return nil, fmt.Errorf("min relevance must be between 0.0 and 1.0")
	}
	return &core.Settings{
		Provider: core.ProviderConfig{
			Host:         strings.TrimSpace(p.inputs[fieldHost].Value()),
			Port:         port,
			HTTPS:        p.httpsOn,
			APIKey:       p.inputs[fieldAPIKey].Value(),
			Model:        p.SelectedResponseModel(),
			KeywordModel: p.SelectedKeywordModel(),
		},
		Search: core.SearchConfig{
			MaxResults:   maxResults,
			MinRelevance: minRel,
		},
		Debug:   p.debug,
		Theme:   p.SelectedTheme(),
		Web:     p.web,
		RootDir: strings.TrimSpace(p.inputs[fieldRootDir].Value()),
	}, nil
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
			Host:         host,
			Port:         port,
			HTTPS:        https,
			APIKey:       apiKey,
			Model:        p.SelectedResponseModel(),
			KeywordModel: p.SelectedKeywordModel(),
		}
		models, err := engine.FetchModels(context.Background(), cfg)
		return ModelsFetchedMsg{Models: models, Err: err}
	}
}

func (p *SettingsPage) isInputField(f settingsField) bool {
	return f == fieldHost || f == fieldPort || f == fieldAPIKey ||
		f == fieldModelFilter || f == fieldRootDir || f == fieldMaxResults || f == fieldMinRelevance
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

func (p *SettingsPage) syncSelectedModel(allowFallback bool) {
	if len(p.models) == 0 {
		return
	}
	current := p.SelectedResponseModel()
	if allowFallback {
		current = p.SelectedKeywordModel()
		if strings.TrimSpace(current) == "" {
			return
		}
	}
	filtered := p.filteredModels()
	if len(filtered) == 0 {
		return
	}
	if slices.Contains(filtered, current) {
		return
	}
	if allowFallback {
		p.keywordModel = filtered[0]
		return
	}
	p.responseModel = filtered[0]
}

func (p *SettingsPage) cycleSelectedModel(dir int, allowFallback bool) {
	options := slices.Clone(p.filteredModels())
	if allowFallback {
		options = append([]string{""}, options...)
	}
	if len(options) == 0 {
		return
	}

	current := p.SelectedResponseModel()
	if allowFallback {
		current = p.SelectedKeywordModel()
	}
	idx := -1
	for i, option := range options {
		if option == current {
			idx = i
			break
		}
	}
	if idx < 0 {
		if dir < 0 {
			idx = len(options) - 1
		} else {
			idx = 0
		}
	} else {
		idx = (idx + dir + len(options)) % len(options)
	}

	if allowFallback {
		p.keywordModel = options[idx]
		return
	}
	p.responseModel = options[idx]
}

func requestSaveSettingsCmd(settings *core.Settings) tea.Cmd {
	return func() tea.Msg {
		return SaveSettingsRequestedMsg{Settings: settings}
	}
}

func (p *SettingsPage) previewDataDir() string {
	root := strings.TrimSpace(p.inputs[fieldRootDir].Value())
	if root == "" {
		return config.DefaultDataDir()
	}
	return filepath.Join(root, "data")
}

func (p *SettingsPage) previewConfigDir() string {
	root := strings.TrimSpace(p.inputs[fieldRootDir].Value())
	if root == "" {
		return config.DefaultConfigDir()
	}
	return filepath.Join(root, "config")
}

func (p *SettingsPage) previewStateDir() string {
	root := strings.TrimSpace(p.inputs[fieldRootDir].Value())
	if root == "" {
		return config.DefaultStateDir()
	}
	return filepath.Join(root, "state")
}

func (p *SettingsPage) recalcViewport() {
	panelWidth := max(24, p.width-4)
	panelFrameWidth := styles.Panel.GetHorizontalFrameSize()
	panelFrameHeight := styles.Panel.GetVerticalFrameSize()
	footerHeight := lipgloss.Height(p.footerView())

	p.viewport.Width = max(1, panelWidth-panelFrameWidth-2)
	p.viewport.Height = max(5, p.height-panelFrameHeight-footerHeight)
}

func (p *SettingsPage) refreshViewport() {
	content, focusStart, focusEnd := p.renderContent()
	p.viewport.SetContent(content)
	p.ensureVisible(focusStart, focusEnd)
}

func (p *SettingsPage) renderContent() (string, int, int) {
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

	mockLLMLabel := "[ ] off"
	if p.debug.MockLLM {
		mockLLMLabel = "[x] on"
	}
	if p.focused == fieldMockLLM {
		mockLLMLabel = styles.Accent.Render(mockLLMLabel) + styles.Muted.Render("  Space to toggle")
	}

	type block struct {
		field   settingsField
		focused bool
		content string
	}

	blocks := []block{
		{content: styles.SectionHeader.Render("LLM Provider")},
		{field: fieldHost, focused: p.focused == fieldHost, content: row("Host / IP", p.inputView(fieldHost))},
		{field: fieldPort, focused: p.focused == fieldPort, content: row("Port", p.inputView(fieldPort))},
		{field: fieldHTTPS, focused: p.focused == fieldHTTPS, content: row("HTTPS", httpsLabel)},
		{field: fieldAPIKey, focused: p.focused == fieldAPIKey, content: row("API Key", p.inputView(fieldAPIKey))},
		{content: ""},
		{field: fieldFetchModels, focused: p.focused == fieldFetchModels, content: row("", fetchBtn)},
		{content: ""},
		{field: fieldModelFilter, focused: p.focused == fieldModelFilter, content: row("Filter", p.inputView(fieldModelFilter))},
		{field: fieldResponseModel, focused: p.focused == fieldResponseModel, content: row("Response model", p.modelSelectorView(p.SelectedResponseModel(), false, p.focused == fieldResponseModel))},
		{field: fieldKeywordModel, focused: p.focused == fieldKeywordModel, content: row("Keyword model", p.modelSelectorView(p.SelectedKeywordModel(), true, p.focused == fieldKeywordModel))},
		{field: fieldTheme, focused: p.focused == fieldTheme, content: row("Theme", p.themeSelectorView())},
		{content: ""},
		{content: styles.SectionHeader.Render("Search")},
		{field: fieldMaxResults, focused: p.focused == fieldMaxResults, content: row("Max ref quotes", p.inputView(fieldMaxResults))},
		{field: fieldMinRelevance, focused: p.focused == fieldMinRelevance, content: row("Min relevance", p.inputView(fieldMinRelevance))},
		{content: styles.Muted.Render("0.0 keeps broad matches. Try 0.3-0.7 for cleaner results; 1.0 is very strict.")},
		{content: ""},
		{content: styles.SectionHeader.Render("Debug")},
		{field: fieldMockLLM, focused: p.focused == fieldMockLLM, content: row("Mock LLM", mockLLMLabel)},
		{content: styles.Muted.Render("Refine returns the original text, keywords split on spaces, and answers combine reference quotes.")},
		{content: ""},
		{content: styles.SectionHeader.Render("Local Storage")},
		{field: fieldRootDir, focused: p.focused == fieldRootDir, content: row("Config root", p.inputView(fieldRootDir))},
		{content: styles.Muted.Render("Leave blank for the default XDG/AppData locations. Saving switches iRecall to that root immediately.")},
		{content: row("Data dir", styles.Muted.Render(p.previewDataDir()))},
		{content: row("Config dir", styles.Muted.Render(p.previewConfigDir()))},
		{content: row("State dir", styles.Muted.Render(p.previewStateDir()))},
	}

	lines := make([]string, 0, len(blocks))
	focusStart := 0
	focusEnd := 0
	lineCount := 0
	for _, block := range blocks {
		blockLines := strings.Split(block.content, "\n")
		if block.focused {
			focusStart = lineCount
			focusEnd = lineCount + len(blockLines) - 1
		}
		lines = append(lines, blockLines...)
		lineCount += len(blockLines)
	}

	return strings.Join(lines, "\n"), focusStart, focusEnd
}

func (p *SettingsPage) ensureVisible(start, end int) {
	if p.viewport.Height <= 0 {
		return
	}
	blockHeight := end - start + 1
	padding := 1
	if blockHeight >= p.viewport.Height {
		p.viewport.SetYOffset(max(0, start-padding))
		return
	}
	if start <= p.viewport.YOffset+padding {
		p.viewport.SetYOffset(max(0, start-padding))
		return
	}
	bottom := p.viewport.YOffset + p.viewport.Height - 1
	if end >= bottom-padding {
		p.viewport.SetYOffset(end - p.viewport.Height + 1 + padding)
	}
}

func (p SettingsPage) settingsViewportView() string {
	scrollbar := p.scrollbarView()
	if scrollbar == "" {
		return p.viewport.View()
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, p.viewport.View(), " ", scrollbar)
}

func (p SettingsPage) scrollbarView() string {
	total := p.viewport.TotalLineCount()
	visible := p.viewport.VisibleLineCount()
	if total <= visible || visible <= 0 {
		return ""
	}

	trackStyle := lipgloss.NewStyle().Foreground(styles.ColorBorder)
	thumbStyle := lipgloss.NewStyle().Foreground(styles.ColorPrimary)
	track := make([]string, visible)
	for i := range track {
		track[i] = trackStyle.Render("│")
	}

	thumbHeight := max(1, visible*visible/total)
	maxThumbTop := visible - thumbHeight
	thumbTop := 0
	if maxThumbTop > 0 {
		thumbTop = p.viewport.YOffset * maxThumbTop / max(1, total-visible)
	}
	for i := thumbTop; i < thumbTop+thumbHeight && i < len(track); i++ {
		track[i] = thumbStyle.Render("█")
	}

	return strings.Join(track, "\n")
}
