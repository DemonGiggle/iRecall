package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gigol/irecall/core"
	"github.com/gigol/irecall/tui/pages"
	"github.com/gigol/irecall/tui/styles"
)

type activePage int

const (
	pageRecall activePage = iota
	pageSettings
)

type overlayKind int

const (
	overlayNone overlayKind = iota
	overlayAddQuote
)

// App is the root Bubbletea model. It owns page routing and global key handling.
type App struct {
	engine   *core.Engine
	settings *core.Settings

	page    activePage
	overlay overlayKind

	recall   pages.RecallPage
	settings_ pages.SettingsPage
	addQuote pages.AddQuotePage

	width  int
	height int
}

func NewApp(engine *core.Engine, settings *core.Settings, width, height int) App {
	return App{
		engine:   engine,
		settings: settings,
		page:     pageRecall,
		overlay:  overlayNone,
		recall:   pages.NewRecallPage(engine, width, height-3),
		settings_: pages.NewSettingsPage(engine, width, height-3, settings),
		addQuote: pages.NewAddQuotePage(engine, width, height),
		width:    width,
		height:   height,
	}
}

func (a App) Init() tea.Cmd {
	return a.recall.Init()
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.recall.SetSize(msg.Width, msg.Height-3)
		a.settings_.SetSize(msg.Width, msg.Height-3)
		a.addQuote.SetSize(msg.Width, msg.Height)
		return a, nil

	case tea.KeyMsg:
		// Global quit — only when no overlay is open.
		if a.overlay == overlayNone {
			switch msg.String() {
			case "ctrl+c":
				return a, tea.Quit
			case "tab":
				// From recall page, Tab switches to settings.
				// From settings page, Tab cycles form fields (handled by settings page).
				if a.page == pageRecall {
					a.page = pageSettings
					return a, nil
				}
			case "esc":
				// From settings page, Esc goes back to recall.
				if a.page == pageSettings {
					a.page = pageRecall
					return a, nil
				}
			}
		}
		if msg.String() == "ctrl+c" {
			return a, tea.Quit
		}

	// Overlay lifecycle.
	case pages.OpenAddQuoteMsg:
		a.overlay = overlayAddQuote
		a.addQuote.Reset()
		return a, a.addQuote.Init()

	case pages.CloseAddQuoteMsg:
		a.overlay = overlayNone
		return a, nil
	}

	// Route messages to the active overlay or page.
	if a.overlay == overlayAddQuote {
		var cmd tea.Cmd
		a.addQuote, cmd = a.addQuote.Update(msg)
		cmds = append(cmds, cmd)
		return a, tea.Batch(cmds...)
	}

	switch a.page {
	case pageRecall:
		var cmd tea.Cmd
		a.recall, cmd = a.recall.Update(msg)
		cmds = append(cmds, cmd)
	case pageSettings:
		var cmd tea.Cmd
		a.settings_, cmd = a.settings_.Update(msg)
		cmds = append(cmds, cmd)
	}

	return a, tea.Batch(cmds...)
}

func (a App) View() string {
	header := a.headerView()
	var body string

	switch a.page {
	case pageRecall:
		body = a.recall.View()
	case pageSettings:
		body = a.settings_.View()
	}

	base := lipgloss.JoinVertical(lipgloss.Left, header, body)

	if a.overlay == overlayAddQuote {
		return lipgloss.Place(a.width, a.height,
			lipgloss.Center, lipgloss.Center,
			a.addQuote.View(),
			lipgloss.WithWhitespaceChars(" "),
			lipgloss.WithWhitespaceForeground(styles.ColorMuted),
		)
	}

	return base
}

func (a App) headerView() string {
	title := styles.TitleBar.Render("iRecall")

	recallTab := styles.TabInactive.Render("Recall")
	settingsTab := styles.TabInactive.Render("Settings")
	if a.page == pageRecall {
		recallTab = styles.TabActive.Render("Recall")
	} else {
		settingsTab = styles.TabActive.Render("Settings")
	}
	tabs := recallTab + styles.Muted.Render(" | ") + settingsTab

	spacer := lipgloss.NewStyle().Width(a.width - lipgloss.Width(title) - lipgloss.Width(tabs) - 4)
	return lipgloss.JoinHorizontal(lipgloss.Top,
		title,
		spacer.Render(""),
		styles.HelpBar.Render(tabs),
	)
}
