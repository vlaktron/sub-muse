package ui

import (
	"fmt"

	"sub-muse/internal/config"
	"sub-muse/internal/subsonic"
	"sub-muse/internal/theme"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	cfg    *config.Config
	client *subsonic.Client
	styles Styles

	width, height int

	activeTab   TabType
	searchInput string

	songs     []subsonic.Song
	artists   []subsonic.Artist   // nolint:unused
	albums    []subsonic.Album    // nolint:unused
	playlists []subsonic.Playlist // nolint:unused

	selectedSong   *subsonic.Song // nolint:unused
	selectedAlbum  *subsonic.Album
	selectedArtist *subsonic.Artist
}

func NewModel(cfg *config.Config, colors theme.Colors) Model {
	return Model{
		cfg:    cfg,
		client: subsonic.NewClient(cfg.ServerURL, cfg.Username, cfg.Password, cfg.ClientName),
		styles: NewStyles(colors.Accent),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	statusBar := m.renderStatusBar()
	contentRow := m.renderContentRow()
	playerBar := m.renderPlayerBar()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		statusBar,
		contentRow,
		playerBar,
	)
}

func (m Model) renderStatusBar() string {
	status := "Connected: " + m.cfg.ServerURL
	return m.styles.StatusBar.Render(status)
}

func (m Model) renderContentRow() string {
	contentHeight := m.height - 1 - 3

	browserW := m.width * 65 / 100
	infoW := m.width - browserW - 2

	browserContent := m.renderBrowser()
	infoContent := m.renderInfoPane()

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.styles.BrowserBorder.Width(browserW).Height(contentHeight).Render(browserContent),
		m.styles.InfoBorder.Width(infoW).Height(contentHeight).Render(infoContent),
	)
}

func (m Model) renderBrowser() string {
	headers := []string{"#", "Title", "Artist", "Yr", "Genre"}
	rows := make([][]string, 0)

	for i, song := range m.songs {
		year := ""
		if song.Year > 0 {
			year = fmt.Sprintf("%d", song.Year)
		}
		rows = append(rows, []string{
			fmt.Sprintf("%d", i+1),
			truncateString(song.Title, 20),
			truncateString(song.Artist, 15),
			year,
			truncateString(song.Genre, 10),
		})
	}

	browser := &Browser{
		Headers:       headers,
		Rows:          rows,
		SelectedIndex: 0,
		TabType:       m.activeTab,
		Filter:        m.searchInput,
		ScrollOffset:  0,
		Styles:        m.styles,
	}

	return browser.Render(m.width*65/100, 20)
}

func (m Model) renderInfoPane() string {
	if m.selectedAlbum != nil {
		return m.renderAlbumInfo()
	}
	if m.selectedArtist != nil {
		return m.renderArtistInfo()
	}
	return "Select an item to view details"
}

func (m Model) renderAlbumInfo() string {
	if m.selectedAlbum == nil {
		return ""
	}

	info := "No cover art available"
	return info
}

func (m Model) renderArtistInfo() string {
	if m.selectedArtist == nil {
		return ""
	}

	info := "No cover art available"
	return info
}

func (m Model) renderPlayerBar() string {
	playerBar := &PlayerBar{
		Progress:    0.0,
		CurrentTime: 0,
		Duration:    0,
		Playing:     false,
		SongTitle:   "No song playing",
		Artist:      "",
		Album:       "",
		Styles:      m.styles,
	}

	return playerBar.Render(m.width)
}
