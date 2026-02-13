package encryption

import (
	"errors"
	"fmt"
	"os"

	"github.com/ProtonMail/gopenpgp/v2/helper"
)

// PGPEncryptor implements Encryptor using PGP (GPG) public-key encryption.
//
// Security properties:
//   - Asymmetric encryption using recipient's public key
//   - Decryption requires the corresponding private key + passphrase
//   - Messages are armored (ASCII-safe) for transport
type PGPEncryptor struct {
	publicKey  string
	privateKey string
	passphrase []byte
}

// NewPGPEncryptor creates a new PGP encryptor from key file paths and passphrase.
func NewPGPEncryptor(publicKeyFile, privateKeyFile, passphrase string) (*PGPEncryptor, error) {
	pubKey, err := os.ReadFile(publicKeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key file: %w", err)
	}

	privKey, err := os.ReadFile(privateKeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}

	if len(pubKey) == 0 {
		return nil, errors.New("public key file is empty")
	}

	if len(privKey) == 0 {
		return nil, errors.New("private key file is empty")
	}

	return &PGPEncryptor{
		publicKey:  string(pubKey),
		privateKey: string(privKey),
		passphrase: []byte(passphrase),
	}, nil
}

// NewPGPEncryptorFromKeys creates a new PGP encryptor directly from key strings.
// This is useful for testing or when keys are provided inline.
func NewPGPEncryptorFromKeys(publicKey, privateKey, passphrase string) (*PGPEncryptor, error) {
	if publicKey == "" {
		return nil, errors.New("public key cannot be empty")
	}

	if privateKey == "" {
		return nil, errors.New("private key cannot be empty")
	}

	return &PGPEncryptor{
		publicKey:  publicKey,
		privateKey: privateKey,
		passphrase: []byte(passphrase),
	}, nil
}

// Encrypt encrypts plaintext using the PGP public key.
func (e *PGPEncryptor) Encrypt(plaintext string) (string, error) {
	armor, err := helper.EncryptMessageArmored(e.publicKey, plaintext)
	if err != nil {
		return "", fmt.Errorf("PGP encryption failed: %w", err)
	}
	return armor, nil
}

// Decrypt decrypts PGP armored ciphertext using the private key and passphrase.
func (e *PGPEncryptor) Decrypt(ciphertext string) (string, error) {
	plaintext, err := helper.DecryptMessageArmored(e.privateKey, e.passphrase, ciphertext)
	if err != nil {
		return "", fmt.Errorf("PGP decryption failed: %w", err)
	}
	return plaintext, nil
}

// Name returns the encryption algorithm name.
func (e *PGPEncryptor) Name() string {
	return "PGP"
}
