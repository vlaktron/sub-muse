package subsonic

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	testBaseURL    = os.Getenv("TEST_SUBSONIC_URL")
	testUsername   = os.Getenv("TEST_SUBSONIC_USERNAME")
	testPassword   = os.Getenv("TEST_SUBSONIC_PASSWORD")
	testClientName = os.Getenv("TEST_SUBSONIC_CLIENT_NAME")
)

func init() {
	if testClientName == "" {
		testClientName = "sub-muse-test"
	}
}

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
	require.Contains(t, url, "http://example.com/rest/test")
}

func TestClient_buildRequest_QueryParams(t *testing.T) {
	client := NewClient("http://example.com/", "testuser", "testpass!", "testclient")

	url, err := client.buildRequest("test", nil)
	require.NoError(t, err)

	require.Contains(t, url, "u=testuser")
	require.Contains(t, url, "p=testpass%21")
	require.Contains(t, url, "v=1.16.1")
	require.Contains(t, url, "c=testclient")
	require.Contains(t, url, "f=json")
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
	mockJSON := `{
		"subsonic-response": {
			"status": "ok",
			"artists": {
				"index": [
					{
						"name": "A",
						"artist": [
							{
								"id": "1",
								"name": "Artist 1",
								"albumCount": 5,
								"coverArt": "1"
							}
						]
					}
				]
			}
		}
	}`

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

	var response struct {
		Status  string `json:"status"`
		Artists struct {
			Index []struct {
				Artist []Artist `json:"artist"`
			} `json:"index"`
		} `json:"artists"`
	}

	err := client.sendRequest("getArtists", nil, &response)
	require.NoError(t, err)
	require.Len(t, response.Artists.Index, 1)
	require.Len(t, response.Artists.Index[0].Artist, 1)
	require.Equal(t, "Artist 1", response.Artists.Index[0].Artist[0].Name)
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

func TestClient_buildRequest_TokenAuth_HasTokenAndSalt(t *testing.T) {
	client := NewClientWithTokenAuth("http://example.com/", "testuser", "testpass", "testclient")

	url, err := client.buildRequest("test", nil)
	require.NoError(t, err)
	require.Contains(t, url, "t=")
	require.Contains(t, url, "s=")
	require.NotContains(t, url, "p=")
}

func TestClient_buildRequest_TokenAuth_NoPassword(t *testing.T) {
	client := NewClientWithTokenAuth("http://example.com/", "testuser", "testpass", "testclient")

	url, err := client.buildRequest("test", nil)
	require.NoError(t, err)
	require.NotContains(t, url, "p=testpass")
}

func TestClient_buildRequest_TokenAuth_ValidMD5(t *testing.T) {
	client := NewClientWithTokenAuth("http://example.com/", "testuser", "testpass", "testclient")

	urlStr, err := client.buildRequest("test", nil)
	require.NoError(t, err)

	require.Contains(t, urlStr, "t=")
	require.Contains(t, urlStr, "s=")

	query, err := url.ParseQuery(urlStr)
	require.NoError(t, err)

	token := query.Get("t")
	salt := query.Get("s")

	expectedToken := md5.Sum([]byte("testpass" + salt))
	expectedTokenHex := hex.EncodeToString(expectedToken[:])

	require.Equal(t, expectedTokenHex, token)
}

func TestClient_buildRequest_PasswordMode_Default(t *testing.T) {
	client := NewClient("http://example.com/", "testuser", "testpass!", "testclient")

	url, err := client.buildRequest("test", nil)
	require.NoError(t, err)
	require.Contains(t, url, "p=testpass%21")
	require.NotContains(t, url, "t=")
	require.NotContains(t, url, "s=")
}

func TestGenerateSalt_Uniqueness(t *testing.T) {
	salt1, err := generateSalt()
	require.NoError(t, err)

	salt2, err := generateSalt()
	require.NoError(t, err)

	require.NotEqual(t, salt1, salt2)
}

func TestGenerateSalt_Length(t *testing.T) {
	salt, err := generateSalt()
	require.NoError(t, err)

	require.Equal(t, 18, len(salt))
}
