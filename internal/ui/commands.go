package ui

import (
	"sub-muse/internal/subsonic"

	tea "github.com/charmbracelet/bubbletea"
)

func loadSongsCmd(client *subsonic.Client) tea.Cmd {
	return func() tea.Msg {
		songs, err := client.GetSongs()
		return songsLoadedMsg{songs: songs, err: err}
	}
}

func loadArtistsCmd(client *subsonic.Client) tea.Cmd {
	return func() tea.Msg {
		artists, err := client.GetArtists()
		return artistsLoadedMsg{artists: artists, err: err}
	}
}

func loadAlbumsCmd(client *subsonic.Client) tea.Cmd {
	return func() tea.Msg {
		albums, err := client.GetAlbums()
		return albumsLoadedMsg{albums: albums, err: err}
	}
}

func loadPlaylistsCmd(client *subsonic.Client) tea.Cmd {
	return func() tea.Msg {
		playlists, err := client.GetPlaylists()
		return playlistsLoadedMsg{playlists: playlists, err: err}
	}
}
