package main

import (
	"github.com/Dafaque/sshaman/internal/credentials"
	"github.com/Dafaque/sshaman/pkg/clients/ssh"
)

func connect(manager *credentials.Manager) error {
	creds, err := manager.Get(flagName)
	if err != nil {
		return err
	}
	err = manager.Done()
	if err != nil {
		return err
	}
	cl, err := ssh.NewSshClient(creds)
	if err != nil {
		return err
	}
	return cl.Loop()
}
