package uberspace

import (
	"bytes"
	"context"
	"fmt"
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

	if err := c.SSHClient.WriteFile(ctx, path, []byte(config)); err != nil {
		return err
	}

	_, err := c.SSHClient.Run("chmod +x " + path)
	if err != nil {
		return err
	}

	_, err = c.SSHClient.Run("supervisorctl reread")
	if err != nil {
		return err
	}

	_, err = c.SSHClient.Run("supervisorctl update")
	if err != nil {
		return err
	}

	_, err = c.SSHClient.Run("supervisorctl start " + name)

	return err
}

func (c *Client) SupervisorServiceExists(name string) (bool, error) {
	out, err := c.SSHClient.Run("supervisorctl status " + name)
	if err != nil {
		if bytes.Contains(out, []byte("no such process")) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (c *Client) SupervisorServiceRemove(user, name string) error {
	_, err := c.SSHClient.Run("supervisorctl stop " + name)
	if err != nil {
		return err
	}

	_, err = c.SSHClient.Run("supervisorctl remove " + name)
	if err != nil {
		return err
	}

	_, err = c.SSHClient.Run(fmt.Sprintf("rm /home/%s/etc/services.d/%s.ini", user, name))

	return err
}
