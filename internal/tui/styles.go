package tui

import "github.com/charmbracelet/lipgloss"

// Color palette
var (
	primaryColor   = lipgloss.Color("99")  // Purple
	secondaryColor = lipgloss.Color("39")  // Cyan
	accentColor    = lipgloss.Color("212") // Pink
	successColor   = lipgloss.Color("82")  // Green
	warningColor   = lipgloss.Color("214") // Orange
	errorColor     = lipgloss.Color("196") // Red
	dimColor       = lipgloss.Color("240") // Gray
	textColor      = lipgloss.Color("252") // Light gray
)

// Style definitions
var (
	// Title style
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			MarginBottom(1)

	// Help style (footer)
	helpStyle = lipgloss.NewStyle().
			Foreground(dimColor)

	// Selected item style
	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(secondaryColor)

	// Normal item style
	normalStyle = lipgloss.NewStyle().
			Foreground(textColor)

	// Disabled/dim style
	disabledStyle = lipgloss.NewStyle().
			Foreground(dimColor)

	// Widget enabled indicator
	enabledIndicator = lipgloss.NewStyle().
				Foreground(successColor).
				SetString("●")

	// Widget disabled indicator
	disabledIndicator = lipgloss.NewStyle().
				Foreground(dimColor).
				SetString("○")

	// Section header style
	sectionStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			MarginTop(1)

	// Preview box style
	previewBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(dimColor).
			Padding(0, 1)

	// Input field style
	inputStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(dimColor).
			Padding(0, 1)

	// Focused input style
	focusedInputStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(secondaryColor).
				Padding(0, 1)

	// Error message style
	errorStyle = lipgloss.NewStyle().
			Foreground(errorColor)

	// Success message style
	successStyle = lipgloss.NewStyle().
			Foreground(successColor)

	// Warning style
	warningStyle = lipgloss.NewStyle().
			Foreground(warningColor)

	// Cursor style
	cursorStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			SetString("▸")

	// Line header style
	lineHeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor)

	// Button style
	buttonStyle = lipgloss.NewStyle().
			Foreground(textColor).
			Background(dimColor).
			Padding(0, 1)

	// Active button style
	activeButtonStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("0")).
				Background(secondaryColor).
				Padding(0, 1)
)
