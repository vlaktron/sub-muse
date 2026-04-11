package subsonic

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

// GetMusicFolders gets the list of music folders
func (c *Client) GetMusicFolders() (*MusicDirectory, error) {
	var response SubsonicResponse
	err := c.sendRequest("getMusicFolders", nil, &response)
	if err != nil {
		return nil, err
	}

	// Parse the response to get music folders
	// This is a simplified example - you'd need to handle the actual XML structure
	return &MusicDirectory{}, nil
}

// GetArtists gets the list of artists
func (c *Client) GetArtists() ([]Artist, error) {
	params := map[string]string{
		"musicFolderId": "",
	}

	var response struct {
		XMLName xml.Name `xml:"subsonic-response"`
		Artists struct {
			Artist []Artist `xml:"artist"`
		} `xml:"artists"`
	}

	err := c.sendRequest("getArtists", params, &response)
	if err != nil {
		return nil, err
	}

	return response.Artists.Artist, nil
}

// GetAlbums gets the list of albums
func (c *Client) GetAlbums() ([]Album, error) {
	var response struct {
		XMLName xml.Name `xml:"subsonic-response"`
		Albums  struct {
			Album []Album `xml:"album"`
		} `xml:"albums"`
	}

	err := c.sendRequest("getAlbumList", map[string]string{
		"type": "alphabeticalByName",
		"size": "1000",
	}, &response)
	if err != nil {
		return nil, err
	}

	return response.Albums.Album, nil
}

// GetSongs gets the list of songs
func (c *Client) GetSongs() ([]Song, error) {
	var response struct {
		XMLName xml.Name `xml:"subsonic-response"`
		Songs   struct {
			Song []Song `xml:"song"`
		} `xml:"songs"`
	}

	err := c.sendRequest("getSongs", map[string]string{
		"size": "1000",
	}, &response)
	if err != nil {
		return nil, err
	}

	return response.Songs.Song, nil
}

// GetAlbum gets a specific album by ID
func (c *Client) GetAlbum(id string) (*Album, error) {
	var response struct {
		XMLName xml.Name `xml:"subsonic-response"`
		Album   struct {
			Album Album `xml:"album"`
		} `xml:"album"`
	}

	err := c.sendRequest("getAlbum", map[string]string{
		"id": id,
	}, &response)
	if err != nil {
		return nil, err
	}

	return &response.Album.Album, nil
}

// GetSongsByArtist gets songs by a specific artist
func (c *Client) GetSongsByArtist(artistID string) ([]Song, error) {
	var response struct {
		XMLName xml.Name `xml:"subsonic-response"`
		Songs   struct {
			Song []Song `xml:"song"`
		} `xml:"songs"`
	}

	err := c.sendRequest("getSongsByArtist", map[string]string{
		"artistId": artistID,
	}, &response)
	if err != nil {
		return nil, err
	}

	return response.Songs.Song, nil
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
