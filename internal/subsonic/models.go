package subsonic

import "time"

type SubsonicResponse struct {
	Status string `json:"status"`
	Error  *Error `json:"error,omitempty"`
}

type MusicFolder struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Song struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Artist      string `json:"artist"`
	Album       string `json:"album"`
	AlbumID     string `json:"albumId"`
	ArtistID    string `json:"artistId"`
	Duration    int    `json:"duration"`
	Track       int    `json:"track"`
	Year        int    `json:"year"`
	CoverArtID  string `json:"coverArt"`
	Size        int64  `json:"size"`
	BitRate     int    `json:"bitRate"`
	ContentType string `json:"contentType"`
	Suffix      string `json:"suffix"`
}

type Album struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Artist     string     `json:"artist"`
	ArtistID   string     `json:"artistId"`
	SongCount  int        `json:"songCount"`
	Duration   int        `json:"duration"`
	CoverArtID string     `json:"coverArt"`
	Year       int        `json:"year"`
	Starred    *time.Time `json:"starred,omitempty"`
	Songs      []Song     `json:"song,omitempty"`
}

type Artist struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	AlbumCount int     `json:"albumCount"`
	CoverArtID string  `json:"coverArt"`
	Albums     []Album `json:"album,omitempty"`
}

type MusicDirectory struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Children []Child `json:"child,omitempty"`
}

type Child struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Artist      string `json:"artist"`
	Album       string `json:"album"`
	Duration    int    `json:"duration"`
	CoverArtID  string `json:"coverArt"`
	Size        int64  `json:"size"`
	BitRate     int    `json:"bitRate"`
	ContentType string `json:"contentType"`
	Suffix      string `json:"suffix"`
	IsDir       bool   `json:"isDir"`
}
