package subsonic

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	baseURL    string
	username   string
	password   string
	clientName string
	httpClient HTTPClient
}

func NewClient(baseURL, username, password, clientName string) *Client {
	return &Client{
		baseURL:    baseURL,
		username:   username,
		password:   password,
		clientName: clientName,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) buildRequest(endpoint string, params map[string]string) (string, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return "", err
	}

	u.Path = fmt.Sprintf("%s/%s", u.Path, endpoint)

	query := u.Query()
	query.Set("u", c.username)
	query.Set("p", c.password)
	query.Set("v", "1.15.0")
	query.Set("c", c.clientName)

	for key, value := range params {
		query.Set(key, value)
	}

	u.RawQuery = query.Encode()
	return u.String(), nil
}

func (c *Client) sendRequest(endpoint string, params map[string]string, result interface{}) error {
	requestURL, err := c.buildRequest(endpoint, params)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	if result == nil {
		return nil
	}

	return xml.NewDecoder(resp.Body).Decode(result)
}
