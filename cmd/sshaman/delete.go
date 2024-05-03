package main

import (
	"fmt"

	"github.com/Dafaque/sshaman/internal/credentials"
)

func deleteCredentials(local credentials.Manager, remote credentials.Manager) error {
	if flagLocal {
		if err := local.Del(flagAlias); err != nil {
			return err
		}
		fmt.Println("local credentials deleted for", flagAlias)
	}
	if flagRemote {
		if remote == nil {
			return errRemoteNotConfigured
		}
		if err := remote.Del(flagAlias); err != nil {
			return err
		}
		fmt.Println("remote credentials deleted for", flagAlias)
	}
	return nil
}
