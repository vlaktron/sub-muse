package ui

import (
	"fmt"
	"time"

	"sub-muse/internal/config"
	"sub-muse/internal/player"
	"sub-muse/internal/subsonic"
	"sub-muse/internal/theme"

	table "github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type PaneType int

const (
	PaneBrowser PaneType = iota
	PaneInfo
)

type Model struct {
	cfg    *config.Config
	client *subsonic.Client
	styles Styles

	width, height int

	activeTab   TabType
	activePane  PaneType
	searchInput string

	songs     []subsonic.Song
	artists   []subsonic.Artist
	albums    []subsonic.Album
	playlists []subsonic.Playlist

	selectedAlbum    *subsonic.Album
	selectedArtist   *subsonic.Artist
	selectedPlaylist *subsonic.Playlist

	browser  *Browser
	infoPane *InfoPane

	player     *player.Player
	nowPlaying *subsonic.Song
	queue      []subsonic.Song
	queuePos   int
	elapsed    time.Duration

	coverArtCache map[string]string
	errorMsg      string
}

func NewModel(cfg *config.Config, colors theme.Colors) Model {
	styles := NewStyles(colors.Accent) // Create styles first

	return Model{
		cfg:           cfg,
		client:        subsonic.NewClient(cfg.ServerURL, cfg.Username, cfg.Password, cfg.ClientName),
		styles:        styles,
		browser:       NewBrowser(styles),
		infoPane:      NewInfoPane(styles, make(map[string]string)),
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

		browserHeight := m.height - 10
		browserWidith := m.width * 65 / 100

		m.browser.table.SetHeight(browserHeight - 4)
		m.browser.table.SetWidth(browserWidith - 2)

		return m, nil
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)
	case songsLoadedMsg:
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
		} else {
			m.songs = msg.songs
			// Convert your songs to the table.Row format
			var rows []table.Row
			for i, song := range m.songs {
				rows = append(rows, table.Row{
					fmt.Sprintf("%d", i+1),
					song.Title,
					song.Artist,
					fmt.Sprintf("%d", song.Year),
					song.Genre,
				})
			}
			// Send the rows to the browser component
			m.browser.UpdateData(TabSongs, rows)
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
	case coverArtLoadedMsg:
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
		} else if msg.id != "" {
			return m, renderCoverArtCmd(msg.id, msg.data, m.infoPane.coverArtWidth, m.infoPane.coverArtHeight)
		}
	case playlistsLoadedMsg:
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
		} else {
			m.playlists = msg.playlists
			m.updateBrowserForTab()
		}
	case coverArtRenderedMsg:
		if msg.id != "" {
			m.coverArtCache[msg.id] = msg.rendered
		}
	case albumDetailMsg:
		if msg.err == nil && msg.album != nil {
			m.selectedAlbum = msg.album
			return m, m.infoPane.SetSelectedAlbum(msg.album)
		}
	case artistDetailMsg:
		if msg.err == nil && msg.artist != nil {
			m.selectedArtist = msg.artist
		}
	case playbackStartedMsg:
		m.nowPlaying = &msg.song
		m.queue = []subsonic.Song{msg.song}
		m.queuePos = 0
		m.infoPane.SetSelectedSong(&msg.song)
		return m, m.playerTickCmd()
	case playbackTickMsg:
		state := m.player.GetState()
		if m.nowPlaying != nil && !state.IsPlaying {
			// Song finished naturally
			m.nowPlaying = nil
			m.elapsed = 0
			return m, nil
		}
		m.elapsed = state.Elapsed
		if m.nowPlaying != nil {
			return m, m.playerTickCmd()
		}
	case playbackStoppedMsg:
		m.nowPlaying = nil
		m.elapsed = 0
	case playbackErrorMsg:
		m.errorMsg = msg.err.Error()
		m.nowPlaying = nil
		m.elapsed = 0
	}
	return m, nil
}

func (m Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Global keys work in any pane
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case " ", "p":
		return m.handlePlayPause()
	case "n":
		return m.handleNext()
	case "N":
		return m.handlePrev()
	case "s":
		return m.handleStop()
	}

	// Pane switching
	switch msg.String() {
	case "l", "right":
		if m.activePane == PaneBrowser {
			m.activePane = PaneInfo
			return m, nil
		}
	case "h", "left", "esc":
		if m.activePane == PaneInfo {
			m.activePane = PaneBrowser
			return m, nil
		}
	}

	// Browser-pane keys
	if m.activePane == PaneBrowser {
		switch msg.String() {
		case "tab":
			m.activeTab = m.activeTab.Next()
			m.updateBrowserForTab()
		case "shift+tab":
			m.activeTab = m.activeTab.Prev()
			m.updateBrowserForTab()
		case "1":
			m.activeTab = TabSongs
			m.updateBrowserForTab()
		case "2":
			m.activeTab = TabArtists
			m.updateBrowserForTab()
		case "3":
			m.activeTab = TabAlbums
			m.updateBrowserForTab()
		case "4":
			m.activeTab = TabPlaylists
			m.updateBrowserForTab()
		case "up", "k", "down", "j":
			_, cmd := m.browser.Update(msg)
			return m, cmd
		case "enter":
			return m.handleEnterKey()
		}
	}

	if m.activePane == PaneInfo {
		switch msg.String() {
		case "up", "k", "down", "j":
			_, cmd := m.infoPane.Update(msg)
			return m, cmd
		case "enter":
			model, cmd := m.handleTrackEnter()
			return model, cmd
		}
	}

	return m, nil
}

func (m Model) updateBrowserForTab() {
	var rows []table.Row

	switch m.activeTab {
	case TabSongs:
		for i, s := range m.songs {
			rows = append(rows, table.Row{
				fmt.Sprintf("%d", i+1),
				s.Title,
				s.Artist,
				fmt.Sprintf("%d", s.Year),
				s.Genre,
			})
		}
	case TabArtists:
		for i, a := range m.artists {
			rows = append(rows, table.Row{
				fmt.Sprintf("%d", i+1),
				a.Name,
			})
		}
	case TabAlbums:
		for i, al := range m.albums {
			year := ""
			if al.Year > 0 {
				year = fmt.Sprintf("%d", al.Year)
			}
			rows = append(rows, table.Row{
				fmt.Sprintf("%d", i+1),
				al.Name,
				al.Artist,
				year,
			})
		}
	case TabPlaylists:
		for i, p := range m.playlists {
			rows = append(rows, table.Row{
				fmt.Sprintf("%d", i+1),
				p.Name,
				fmt.Sprintf("%d", p.SongCount),
			})
		}
	}

	// Pass the rows to the browser component
	m.browser.UpdateData(m.activeTab, rows)
}

func (m Model) selectedDataIndex() int {
	// Get the currently selected row from the table component
	row := m.browser.table.SelectedRow()
	if len(row) == 0 {
		return -1
	}

	// Your logic to parse the "#" column (row[0]) back to an int
	var n int
	if _, err := fmt.Sscanf(row[0], "%d", &n); err != nil {
		return -1
	}
	return n - 1
}

func (m Model) handleEnterKey() (tea.Model, tea.Cmd) {
	idx := m.selectedDataIndex()
	if idx < 0 {
		return m, nil
	}
	switch m.activeTab {
	case TabSongs:
		if idx < len(m.songs) {
			return m, tea.Batch(playSongCmd(m.client, m.player, m.songs[idx]), m.infoPane.SetSelectedSong(&m.songs[idx]))
		}
	case TabArtists:
		if idx < len(m.artists) {
			m.selectedArtist = &m.artists[idx]
			return m, loadArtistDetailCmd(m.client, m.selectedArtist.ID)
		}
	case TabAlbums:
		if idx < len(m.albums) {
			m.selectedAlbum = &m.albums[idx]
			cmds := []tea.Cmd{loadAlbumDetailCmd(m.client, m.selectedAlbum.ID)}
			if m.selectedAlbum.CoverArtID != "" {
				cmds = append(cmds, loadCoverArtCmd(m.client, m.selectedAlbum.CoverArtID))
			}
			cmds = append(cmds, m.infoPane.SetSelectedAlbum(m.selectedAlbum))
			return m, tea.Batch(cmds...)
		}
	case TabPlaylists:
		if idx < len(m.playlists) {
			m.selectedPlaylist = &m.playlists[idx]
			return m, m.infoPane.SetSelectedPlaylist(m.selectedPlaylist)
		}
	}
	return m, nil
}

func (m Model) playerTickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return playbackTickMsg{}
	})
}

func (m Model) handlePlayPause() (tea.Model, tea.Cmd) {
	state := m.player.GetState()
	if state.IsPlaying {
		return m, stopSongCmd(m.player)
	}
	// Restart current song if one is queued
	if m.nowPlaying != nil {
		return m, playSongCmd(m.client, m.player, *m.nowPlaying)
	}
	if len(m.queue) > 0 {
		return m, playSongCmd(m.client, m.player, m.queue[m.queuePos])
	}
	return m, nil
}

func (m Model) handleNext() (tea.Model, tea.Cmd) {
	if len(m.queue) == 0 {
		return m, nil
	}

	m.queuePos = (m.queuePos + 1) % len(m.queue)
	song := m.queue[m.queuePos]
	return m, playSongCmd(m.client, m.player, song)
}

func (m Model) handlePrev() (tea.Model, tea.Cmd) {
	if len(m.queue) == 0 {
		return m, nil
	}

	m.queuePos = (m.queuePos - 1 + len(m.queue)) % len(m.queue)
	song := m.queue[m.queuePos]
	return m, playSongCmd(m.client, m.player, song)
}

func (m Model) handleStop() (tea.Model, tea.Cmd) {
	return m, stopSongCmd(m.player)
}

func (m Model) handleTrackEnter() (tea.Model, tea.Cmd) {
	row := m.infoPane.GetTrackTable().SelectedRow()
	if len(row) == 0 {
		return m, nil
	}

	var trackNum int
	if _, err := fmt.Sscanf(row[0], "%d", &trackNum); err != nil {
		return m, nil
	}
	trackNum--

	if m.selectedAlbum != nil && trackNum >= 0 && trackNum < len(m.selectedAlbum.Songs) {
		return m, playSongCmd(m.client, m.player, m.selectedAlbum.Songs[trackNum])
	}
	if m.selectedPlaylist != nil && trackNum >= 0 && trackNum < len(m.selectedPlaylist.Songs) {
		return m, playSongCmd(m.client, m.player, m.selectedPlaylist.Songs[trackNum])
	}

	return m, nil
}

func (m Model) View() string {
	statusHeight := 5
	playerHeight := 5
	contentHeight := m.height - statusHeight - playerHeight

	if contentHeight < 5 {
		return "Terminal too small..."
	}

	statusBar := m.renderStatusBar()
	contentRow := m.renderContentRow(contentHeight)
	playerBar := m.renderPlayerBar()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		statusBar,
		contentRow,
		playerBar,
	)
}

func (m Model) renderContentRow(h int) string {
	browserW := m.width * 65 / 100
	infoW := max(0, m.width-browserW-2)

	browserView := m.browser.Render(m.activeTab, m.searchInput, browserW, h)

	browserStyle := m.styles.BrowserBorder
	infoStyle := m.styles.InfoBorder
	if m.activePane == PaneInfo {
		browserStyle = browserStyle.BorderForeground(lipgloss.Color("#444444"))
	}

	m.infoPane.SetCoverArtDimensions(infoW, h)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		browserStyle.Width(browserW).Height(h).Render(browserView),
		infoStyle.Width(infoW).Height(h).Render(m.infoPane.Render()),
	)
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

func (m Model) renderStatusBar() string {
	status := "Connected: " + m.cfg.ServerURL
	return m.styles.StatusBar.Render(status)
}
