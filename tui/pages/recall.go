package pages

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gigol/irecall/core"
	"github.com/gigol/irecall/tui/styles"
)

// --- Messages ---

// TokenMsg carries a single streamed token from the LLM.
type TokenMsg struct{ Token string }

// RecallDoneMsg signals that streaming has finished.
type RecallDoneMsg struct{ Err error }

// QuotesReadyMsg carries the retrieved reference quotes.
type QuotesReadyMsg struct{ Quotes []core.Quote }

// KeywordsReadyMsg carries the extracted search keywords.
type KeywordsReadyMsg struct{ Keywords []string }

// --- RecallPage ---

// RecallPage is the main Q&A page.
type RecallPage struct {
	engine *core.Engine

	input     textinput.Model
	response  viewport.Model
	refPanel  viewport.Model
	spinner   spinner.Model
	busy      bool
	statusMsg string

	quotes   []core.Quote
	keywords []string
	respBuf  string
	quoteFns quoteSelection

	width  int
	height int
}

func NewRecallPage(engine *core.Engine, width, height int) RecallPage {
	ti := textinput.New()
	ti.Placeholder = "Ask anything..."
	ti.Focus()
	ti.CharLimit = 2000

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(styles.ColorAccent)

	p := RecallPage{
		engine:   engine,
		input:    ti,
		spinner:  sp,
		width:    width,
		height:   height,
		quoteFns: newQuoteSelection(),
	}
	p.recalcLayout()
	return p
}

func (p RecallPage) Init() tea.Cmd {
	return textinput.Blink
}

func (p RecallPage) Update(msg tea.Msg) (RecallPage, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+n":
			return p, func() tea.Msg { return OpenQuoteEditorMsg{Mode: QuoteEditorModeAdd} }

		case "enter":
			if p.busy || strings.TrimSpace(p.input.Value()) == "" {
				break
			}
			question := strings.TrimSpace(p.input.Value())
			p.input.SetValue("")
			p.respBuf = ""
			p.response.SetContent("")
			p.refPanel.SetContent("")
			p.quotes = nil
			p.keywords = nil
			p.quoteFns.clear()
			p.busy = true
			p.statusMsg = ""
			return p, tea.Batch(p.spinner.Tick, p.runRecall(question))
		case "up":
			p.quoteFns.move(-1, p.quotes)
			p.refreshReferencePanel()
			return p, nil
		case "down":
			p.quoteFns.move(1, p.quotes)
			p.refreshReferencePanel()
			return p, nil
		case " ":
			p.quoteFns.toggleCurrent(p.quotes)
			p.refreshReferencePanel()
			return p, nil
		case "e":
			if q := p.quoteFns.current(p.quotes); q != nil {
				quote := *q
				return p, func() tea.Msg {
					return OpenQuoteEditorMsg{Mode: QuoteEditorModeEdit, Quote: &quote}
				}
			}
		case "d":
			selected := p.quoteFns.selectedQuotes(p.quotes)
			if len(selected) > 0 {
				return p, func() tea.Msg { return OpenDeleteQuotesMsg{Quotes: selected} }
			}
			return p, nil
		}

	case TokenMsg:
		p.respBuf += msg.Token
		p.response.SetContent(p.respBuf)
		p.response.GotoBottom()

	case KeywordsReadyMsg:
		p.keywords = msg.Keywords

	case QuotesReadyMsg:
		p.quotes = msg.Quotes
		p.quoteFns.clamp(p.quotes)
		p.refreshReferencePanel()

	case RecallDoneMsg:
		p.busy = false
		if msg.Err != nil {
			p.statusMsg = "Error: " + msg.Err.Error()
		}

	case quotesAndStreamMsg:
		return p.handleQuotesAndStream(msg)

	case tokenWithChannel:
		p.respBuf += msg.token
		p.response.SetContent(p.respBuf)
		p.response.GotoBottom()
		ch := msg.ch
		cmds = append(cmds, func() tea.Msg {
			return drainNext(ch)
		})

	case spinner.TickMsg:
		if p.busy {
			var cmd tea.Cmd
			p.spinner, cmd = p.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	// Delegate key events to sub-components when not intercepted above.
	var cmd tea.Cmd
	p.input, cmd = p.input.Update(msg)
	cmds = append(cmds, cmd)
	p.response, cmd = p.response.Update(msg)
	cmds = append(cmds, cmd)
	p.refPanel, cmd = p.refPanel.Update(msg)
	cmds = append(cmds, cmd)

	return p, tea.Batch(cmds...)
}

func (p RecallPage) View() string {
	helpLine := styles.HelpBar.Render(
		"Enter: Ask   Ctrl+N: Add Quote   ↑/↓: Move quote   Space: Select   E: Edit   D: Delete   Tab: Settings",
	)
	if p.busy {
		helpLine = styles.HelpBar.Render(p.spinner.View() + " Thinking...")
	}

	inputBox := styles.PanelActive.Width(p.width - 4).Render(p.input.View())

	var keywordsLine string
	if len(p.keywords) > 0 {
		keywordsLine = styles.Muted.Render("Keywords: ") + styles.Accent.Render(strings.Join(p.keywords, "  ·  "))
	} else {
		keywordsLine = styles.Muted.Render("Keywords: —")
	}

	responseLabel := styles.Accent.Render("Response")
	responseBox := styles.Panel.Width(p.width - 4).
		Height(p.response.Height + 3).
		Render(responseLabel + "\n" + p.response.View())

	refLabel := styles.Accent.Render("Reference Quotes")
	refBox := styles.Panel.Width(p.width - 4).
		Height(p.refPanel.Height + 3).
		Render(refLabel + "\n" + p.refPanel.View())

	var status string
	if p.statusMsg != "" {
		status = styles.StatusErr.Render(p.statusMsg)
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		inputBox,
		keywordsLine,
		helpLine,
		responseBox,
		refBox,
		status,
	)
}

func (p *RecallPage) SetSize(width, height int) {
	p.width = width
	p.height = height
	p.recalcLayout()
}

func (p *RecallPage) recalcLayout() {
	innerW := p.width - 6 // account for panel borders + padding
	// Divide remaining vertical space between response and ref panels.
	// Height() is inner (before borders), so Panel.Height(n) renders n+2 outer lines.
	// Fixed overhead: 3 (input) + 1 (keywords) + 1 (help) + 5 (resp panel) + 5 (ref panel) + 1 (status) = 16
	// Target body = p.height so total app = header(1) + p.height = T-2, no overflow.
	remaining := p.height - 16
	responseH := remaining * 2 / 3
	refH := remaining - responseH
	if responseH < 3 {
		responseH = 3
	}
	if refH < 3 {
		refH = 3
	}
	p.response = viewport.New(innerW, responseH)
	p.refPanel = viewport.New(innerW, refH)
}

func (p *RecallPage) refreshReferencePanel() {
	p.quoteFns.clamp(p.quotes)
	p.refPanel.SetContent(renderQuoteFunctionList(p.quotes, p.quoteFns, p.refPanel.Width, false))
}

// runRecall starts the full recall pipeline as a tea.Cmd.
func (p *RecallPage) runRecall(question string) tea.Cmd {
	engine := p.engine
	return func() tea.Msg {
		ctx := context.Background()

		keywords, err := engine.ExtractKeywords(ctx, question)
		if err != nil {
			return RecallDoneMsg{Err: fmt.Errorf("keyword extraction: %w", err)}
		}

		quotes, err := engine.SearchQuotes(ctx, keywords)
		if err != nil {
			return RecallDoneMsg{Err: fmt.Errorf("search: %w", err)}
		}

		// Return quotes immediately so the TUI can render them before streaming starts.
		// We can't send two messages from one Cmd, so we return QuotesReadyMsg here
		// and chain streaming via a subsequent Cmd.
		_ = quotes // passed via closure to streaming goroutine below

		return quotesAndStreamMsg{question: question, quotes: quotes, keywords: keywords}
	}
}

// quotesAndStreamMsg is an internal message that carries quotes + triggers streaming.
type quotesAndStreamMsg struct {
	question string
	quotes   []core.Quote
	keywords []string
}

// We handle this internal message in Update by dispatching two effects.
func (p RecallPage) handleQuotesAndStream(msg quotesAndStreamMsg) (RecallPage, tea.Cmd) {
	p.quotes = msg.quotes
	p.keywords = msg.keywords
	p.quoteFns.clamp(msg.quotes)
	p.refreshReferencePanel()

	engine := p.engine
	return p, func() tea.Msg {
		ctx := context.Background()
		tokenCh := make(chan string, 64)

		if err := engine.GenerateResponse(ctx, msg.question, msg.quotes, tokenCh); err != nil {
			return RecallDoneMsg{Err: err}
		}

		// Drain the channel and send each token as a TokenMsg.
		// Bubbletea doesn't support sending multiple messages from one Cmd,
		// so we use a recursive Cmd pattern: each Cmd reads one token and
		// returns either a TokenMsg or RecallDoneMsg.
		return drainNext(tokenCh)
	}
}

func drainNext(ch <-chan string) tea.Msg {
	tok, ok := <-ch
	if !ok {
		return RecallDoneMsg{}
	}
	return tokenWithChannel{token: tok, ch: ch}
}

type tokenWithChannel struct {
	token string
	ch    <-chan string
}

func (p *RecallPage) ApplyQuoteUpdate(updated core.Quote) {
	for i := range p.quotes {
		if p.quotes[i].ID == updated.ID {
			p.quotes[i] = updated
			p.refreshReferencePanel()
			return
		}
	}
}

func (p *RecallPage) RemoveQuotes(ids []int64) {
	if len(ids) == 0 || len(p.quotes) == 0 {
		return
	}
	remove := idsSet(ids)
	filtered := p.quotes[:0]
	for _, q := range p.quotes {
		if !remove[q.ID] {
			filtered = append(filtered, q)
		}
	}
	p.quotes = filtered
	p.quoteFns.clamp(p.quotes)
	p.refreshReferencePanel()
}
