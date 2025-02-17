package uberspace

import (
	"bytes"
	"context"
	"os/exec"
)

func (c *Client) WebDomainAdd(ctx context.Context, domain string) error {
	cmd := exec.CommandContext(ctx, "uberspace", "web", "domain", "add", domain)

	_, err := c.Runner.Run(cmd)

	return err
}

func (c *Client) WebDomainRead(ctx context.Context, domain string) (bool, error) {
	cmd := exec.CommandContext(ctx, "uberspace", "web", "domain", "list")

	out, err := c.Runner.Run(cmd)
	if err != nil {
		return false, err
	}

	return bytes.Contains(out, []byte(domain)), nil
}

func (c *Client) WebDomainDelete(ctx context.Context, domain string) (bool, error) {
	cmd := exec.CommandContext(ctx, "uberspace", "web", "domain", "del", domain)

	_, err := c.Runner.Run(cmd)
	if err != nil {
		return false, err
	}

	return true, nil
}
