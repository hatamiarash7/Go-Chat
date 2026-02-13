package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"

	"golang.org/x/crypto/argon2"
)

const (
	// Argon2id parameters for key derivation.
	argon2Time    = 3
	argon2Memory  = 64 * 1024 // 64 MB
	argon2Threads = 4
	argon2KeyLen  = 32 // AES-256

	// Fixed salt for deterministic key derivation from passphrase.
	// Each message still gets a unique random nonce, ensuring ciphertext uniqueness.
	aesKeySalt = "go-chat-aes256gcm-v1"
)

// AESEncryptor implements Encryptor using AES-256-GCM with Argon2id key derivation.
//
// Security properties:
//   - AES-256-GCM provides authenticated encryption (confidentiality + integrity)
//   - Argon2id derives the encryption key from a passphrase (memory-hard, side-channel resistant)
//   - Each message uses a unique random 12-byte nonce
//   - Output format: base64(nonce || ciphertext || tag)
type AESEncryptor struct {
	gcm cipher.AEAD
}

// NewAESEncryptor creates a new AES-256-GCM encryptor with a key derived from the passphrase.
func NewAESEncryptor(passphrase string) (*AESEncryptor, error) {
	if passphrase == "" {
		return nil, errors.New("passphrase cannot be empty")
	}

	key := argon2.IDKey(
		[]byte(passphrase),
		[]byte(aesKeySalt),
		argon2Time,
		argon2Memory,
		argon2Threads,
		argon2KeyLen,
	)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	return &AESEncryptor{gcm: gcm}, nil
}

// Encrypt encrypts plaintext using AES-256-GCM with a random nonce.
func (e *AESEncryptor) Encrypt(plaintext string) (string, error) {
	nonce := make([]byte, e.gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Seal appends ciphertext+tag to nonce, so result = nonce || ciphertext || tag
	ciphertext := e.gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts an AES-256-GCM ciphertext encoded in base64.
func (e *AESEncryptor) Decrypt(encoded string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	nonceSize := e.gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := e.gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	return string(plaintext), nil
}

// Name returns the encryption algorithm name.
func (e *AESEncryptor) Name() string {
	return "AES-256-GCM"
}
