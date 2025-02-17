package uberspace

import (
	"github.com/cugu/terraform-provider-uberspace/ssh"
)

type Client struct {
	User   string
	Runner *ssh.Client
}
