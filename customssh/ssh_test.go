package customssh

import "testing"

func TestSSH(t *testing.T) {
	c := CustomSSH{
		Address: "127.0.0.1",
		Port:    22,
		User:    "root",
		Auth:    "root",
	}
	err := c.CheckSSHClient()
	if err != nil {
		t.Error(err)
	}
	t.Log("success")
}
