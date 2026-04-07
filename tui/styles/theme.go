package styles

import (
	"slices"

	"github.com/charmbracelet/lipgloss"
)

type Theme struct {
	Name     string
	Primary  string
	Accent   string
	Muted    string
	Success  string
	Error    string
	Warning  string
	Fg       string
	Border   string
	Selected string
}

var (
	availableThemes = []Theme{
		{
			Name:     "violet",
			Primary:  "#7C3AED",
			Accent:   "#A78BFA",
			Muted:    "#6B7280",
			Success:  "#10B981",
			Error:    "#EF4444",
			Warning:  "#F59E0B",
			Fg:       "#F9FAFB",
			Border:   "#374151",
			Selected: "#1F2937",
		},
		{
			Name:     "forest",
			Primary:  "#0F766E",
			Accent:   "#2DD4BF",
			Muted:    "#6B7280",
			Success:  "#22C55E",
			Error:    "#EF4444",
			Warning:  "#F59E0B",
			Fg:       "#ECFDF5",
			Border:   "#334155",
			Selected: "#0F172A",
		},
		{
			Name:     "sunset",
			Primary:  "#C2410C",
			Accent:   "#FB923C",
			Muted:    "#78716C",
			Success:  "#16A34A",
			Error:    "#DC2626",
			Warning:  "#F59E0B",
			Fg:       "#FFFBEB",
			Border:   "#44403C",
			Selected: "#292524",
		},
		{
			Name:     "ocean",
			Primary:  "#0369A1",
			Accent:   "#38BDF8",
			Muted:    "#64748B",
			Success:  "#10B981",
			Error:    "#EF4444",
			Warning:  "#F59E0B",
			Fg:       "#F8FAFC",
			Border:   "#334155",
			Selected: "#0F172A",
		},
		{
			Name:     "paper",
			Primary:  "#1D4ED8",
			Accent:   "#0F766E",
			Muted:    "#6B7280",
			Success:  "#15803D",
			Error:    "#B91C1C",
			Warning:  "#B45309",
			Fg:       "#111827",
			Border:   "#CBD5E1",
			Selected: "#E2E8F0",
		},
	}
	currentThemeName string

	// Colors
	ColorPrimary  lipgloss.Color
	ColorAccent   lipgloss.Color
	ColorMuted    lipgloss.Color
	ColorSuccess  lipgloss.Color
	ColorError    lipgloss.Color
	ColorWarning  lipgloss.Color
	ColorFg       lipgloss.Color
	ColorBorder   lipgloss.Color
	ColorSelected lipgloss.Color

	// Base styles
	Bold   lipgloss.Style
	Muted  lipgloss.Style
	Accent lipgloss.Style

	// App title bar
	TitleBar lipgloss.Style

	// Active / inactive tab indicators
	TabActive   lipgloss.Style
	TabInactive lipgloss.Style

	// Panel borders
	Panel       lipgloss.Style
	PanelActive lipgloss.Style

	// Section headers inside panels
	SectionHeader lipgloss.Style

	// Status / help bar at bottom
	HelpBar lipgloss.Style

	// Status messages
	StatusOK  lipgloss.Style
	StatusErr lipgloss.Style

	// Quote items in reference list
	QuoteItem   lipgloss.Style
	QuoteNumber lipgloss.Style

	// Modal overlay
	Modal lipgloss.Style

	// Form labels
	FormLabel lipgloss.Style

	// Button (focused)
	ButtonFocused lipgloss.Style

	// Button (unfocused)
	ButtonNormal lipgloss.Style
)

func init() {
	ApplyTheme("violet")
}

func ThemeNames() []string {
	names := make([]string, 0, len(availableThemes))
	for _, theme := range availableThemes {
		names = append(names, theme.Name)
	}
	return names
}

func CurrentThemeName() string {
	return currentThemeName
}

func ApplyTheme(name string) {
	theme := resolveTheme(name)
	currentThemeName = theme.Name

	ColorPrimary = lipgloss.Color(theme.Primary)
	ColorAccent = lipgloss.Color(theme.Accent)
	ColorMuted = lipgloss.Color(theme.Muted)
	ColorSuccess = lipgloss.Color(theme.Success)
	ColorError = lipgloss.Color(theme.Error)
	ColorWarning = lipgloss.Color(theme.Warning)
	ColorFg = lipgloss.Color(theme.Fg)
	ColorBorder = lipgloss.Color(theme.Border)
	ColorSelected = lipgloss.Color(theme.Selected)

	Bold = lipgloss.NewStyle().Bold(true)
	Muted = lipgloss.NewStyle().Foreground(ColorMuted)
	Accent = lipgloss.NewStyle().Foreground(ColorAccent)

	TitleBar = lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Bold(true).
		Padding(0, 1)

	TabActive = lipgloss.NewStyle().
		Foreground(ColorAccent).
		Bold(true).
		Underline(true)

	TabInactive = lipgloss.NewStyle().
		Foreground(ColorMuted)

	Panel = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorder).
		Padding(0, 1)

	PanelActive = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorPrimary).
		Padding(0, 1)

	SectionHeader = lipgloss.NewStyle().
		Foreground(ColorAccent).
		Bold(true).
		MarginBottom(1)

	HelpBar = lipgloss.NewStyle().
		Foreground(ColorMuted).
		Padding(0, 1)

	StatusOK = lipgloss.NewStyle().Foreground(ColorSuccess)
	StatusErr = lipgloss.NewStyle().Foreground(ColorError)

	QuoteItem = lipgloss.NewStyle().
		Foreground(ColorFg).
		PaddingLeft(1)

	QuoteNumber = lipgloss.NewStyle().
		Foreground(ColorAccent).
		Bold(true)

	Modal = lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(ColorPrimary).
		Padding(1, 2)

	FormLabel = lipgloss.NewStyle().
		Foreground(ColorAccent).
		Width(20)

	ButtonFocused = lipgloss.NewStyle().
		Foreground(ColorFg).
		Border(lipgloss.NormalBorder()).
		BorderForeground(ColorPrimary).
		Background(ColorPrimary).
		Padding(0, 2).
		Bold(true)

	ButtonNormal = lipgloss.NewStyle().
		Foreground(ColorMuted).
		Border(lipgloss.NormalBorder()).
		BorderForeground(ColorBorder).
		Padding(0, 2)
}

func resolveTheme(name string) Theme {
	names := ThemeNames()
	if !slices.Contains(names, name) {
		return availableThemes[0]
	}
	for _, theme := range availableThemes {
		if theme.Name == name {
			return theme
		}
	}
	return availableThemes[0]
}
