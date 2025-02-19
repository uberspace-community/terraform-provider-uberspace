package ssh

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"

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

func (c *Client) Run(cmd string) ([]byte, error) {
	return c.RunWithStdin(cmd, nil)
}

func (c *Client) RunWithStdin(cmd string, stdin io.Reader) ([]byte, error) {
	session, err := c.sshClient.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	session.Stdin = stdin

	out, err := session.CombinedOutput(cmd)
	if err != nil {
		if out != nil {
			err = fmt.Errorf("%w: %s", err, string(out))
		}

		return nil, err
	}

	return out, nil
}

func (c *Client) WriteFile(ctx context.Context, path string, content []byte) error {
	scpClient, err := scp.NewClientBySSH(c.sshClient)
	if err != nil {
		return err
	}

	if err := scpClient.CopyFile(ctx, bytes.NewReader(content), path, "0655"); err != nil {
		return fmt.Errorf("copying file: %w", err)
	}

	return nil
}
