package uberspace

import (
	"context"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func (c *Client) RemoteFileCreate(ctx context.Context, src, dst string, executable bool) error {
	var data []byte

	var err error

	if strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://") {
		data, err = httpReadFile(ctx, src)
		if err != nil {
			return err
		}
	} else {
		data, err = os.ReadFile(src)
		if err != nil {
			return err
		}
	}

	if err := c.Runner.WriteFile(ctx, dst, data); err != nil {
		return err
	}

	if executable {
		_, err := c.Runner.Run(exec.CommandContext(ctx, "chmod", "+x", dst))
		if err != nil {
			return err
		}
	}

	return nil
}

func httpReadFile(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (c *Client) RemoteFileExists(ctx context.Context, dst string) (bool, error) {
	_, err := c.Runner.Run(exec.CommandContext(ctx, "test", "-e", dst))

	return err == nil, err
}

func (c *Client) RemoteFileDelete(ctx context.Context, dst string) (bool, error) {
	_, err := c.Runner.Run(exec.CommandContext(ctx, "rm", dst))

	return err == nil, err
}
