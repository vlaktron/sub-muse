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
	mockXML := `<?xml version="1.0"?>
<subsonic-response status="ok">
	<artists>
		<artist id="1" name="Artist 1" albumCount="5" coverArt="1"/>
		<artist id="2" name="Artist 2" albumCount="3" coverArt="2"/>
	</artists>
</subsonic-response>`

	client := &Client{
		httpClient: &mockHTTPClient{
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(mockXML))),
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
	mockXML := `<?xml version="1.0"?>
<subsonic-response status="ok">
	<artists>
	</artists>
</subsonic-response>`

	client := &Client{
		httpClient: &mockHTTPClient{
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(mockXML))),
				}, nil
			},
		},
	}

	artists, err := client.GetArtists()
	require.NoError(t, err)
	require.Empty(t, artists)
}

func TestGetAlbums_Success(t *testing.T) {
	mockXML := `<?xml version="1.0"?>
<subsonic-response status="ok">
	<albums>
		<album id="1" name="Album 1" artist="Artist 1" songCount="10" duration="300" coverArt="1"/>
		<album id="2" name="Album 2" artist="Artist 2" songCount="12" duration="350" coverArt="2"/>
	</albums>
</subsonic-response>`

	client := &Client{
		httpClient: &mockHTTPClient{
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(mockXML))),
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
	mockXML := `<?xml version="1.0"?>
<subsonic-response status="ok">
	<songs>
		<song id="1" title="Song 1" artist="Artist 1" album="Album 1" duration="180" track="1" year="2020" coverArt="1" size="5000000" bitRate="128" contentType="audio/mpeg" suffix="mp3"/>
		<song id="2" title="Song 2" artist="Artist 2" album="Album 2" duration="240" track="2" year="2021" coverArt="2" size="6000000" bitRate="192" contentType="audio/mpeg" suffix="mp3"/>
	</songs>
</subsonic-response>`

	client := &Client{
		httpClient: &mockHTTPClient{
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(mockXML))),
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
	mockXML := `<?xml version="1.0"?>
<subsonic-response status="ok">
</subsonic-response>`

	client := &Client{
		httpClient: &mockHTTPClient{
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(mockXML))),
				}, nil
			},
		},
	}

	var response struct{}
	err := client.sendRequest("ping", nil, &response)
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
