package uberspace

import (
	"bytes"
)

func (c *Client) WebDomainAdd(domain string) error {
	_, err := c.SSHClient.Run("uberspace web domain add " + domain)

	return err
}

func (c *Client) WebDomainExists(domain string) (bool, error) {
	out, err := c.SSHClient.Run("uberspace web domain list")
	if err != nil {
		return false, err
	}

	return bytes.Contains(out, []byte(domain)), nil
}

func (c *Client) WebDomainDelete(domain string) error {
	_, err := c.SSHClient.Run("uberspace web domain del " + domain)

	return err
}
