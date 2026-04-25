package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Browser struct {
	table  table.Model
	styles Styles
}

func NewBrowser(s Styles) *Browser {
	t := table.New(table.WithFocused(true))

	// Apply the industrial styles from your styles.go
	ts := table.DefaultStyles()
	ts.Header = ts.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(s.Accent).
		BorderBottom(true).
		Bold(true)

	ts.Selected = ts.Selected.
		Foreground(s.SelectionForeground).
		Background(s.SelectionBackground).
		Bold(false)

	t.SetStyles(ts)

	return &Browser{
		table:  t,
		styles: s,
	}
}

// UpdateData refreshes the table content when tabs switch or search results change
func (b *Browser) UpdateData(tab TabType, rows []table.Row) {
	var columns []table.Column

	switch tab {
	case TabSongs:
		columns = []table.Column{
			{Title: "#", Width: 4},
			{Title: "Title", Width: 25},
			{Title: "Artist", Width: 20},
			{Title: "Yr", Width: 6},
			{Title: "Genre", Width: 15},
		}
	case TabArtists:
		columns = []table.Column{
			{Title: "#", Width: 4},
			{Title: "Artist", Width: 50},
		}
	case TabAlbums:
		columns = []table.Column{
			{Title: "#", Width: 4},
			{Title: "Album", Width: 30},
			{Title: "Artist", Width: 20},
			{Title: "Yr", Width: 6},
		}
	case TabPlaylists:
		columns = []table.Column{
			{Title: "#", Width: 4},
			{Title: "Playlist", Width: 40},
			{Title: "Tracks", Width: 10},
		}
	}

	b.table.SetColumns(columns)
	b.table.SetRows(rows)
}

func (b *Browser) Render(activeTab TabType, searchStr string, width, height int) string {
	// 1. Calculate how many rows our header uses
	// Tab row (1) + Spacer (1) + Search bar (1) + Spacer (1) = 4 rows
	headerHeight := 4
	b.table.SetWidth(width - 2)
	b.table.SetHeight(height - headerHeight)

	// 2. Render Tab Bar (using your TabLabels map)
	var tabs []string
	for i := 0; i < 4; i++ {
		t := TabType(i)
		label := strings.ToUpper(TabLabels[t]) // 90s industrial look likes All-Caps
		style := b.styles.InactiveTab
		if t == activeTab {
			style = b.styles.ActiveTab
		}
		tabs = append(tabs, style.Render(fmt.Sprintf(" %d:%s ", i+1, label)))
	}
	tabRow := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)

	// 3. Render Search Bar
	searchBar := b.styles.SearchInput.Render(fmt.Sprintf(" [ / ] SEARCH: %s ", searchStr))

	// 4. Join it all vertically
	// The table.View() already includes its own headers
	return lipgloss.JoinVertical(
		lipgloss.Left,
		tabRow,
		"", // Empty string creates a spacer line
		searchBar,
		"",
		b.table.View(),
	)
}

// GetSelectedIndex returns the index of the current row in the data slice
func (b *Browser) GetSelectedIndex() int {
	return b.table.Cursor()
}

// Update passes messages (like keyboard events) to the internal table
func (b *Browser) Update(msg tea.Msg) (table.Model, tea.Cmd) {
	var cmd tea.Cmd
	b.table, cmd = b.table.Update(msg)
	return b.table, cmd
}
