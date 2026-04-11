package subsonic

func (c *Client) Ping() error {
	err := c.sendRequest("ping", nil, nil)
	return err
}
