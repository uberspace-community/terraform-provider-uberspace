package uberspace

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

const (
	supervisorServiceTemplate = `[program:%s]
command=%s`
	supervisorServiceTemplateWithEnv = `[program:%s]
command=%s
environment=%s`
)

func (c *Client) SupervisorServiceCreate(ctx context.Context, user, name, command string, env map[string]string) error {
	var config string

	if env != nil {
		var envs []string
		for k, v := range env {
			envs = append(envs, fmt.Sprintf("%s=%s", k, v))
		}

		config = fmt.Sprintf(supervisorServiceTemplateWithEnv, name, command, strings.Join(envs, ","))
	} else {
		config = fmt.Sprintf(supervisorServiceTemplate, name, command)
	}

	path := fmt.Sprintf("/home/%s/etc/services.d/%s.ini", user, name)

	if err := c.Runner.WriteFile(ctx, path, []byte(config)); err != nil {
		return err
	}

	_, err := c.Runner.Run(exec.CommandContext(ctx, "chmod", "+x", path))
	if err != nil {
		return err
	}

	_, err = c.Runner.Run(exec.CommandContext(ctx, "supervisorctl", "reread"))
	if err != nil {
		return err
	}

	_, err = c.Runner.Run(exec.CommandContext(ctx, "supervisorctl", "update"))
	if err != nil {
		return err
	}

	_, err = c.Runner.Run(exec.CommandContext(ctx, "supervisorctl", "start", name))

	return err
}

func (c *Client) SupervisorServiceRead(ctx context.Context, name string) (bool, error) {
	cmd := exec.CommandContext(ctx, "supervisorctl", "status", name)

	out, err := c.Runner.Run(cmd)
	if err != nil {
		if bytes.Contains(out, []byte("no such process")) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (c *Client) SupervisorServiceDrop(ctx context.Context, user, name string) (bool, error) {
	_, err := c.Runner.Run(exec.CommandContext(ctx, "supervisorctl", "stop", name))
	if err != nil {
		return false, err
	}

	_, err = c.Runner.Run(exec.CommandContext(ctx, "supervisorctl", "remove", name))
	if err != nil {
		return false, err
	}

	_, err = c.Runner.Run(exec.CommandContext(ctx, "rm", fmt.Sprintf("/home/%s/etc/services.d/%s.ini", user, name))) //nolint: gosec
	if err != nil {
		return false, err
	}

	return true, nil
}
