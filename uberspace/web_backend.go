package uberspace

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
)

func (c *Client) WebBackendSet(ctx context.Context, uri string, port int32) error {
	cmd := exec.CommandContext(ctx, "uberspace", "web", "backend", "set", uri, "--http", "--port", fmt.Sprint(port)) //nolint: gosec

	_, err := c.Runner.Run(cmd)

	return err
}

func (c *Client) WebBackendRead(ctx context.Context, uri string) (bool, error) {
	cmd := exec.CommandContext(ctx, "uberspace", "web", "backend", "list")

	out, err := c.Runner.Run(cmd)
	if err != nil {
		return false, err
	}

	return bytes.Contains(out, []byte(uri)), nil
}

func (c *Client) WebBackendDelete(ctx context.Context, uri string) (bool, error) {
	cmd := exec.CommandContext(ctx, "uberspace", "web", "backend", "del", uri)

	_, err := c.Runner.Run(cmd)
	if err != nil {
		return false, err
	}

	return true, nil
}
