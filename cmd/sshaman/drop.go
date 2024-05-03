package main

import (
	"errors"
	"fmt"

	"github.com/Dafaque/sshaman/internal/credentials"
)

func dropCredentials(local credentials.Manager, remote credentials.Manager) error {
	if !flagForce {
		return errors.New("this operation will delete all your data. if you are sure of what you are doing, use the flag -force")
	}
	if flagLocal {
		if err := local.Drop(); err != nil {
			return err
		}
		fmt.Println("local credentials cleared")
	}
	if flagRemote {
		if remote == nil {
			return errRemoteNotConfigured
		}
		if err := remote.Drop(); err != nil {
			return err
		}
		fmt.Println("remote credentials cleared")
	}
	return nil
}
