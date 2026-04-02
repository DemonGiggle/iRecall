package pages

import (
	"context"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gigol/irecall/core"
	"github.com/gigol/irecall/tui/styles"
)

type CloseUserProfilePromptMsg struct {
	Profile *core.UserProfile
}

type SaveUserProfileDoneMsg struct {
	Profile *core.UserProfile
	Err     error
}

type UserProfilePromptPage struct {
	engine    *core.Engine
	profile   *core.UserProfile
	input     textinput.Model
	statusMsg string
	isErr     bool
	busy      bool
	width     int
	height    int
}

func NewUserProfilePromptPage(engine *core.Engine, width, height int, profile *core.UserProfile) UserProfilePromptPage {
	input := textinput.New()
	input.Placeholder = "Your name"
	input.CharLimit = 120
	input.Focus()
	if profile != nil {
		input.SetValue(profile.DisplayName)
	}
	return UserProfilePromptPage{
		engine:    engine,
		profile:   profile,
		input:     input,
		width:     width,
		height:    height,
		statusMsg: "",
	}
}

func (p UserProfilePromptPage) Init() tea.Cmd {
	return textinput.Blink
}

func (p UserProfilePromptPage) Update(msg tea.Msg) (UserProfilePromptPage, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "ctrl+s":
			if p.busy {
				return p, nil
			}
			name := strings.TrimSpace(p.input.Value())
			if name == "" {
				p.statusMsg = "Please enter a name to continue."
				p.isErr = true
				return p, nil
			}
			p.busy = true
			p.statusMsg = ""
			return p, p.save(name)
		}
	case SaveUserProfileDoneMsg:
		p.busy = false
		if msg.Err != nil {
			p.statusMsg = "Error: " + msg.Err.Error()
			p.isErr = true
			return p, nil
		}
		return p, func() tea.Msg { return CloseUserProfilePromptMsg{Profile: msg.Profile} }
	}

	var cmd tea.Cmd
	p.input, cmd = p.input.Update(msg)
	return p, cmd
}

func (p UserProfilePromptPage) View() string {
	title := styles.Bold.Foreground(styles.ColorPrimary).Render(" Set Your Name ")
	description := "Your name is attached to quotes you share and shown when other users receive your quotes."
	help := "  enter: Save name and continue"
	if p.busy {
		help = "  Saving..."
	}

	var status string
	if p.statusMsg != "" {
		if p.isErr {
			status = styles.StatusErr.Render(p.statusMsg)
		} else {
			status = styles.StatusOK.Render(p.statusMsg)
		}
	}

	body := lipgloss.JoinVertical(lipgloss.Left,
		styles.QuoteItem.Render(description),
		"",
		styles.SectionHeader.Render("Display Name"),
		styles.PanelActive.Width(max(36, p.width-30)).Render(p.input.View()),
		"",
		status,
		styles.HelpBar.Render(help),
	)

	return lipgloss.Place(p.width, p.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Left,
			title,
			styles.Modal.Width(max(50, p.width-20)).Render(body),
		),
	)
}

func (p *UserProfilePromptPage) SetSize(width, height int) {
	p.width = width
	p.height = height
}

func (p UserProfilePromptPage) save(name string) tea.Cmd {
	engine := p.engine
	profile := core.UserProfile{}
	if p.profile != nil {
		profile = *p.profile
	}
	profile.DisplayName = name
	return func() tea.Msg {
		err := engine.SaveUserProfile(context.Background(), &profile)
		return SaveUserProfileDoneMsg{Profile: &profile, Err: err}
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
