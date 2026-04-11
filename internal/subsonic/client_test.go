package subsonic

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func getTestEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

var (
	testBaseURL    = getTestEnv("TEST_SUBSONIC_URL", "http://192.168.1.5:4533/")
	testUsername   = getTestEnv("TEST_SUBSONIC_USERNAME", "sub-muse")
	testPassword   = getTestEnv("TEST_SUBSONIC_PASSWORD", "Test12345!")
	testClientName = getTestEnv("TEST_SUBSONIC_CLIENT_NAME", "sub-muse-test")
)

func TestNewClient_Success(t *testing.T) {
	client := NewClient(testBaseURL, testUsername, testPassword, testClientName)

	require.Equal(t, testBaseURL, client.baseURL)
	require.Equal(t, testUsername, client.username)
	require.Equal(t, testPassword, client.password)
	require.Equal(t, testClientName, client.clientName)
	require.NotNil(t, client.httpClient)
}

func TestClient_buildRequest_BasicURL(t *testing.T) {
	client := NewClient("http://example.com/", testUsername, testPassword, testClientName)

	url, err := client.buildRequest("test", nil)
	require.NoError(t, err)
	require.Contains(t, url, "http://example.com//test")
}

func TestClient_buildRequest_QueryParams(t *testing.T) {
	client := NewClient("http://example.com/", testUsername, testPassword, testClientName)

	url, err := client.buildRequest("test", nil)
	require.NoError(t, err)

	require.Contains(t, url, "u="+testUsername)
	require.Contains(t, url, "p=Test12345%21")
	require.Contains(t, url, "v=1.15.0")
	require.Contains(t, url, "c="+testClientName)
}

func TestClient_buildRequest_CustomParams(t *testing.T) {
	client := NewClient("http://example.com/", testUsername, testPassword, testClientName)

	params := map[string]string{
		"musicFolderId": "123",
		"size":          "50",
	}

	url, err := client.buildRequest("test", params)
	require.NoError(t, err)

	require.Contains(t, url, "musicFolderId=123")
	require.Contains(t, url, "size=50")
}

func TestClient_buildRequest_Error(t *testing.T) {
	client := &Client{
		baseURL:    ":invalid",
		username:   testUsername,
		password:   testPassword,
		clientName: testClientName,
	}

	_, err := client.buildRequest("test", nil)
	require.Error(t, err)
}

func TestClient_sendRequest_Success(t *testing.T) {
	mockXML := `<?xml version="1.0"?>
<subsonic-response status="ok">
	<artists>
		<artist id="1" name="Artist 1" albumCount="5" coverArt="1"/>
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

	var response struct {
		XMLName xml.Name `xml:"subsonic-response"`
		Artists struct {
			Artist []Artist `xml:"artist"`
		} `xml:"artists"`
	}

	err := client.sendRequest("getArtists", nil, &response)
	require.NoError(t, err)
	require.Len(t, response.Artists.Artist, 1)
	require.Equal(t, "Artist 1", response.Artists.Artist[0].Name)
}

func TestClient_sendRequest_HTTPError(t *testing.T) {
	client := &Client{
		httpClient: &mockHTTPClient{
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusUnauthorized,
					Status:     "401 Unauthorized",
					Body:       io.NopCloser(bytes.NewReader([]byte("Unauthorized"))),
				}, nil
			},
		},
	}

	err := client.sendRequest("getArtists", nil, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "HTTP 401")
}

func TestClient_sendRequest_RequestError(t *testing.T) {
	client := &Client{
		httpClient: &mockHTTPClient{
			doFunc: func(req *http.Request) (*http.Response, error) {
				return nil, fmt.Errorf("network error")
			},
		},
	}

	err := client.sendRequest("getArtists", nil, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "network error")
}
