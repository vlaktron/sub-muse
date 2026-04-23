package subsonic

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type AuthMode int

const (
	AuthPassword AuthMode = iota
	AuthToken
)

type Client struct {
	baseURL    string
	username   string
	password   string
	clientName string
	authMode   AuthMode
	httpClient HTTPClient
}

func NewClient(baseURL, username, password, clientName string) *Client {
	return &Client{
		baseURL:    baseURL,
		username:   username,
		password:   password,
		clientName: clientName,
		authMode:   AuthPassword,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func NewClientWithTokenAuth(baseURL, username, password, clientName string) *Client {
	c := NewClient(baseURL, username, password, clientName)
	c.authMode = AuthToken
	return c
}

func (c *Client) buildRequest(endpoint string, params map[string]string) (string, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return "", err
	}

	u.Path = "/rest/" + endpoint

	query := u.Query()
	query.Set("u", c.username)

	switch c.authMode {
	case AuthToken:
		salt, err := generateSalt()
		if err != nil {
			return "", err
		}
		hash := md5.Sum([]byte(c.password + salt))
		query.Set("t", hex.EncodeToString(hash[:]))
		query.Set("s", salt)
	default:
		query.Set("p", c.password)
	}

	query.Set("v", "1.16.1")
	query.Set("c", c.clientName)
	query.Set("f", "json")

	for key, value := range params {
		query.Set(key, value)
	}

	u.RawQuery = query.Encode()
	return u.String(), nil
}

func generateSalt() (string, error) {
	b := make([]byte, 9)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
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

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Unwrap the outer "subsonic-response" envelope
	var outer struct {
		Inner json.RawMessage `json:"subsonic-response"`
	}
	if err := json.Unmarshal(bodyBytes, &outer); err != nil {
		return err
	}

	// Check status
	var statusCheck struct {
		Status string `json:"status"`
		Error  *Error `json:"error"`
	}
	if err := json.Unmarshal(outer.Inner, &statusCheck); err != nil {
		return err
	}

	if statusCheck.Status != "ok" {
		if statusCheck.Error != nil {
			return fmt.Errorf("API error %d: %s", statusCheck.Error.Code, statusCheck.Error.Message)
		}
		return fmt.Errorf("API error: status=%s", statusCheck.Status)
	}

	if result == nil {
		return nil
	}

	return json.Unmarshal(outer.Inner, result)
}

func (c *Client) buildRequestWithValues(endpoint string, params url.Values) (string, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return "", err
	}

	u.Path = "/rest/" + endpoint

	query := u.Query()
	query.Set("u", c.username)

	switch c.authMode {
	case AuthToken:
		salt, err := generateSalt()
		if err != nil {
			return "", err
		}
		hash := md5.Sum([]byte(c.password + salt))
		query.Set("t", hex.EncodeToString(hash[:]))
		query.Set("s", salt)
	default:
		query.Set("p", c.password)
	}

	query.Set("v", "1.16.1")
	query.Set("c", c.clientName)
	query.Set("f", "json")

	for key, values := range params {
		for _, v := range values {
			query.Add(key, v)
		}
	}

	u.RawQuery = query.Encode()
	return u.String(), nil
}

func (c *Client) sendRequestWithValues(endpoint string, params url.Values, result interface{}) error {
	requestURL, err := c.buildRequestWithValues(endpoint, params)
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

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Unwrap the outer "subsonic-response" envelope
	var outer struct {
		Inner json.RawMessage `json:"subsonic-response"`
	}
	if err := json.Unmarshal(bodyBytes, &outer); err != nil {
		return err
	}

	// Check status
	var statusCheck struct {
		Status string `json:"status"`
		Error  *Error `json:"error"`
	}
	if err := json.Unmarshal(outer.Inner, &statusCheck); err != nil {
		return err
	}

	if statusCheck.Status != "ok" {
		if statusCheck.Error != nil {
			return fmt.Errorf("API error %d: %s", statusCheck.Error.Code, statusCheck.Error.Message)
		}
		return fmt.Errorf("API error: status=%s", statusCheck.Status)
	}

	if result == nil {
		return nil
	}

	return json.Unmarshal(outer.Inner, result)
}
