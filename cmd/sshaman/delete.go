package main

import (
	"fmt"

	"github.com/Dafaque/sshaman/v2/internal/credentials"
)

func removeCredentials(manager *credentials.Manager) error {
	if err := manager.Del(flagName); err != nil {
		return err
	}
	fmt.Println("credentials removed for", flagName)
	return nil
}
