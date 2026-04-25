package ui

import (
	"fmt"
	"strings"

	"sub-muse/internal/subsonic"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type InfoPane struct {
	styles           Styles
	coverArtCache    map[string]string
	coverArtWidth    int
	coverArtHeight   int
	trackTable       table.Model
	selectedSong     *subsonic.Song
	selectedAlbum    *subsonic.Album
	selectedPlaylist *subsonic.Playlist
}

func NewInfoPane(s Styles, cache map[string]string) *InfoPane {
	t := table.New(table.WithFocused(true))

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

	return &InfoPane{
		styles:         s,
		coverArtCache:  cache,
		coverArtWidth:  40,
		coverArtHeight: 20,
		trackTable:     t,
	}
}

func (p *InfoPane) SetCoverArtDimensions(width, height int) {
	p.coverArtWidth = width
	p.coverArtHeight = height
}

func (p *InfoPane) SetSelectedSong(song *subsonic.Song) tea.Cmd {
	p.selectedSong = song
	p.selectedAlbum = nil
	p.selectedPlaylist = nil
	p.trackTable.SetRows([]table.Row{})
	return nil
}

func (p *InfoPane) SetSelectedAlbum(album *subsonic.Album) tea.Cmd {
	p.selectedAlbum = album
	p.selectedSong = nil
	p.selectedPlaylist = nil

	var rows []table.Row
	for i, track := range album.Songs {
		rows = append(rows, table.Row{
			fmt.Sprintf("%d", i+1),
			track.Title,
			formatDuration(track.Duration),
		})
	}
	p.trackTable.SetColumns([]table.Column{
		{Title: "#", Width: 4},
		{Title: "Title", Width: 30},
		{Title: "Duration", Width: 8},
	})
	p.trackTable.SetRows(rows)

	return nil
}

func (p *InfoPane) SetSelectedPlaylist(playlist *subsonic.Playlist) tea.Cmd {
	p.selectedPlaylist = playlist
	p.selectedSong = nil
	p.selectedAlbum = nil

	var rows []table.Row
	for i, track := range playlist.Songs {
		rows = append(rows, table.Row{
			fmt.Sprintf("%d", i+1),
			track.Title,
			formatDuration(track.Duration),
		})
	}
	p.trackTable.SetColumns([]table.Column{
		{Title: "#", Width: 4},
		{Title: "Title", Width: 30},
		{Title: "Duration", Width: 8},
	})
	p.trackTable.SetRows(rows)

	return nil
}

func formatDuration(seconds int) string {
	minutes := seconds / 60
	secs := seconds % 60
	return fmt.Sprintf("%d:%02d", minutes, secs)
}

func (p *InfoPane) Render() string {
	if p.selectedSong != nil {
		return p.renderSong()
	}
	if p.selectedAlbum != nil {
		return p.renderAlbum()
	}
	if p.selectedPlaylist != nil {
		return p.renderPlaylist()
	}
	return "Select an item to view details"
}

func (p *InfoPane) renderSong() string {
	song := p.selectedSong
	if song == nil {
		return ""
	}

	var lines []string

	if art, ok := p.coverArtCache[song.CoverArtID]; ok && art != "" {
		lines = append(lines, art)
	}

	if song.Artist != "" {
		lines = append(lines, p.styles.InfoTitle.Render(song.Artist))
	}
	if song.Album != "" {
		lines = append(lines, p.styles.InfoTitle.Render(song.Album))
	}
	if song.Duration > 0 {
		lines = append(lines, fmt.Sprintf("Duration: %s", formatDuration(song.Duration)))
	}

	return strings.Join(lines, "\n")
}

func (p *InfoPane) renderAlbum() string {
	album := p.selectedAlbum
	if album == nil {
		return ""
	}

	var lines []string

	if art, ok := p.coverArtCache[album.CoverArtID]; ok && art != "" {
		lines = append(lines, art)
	}

	lines = append(lines, p.styles.InfoTitle.Render(album.Name))
	if album.Artist != "" {
		lines = append(lines, p.styles.InfoTitle.Render(album.Artist))
	}
	if album.Year > 0 {
		lines = append(lines, fmt.Sprintf("Year: %d", album.Year))
	}
	if album.SongCount > 0 {
		lines = append(lines, fmt.Sprintf("Tracks: %d", album.SongCount))
	}

	tableHeight := p.coverArtHeight - len(lines)
	if tableHeight < 5 {
		tableHeight = 5
	}
	p.trackTable.SetHeight(tableHeight)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		strings.Join(lines, "\n"),
		"",
		p.trackTable.View(),
	)
}

func (p *InfoPane) renderPlaylist() string {
	playlist := p.selectedPlaylist
	if playlist == nil {
		return ""
	}

	var lines []string

	lines = append(lines, p.styles.InfoTitle.Render(playlist.Name))
	if playlist.SongCount > 0 {
		lines = append(lines, fmt.Sprintf("Tracks: %d", playlist.SongCount))
	}
	if playlist.Duration > 0 {
		lines = append(lines, fmt.Sprintf("Duration: %s", formatDuration(playlist.Duration)))
	}

	tableHeight := p.coverArtHeight - len(lines)
	if tableHeight < 5 {
		tableHeight = 5
	}
	p.trackTable.SetHeight(tableHeight)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		strings.Join(lines, "\n"),
		"",
		p.trackTable.View(),
	)
}

func (p *InfoPane) Update(msg tea.Msg) (table.Model, tea.Cmd) {
	var cmd tea.Cmd
	p.trackTable, cmd = p.trackTable.Update(msg)
	return p.trackTable, cmd
}

func (p *InfoPane) GetTrackTable() table.Model {
	return p.trackTable
}
