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

type RecallHistorySavedMsg struct{ Err error }
type RecallQuoteSavedMsg struct {
	Quote *core.Quote
	Err   error
}

// --- RecallPage ---

// RecallPage is the main Q&A page.
type RecallPage struct {
	engine *core.Engine

	input     textinput.Model
	response  viewport.Model
	refQuotes quoteListWidget
	spinner   spinner.Model
	busy      bool
	statusMsg string
	statusErr bool
	focus     recallFocus

	keywords []string
	question string
	respBuf  string

	width  int
	height int
}

type recallFocus int

const (
	focusInput recallFocus = iota
	focusReferenceQuotes
)

func NewRecallPage(engine *core.Engine, width, height int) RecallPage {
	ti := textinput.New()
	ti.Placeholder = "Ask anything..."
	ti.Focus()
	ti.CharLimit = 2000

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(styles.ColorAccent)

	p := RecallPage{
		engine:    engine,
		input:     ti,
		spinner:   sp,
		width:     width,
		height:    height,
		focus:     focusInput,
		refQuotes: newQuoteListWidget("Reference Quotes", width-4, max(3, height/3)),
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
		case "ctrl+j":
			p.toggleFocus()
			return p, nil
		}

		switch msg.String() {
		case "ctrl+n":
			return p, func() tea.Msg { return OpenQuoteEditorMsg{Mode: QuoteEditorModeAdd} }
		case "ctrl+s":
			if p.busy || strings.TrimSpace(p.question) == "" || strings.TrimSpace(p.respBuf) == "" {
				break
			}
			p.statusMsg = ""
			p.statusErr = false
			return p, p.saveRecallAsQuote()

		case "enter":
			if p.focus != focusInput {
				break
			}
			if p.busy || strings.TrimSpace(p.input.Value()) == "" {
				break
			}
			question := strings.TrimSpace(p.input.Value())
			p.input.SetValue("")
			p.question = question
			p.respBuf = ""
			p.updateResponsePanel()
			p.refQuotes.ClearQuotes()
			p.keywords = nil
			p.busy = true
			p.statusMsg = ""
			p.statusErr = false
			return p, tea.Batch(p.spinner.Tick, p.runRecall(question))
		}

	case TokenMsg:
		p.respBuf += msg.Token
		p.updateResponsePanel()
		p.response.GotoBottom()

	case KeywordsReadyMsg:
		p.keywords = msg.Keywords

	case QuotesReadyMsg:
		p.refQuotes.SetQuotes(msg.Quotes)

	case RecallDoneMsg:
		p.busy = false
		if msg.Err != nil {
			p.statusMsg = "Error: " + msg.Err.Error()
			p.statusErr = true
			break
		}
		return p, p.saveHistory()

	case RecallHistorySavedMsg:
		if msg.Err != nil {
			p.statusMsg = "Error saving history: " + msg.Err.Error()
			p.statusErr = true
		}

	case RecallQuoteSavedMsg:
		if msg.Err != nil {
			p.statusMsg = "Error saving recall quote: " + msg.Err.Error()
			p.statusErr = true
			break
		}
		p.statusMsg = "Saved recall as quote."
		p.statusErr = false
		return p, func() tea.Msg {
			return OpenNoticeMsg{
				Title:   "Recall Saved as Quote",
				Message: "The current question and grounded response were saved as a quote with generated tags.",
			}
		}

	case quotesAndStreamMsg:
		return p.handleQuotesAndStream(msg)

	case tokenWithChannel:
		p.respBuf += msg.token
		p.updateResponsePanel()
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
	if p.focus == focusInput {
		p.input, cmd = p.input.Update(msg)
		cmds = append(cmds, cmd)
	}
	p.response, cmd = p.response.Update(msg)
	cmds = append(cmds, cmd)
	if p.focus == focusReferenceQuotes {
		var action quoteListAction
		action, cmd = p.refQuotes.Update(msg)
		cmds = append(cmds, cmd)
		switch action.kind {
		case quoteListActionEdit:
			return p, tea.Batch(append(cmds, func() tea.Msg {
				return OpenQuoteEditorMsg{Mode: QuoteEditorModeEdit, Quote: action.quote}
			})...)
		case quoteListActionDelete:
			return p, tea.Batch(append(cmds, func() tea.Msg {
				return OpenDeleteQuotesMsg{Quotes: action.quotes}
			})...)
		case quoteListActionShare:
			return p, tea.Batch(append(cmds, func() tea.Msg {
				return OpenQuoteShareMsg{Quotes: action.quotes}
			})...)
		}
	}

	return p, tea.Batch(cmds...)
}

func (p RecallPage) View() string {
	helpLine := styles.HelpBar.Render(
		"enter: Ask   ctrl+n: Add Quote   ctrl+s: Save Q/A as Quote   ctrl+j: Jump focus   tab/shift+tab: Switch Page",
	)
	if p.busy {
		helpLine = styles.HelpBar.Render(p.spinner.View() + " Thinking...")
	}

	inputStyle := styles.Panel
	if p.focus == focusInput {
		inputStyle = styles.PanelActive
	}
	inputBox := inputStyle.Width(p.width - 4).Render(p.input.View())

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

	refBox := p.refQuotes.View(
		p.focus == focusReferenceQuotes,
		"ctrl+j: Focus input",
		"ctrl+j: Focus input",
	)

	var status string
	if p.statusMsg != "" {
		if p.statusErr {
			status = styles.StatusErr.Render(p.statusMsg)
		} else {
			status = styles.StatusOK.Render(p.statusMsg)
		}
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
	// Fixed overhead: 3 (input) + 1 (keywords) + 1 (help) + 5 (resp panel) + 7 (ref panel incl. local help) + 1 (status) = 18
	// Target body = p.height so total app = header(1) + p.height = T-2, no overflow.
	remaining := p.height - 18
	responseH := remaining * 2 / 3
	refH := remaining - responseH
	if responseH < 3 {
		responseH = 3
	}
	if refH < 3 {
		refH = 3
	}
	p.response = viewport.New(innerW, responseH)
	p.refQuotes.SetSize(p.width-4, refH)
}

func (p *RecallPage) updateResponsePanel() {
	var lines []string
	if strings.TrimSpace(p.question) != "" {
		lines = append(lines, styles.Muted.Render("Question: ")+styles.Accent.Render(p.question))
		lines = append(lines, "")
	}
	if p.respBuf != "" {
		lines = append(lines, p.respBuf)
	}
	p.response.SetContent(strings.Join(lines, "\n"))
}

func (p *RecallPage) toggleFocus() {
	if p.focus == focusInput {
		p.focus = focusReferenceQuotes
		p.input.Blur()
		return
	}
	p.focus = focusInput
	p.input.Focus()
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
	p.keywords = msg.keywords
	p.refQuotes.SetQuotes(msg.quotes)

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
	p.refQuotes.ApplyQuoteUpdate(updated)
}

func (p *RecallPage) RemoveQuotes(ids []int64) {
	p.refQuotes.RemoveQuotes(ids)
}

func (p *RecallPage) saveHistory() tea.Cmd {
	engine := p.engine
	question := p.question
	response := p.respBuf
	quotes := append([]core.Quote(nil), p.refQuotes.quotes...)
	return func() tea.Msg {
		_, err := engine.SaveRecallHistory(context.Background(), question, response, quotes)
		return RecallHistorySavedMsg{Err: err}
	}
}

func (p *RecallPage) saveRecallAsQuote() tea.Cmd {
	engine := p.engine
	question := p.question
	response := p.respBuf
	keywords := append([]string(nil), p.keywords...)
	return func() tea.Msg {
		quote, err := engine.SaveRecallAsQuote(context.Background(), question, response, keywords)
		return RecallQuoteSavedMsg{Quote: quote, Err: err}
	}
}
