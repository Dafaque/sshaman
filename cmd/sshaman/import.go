package main

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
	"hash/fnv"
	"os"
	"path/filepath"

	"github.com/vmihailenco/msgpack"
	"golang.org/x/term"

	"github.com/Dafaque/sshaman/internal/credentials"
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
	key := fnv.New128a()
	key.Write(password)
	c, err := aes.NewCipher(key.Sum(nil))
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return err
	}
	nonce, ciphertext := encrypted[:gcm.NonceSize()], encrypted[gcm.NonceSize():]

	decrypted, err := gcm.Open(nil, nonce, ciphertext, nil)
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
