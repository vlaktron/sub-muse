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
	artists   []subsonic.Artist
	albums    []subsonic.Album
	playlists []subsonic.Playlist

	selectedSong   *subsonic.Song
	selectedAlbum  *subsonic.Album
	selectedArtist *subsonic.Artist

	browser *Browser
}

func NewModel(cfg *config.Config, colors theme.Colors) Model {
	return Model{
		cfg:    cfg,
		client: subsonic.NewClient(cfg.ServerURL, cfg.Username, cfg.Password, cfg.ClientName),
		styles: NewStyles(colors.Accent),
		browser: &Browser{
			Headers:       []string{"#", "Title", "Artist", "Yr", "Genre"},
			SelectedIndex: 0,
			TabType:       TabSongs,
			Filter:        "",
			ScrollOffset:  0,
			Styles:        NewStyles(colors.Accent),
		},
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		loadSongsCmd(m.client),
		loadArtistsCmd(m.client),
		loadAlbumsCmd(m.client),
		loadPlaylistsCmd(m.client),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)
	case songsLoadedMsg:
		if msg.err == nil {
			m.songs = msg.songs
			m.updateBrowserForTab()
		}
	case artistsLoadedMsg:
		if msg.err == nil {
			m.artists = msg.artists
			m.updateBrowserForTab()
		}
	case albumsLoadedMsg:
		if msg.err == nil {
			m.albums = msg.albums
			m.updateBrowserForTab()
		}
	case playlistsLoadedMsg:
		if msg.err == nil {
			m.playlists = msg.playlists
			m.updateBrowserForTab()
		}
	case albumDetailMsg:
		if msg.err == nil && msg.album != nil {
			m.selectedAlbum = msg.album
		}
	case artistDetailMsg:
		if msg.err == nil && msg.artist != nil {
			m.selectedArtist = msg.artist
		}
	case playbackStartedMsg:
		m.selectedSong = &msg.song
	}
	return m, nil
}

func (m Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "tab":
		m.activeTab = m.activeTab.Next()
		m.updateBrowserForTab()
		return m, nil
	case "shift+tab":
		m.activeTab = m.activeTab.Prev()
		m.updateBrowserForTab()
		return m, nil
	case "1":
		m.activeTab = TabSongs
		m.updateBrowserForTab()
		return m, nil
	case "2":
		m.activeTab = TabArtists
		m.updateBrowserForTab()
		return m, nil
	case "3":
		m.activeTab = TabAlbums
		m.updateBrowserForTab()
		return m, nil
	case "4":
		m.activeTab = TabPlaylists
		m.updateBrowserForTab()
		return m, nil
	case "up", "k":
		m.browser.Up()
		return m, nil
	case "down", "j":
		m.browser.Down()
		return m, nil
	case "enter":
		return m.handleEnterKey()
	}
	return m, nil
}

func (m Model) updateBrowserForTab() {
	m.browser.TabType = m.activeTab
	switch m.activeTab {
	case TabSongs:
		m.browser.Rows = m.formatSongsForTable()
		if len(m.browser.Rows) > 0 && m.selectedSong != nil {
			idx := m.findSongIndex(*m.selectedSong)
			if idx < len(m.browser.Rows) {
				m.browser.SelectedIndex = idx
			}
		}
	case TabArtists:
		m.browser.Rows = m.formatArtistsForTable()
	case TabAlbums:
		m.browser.Rows = m.formatAlbumsForTable()
	case TabPlaylists:
		m.browser.Rows = m.formatPlaylistsForTable()
	}
	m.browser.ScrollToCenter()
}

func (m Model) formatSongsForTable() [][]string {
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
	return rows
}

func (m Model) formatArtistsForTable() [][]string {
	rows := make([][]string, 0)
	for i, artist := range m.artists {
		rows = append(rows, []string{
			fmt.Sprintf("%d", i+1),
			truncateString(artist.Name, 30),
			"",
			"",
			"",
		})
	}
	return rows
}

func (m Model) formatAlbumsForTable() [][]string {
	rows := make([][]string, 0)
	for i, album := range m.albums {
		year := ""
		if album.Year > 0 {
			year = fmt.Sprintf("%d", album.Year)
		}
		rows = append(rows, []string{
			fmt.Sprintf("%d", i+1),
			truncateString(album.Name, 20),
			truncateString(album.Artist, 15),
			year,
			"",
		})
	}
	return rows
}

func (m Model) formatPlaylistsForTable() [][]string {
	rows := make([][]string, 0)
	for i, playlist := range m.playlists {
		rows = append(rows, []string{
			fmt.Sprintf("%d", i+1),
			truncateString(playlist.Name, 30),
			fmt.Sprintf("%d", playlist.SongCount),
			"",
			"",
		})
	}
	return rows
}

func (m Model) handleEnterKey() (tea.Model, tea.Cmd) {
	if m.activeTab == TabAlbums && m.selectedAlbum != nil {
		return m, loadAlbumDetailCmd(m.client, m.selectedAlbum.ID)
	}
	if m.activeTab == TabArtists && m.selectedArtist != nil {
		return m, loadArtistDetailCmd(m.client, m.selectedArtist.ID)
	}
	if m.activeTab == TabSongs && m.selectedSong != nil {
		return m, playSongCmd(m.client, *m.selectedSong)
	}
	return m, nil
}

func (m Model) findSongIndex(song subsonic.Song) int {
	for i, s := range m.songs {
		if s.ID == song.ID {
			return i
		}
	}
	return 0
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
	m.browser.Headers = []string{"#", "Title", "Artist", "Yr", "Genre"}
	m.browser.Filter = m.searchInput

	switch m.activeTab {
	case TabSongs:
		m.browser.Rows = m.formatSongsForTable()
	case TabArtists:
		m.browser.Rows = m.formatArtistsForTable()
	case TabAlbums:
		m.browser.Rows = m.formatAlbumsForTable()
	case TabPlaylists:
		m.browser.Rows = m.formatPlaylistsForTable()
	}

	if m.browser.SelectedIndex >= len(m.browser.Rows) {
		m.browser.SelectedIndex = 0
	}

	return m.browser.Render(m.width*65/100, m.height-4)
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
