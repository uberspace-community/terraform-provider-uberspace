package uberspace

import (
	"bytes"
	"strconv"
)

func (c *Client) MySQLDatabaseCreate(name string) error {
	_, err := c.SSHClient.Run("mysql -e " + strconv.Quote("CREATE DATABASE "+name))

	return err
}

func (c *Client) MySQLDatabaseExists(name string) (bool, error) {
	out, err := c.SSHClient.Run("mysql -e " + strconv.Quote("SHOW DATABASES LIKE '"+name+"'"))
	if err != nil {
		return false, err
	}

	return bytes.Contains(out, []byte(name)), nil
}

func (c *Client) MySQLDatabaseDrop(name string) error {
	_, err := c.SSHClient.Run("mysql -e " + strconv.Quote("DROP DATABASE "+name))

	return err
}
