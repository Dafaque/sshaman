package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/vmihailenco/msgpack"
	"golang.org/x/term"

	"github.com/Dafaque/sshaman/internal/credentials"
	"github.com/Dafaque/sshaman/internal/encryption"
)

func importCredentials(local credentials.Manager, remote credentials.Manager) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	encrypted, err := os.ReadFile(filepath.Join(wd, defaultDumpFileName))
	if err != nil {
		return err
	}
	fd := int(os.Stdin.Fd())
	var password []byte
	if !flagSkipPassword {
		fmt.Printf("Enter password for dump file: ")
		password, err = term.ReadPassword(fd)
		if err != nil {
			return err
		}
		println()
	}
	decrypted, err := encryption.DecryptWithSecret(password, encrypted)
	if err != nil {
		return err
	}

	var creds []*credentials.Credentials
	err = msgpack.Unmarshal(decrypted, &creds)
	if err != nil {
		return err
	}
	if flagDryRun {
		fmt.Println("Dry run, not importing credentials:")
		displayListCredentials(creds)
		return nil
	}
	for _, cred := range creds {
		if cred.Source != nil {
			continue
		}
		err = local.Set(cred, flagForce)
		if err != nil {
			if errors.Is(err, credentials.ErrCredentialsExist) {
				fmt.Println("Credentials already exists, skipping: ", cred.Alias)
				continue
			}
			return err
		}
	}
	println("Credentials imported")
	return listCredentials(local, remote)
}
