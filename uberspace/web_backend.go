package uberspace

import (
	"bytes"
	"fmt"
)

func (c *Client) WebBackendSet(uri string, port int32) error {
	_, err := c.SSHClient.Run(fmt.Sprintf("uberspace web backend set %s --http --port %d", uri, port))

	return err
}

func (c *Client) WebBackendExists(uri string) (bool, error) {
	out, err := c.SSHClient.Run("uberspace web backend list")
	if err != nil {
		return false, err
	}

	return bytes.Contains(out, []byte(uri)), nil
}

func (c *Client) WebBackendDelete(uri string) error {
	_, err := c.SSHClient.Run("uberspace web backend del " + uri)

	return err
}
