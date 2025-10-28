package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	primaryColor   = lipgloss.Color("#7D56F4")
	secondaryColor = lipgloss.Color("#FF79C6")
	accentColor    = lipgloss.Color("#50FA7B")
	warningColor   = lipgloss.Color("#FFB86C")
	errorColor     = lipgloss.Color("#FF5555")
	mutedColor     = lipgloss.Color("#6272A4")
	bgColor        = lipgloss.Color("#282A36")
	fgColor        = lipgloss.Color("#F8F8F2")

	// Base styles
	baseStyle = lipgloss.NewStyle().
			Foreground(fgColor)

	// Panel styles
	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 2)

	activePanelStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(accentColor).
				Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Padding(0, 1)

	// List item styles
	selectedItemStyle = lipgloss.NewStyle().
				Foreground(accentColor).
				Bold(true).
				PaddingLeft(2)

	normalItemStyle = lipgloss.NewStyle().
			Foreground(fgColor).
			PaddingLeft(2)

	currentWorktreeStyle = lipgloss.NewStyle().
				Foreground(secondaryColor).
				Bold(true).
				PaddingLeft(2)

	// Detail styles
	detailKeyStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)

	detailValueStyle = lipgloss.NewStyle().
				Foreground(fgColor)

	// Help/Status bar
	helpStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Padding(0, 1)

	statusStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true)

	// Modal styles
	modalStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 2).
			Width(60)

	modalTitleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Align(lipgloss.Center)

	inputLabelStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Bold(true)

	buttonStyle = lipgloss.NewStyle().
			Foreground(fgColor).
			Background(primaryColor).
			Padding(0, 2).
			MarginRight(1)

	selectedButtonStyle = lipgloss.NewStyle().
				Foreground(fgColor).
				Background(accentColor).
				Padding(0, 2).
				MarginRight(1).
				Bold(true)

	cancelButtonStyle = lipgloss.NewStyle().
				Foreground(fgColor).
				Background(mutedColor).
				Padding(0, 2)

	selectedCancelButtonStyle = lipgloss.NewStyle().
					Foreground(fgColor).
					Background(errorColor).
					Padding(0, 2).
					Bold(true)
)
