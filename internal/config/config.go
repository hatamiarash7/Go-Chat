// Package config provides configuration management for the Go-Chat application.
package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// Mode represents the application start mode.
type Mode string

const (
	ModeServer Mode = "server"
	ModeClient Mode = "client"
)

// Encryption represents the encryption algorithm to use.
type Encryption string

const (
	EncryptionPGP Encryption = "pgp"
	EncryptionAES Encryption = "aes"
)

// Config holds the application configuration.
type Config struct {
	Mode           Mode
	Host           string
	Port           string
	PublicKeyFile  string
	PrivateKeyFile string
	Passphrase     string
	Encryption     Encryption
}

// Load reads configuration from environment variables and returns a Config.
func Load() (*Config, error) {
	mode, ok := os.LookupEnv("START_MODE")
	if !ok || mode == "" {
		return nil, errors.New("START_MODE is not set (must be 'server' or 'client')")
	}

	m := Mode(strings.ToLower(mode))
	if m != ModeServer && m != ModeClient {
		return nil, fmt.Errorf("invalid START_MODE: %q (must be 'server' or 'client')", mode)
	}

	cfg := &Config{
		Mode: m,
		Host: getEnvDefault("HOST", "localhost"),
		Port: getEnvDefault("PORT", "12345"),
	}

	enc := strings.ToLower(getEnvDefault("ENCRYPTION", "pgp"))
	switch Encryption(enc) {
	case EncryptionPGP:
		cfg.Encryption = EncryptionPGP
	case EncryptionAES:
		cfg.Encryption = EncryptionAES
	default:
		return nil, fmt.Errorf("invalid ENCRYPTION: %q (must be 'pgp' or 'aes')", enc)
	}

	if m == ModeClient {
		if err := cfg.loadClientConfig(); err != nil {
			return nil, err
		}
	}

	return cfg, nil
}

// loadClientConfig loads client-specific configuration.
func (c *Config) loadClientConfig() error {
	c.Passphrase = os.Getenv("PASSPHRASE")

	if c.Encryption == EncryptionPGP {
		pubFile, ok := os.LookupEnv("PUBLIC_KEY_FILE")
		if !ok || pubFile == "" {
			return errors.New("PUBLIC_KEY_FILE is required for PGP encryption")
		}
		c.PublicKeyFile = pubFile

		privFile, ok := os.LookupEnv("PRIVATE_KEY_FILE")
		if !ok || privFile == "" {
			return errors.New("PRIVATE_KEY_FILE is required for PGP encryption")
		}
		c.PrivateKeyFile = privFile

		if c.Passphrase == "" {
			return errors.New("PASSPHRASE is required for PGP encryption")
		}
	}

	if c.Encryption == EncryptionAES {
		if c.Passphrase == "" {
			return errors.New("PASSPHRASE is required for AES encryption")
		}
	}

	return nil
}

// Address returns the host:port address string.
func (c *Config) Address() string {
	return c.Host + ":" + c.Port
}

// getEnvDefault returns the environment variable value or a default.
func getEnvDefault(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultVal
}
