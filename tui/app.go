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
	pageQuotes
	pageSettings
	pageCount // sentinel
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

	recall    pages.RecallPage
	quotes    pages.QuotesPage
	settings_ pages.SettingsPage
	addQuote  pages.AddQuotePage

	width  int
	height int
}

func NewApp(engine *core.Engine, settings *core.Settings, width, height int) App {
	return App{
		engine:    engine,
		settings:  settings,
		page:      pageRecall,
		overlay:   overlayNone,
		recall:    pages.NewRecallPage(engine, width, height-3),
		quotes:    pages.NewQuotesPage(engine, width, height-3),
		settings_: pages.NewSettingsPage(engine, width, height-3, settings),
		addQuote:  pages.NewAddQuotePage(engine, width, height),
		width:     width,
		height:    height,
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
		a.quotes.SetSize(msg.Width, msg.Height-3)
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
				// Cycle forward: Recall → Quotes → Settings → Recall.
				// Settings page uses Tab internally for field cycling, so only
				// advance from non-settings pages here.
				if a.page != pageSettings {
					next := activePage((int(a.page) + 1) % int(pageCount))
					a.page = next
					if next == pageQuotes {
						cmds = append(cmds, a.quotes.Reload())
					}
					return a, tea.Batch(cmds...)
				}
			case "esc":
				// From settings page, Esc returns to recall.
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
		// Reload quotes if the quotes page is visible.
		if a.page == pageQuotes {
			cmds = append(cmds, a.quotes.Reload())
		}
		return a, tea.Batch(cmds...)

	case pages.QuotesLoadedMsg:
		var cmd tea.Cmd
		a.quotes, cmd = a.quotes.Update(msg)
		cmds = append(cmds, cmd)
		return a, tea.Batch(cmds...)
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
	case pageQuotes:
		var cmd tea.Cmd
		a.quotes, cmd = a.quotes.Update(msg)
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
	case pageQuotes:
		body = a.quotes.View()
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

	tab := func(label string, p activePage) string {
		if a.page == p {
			return styles.TabActive.Render(label)
		}
		return styles.TabInactive.Render(label)
	}
	sep := styles.Muted.Render(" | ")
	tabs := tab("Recall", pageRecall) + sep +
		tab("Quotes", pageQuotes) + sep +
		tab("Settings", pageSettings)

	spacer := lipgloss.NewStyle().Width(a.width - lipgloss.Width(title) - lipgloss.Width(tabs) - 4)
	return lipgloss.JoinHorizontal(lipgloss.Top,
		title,
		spacer.Render(""),
		styles.HelpBar.Render(tabs),
	)
}
