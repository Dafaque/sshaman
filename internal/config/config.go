package config

import (
	"fmt"
	"os"
	"path"
)

const (
	appHome            string = ".sshaman"
	homeOverrideEnvKey string = "SSHAMAN_HOME"
)

type Config struct {
	Home          string
	EncryptionKey []byte
}

func NewConfig() (*Config, error) {
	var cfg Config
	var home string = os.Getenv(homeOverrideEnvKey)
	if home == "" {
		userHome, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		home = userHome
	}
	home = path.Join(home, appHome)
	if _, err := os.Stat(home); os.IsNotExist(err) {
		if err := os.Mkdir(home, os.ModePerm); err != nil {
			return nil, fmt.Errorf("failed to init app home directory: %w", err)
		}
	}
	cfg.Home = home
	return &cfg, nil
}
