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
	mockJSON := `{"subsonic-response":{"status":"ok","artists":{"index":[{"name":"A","artist":[{"id":"1","name":"Artist 1","albumCount":5,"coverArt":"1"}]},{"name":"B","artist":[{"id":"2","name":"Artist 2","albumCount":3,"coverArt":"2"}]}]}}}`

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
	mockJSON := `{"subsonic-response":{"status":"ok","artists":{"index":[]}}}`

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

func TestStream_Success(t *testing.T) {
	mockAudio := []byte{0x00, 0x01, 0x02, 0x03, 0x04}

	client := &Client{
		httpClient: &mockHTTPClient{
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(mockAudio)),
				}, nil
			},
		},
	}

	data, err := client.Stream(WithID("123"))
	require.NoError(t, err)
	require.Equal(t, mockAudio, data)
}

func TestStream_WithFormat(t *testing.T) {
	var capturedURL string
	mockAudio := []byte{0x00, 0x01, 0x02}

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
					Body:       io.NopCloser(bytes.NewReader(mockAudio)),
				}, nil
			},
		},
	}

	_, err := client.Stream(WithID("456"), WithFormat("mp3"))
	require.NoError(t, err)
	require.Contains(t, capturedURL, "id=456")
	require.Contains(t, capturedURL, "format=mp3")
}

func TestStream_WithMaxBitRate(t *testing.T) {
	var capturedURL string
	mockAudio := []byte{0x00, 0x01, 0x02}

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
					Body:       io.NopCloser(bytes.NewReader(mockAudio)),
				}, nil
			},
		},
	}

	_, err := client.Stream(WithID("789"), WithMaxBitRate(128))
	require.NoError(t, err)
	require.Contains(t, capturedURL, "maxBitRate=128")
}

func TestStream_MultipleOptions(t *testing.T) {
	var capturedURL string
	mockAudio := []byte{0x00, 0x01, 0x02}

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
					Body:       io.NopCloser(bytes.NewReader(mockAudio)),
				}, nil
			},
		},
	}

	_, err := client.Stream(WithID("111"), WithFormat("flac"), WithMaxBitRate(256))
	require.NoError(t, err)
	require.Contains(t, capturedURL, "id=111")
	require.Contains(t, capturedURL, "format=flac")
	require.Contains(t, capturedURL, "maxBitRate=256")
}

func TestStream_HTTPError(t *testing.T) {
	client := &Client{
		httpClient: &mockHTTPClient{
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Status:     "404 Not Found",
					Body:       io.NopCloser(bytes.NewReader([]byte("not found"))),
				}, nil
			},
		},
	}

	_, err := client.Stream(WithID("999"))
	require.Error(t, err)
	require.Contains(t, err.Error(), "HTTP 404")
}

func TestStream_NetworkError(t *testing.T) {
	client := &Client{
		httpClient: &mockHTTPClient{
			doFunc: func(req *http.Request) (*http.Response, error) {
				return nil, fmt.Errorf("network error")
			},
		},
	}

	_, err := client.Stream(WithID("999"))
	require.Error(t, err)
	require.Equal(t, "network error", err.Error())
}

func TestStream_VerifyParams(t *testing.T) {
	var capturedURL string
	mockAudio := []byte{0x00, 0x01, 0x02}

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
					Body:       io.NopCloser(bytes.NewReader(mockAudio)),
				}, nil
			},
		},
	}

	_, err := client.Stream(WithID("42"))
	require.NoError(t, err)
	require.Contains(t, capturedURL, "id=42")
	require.Contains(t, capturedURL, "f=json")
}

func TestStream_WithTimeOffset(t *testing.T) {
	var capturedURL string
	mockAudio := []byte{0x00, 0x01, 0x02}

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
					Body:       io.NopCloser(bytes.NewReader(mockAudio)),
				}, nil
			},
		},
	}

	_, err := client.Stream(WithID("123"), WithTimeOffset(30))
	require.NoError(t, err)
	require.Contains(t, capturedURL, "id=123")
	require.Contains(t, capturedURL, "timeOffset=30")
}

func TestStream_EmptyID(t *testing.T) {
	client := &Client{
		httpClient: &mockHTTPClient{
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusBadRequest,
					Status:     "400 Bad Request",
					Body:       io.NopCloser(bytes.NewReader([]byte("empty id"))),
				}, nil
			},
		},
	}

	_, err := client.Stream(WithID(""))
	require.Error(t, err)
	require.Contains(t, err.Error(), "HTTP 400")
}

func TestStream_ServerError(t *testing.T) {
	client := &Client{
		httpClient: &mockHTTPClient{
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Status:     "500 Internal Server Error",
					Body:       io.NopCloser(bytes.NewReader([]byte("server error"))),
				}, nil
			},
		},
	}

	_, err := client.Stream(WithID("999"))
	require.Error(t, err)
	require.Contains(t, err.Error(), "HTTP 500")
}

func TestGetMusicFolders_Success(t *testing.T) {
	mockJSON := `{"subsonic-response":{"status":"ok","musicFolders":{"musicFolder":[{"id":1,"name":"music"},{"id":4,"name":"upload"}]}}}`

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

	folders, err := client.GetMusicFolders()
	require.NoError(t, err)
	require.Len(t, folders, 2)
	require.Equal(t, 1, folders[0].ID)
	require.Equal(t, "music", folders[0].Name)
	require.Equal(t, 4, folders[1].ID)
	require.Equal(t, "upload", folders[1].Name)
}

func TestGetMusicFolders_Empty(t *testing.T) {
	mockJSON := `{"subsonic-response":{"status":"ok","musicFolders":{"musicFolder":[]}}}`

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

	folders, err := client.GetMusicFolders()
	require.NoError(t, err)
	require.Empty(t, folders)
}

func TestGetMusicFolders_VerifyParams(t *testing.T) {
	var capturedURL string
	mockJSON := `{"subsonic-response":{"status":"ok","musicFolders":{"musicFolder":[]}}}`

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

	_, err := client.GetMusicFolders()
	require.NoError(t, err)
	require.Contains(t, capturedURL, "getMusicFolders")
	require.Contains(t, capturedURL, "f=json")
}
