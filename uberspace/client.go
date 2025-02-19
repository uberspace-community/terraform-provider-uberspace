package uberspace

import (
	"github.com/uberspace-community/terraform-provider-uberspace/ssh"
)

type Client struct {
	User      string
	SSHClient *ssh.Client
}
