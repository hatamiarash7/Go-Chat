// Package encryption provides encryption interfaces and implementations for Go-Chat.
package encryption

// Encryptor defines the interface for message encryption and decryption.
type Encryptor interface {
	// Encrypt encrypts a plaintext message and returns the encoded ciphertext.
	Encrypt(plaintext string) (string, error)

	// Decrypt decrypts an encoded ciphertext and returns the plaintext message.
	Decrypt(ciphertext string) (string, error)

	// Name returns the name of the encryption algorithm.
	Name() string
}
