package subsonic

import (
	"fmt"
	"io"
	"net/http"
)

// Ping checks connectivity to the server
func (c *Client) Ping() error {
	return c.sendRequest("ping", nil, nil)
}

// GetMusicFolders gets the list of configured music folders
func (c *Client) GetMusicFolders() ([]MusicFolder, error) {
	var envelope struct {
		MusicFolders struct {
			Folder []MusicFolder `json:"musicFolder"`
		} `json:"musicFolders"`
	}

	if err := c.sendRequest("getMusicFolders", nil, &envelope); err != nil {
		return nil, err
	}

	return envelope.MusicFolders.Folder, nil
}

// GetArtists gets the list of artists
func (c *Client) GetArtists() ([]Artist, error) {
	var envelope struct {
		Artists struct {
			Index []struct {
				Artist []Artist `json:"artist"`
			} `json:"index"`
		} `json:"artists"`
	}

	if err := c.sendRequest("getArtists", nil, &envelope); err != nil {
		return nil, err
	}

	var artists []Artist
	for _, idx := range envelope.Artists.Index {
		artists = append(artists, idx.Artist...)
	}
	return artists, nil
}

// GetArtist gets an artist by ID including their albums
func (c *Client) GetArtist(artistID string) (*Artist, error) {
	var envelope struct {
		Artist Artist `json:"artist"`
	}

	if err := c.sendRequest("getArtist", map[string]string{"id": artistID}, &envelope); err != nil {
		return nil, err
	}

	return &envelope.Artist, nil
}

// GetAlbums gets the list of albums alphabetically
func (c *Client) GetAlbums() ([]Album, error) {
	var envelope struct {
		AlbumList struct {
			Albums []Album `json:"album"`
		} `json:"albumList"`
	}

	if err := c.sendRequest("getAlbumList", map[string]string{
		"type": "alphabeticalByName",
		"size": "1000",
	}, &envelope); err != nil {
		return nil, err
	}

	return envelope.AlbumList.Albums, nil
}

// GetAlbum gets a specific album by ID including its songs
func (c *Client) GetAlbum(id string) (*Album, error) {
	var envelope struct {
		Album Album `json:"album"`
	}

	if err := c.sendRequest("getAlbum", map[string]string{"id": id}, &envelope); err != nil {
		return nil, err
	}

	return &envelope.Album, nil
}

// GetSongs gets a random list of songs
func (c *Client) GetSongs() ([]Song, error) {
	var envelope struct {
		RandomSongs struct {
			Song []Song `json:"song"`
		} `json:"randomSongs"`
	}

	if err := c.sendRequest("getRandomSongs", map[string]string{
		"size": "1000",
	}, &envelope); err != nil {
		return nil, err
	}

	return envelope.RandomSongs.Song, nil
}

// GetCoverArt gets cover art for an album
func (c *Client) GetCoverArt(coverArtID string, size int) ([]byte, error) {
	params := map[string]string{
		"id": coverArtID,
	}
	if size > 0 {
		params["size"] = fmt.Sprintf("%d", size)
	}

	requestURL, err := c.buildRequest("getCoverArt", params)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	return io.ReadAll(resp.Body)
}

type StreamOptions struct {
	ID         string
	Format     string
	MaxBitRate int
	TimeOffset int
}

type StreamOption func(*StreamOptions)

func WithID(id string) StreamOption {
	return func(o *StreamOptions) {
		o.ID = id
	}
}

func WithFormat(format string) StreamOption {
	return func(o *StreamOptions) {
		o.Format = format
	}
}

func WithMaxBitRate(bitRate int) StreamOption {
	return func(o *StreamOptions) {
		o.MaxBitRate = bitRate
	}
}

func WithTimeOffset(offset int) StreamOption {
	return func(o *StreamOptions) {
		o.TimeOffset = offset
	}
}

func (c *Client) Stream(opts ...StreamOption) ([]byte, error) {
	options := &StreamOptions{}
	for _, opt := range opts {
		opt(options)
	}

	params := map[string]string{
		"id": options.ID,
	}
	if options.Format != "" {
		params["format"] = options.Format
	}
	if options.MaxBitRate > 0 {
		params["maxBitRate"] = fmt.Sprintf("%d", options.MaxBitRate)
	}
	if options.TimeOffset > 0 {
		params["timeOffset"] = fmt.Sprintf("%d", options.TimeOffset)
	}

	requestURL, err := c.buildRequest("stream", params)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	return io.ReadAll(resp.Body)
}
