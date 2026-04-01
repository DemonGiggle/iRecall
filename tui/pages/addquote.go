package pages

import (
	"context"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gigol/irecall/core"
	"github.com/gigol/irecall/tui/styles"
)

// --- Messages ---

// CloseAddQuoteMsg tells the app to dismiss the overlay.
type CloseAddQuoteMsg struct{}

// AddQuoteDoneMsg signals the result of an add-quote operation.
type AddQuoteDoneMsg struct {
	Quote *core.Quote
	Err   error
}

// --- AddQuotePage ---

// AddQuotePage is a modal overlay for capturing a new quote.
type AddQuotePage struct {
	engine    *core.Engine
	textarea  textarea.Model
	spinner   spinner.Model
	busy      bool
	statusMsg string
	isErr     bool
	clearAt   time.Time

	width  int
	height int
}

func NewAddQuotePage(engine *core.Engine, width, height int) AddQuotePage {
	ta := textarea.New()
	ta.Placeholder = "Type or paste your note here. Multi-line input supported."
	ta.Focus()
	ta.SetWidth(width - 12)
	ta.SetHeight(6)
	ta.CharLimit = 10000

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(styles.ColorAccent)

	return AddQuotePage{
		engine:   engine,
		textarea: ta,
		spinner:  sp,
		width:    width,
		height:   height,
	}
}

func (p AddQuotePage) Init() tea.Cmd {
	return textarea.Blink
}

func (p AddQuotePage) Update(msg tea.Msg) (AddQuotePage, tea.Cmd) {
	var cmds []tea.Cmd

	// Clear status message after timeout.
	if !p.clearAt.IsZero() && time.Now().After(p.clearAt) {
		p.statusMsg = ""
		p.clearAt = time.Time{}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if !p.busy {
				return p, func() tea.Msg { return CloseAddQuoteMsg{} }
			}

		case "ctrl+s":
			if p.busy {
				break
			}
			content := strings.TrimSpace(p.textarea.Value())
			if content == "" {
				p.statusMsg = "Nothing to save."
				p.isErr = true
				p.clearAt = time.Now().Add(2 * time.Second)
				break
			}
			p.busy = true
			p.statusMsg = ""
			cmds = append(cmds, p.spinner.Tick, p.doAdd(content))
		}

	case AddQuoteDoneMsg:
		p.busy = false
		if msg.Err != nil {
			p.statusMsg = "Error: " + msg.Err.Error()
			p.isErr = true
			p.clearAt = time.Now().Add(4 * time.Second)
		} else {
			p.statusMsg = "Saved."
			p.isErr = false
			p.clearAt = time.Now().Add(2 * time.Second)
			p.textarea.Reset()
			// Close after a brief moment so the user sees "Saved."
			cmds = append(cmds, tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
				return CloseAddQuoteMsg{}
			}))
		}

	case spinner.TickMsg:
		if p.busy {
			var cmd tea.Cmd
			p.spinner, cmd = p.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	var cmd tea.Cmd
	p.textarea, cmd = p.textarea.Update(msg)
	cmds = append(cmds, cmd)

	return p, tea.Batch(cmds...)
}

func (p AddQuotePage) View() string {
	helpLine := "  Ctrl+S: Save   Esc: Cancel"
	if p.busy {
		helpLine = "  " + p.spinner.View() + " Saving..."
	}

	var statusLine string
	if p.statusMsg != "" {
		if p.isErr {
			statusLine = "\n  " + styles.StatusErr.Render(p.statusMsg)
		} else {
			statusLine = "\n  " + styles.StatusOK.Render(p.statusMsg)
		}
	}

	hint := styles.Muted.Render("  Tags will be extracted automatically by the LLM.")

	inner := lipgloss.JoinVertical(lipgloss.Left,
		"\n",
		p.textarea.View(),
		"\n",
		hint,
		statusLine,
		"\n",
		styles.HelpBar.Render(helpLine),
	)

	modalW := p.width - 8
	if modalW < 40 {
		modalW = 40
	}

	modal := styles.Modal.Width(modalW).Render(inner)

	// Center the modal on screen.
	return lipgloss.Place(p.width, p.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Left,
			styles.Bold.Foreground(styles.ColorPrimary).Render(" Add Quote "),
			modal,
		),
	)
}

func (p *AddQuotePage) SetSize(width, height int) {
	p.width = width
	p.height = height
	p.textarea.SetWidth(width - 12)
}

func (p *AddQuotePage) Reset() {
	p.textarea.Reset()
	p.textarea.Focus()
	p.statusMsg = ""
	p.isErr = false
	p.busy = false
	p.clearAt = time.Time{}
}

func (p *AddQuotePage) doAdd(content string) tea.Cmd {
	engine := p.engine
	return func() tea.Msg {
		q, err := engine.AddQuote(context.Background(), content)
		return AddQuoteDoneMsg{Quote: q, Err: err}
	}
}
