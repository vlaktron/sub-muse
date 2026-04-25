/*
* Package for lipgloss styles
 */

package ui

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	// Accent color from theme
	Accent lipgloss.Color

	// Text colors
	Foreground lipgloss.Color
	Background lipgloss.Color

	// Selection colors
	SelectionForeground lipgloss.Color
	SelectionBackground lipgloss.Color

	// Tab Styles
	ActiveTab   lipgloss.Style
	InactiveTab lipgloss.Style

	// Browser styles
	BrowserBorder lipgloss.Style
	BrowserTitle  lipgloss.Style
	Table         lipgloss.Style

	// Info pane styles
	InfoBorder lipgloss.Style
	InfoTitle  lipgloss.Style

	// Player bar styles
	PlayerBar      lipgloss.Style
	ProgressFilled lipgloss.Style
	ProgressEmpty  lipgloss.Style

	// Status bar styles
	StatusBar lipgloss.Style

	// Search styles
	SearchInput lipgloss.Style
	SearchLabel lipgloss.Style
}

func NewStyles(accent string) Styles {
	accentColor := lipgloss.Color(accent)

	return Styles{
		Accent: accentColor,

		Foreground: lipgloss.Color("#FFFFFF"),
		Background: lipgloss.Color("#000000"),

		SelectionForeground: lipgloss.Color("#000000"),
		SelectionBackground: accentColor,

		ActiveTab: lipgloss.NewStyle().
			Background(accentColor).
			Foreground(lipgloss.Color("#000000")). // Black text on accent background
			Padding(0, 1).
			Bold(true).
			MarginRight(1),

		InactiveTab: lipgloss.NewStyle().
			Foreground(accentColor).
			Padding(0, 1).
			MarginRight(1),

		BrowserBorder: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(accentColor).
			Padding(0, 1),

		BrowserTitle: lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true),

		Table: lipgloss.NewStyle(),

		InfoBorder: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(accentColor).
			Padding(1),

		InfoTitle: lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true),

		PlayerBar: lipgloss.NewStyle().
			Background(accentColor).
			Foreground(lipgloss.Color("#000000")).
			Padding(1),

		ProgressFilled: lipgloss.NewStyle().
			Foreground(accentColor),

		ProgressEmpty: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#333333")),

		StatusBar: lipgloss.NewStyle().
			Background(accentColor).
			Foreground(lipgloss.Color("#000000")).
			Padding(0, 1),

		SearchInput: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(accentColor),
	}
}
