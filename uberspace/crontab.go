package uberspace

import (
	"bytes"
	"fmt"
)

func (c *Client) CrontabEntryAdd(entry string) error {
	out, err := c.crontabL()
	if err != nil {
		return fmt.Errorf("failed to read crontab: %w", err)
	}

	out = bytes.TrimSpace(out)
	out = append(out, '\n')
	out = append(out, []byte(entry)...)

	return c.crontabE(out)
}

func (c *Client) CrontabEntryExists(entry string) (bool, error) {
	out, err := c.crontabL()
	if err != nil {
		return false, fmt.Errorf("failed to read crontab: %w", err)
	}

	for line := range bytes.Lines(out) {
		if string(bytes.TrimSpace(line)) == entry {
			return true, nil
		}
	}

	return false, nil
}

func (c *Client) CrontabEntryRemove(entry string) (bool, error) {
	out, err := c.crontabL()
	if err != nil {
		return false, fmt.Errorf("failed to read crontab: %w", err)
	}

	var newLines [][]byte

	found := false

	for line := range bytes.Lines(out) {
		if string(bytes.TrimSpace(line)) == entry {
			found = true
			continue
		}

		newLines = append(newLines, line)
	}

	if !found {
		return false, nil
	}

	newCrontab := bytes.Join(newLines, []byte("\n"))

	return true, c.crontabE(newCrontab)
}

// if the user does not have a crontab, it returns an empty output and no error.
func (c *Client) crontabL() ([]byte, error) {
	out, err := c.SSHClient.Run("crontab -l")

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
func (c *Client) crontabE(crontab []byte) error {
	// ensure crontab ends with a newline
	crontab = append(bytes.TrimSpace(crontab), '\n')

	_, err := c.SSHClient.RunWithStdin("crontab -", bytes.NewReader(crontab))

	return err
}
