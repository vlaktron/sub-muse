package subsonic

import (
	"encoding/xml"
	"time"
)

type SubsonicResponse struct {
	XMLName xml.Name `xml:"subsonic-response"`
	Status  string   `xml:"status,attr"`
	Error   *Error   `xml:"error,omitempty"`
	// Add other response fields as needed
}

type Error struct {
	Code    int    `xml:"code,attr"`
	Message string `xml:"message,attr"`
}

type Song struct {
	XMLName     xml.Name `xml:"song"`
	ID          string   `xml:"id,attr"`
	Title       string   `xml:"title,attr"`
	Artist      string   `xml:"artist,attr"`
	Album       string   `xml:"album,attr"`
	AlbumID     string   `xml:"albumId,attr"`
	ArtistID    string   `xml:"artistId,attr"`
	Duration    int      `xml:"duration,attr"`
	Track       int      `xml:"track,attr"`
	Year        int      `xml:"year,attr"`
	CoverArtID  string   `xml:"coverArt,attr"`
	Size        int64    `xml:"size,attr"`
	BitRate     int      `xml:"bitRate,attr"`
	ContentType string   `xml:"contentType,attr"`
	Suffix      string   `xml:"suffix,attr"`
}

type Album struct {
	XMLName    xml.Name  `xml:"album"`
	ID         string    `xml:"id,attr"`
	Name       string    `xml:"name,attr"`
	Artist     string    `xml:"artist,attr"`
	ArtistID   string    `xml:"artistId,attr"`
	SongCount  int       `xml:"songCount,attr"`
	Duration   int       `xml:"duration,attr"`
	CoverArtID string    `xml:"coverArt,attr"`
	Year       int       `xml:"year,attr"`
	Starred    time.Time `xml:"starred,attr"`
}

type Artist struct {
	XMLName    xml.Name `xml:"artist"`
	ID         string   `xml:"id,attr"`
	Name       string   `xml:"name,attr"`
	AlbumCount int      `xml:"albumCount,attr"`
	CoverArtID string   `xml:"coverArt,attr"`
}

type MusicDirectory struct {
	XMLName  xml.Name `xml:"directory"`
	ID       string   `xml:"id,attr"`
	Name     string   `xml:"name,attr"`
	Children []Child  `xml:"child,omitempty"`
}

type Child struct {
	XMLName     xml.Name `xml:"child"`
	ID          string   `xml:"id,attr"`
	Title       string   `xml:"title,attr"`
	Artist      string   `xml:"artist,attr"`
	Album       string   `xml:"album,attr"`
	Duration    int      `xml:"duration,attr"`
	CoverArtID  string   `xml:"coverArt,attr"`
	Size        int64    `xml:"size,attr"`
	BitRate     int      `xml:"bitRate,attr"`
	ContentType string   `xml:"contentType,attr"`
	Suffix      string   `xml:"suffix,attr"`
	IsDir       bool     `xml:"isDir,attr"`
}
