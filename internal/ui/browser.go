package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Browser struct {
	Headers       []string
	Rows          [][]string
	SelectedIndex int
	TabType       TabType
	Filter        string
	ScrollOffset  int
	Styles        Styles
	lastHeight    int
}

func (b *Browser) Render(width, height int) string {
	if width < 10 || height < 5 {
		return ""
	}

	b.lastHeight = height

	var sb strings.Builder

	tabLabel := TabLabels[b.TabType]
	tabStyle := lipgloss.NewStyle().
		Foreground(b.Styles.Accent).
		Bold(true).
		PaddingLeft(1)
	sb.WriteString(tabStyle.Render("[" + tabLabel + "]"))
	sb.WriteString("\n")

	if b.Filter != "" {
		filterStyle := lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(b.Styles.Accent).
			Padding(0, 1)
		sb.WriteString(filterStyle.Render("Filter: " + b.Filter))
		sb.WriteString("\n")
	}

	filteredRows := b.FilterRows()
	visibleRows := b.VisibleRows()

	if len(filteredRows) == 0 {
		sb.WriteString(lipgloss.NewStyle().Padding(1).Render("No results"))
		return sb.String()
	}

	headerStyle := lipgloss.NewStyle().
		Foreground(b.Styles.Accent).
		Bold(true)

	headerRow := headerStyle.Render(lipgloss.JoinHorizontal(
		lipgloss.Left,
		b.Headers...,
	))
	sb.WriteString(headerRow)
	sb.WriteString("\n")

	startIndex := b.ScrollOffset
	endIndex := startIndex + visibleRows

	if endIndex > len(filteredRows) {
		endIndex = len(filteredRows)
	}

	for i := startIndex; i < endIndex; i++ {
		row := filteredRows[i]
		rowStr := lipgloss.JoinHorizontal(lipgloss.Left, row...)

		if i == b.SelectedIndex {
			rowStr = b.Styles.Table.
				Background(b.Styles.SelectionBackground).
				Foreground(b.Styles.SelectionForeground).
				Render(rowStr)
		} else {
			rowStr = b.Styles.Table.Render(rowStr)
		}

		sb.WriteString(rowStr)
		sb.WriteString("\n")
	}

	return sb.String()
}

func (b *Browser) FilterRows() [][]string {
	if b.Filter == "" {
		return b.Rows
	}

	filterLower := strings.ToLower(b.Filter)
	var filtered [][]string

	for _, row := range b.Rows {
		match := false
		for _, cell := range row {
			if strings.Contains(strings.ToLower(cell), filterLower) {
				match = true
				break
			}
		}
		if match {
			filtered = append(filtered, row)
		}
	}

	return filtered
}

func (b *Browser) VisibleRows() int {
	visible := 1
	if b.Filter != "" {
		visible++
	}
	return max(0, b.lastHeight-3-visible)
}

func (b *Browser) ScrollToCenter() {
	visibleRows := b.VisibleRows()
	if visibleRows <= 0 {
		return
	}

	if b.SelectedIndex >= b.ScrollOffset+visibleRows {
		b.ScrollOffset = b.SelectedIndex - visibleRows + 1
	} else if b.SelectedIndex < b.ScrollOffset {
		b.ScrollOffset = b.SelectedIndex
	}
}

func (b *Browser) Up() {
	if len(b.Rows) == 0 {
		return
	}

	if b.SelectedIndex > 0 {
		b.SelectedIndex--
	} else {
		b.SelectedIndex = len(b.Rows) - 1
	}
	b.ScrollToCenter()
}

func (b *Browser) Down() {
	if len(b.Rows) == 0 {
		return
	}

	if b.SelectedIndex < len(b.Rows)-1 {
		b.SelectedIndex++
	} else {
		b.SelectedIndex = 0
	}
	b.ScrollToCenter()
}
