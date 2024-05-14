package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/vmihailenco/msgpack"
	"golang.org/x/term"

	"github.com/Dafaque/sshaman/internal/credentials"
	"github.com/Dafaque/sshaman/internal/encryption"
)

const (
	defaultDumpFileName = "sshaman.enc"
)

func exportCredentials(local credentials.Manager, remote credentials.Manager) error {
	creds, err := getCredentialsList(local, nil)
	if err != nil {
		return err
	}
	for _, cred := range creds {
		if cred.Local {
			cred.Source = nil
		}
	}

	data, err := msgpack.Marshal(creds)
	if err != nil {
		return err
	}
	// MARK: - Encrypt
	var password []byte
	if !flagSkipPassword {
		fd := int(os.Stdin.Fd())
		fmt.Printf("Enter password for dump file: ")
		password, err = term.ReadPassword(fd)
		if err != nil {
			return err
		}
		println()
	}
	encrypted, err := encryption.EncryptWithSecret(password, data)
	if err != nil {
		return err
	}

	// MARK: - Write
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(pwd, defaultDumpFileName))
	if err != nil {
		return err
	}

	defer file.Close()
	_, err = file.Write(encrypted)
	if err != nil {
		return err
	}
	println("Dump file written to", file.Name())
	return nil
}
