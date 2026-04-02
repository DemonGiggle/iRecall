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

type QuoteEditorMode int

const (
	QuoteEditorModeAdd QuoteEditorMode = iota
	QuoteEditorModeEdit
)

// OpenQuoteEditorMsg tells the app router to show the quote editor overlay.
type OpenQuoteEditorMsg struct {
	Mode  QuoteEditorMode
	Quote *core.Quote
}

// CloseQuoteEditorMsg tells the app router to dismiss the quote editor overlay.
type CloseQuoteEditorMsg struct {
	SavedQuote *core.Quote
}

// QuoteEditorDoneMsg signals the result of an add/edit operation.
type QuoteEditorDoneMsg struct {
	Quote *core.Quote
	Err   error
}

// QuoteRefineDoneMsg signals the result of a draft refinement request.
type QuoteRefineDoneMsg struct {
	Refined string
	Err     error
}

// QuoteEditorPage is a modal overlay for adding or editing a quote.
type QuoteEditorPage struct {
	engine    *core.Engine
	textarea  textarea.Model
	spinner   spinner.Model
	mode      QuoteEditorMode
	editingID int64
	busy      bool
	preview   bool
	statusMsg string
	isErr     bool
	clearAt   time.Time
	refined   string

	width  int
	height int
}

func NewQuoteEditorPage(engine *core.Engine, width, height int) QuoteEditorPage {
	ta := textarea.New()
	ta.Placeholder = "Type or paste your note here. Multi-line input supported."
	ta.Focus()
	ta.SetWidth(width - 12)
	ta.SetHeight(6)
	ta.CharLimit = 10000

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(styles.ColorAccent)

	return QuoteEditorPage{
		engine:   engine,
		textarea: ta,
		spinner:  sp,
		width:    width,
		height:   height,
		mode:     QuoteEditorModeAdd,
	}
}

func (p QuoteEditorPage) Init() tea.Cmd {
	return textarea.Blink
}

func (p QuoteEditorPage) Update(msg tea.Msg) (QuoteEditorPage, tea.Cmd) {
	var cmds []tea.Cmd

	if !p.clearAt.IsZero() && time.Now().After(p.clearAt) {
		p.statusMsg = ""
		p.clearAt = time.Time{}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if p.preview {
			switch msg.String() {
			case "enter":
				p.textarea.SetValue(p.refined)
				p.preview = false
				p.refined = ""
				p.statusMsg = "Refined draft applied. Review and keep editing."
				p.isErr = false
				p.clearAt = time.Now().Add(3 * time.Second)
				p.textarea.Focus()
				return p, nil
			case "esc":
				p.preview = false
				p.refined = ""
				p.statusMsg = "Refined draft discarded."
				p.isErr = false
				p.clearAt = time.Now().Add(2 * time.Second)
				p.textarea.Focus()
				return p, nil
			}
		}

		switch msg.String() {
		case "esc":
			if !p.busy {
				return p, func() tea.Msg { return CloseQuoteEditorMsg{} }
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
			cmds = append(cmds, p.spinner.Tick, p.persistQuote(content))
		case "ctrl+r":
			if p.busy {
				break
			}
			content := strings.TrimSpace(p.textarea.Value())
			if content == "" {
				p.statusMsg = "Nothing to refine."
				p.isErr = true
				p.clearAt = time.Now().Add(2 * time.Second)
				break
			}
			p.busy = true
			p.statusMsg = ""
			cmds = append(cmds, p.spinner.Tick, p.refineQuote(content))
		}

	case QuoteEditorDoneMsg:
		p.busy = false
		if msg.Err != nil {
			p.statusMsg = "Error: " + msg.Err.Error()
			p.isErr = true
			p.clearAt = time.Now().Add(4 * time.Second)
		} else {
			p.statusMsg = "Saved."
			p.isErr = false
			p.clearAt = time.Now().Add(2 * time.Second)
			cmds = append(cmds, tea.Tick(2*time.Second, func(time.Time) tea.Msg {
				return CloseQuoteEditorMsg{SavedQuote: msg.Quote}
			}))
		}

	case QuoteRefineDoneMsg:
		p.busy = false
		if msg.Err != nil {
			p.statusMsg = "Error: " + msg.Err.Error()
			p.isErr = true
			p.clearAt = time.Now().Add(4 * time.Second)
		} else {
			p.preview = true
			p.refined = msg.Refined
			p.statusMsg = ""
			p.isErr = false
			p.clearAt = time.Time{}
			p.textarea.Blur()
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

func (p QuoteEditorPage) View() string {
	helpLine := "  Ctrl+S: Save   Ctrl+R: Refine   Esc: Cancel"
	if p.busy {
		helpLine = "  " + p.spinner.View() + " Working..."
	} else if p.preview {
		helpLine = "  Enter: Accept Refined Draft   Esc: Reject and Continue Editing"
	}

	var statusLine string
	if p.statusMsg != "" {
		if p.isErr {
			statusLine = "\n  " + styles.StatusErr.Render(p.statusMsg)
		} else {
			statusLine = "\n  " + styles.StatusOK.Render(p.statusMsg)
		}
	}

	hint := styles.Muted.Render("  Tags will be regenerated automatically by the LLM.")
	title := " Add Quote "
	if p.mode == QuoteEditorModeEdit {
		title = " Edit Quote "
	}

	body := p.textarea.View()
	if p.preview {
		body = styles.Panel.Width(p.width - 16).Render(p.refined)
		hint = styles.Muted.Render("  Preview the refined draft before applying it to your note.")
	}

	inner := lipgloss.JoinVertical(lipgloss.Left,
		"\n",
		body,
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

	return lipgloss.Place(p.width, p.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Left,
			styles.Bold.Foreground(styles.ColorPrimary).Render(title),
			modal,
		),
	)
}

func (p *QuoteEditorPage) SetSize(width, height int) {
	p.width = width
	p.height = height
	p.textarea.SetWidth(width - 12)
}

func (p *QuoteEditorPage) Reset(mode QuoteEditorMode, quote *core.Quote) {
	p.mode = mode
	p.textarea.Reset()
	p.textarea.Focus()
	p.statusMsg = ""
	p.isErr = false
	p.busy = false
	p.preview = false
	p.refined = ""
	p.clearAt = time.Time{}
	p.editingID = 0
	if quote != nil {
		p.editingID = quote.ID
		p.textarea.SetValue(quote.Content)
	}
}

func (p *QuoteEditorPage) persistQuote(content string) tea.Cmd {
	engine := p.engine
	mode := p.mode
	editingID := p.editingID
	return func() tea.Msg {
		if mode == QuoteEditorModeEdit {
			q, err := engine.UpdateQuote(context.Background(), editingID, content)
			return QuoteEditorDoneMsg{Quote: q, Err: err}
		}
		q, err := engine.AddQuote(context.Background(), content)
		return QuoteEditorDoneMsg{Quote: q, Err: err}
	}
}

func (p *QuoteEditorPage) refineQuote(content string) tea.Cmd {
	engine := p.engine
	return func() tea.Msg {
		refined, err := engine.RefineQuoteDraft(context.Background(), content)
		return QuoteRefineDoneMsg{Refined: refined, Err: err}
	}
}
