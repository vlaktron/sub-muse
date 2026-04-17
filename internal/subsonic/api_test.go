package subsonic

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetArtists_Success(t *testing.T) {
	mockJSON := `{"subsonic-response":{"status":"ok","artists":{"artist":[{"id":"1","name":"Artist 1","albumCount":5,"coverArt":"1"},{"id":"2","name":"Artist 2","albumCount":3,"coverArt":"2"}]}}}`

	client := &Client{
		httpClient: &mockHTTPClient{
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(mockJSON))),
				}, nil
			},
		},
	}

	artists, err := client.GetArtists()
	require.NoError(t, err)
	require.Len(t, artists, 2)
	require.Equal(t, "Artist 1", artists[0].Name)
	require.Equal(t, "Artist 2", artists[1].Name)
}

func TestGetArtists_Empty(t *testing.T) {
	mockJSON := `{"subsonic-response":{"status":"ok","artists":{}}}`

	client := &Client{
		httpClient: &mockHTTPClient{
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(mockJSON))),
				}, nil
			},
		},
	}

	artists, err := client.GetArtists()
	require.NoError(t, err)
	require.Empty(t, artists)
}

func TestGetAlbums_Success(t *testing.T) {
	mockJSON := `{"subsonic-response":{"status":"ok","albumList":{"album":[{"id":"1","name":"Album 1","artist":"Artist 1","songCount":10,"duration":300,"coverArt":"1"},{"id":"2","name":"Album 2","artist":"Artist 2","songCount":12,"duration":350,"coverArt":"2"}]}}}`

	client := &Client{
		httpClient: &mockHTTPClient{
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(mockJSON))),
				}, nil
			},
		},
	}

	albums, err := client.GetAlbums()
	require.NoError(t, err)
	require.Len(t, albums, 2)
	require.Equal(t, "Album 1", albums[0].Name)
	require.Equal(t, "Album 2", albums[1].Name)
}

func TestGetSongs_Success(t *testing.T) {
	mockJSON := `{"subsonic-response":{"status":"ok","randomSongs":{"song":[{"id":"1","title":"Song 1","artist":"Artist 1","album":"Album 1","duration":180,"track":1,"year":2020,"coverArt":"1","size":5000000,"bitRate":128,"contentType":"audio/mpeg","suffix":"mp3"},{"id":"2","title":"Song 2","artist":"Artist 2","album":"Album 2","duration":240,"track":2,"year":2021,"coverArt":"2","size":6000000,"bitRate":192,"contentType":"audio/mpeg","suffix":"mp3"}]}}}`

	client := &Client{
		httpClient: &mockHTTPClient{
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(mockJSON))),
				}, nil
			},
		},
	}

	songs, err := client.GetSongs()
	require.NoError(t, err)
	require.Len(t, songs, 2)
	require.Equal(t, "Song 1", songs[0].Title)
	require.Equal(t, "Song 2", songs[1].Title)
}

func TestPing_Success(t *testing.T) {
	mockJSON := `{"subsonic-response":{"status":"ok"}}`

	client := &Client{
		httpClient: &mockHTTPClient{
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(mockJSON))),
				}, nil
			},
		},
	}

	err := client.Ping()
	require.NoError(t, err)
}

func TestPing_Failure(t *testing.T) {
	client := &Client{
		httpClient: &mockHTTPClient{
			doFunc: func(req *http.Request) (*http.Response, error) {
				return nil, fmt.Errorf("network error")
			},
		},
	}

	err := client.Ping()
	require.Error(t, err)
	require.Equal(t, "network error", err.Error())
}

func TestGetArtist_Success(t *testing.T) {
	mockJSON := `{"subsonic-response":{"status":"ok","artist":{"id":"1","name":"Artist 1","albumCount":3,"coverArt":"1","album":[{"id":"1","name":"Album 1","artist":"Artist 1","songCount":10,"duration":300,"coverArt":"1"}]}}}`

	client := &Client{
		httpClient: &mockHTTPClient{
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(mockJSON))),
				}, nil
			},
		},
	}

	artist, err := client.GetArtist("1")
	require.NoError(t, err)
	require.Equal(t, "Artist 1", artist.Name)
	require.Equal(t, 3, artist.AlbumCount)
	require.Len(t, artist.Albums, 1)
	require.Equal(t, "Album 1", artist.Albums[0].Name)
}

func TestGetArtist_Empty(t *testing.T) {
	mockJSON := `{"subsonic-response":{"status":"ok","artist":{"id":"1","name":"Artist 1","albumCount":0,"coverArt":"1"}}}`

	client := &Client{
		httpClient: &mockHTTPClient{
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(mockJSON))),
				}, nil
			},
		},
	}

	artist, err := client.GetArtist("1")
	require.NoError(t, err)
	require.Equal(t, "Artist 1", artist.Name)
	require.Equal(t, 0, artist.AlbumCount)
	require.Empty(t, artist.Albums)
}

func TestGetArtist_ServerError(t *testing.T) {
	client := &Client{
		httpClient: &mockHTTPClient{
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Status:     "500 Internal Server Error",
					Body:       io.NopCloser(bytes.NewReader([]byte("error"))),
				}, nil
			},
		},
	}

	_, err := client.GetArtist("1")
	require.Error(t, err)
	require.Contains(t, err.Error(), "HTTP 500")
}

func TestGetArtist_VerifyParams(t *testing.T) {
	var capturedURL string
	mockJSON := `{"subsonic-response":{"status":"ok","artist":{"id":"42","name":"Test Artist","albumCount":0}}}`

	client := &Client{
		baseURL:    "http://example.com",
		username:   "user",
		password:   "pass",
		clientName: "test",
		httpClient: &mockHTTPClient{
			doFunc: func(req *http.Request) (*http.Response, error) {
				capturedURL = req.URL.String()
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(mockJSON))),
				}, nil
			},
		},
	}

	_, err := client.GetArtist("42")
	require.NoError(t, err)
	require.Contains(t, capturedURL, "id=42")
	require.Contains(t, capturedURL, "f=json")
}
