package ssh

import (
	"golang.org/x/crypto/ssh"
)

type authMethodConfig struct {
	key        []byte
	passphrase []byte
	password   *string
}

func NewAuthMethodConfig() *authMethodConfig {
	return new(authMethodConfig)
}

func (amc *authMethodConfig) WithKey(key []byte) *authMethodConfig {
	return amc.WithKeyPassphrase(key, nil)
}

func (amc *authMethodConfig) WithKeyPassphrase(key, passphrase []byte) *authMethodConfig {
	amc.key = key
	amc.passphrase = passphrase
	return amc
}

func (amc *authMethodConfig) WithPassword(password string) *authMethodConfig {
	amc.password = &password
	return amc
}

func (amc *authMethodConfig) signerFromPem() (ssh.Signer, error) {
	if len(amc.passphrase) > 0 {
		return ssh.ParsePrivateKeyWithPassphrase(amc.key, amc.passphrase)
	}
	return ssh.ParsePrivateKey(amc.key)
}

func (amc *authMethodConfig) AuthMethods() ([]ssh.AuthMethod, error) {
	var authMethods []ssh.AuthMethod = make([]ssh.AuthMethod, 0)
	if len(amc.key) > 0 {
		amSigner, err := amc.signerFromPem()
		if err != nil {
			return nil, err
		}
		authMethods = append(authMethods, ssh.PublicKeys(amSigner))
	}
	if amc.password != nil && len(*amc.password) > 0 {
		authMethods = append(authMethods, ssh.Password(*amc.password))
	}
	return authMethods, nil
}
