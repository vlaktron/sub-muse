package ui

import (
	"fmt"
	"time"

	"sub-muse/internal/config"
	"sub-muse/internal/player"
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

	player     *player.Player
	nowPlaying *subsonic.Song
	queue      []subsonic.Song
	queuePos   int
	elapsed    time.Duration

	coverArtCache map[string]string
	errorMsg      string
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
		player:        player.NewPlayer(),
		coverArtCache: make(map[string]string),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		loadSongsCmd(m.client),
		loadArtistsCmd(m.client),
		loadAlbumsCmd(m.client),
		loadPlaylistsCmd(m.client),
		m.playerTickCmd(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)
	case songsLoadedMsg:
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
		} else {
			m.songs = msg.songs
			m.updateBrowserForTab()
		}
	case artistsLoadedMsg:
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
		} else {
			m.artists = msg.artists
			m.updateBrowserForTab()
		}
	case albumsLoadedMsg:
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
		} else {
			m.albums = msg.albums
			m.updateBrowserForTab()
			if m.selectedAlbum != nil && m.selectedAlbum.CoverArtID != "" {
				return m, loadCoverArtCmd(m.client, m.selectedAlbum.CoverArtID)
			}
		}
	case playlistsLoadedMsg:
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
		} else {
			m.playlists = msg.playlists
			m.updateBrowserForTab()
		}
	case coverArtLoadedMsg:
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
		} else if msg.id != "" {
			m.coverArtCache[msg.id] = string(msg.data)
			return m, renderCoverArtCmd(msg.data, 200, 200)
		}
	case coverArtRenderedMsg:
		if m.selectedAlbum != nil && m.selectedAlbum.CoverArtID != "" {
			if rendered, ok := m.coverArtCache[m.selectedAlbum.CoverArtID]; ok {
				m.selectedAlbum.CoverArtID = rendered
			}
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
		m.nowPlaying = &msg.song
		m.queue = []subsonic.Song{msg.song}
		m.queuePos = 0
	case playbackTickMsg:
		m.elapsed = m.player.GetState().Elapsed
	case playbackStoppedMsg:
		m.nowPlaying = nil
		m.elapsed = 0
	case playbackErrorMsg:
		m.nowPlaying = nil
		m.elapsed = 0
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
	case " ", "p":
		return m.handlePlayPause()
	case "n":
		return m.handleNext()
	case "N":
		return m.handlePrev()
	case "s":
		return m.handleStop()
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
		cmds := []tea.Cmd{
			loadAlbumDetailCmd(m.client, m.selectedAlbum.ID),
		}
		if m.selectedAlbum.CoverArtID != "" {
			cmds = append(cmds, loadCoverArtCmd(m.client, m.selectedAlbum.CoverArtID))
		}
		return m, tea.Batch(cmds...)
	}
	if m.activeTab == TabArtists && m.selectedArtist != nil {
		cmds := []tea.Cmd{
			loadArtistDetailCmd(m.client, m.selectedArtist.ID),
		}
		if m.selectedArtist.CoverArtID != "" {
			cmds = append(cmds, loadCoverArtCmd(m.client, m.selectedArtist.CoverArtID))
		}
		return m, tea.Batch(cmds...)
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

func (m Model) playerTickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return playbackTickMsg{}
	})
}

func (m Model) handlePlayPause() (tea.Model, tea.Cmd) {
	if m.nowPlaying == nil {
		if len(m.queue) > 0 {
			song := m.queue[m.queuePos]
			return m, playSongCmd(m.client, song)
		}
		return m, nil
	}

	state := m.player.GetState()
	if state.IsPlaying {
		return m, stopSongCmd()
	}
	return m, nil
}

func (m Model) handleNext() (tea.Model, tea.Cmd) {
	if len(m.queue) == 0 {
		return m, nil
	}

	m.queuePos = (m.queuePos + 1) % len(m.queue)
	song := m.queue[m.queuePos]
	return m, playSongCmd(m.client, song)
}

func (m Model) handlePrev() (tea.Model, tea.Cmd) {
	if len(m.queue) == 0 {
		return m, nil
	}

	m.queuePos = (m.queuePos - 1 + len(m.queue)) % len(m.queue)
	song := m.queue[m.queuePos]
	return m, playSongCmd(m.client, song)
}

func (m Model) handleStop() (tea.Model, tea.Cmd) {
	return m, stopSongCmd()
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
	if m.errorMsg != "" {
		status = "Error: " + m.errorMsg
	}
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

	info := m.selectedAlbum.CoverArtID
	if info == "" {
		info = "No cover art available"
	}
	return info
}

func (m Model) renderArtistInfo() string {
	if m.selectedArtist == nil {
		return ""
	}

	info := m.selectedArtist.CoverArtID
	if info == "" {
		info = "No cover art available"
	}
	return info
}

func (m Model) renderPlayerBar() string {
	if m.nowPlaying == nil {
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

	song := *m.nowPlaying
	state := m.player.GetState()
	elapsed := int(state.Elapsed.Seconds())
	duration := song.Duration

	var progress float64
	if duration > 0 {
		progress = float64(elapsed) / float64(duration)
	}

	playerBar := &PlayerBar{
		Progress:    progress,
		CurrentTime: elapsed,
		Duration:    duration,
		Playing:     state.IsPlaying,
		SongTitle:   song.Title,
		Artist:      song.Artist,
		Album:       song.Album,
		Styles:      m.styles,
	}

	return playerBar.Render(m.width)
}
