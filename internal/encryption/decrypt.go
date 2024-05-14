package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"hash/fnv"
)

func DecryptWithSecret(secret []byte, data []byte) ([]byte, error) {
	key := fnv.New128a()
	key.Write(secret)
	c, err := aes.NewCipher(key.Sum(nil))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}
	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]

	decrypted, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return decrypted, nil
}
