package ui

import (
	"sub-muse/internal/subsonic"
)

type songsLoadedMsg struct {
	songs []subsonic.Song
	err   error
}

type artistsLoadedMsg struct {
	artists []subsonic.Artist
	err     error
}

type albumsLoadedMsg struct {
	albums []subsonic.Album
	err    error
}

type playlistsLoadedMsg struct {
	playlists []subsonic.Playlist
	err       error
}

type albumDetailMsg struct {
	album *subsonic.Album
	err   error
}

type artistDetailMsg struct {
	artist *subsonic.Artist
	err    error
}

type coverArtLoadedMsg struct {
	id   string
	data []byte
	err  error
}

type coverArtRenderedMsg struct {
	rendered string
}

type playbackStartedMsg struct {
	song subsonic.Song
}

type playbackStoppedMsg struct{}

type playbackTickMsg struct{}

type playbackErrorMsg struct{}
