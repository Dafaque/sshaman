package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"hash/fnv"
	"io"
)

func EncryptWithSecret(secret []byte, data []byte) ([]byte, error) {
	key := fnv.New128a()
	key.Write(secret)

	c, err := aes.NewCipher(key.Sum(nil))
	if err != nil {
		return nil, err
	}
	key.Reset()

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}
	encrypted := gcm.Seal(nonce, nonce, data, nil)
	return encrypted, nil
}
