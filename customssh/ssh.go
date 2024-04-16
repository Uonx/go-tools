package customssh

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

type CustomSSH struct {
	Address string
	Port    int
	User    string
	Auth    string
}

func (c *CustomSSH) CheckSSHClient() error {
	config := ssh.ClientConfig{
		User:            c.User,
		Auth:            []ssh.AuthMethod{ssh.Password(c.Auth)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", c.Address, c.Port), &config)
	if err != nil {
		return fmt.Errorf("failed to dial: %#v", err)
	}
	defer client.Close()
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %#v", err)
	}
	defer session.Close()
	return nil
}
