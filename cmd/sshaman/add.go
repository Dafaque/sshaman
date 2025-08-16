package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"golang.org/x/term"

	"github.com/Dafaque/sshaman/v2/internal/credentials"
)

func addConnection(manager *credentials.Manager) error {
	creds, err := makeNewCredentials()
	if err != nil {
		return err
	}

	if err := manager.Set(creds, flagForce); err != nil {
		return err
	}
	fmt.Println("local credentials added for", flagName)

	return nil
}

func makeNewCredentials() (*credentials.Credentials, error) {
	if flagHost == emptyString {
		return nil, errors.New("host required")
	}

	if flagUser == emptyString {
		return nil, errors.New("user required")
	}

	if flagName == emptyString {
		return nil, errors.New("alias required")
	}

	var creds credentials.Credentials = credentials.Credentials{
		Name:     flagName,
		Host:     flagHost,
		Port:     flagPort,
		UserName: flagUser,
	}

	fd := int(os.Stdin.Fd())
	var password string
	var passphrase []byte
	var key []byte
	if flagKeyFilePath != emptyString {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		if flagKeyFilePath == "~" {
			// In case of "~", which won't be caught by the "else if"
			flagKeyFilePath = home
		} else if strings.HasPrefix(flagKeyFilePath, "~/") {
			flagKeyFilePath = path.Join(home, flagKeyFilePath[2:])
		}
		file, err := os.ReadFile(flagKeyFilePath)
		if err != nil {
			return nil, err
		}
		key = file
		if !flagSkipPassphrase {
			fmt.Printf("Enter %s key passphrase: ", flagKeyFilePath)
			pp, err := term.ReadPassword(fd)
			println()
			if err != nil {
				return nil, err
			}
			passphrase = pp
		}
	}
	if !flagSkipPassword {
		fmt.Printf("Enter %s's password for %s: ", flagUser, flagHost)
		pass, err := term.ReadPassword(fd)
		println()
		if err != nil {
			return nil, err
		}
		password = string(pass)
	}

	if len(password) > 0 {
		creds.Password = &password
	}
	if len(passphrase) > 0 {
		creds.Passphrase = passphrase
	}
	if len(key) > 0 {
		creds.Key = key
	}

	return &creds, nil
}
