package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"hash/fnv"
	"io"
	"os"
	"path/filepath"

	"github.com/vmihailenco/msgpack"
	"golang.org/x/term"

	"github.com/Dafaque/sshaman/internal/credentials"
)

const (
	defaultDumpFileName = "sshaman.enc"
)

func exportCredentials(manager *credentials.Manager) error {
	creds, err := getCredentialsList(manager)
	if err != nil {
		return err
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
	key := fnv.New128a()
	key.Write(password)
	c, err := aes.NewCipher(key.Sum(nil))
	if err != nil {
		return err
	}
	key.Reset()
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return err
	}
	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return err
	}
	encrypted := gcm.Seal(nonce, nonce, data, nil)

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
