package tui

import (
	"strings"

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
	overlayUserProfilePrompt
	overlayQuoteEditor
	overlayDeleteQuotes
	overlayQuoteShare
)

// App is the root Bubbletea model. It owns page routing and global key handling.
type App struct {
	engine   *core.Engine
	settings *core.Settings
	profile  *core.UserProfile

	page    activePage
	overlay overlayKind

	userProfile  pages.UserProfilePromptPage
	recall       pages.RecallPage
	quotes       pages.QuotesPage
	settings_    pages.SettingsPage
	quoteEditor  pages.QuoteEditorPage
	deleteQuotes pages.DeleteQuotesPage
	quoteShare   pages.QuoteSharePage

	width  int
	height int
}

func NewApp(engine *core.Engine, settings *core.Settings, profile *core.UserProfile, width, height int) App {
	overlay := overlayNone
	if profile == nil || profile.DisplayName == "" {
		overlay = overlayUserProfilePrompt
	}
	return App{
		engine:       engine,
		settings:     settings,
		profile:      profile,
		page:         pageRecall,
		overlay:      overlay,
		userProfile:  pages.NewUserProfilePromptPage(engine, width, height, profile),
		recall:       pages.NewRecallPage(engine, width, height-3),
		quotes:       pages.NewQuotesPage(engine, width, height-3),
		settings_:    pages.NewSettingsPage(engine, width, height-3, settings),
		quoteEditor:  pages.NewQuoteEditorPage(engine, width, height),
		deleteQuotes: pages.NewDeleteQuotesPage(engine, width, height),
		quoteShare:   pages.NewQuoteSharePage(engine, width, height),
		width:        width,
		height:       height,
	}
}

func (a App) Init() tea.Cmd {
	if a.overlay == overlayUserProfilePrompt {
		return a.userProfile.Init()
	}
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
		a.userProfile.SetSize(msg.Width, msg.Height)
		a.quoteEditor.SetSize(msg.Width, msg.Height)
		a.deleteQuotes.SetSize(msg.Width, msg.Height)
		a.quoteShare.SetSize(msg.Width, msg.Height)
		return a, nil

	case tea.KeyMsg:
		// Global quit — only when no overlay is open.
		if a.overlay == overlayNone {
			switch msg.String() {
			case "ctrl+c":
				return a, tea.Quit
			case "tab":
				// Cycle forward: Recall → Quotes → Settings → Recall.
				next := activePage((int(a.page) + 1) % int(pageCount))
				a.page = next
				if next == pageQuotes {
					cmds = append(cmds, a.quotes.Reload())
				}
				return a, tea.Batch(cmds...)
			case "shift+tab":
				// Cycle backward: Recall ← Quotes ← Settings ← Recall.
				next := activePage((int(a.page) - 1 + int(pageCount)) % int(pageCount))
				a.page = next
				if next == pageQuotes {
					cmds = append(cmds, a.quotes.Reload())
				}
				return a, tea.Batch(cmds...)
			}
		}
		if msg.String() == "ctrl+c" {
			return a, tea.Quit
		}

	// Overlay lifecycle.
	case pages.CloseUserProfilePromptMsg:
		a.overlay = overlayNone
		a.profile = msg.Profile
		return a, nil

	case pages.OpenQuoteEditorMsg:
		a.overlay = overlayQuoteEditor
		a.quoteEditor.Reset(msg.Mode, msg.Quote)
		return a, a.quoteEditor.Init()

	case pages.CloseQuoteEditorMsg:
		a.overlay = overlayNone
		if msg.SavedQuote != nil {
			a.recall.ApplyQuoteUpdate(*msg.SavedQuote)
			a.quotes.ApplyQuoteUpdate(*msg.SavedQuote)
		}
		cmds = append(cmds, a.quotes.Reload())
		return a, tea.Batch(cmds...)

	case pages.OpenDeleteQuotesMsg:
		a.overlay = overlayDeleteQuotes
		a.deleteQuotes.Reset(msg.Quotes)
		return a, a.deleteQuotes.Init()

	case pages.OpenQuoteShareMsg:
		a.overlay = overlayQuoteShare
		a.quoteShare.Reset(msg.Quotes)
		return a, a.quoteShare.Init()

	case pages.CloseDeleteQuotesMsg:
		a.overlay = overlayNone
		if len(msg.DeletedIDs) > 0 {
			a.recall.RemoveQuotes(msg.DeletedIDs)
			a.quotes.RemoveQuotes(msg.DeletedIDs)
		}
		cmds = append(cmds, a.quotes.Reload())
		return a, tea.Batch(cmds...)

	case pages.CloseQuoteShareMsg:
		a.overlay = overlayNone
		return a, nil

	case pages.QuotesLoadedMsg:
		var cmd tea.Cmd
		a.quotes, cmd = a.quotes.Update(msg)
		cmds = append(cmds, cmd)
		return a, tea.Batch(cmds...)
	}

	// Route messages to the active overlay or page.
	if a.overlay == overlayUserProfilePrompt {
		var cmd tea.Cmd
		a.userProfile, cmd = a.userProfile.Update(msg)
		cmds = append(cmds, cmd)
		return a, tea.Batch(cmds...)
	}
	if a.overlay == overlayQuoteEditor {
		var cmd tea.Cmd
		a.quoteEditor, cmd = a.quoteEditor.Update(msg)
		cmds = append(cmds, cmd)
		return a, tea.Batch(cmds...)
	}
	if a.overlay == overlayDeleteQuotes {
		var cmd tea.Cmd
		a.deleteQuotes, cmd = a.deleteQuotes.Update(msg)
		cmds = append(cmds, cmd)
		return a, tea.Batch(cmds...)
	}
	if a.overlay == overlayQuoteShare {
		var cmd tea.Cmd
		a.quoteShare, cmd = a.quoteShare.Update(msg)
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

	if a.overlay == overlayUserProfilePrompt {
		return lipgloss.Place(a.width, a.height,
			lipgloss.Center, lipgloss.Center,
			a.userProfile.View(),
			lipgloss.WithWhitespaceChars(" "),
			lipgloss.WithWhitespaceForeground(styles.ColorMuted),
		)
	}
	if a.overlay == overlayQuoteEditor {
		return lipgloss.Place(a.width, a.height,
			lipgloss.Center, lipgloss.Center,
			a.quoteEditor.View(),
			lipgloss.WithWhitespaceChars(" "),
			lipgloss.WithWhitespaceForeground(styles.ColorMuted),
		)
	}
	if a.overlay == overlayDeleteQuotes {
		return lipgloss.Place(a.width, a.height,
			lipgloss.Center, lipgloss.Center,
			a.deleteQuotes.View(),
			lipgloss.WithWhitespaceChars(" "),
			lipgloss.WithWhitespaceForeground(styles.ColorMuted),
		)
	}
	if a.overlay == overlayQuoteShare {
		return lipgloss.Place(a.width, a.height,
			lipgloss.Center, lipgloss.Center,
			a.quoteShare.View(),
			lipgloss.WithWhitespaceChars(" "),
			lipgloss.WithWhitespaceForeground(styles.ColorMuted),
		)
	}

	return base
}

func (a App) headerView() string {
	title := styles.TitleBar.Render("iRecall")
	greeting := ""
	if a.profile != nil && strings.TrimSpace(a.profile.DisplayName) != "" {
		greeting = styles.HelpBar.Render("Hi! " + strings.TrimSpace(a.profile.DisplayName))
	}

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

	right := tabs
	if greeting != "" {
		right = greeting + "  " + tabs
	}
	spacer := lipgloss.NewStyle().Width(a.width - lipgloss.Width(title) - lipgloss.Width(right) - 4)
	return lipgloss.JoinHorizontal(lipgloss.Top,
		title,
		spacer.Render(""),
		right,
	)
}
