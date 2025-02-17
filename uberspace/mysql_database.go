package uberspace

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
)

func (c *Client) MySQLDatabaseCreate(ctx context.Context, name string) error {
	mySQLQuery := fmt.Sprintf("%q", "CREATE DATABASE "+name)
	cmd := exec.CommandContext(ctx, "mysql", "-e", mySQLQuery)

	_, err := c.Runner.Run(cmd)

	return err
}

func (c *Client) MySQLDatabaseRead(ctx context.Context, name string) (bool, error) {
	mySQLQuery := fmt.Sprintf("%q", "SHOW DATABASES LIKE '"+name+"'")
	cmd := exec.CommandContext(ctx, "mysql", "-e", mySQLQuery)

	out, err := c.Runner.Run(cmd)
	if err != nil {
		return false, err
	}

	return bytes.Contains(out, []byte(name)), nil
}

func (c *Client) MySQLDatabaseDrop(ctx context.Context, name string) (bool, error) {
	mySQLQuery := fmt.Sprintf("%q", "DROP DATABASE "+name)
	cmd := exec.CommandContext(ctx, "mysql", "-e", mySQLQuery)

	_, err := c.Runner.Run(cmd)
	if err != nil {
		return false, err
	}

	return true, nil
}
