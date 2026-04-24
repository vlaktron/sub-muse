package ui

type TabType int

const (
	TabSongs TabType = iota
	TabArtists
	TabAlbums
	TabPlaylists
)

var TabLabels = map[TabType]string{
	TabSongs:     "Songs",
	TabArtists:   "Artists",
	TabAlbums:    "Albums",
	TabPlaylists: "Playlists",
}

func (t TabType) Next() TabType {
	return TabType((int(t) + 1) % 4)
}

func (t TabType) Prev() TabType {
	return TabType((int(t) + 3) % 4)
}
