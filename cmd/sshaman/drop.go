package main

import (
	"errors"
	"fmt"

	"github.com/Dafaque/sshaman/internal/credentials"
)

func dropCredentials(manager *credentials.Manager) error {
	if !flagForce {
		return errors.New("this operation will delete all your data. if you are sure of what you are doing, use the flag -force")
	}
	if err := manager.Drop(); err != nil {
		return err
	}
	fmt.Println("local credentials cleared")

	return nil
}
