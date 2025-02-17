package ssh

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"os/exec"
	"strings"

	"github.com/bramvdbogaerde/go-scp"
	"golang.org/x/crypto/ssh"
)

type Config struct {
	User       string
	Host       string
	PrivateKey string
	Password   string
}

type Client struct {
	sshClient *ssh.Client
}

func NewClient(config *Config) (*Client, error) {
	sshConfig := &ssh.ClientConfig{
		User:            config.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //nolint: gosec
	}

	switch {
	case config.PrivateKey != "":
		signer, err := ssh.ParsePrivateKey([]byte(config.PrivateKey))
		if err != nil {
			return nil, err
		}

		sshConfig.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
	case config.Password != "":
		sshConfig.Auth = []ssh.AuthMethod{ssh.Password(config.Password)}
	default:
		return nil, fmt.Errorf("either private_key or password must be set")
	}

	client, err := ssh.Dial("tcp", net.JoinHostPort(config.Host, "22"), sshConfig)
	if err != nil {
		return nil, fmt.Errorf("dialing SSH: %w", err)
	}

	return &Client{sshClient: client}, nil
}

func (c *Client) Run(cmd *exec.Cmd) ([]byte, error) {
	session, err := c.sshClient.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	session.Stdin = cmd.Stdin

	command := strings.Join(cmd.Args, " ")

	out, err := session.CombinedOutput(command)
	if err != nil {
		return out, fmt.Errorf("%w: %s", err, string(out))
	}

	return out, nil
}

func (c *Client) WriteFile(ctx context.Context, path string, content []byte) error {
	scpClient, err := scp.NewClientBySSH(c.sshClient)
	if err != nil {
		return err
	}

	src := bytes.NewReader(content)

	if err := scpClient.CopyFile(ctx, src, path, "0655"); err != nil {
		return fmt.Errorf("copying file: %w", err)
	}

	return nil
}
