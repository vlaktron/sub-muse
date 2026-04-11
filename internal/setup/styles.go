package setup

import "github.com/charmbracelet/lipgloss"

var (
	primaryColor = lipgloss.Color("#7D56F4")
	errorColor   = lipgloss.Color("#FF5370")
	successColor = lipgloss.Color("#00C853")
	textColor    = lipgloss.Color("#CCCCCC")

	baseStyle = lipgloss.NewStyle().
		Width(60).
		Margin(1).
		Padding(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Align(lipgloss.Center)

	titleStyle = baseStyle.Copy().
		BorderForeground(successColor)

	errorStyle = baseStyle.Copy().
		BorderForeground(errorColor)

	wizardTitleStyle = lipgloss.NewStyle().
		Width(60).
		MarginBottom(1).
		Align(lipgloss.Center).
		Bold(true).
		Foreground(primaryColor)

	inputStyle = lipgloss.NewStyle().
		PaddingLeft(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#666"))

	focusedInputStyle = inputStyle.Copy().
		BorderForeground(primaryColor)

	buttonStyle = lipgloss.NewStyle().
		MarginRight(2).
		PaddingLeft(2).
		PaddingRight(2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Background(primaryColor).
		Foreground(lipgloss.Color("#000"))

	focusedButtonStyle = buttonStyle.Copy().
		Background(lipgloss.Color("#9D7DFF"))

	helpStyle = lipgloss.NewStyle().
		MarginTop(1).
		Foreground(lipgloss.Color("#888")).
		Align(lipgloss.Center)

	spinnerStyle = lipgloss.NewStyle().
		MarginTop(1).
		Align(lipgloss.Center)
)
