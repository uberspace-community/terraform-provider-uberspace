package uberspace

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
)

func (c *Client) CrontabEntryExists(ctx context.Context, entry string) (bool, error) {
	out, err := c.crontabL(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to read crontab: %w", err)
	}

	for line := range bytes.Lines(out) {
		if string(line) == entry {
			return true, nil
		}
	}

	return false, nil
}

func (c *Client) CrontabEntryAdd(ctx context.Context, entry string) error {
	out, err := c.crontabL(ctx)
	if err != nil {
		return fmt.Errorf("failed to read crontab: %w", err)
	}

	out = bytes.TrimSpace(out)
	out = append(out, '\n')
	out = append(out, []byte(entry)...)

	return c.crontabE(ctx, out)
}

func (c *Client) CrontabEntryRemove(ctx context.Context, entry string) (bool, error) {
	out, err := c.crontabL(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to read crontab: %w", err)
	}

	var newLines [][]byte

	found := false

	for line := range bytes.Lines(out) {
		if string(line) == entry {
			found = true
			continue
		}

		newLines = append(newLines, line)
	}

	if !found {
		return false, nil
	}

	newCrontab := bytes.Join(newLines, []byte("\n"))

	return true, c.crontabE(ctx, newCrontab)
}

// if the user does not have a crontab, it returns an empty output and no error.
func (c *Client) crontabL(ctx context.Context) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "crontab", "-l")

	out, err := c.Runner.Run(cmd)

	switch {
	case err != nil && bytes.Contains(out, []byte("no crontab for")):
		return []byte{}, nil
	case err != nil:
		return nil, err
	default:
		return out, nil
	}
}

// crontabE sets the crontab.
func (c *Client) crontabE(ctx context.Context, crontab []byte) error {
	// ensure crontab ends with a newline
	crontab = bytes.TrimSpace(crontab)
	crontab = append(crontab, '\n')

	cmd := exec.CommandContext(ctx, "crontab", "-")
	cmd.Stdin = bytes.NewReader(crontab)

	_, err := c.Runner.Run(cmd)

	return err
}
