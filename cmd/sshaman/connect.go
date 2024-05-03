package main

import (
	"errors"

	"github.com/Dafaque/sshaman/internal/credentials"
	"github.com/Dafaque/sshaman/pkg/clients/ssh"
)

func connect(local credentials.Manager, remote credentials.Manager) error {
	var creds *credentials.Credentials
	if flagRemote {
		if remote == nil {
			return errRemoteNotConfigured
		}
		remoteCreds, err := remote.Get(flagAlias)
		if err != nil {
			return err
		}
		creds = remoteCreds
	}
	if flagLocal && creds == nil {
		localCreds, err := local.Get(flagAlias)
		if err != nil {
			return err
		}
		creds = localCreds
	}
	if creds == nil {
		return errors.New("no credentials source given")
	}
	cl, err := ssh.NewSshClient(creds)
	if err != nil {
		return err
	}
	return cl.Loop()
}
