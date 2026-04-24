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

func loadAlbumDetailCmd(client *subsonic.Client, albumID string) tea.Cmd {
	return func() tea.Msg {
		album, err := client.GetAlbum(albumID)
		return albumDetailMsg{album: album, err: err}
	}
}

func loadArtistDetailCmd(client *subsonic.Client, artistID string) tea.Cmd {
	return func() tea.Msg {
		artist, err := client.GetArtist(artistID)
		return artistDetailMsg{artist: artist, err: err}
	}
}

func playSongCmd(client *subsonic.Client, song subsonic.Song) tea.Cmd {
	return func() tea.Msg {
		return playbackStartedMsg{song: song}
	}
}

func stopSongCmd() tea.Cmd {
	return func() tea.Msg {
		return playbackStoppedMsg{}
	}
}
