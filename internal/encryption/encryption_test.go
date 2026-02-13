package encryption

import (
	"strings"
	"testing"
)

// --- AES-256-GCM Tests ---

func TestAESEncryptor_RoundTrip(t *testing.T) {
	enc, err := NewAESEncryptor("my-secret-passphrase")
	if err != nil {
		t.Fatalf("failed to create encryptor: %v", err)
	}

	tests := []string{
		"Hello, World!",
		"",
		"Short",
		strings.Repeat("Long message ", 100),
		"Special chars: !@#$%^&*()_+-=[]{}|;':\",./<>?",
		"Unicode: こんにちは世界",
	}

	for _, plaintext := range tests {
		name := plaintext
		if len(name) > 20 {
			name = name[:20]
		}
		if name == "" {
			name = "(empty)"
		}
		t.Run(name, func(t *testing.T) {
			ciphertext, err := enc.Encrypt(plaintext)
			if err != nil {
				t.Fatalf("encrypt failed: %v", err)
			}

			if plaintext != "" && ciphertext == plaintext {
				t.Fatal("ciphertext should differ from plaintext")
			}

			decrypted, err := enc.Decrypt(ciphertext)
			if err != nil {
				t.Fatalf("decrypt failed: %v", err)
			}

			if decrypted != plaintext {
				t.Errorf("expected %q, got %q", plaintext, decrypted)
			}
		})
	}
}

func TestAESEncryptor_UniqueNonces(t *testing.T) {
	enc, err := NewAESEncryptor("test-passphrase")
	if err != nil {
		t.Fatalf("failed to create encryptor: %v", err)
	}

	ct1, err := enc.Encrypt("same message")
	if err != nil {
		t.Fatalf("encrypt 1 failed: %v", err)
	}

	ct2, err := enc.Encrypt("same message")
	if err != nil {
		t.Fatalf("encrypt 2 failed: %v", err)
	}

	if ct1 == ct2 {
		t.Fatal("two encryptions of the same message should produce different ciphertexts")
	}
}

func TestAESEncryptor_WrongPassphrase(t *testing.T) {
	enc1, _ := NewAESEncryptor("correct-passphrase")
	enc2, _ := NewAESEncryptor("wrong-passphrase")

	ciphertext, err := enc1.Encrypt("secret message")
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	_, err = enc2.Decrypt(ciphertext)
	if err == nil {
		t.Fatal("decryption with wrong passphrase should fail")
	}
}

func TestAESEncryptor_EmptyPassphrase(t *testing.T) {
	_, err := NewAESEncryptor("")
	if err == nil {
		t.Fatal("expected error for empty passphrase")
	}
}

func TestAESEncryptor_InvalidCiphertext(t *testing.T) {
	enc, _ := NewAESEncryptor("test")

	tests := []struct {
		name  string
		input string
	}{
		{"not base64", "not-valid-base64!!!"},
		{"too short", "AAAA"},
		{"empty", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := enc.Decrypt(tt.input)
			if err == nil {
				t.Fatal("expected error for invalid ciphertext")
			}
		})
	}
}

func TestAESEncryptor_Name(t *testing.T) {
	enc, _ := NewAESEncryptor("test")
	if enc.Name() != "AES-256-GCM" {
		t.Errorf("expected name %q, got %q", "AES-256-GCM", enc.Name())
	}
}

// --- PGP Tests ---

func TestPGPEncryptor_EmptyKeys(t *testing.T) {
	_, err := NewPGPEncryptorFromKeys("", "private", "pass")
	if err == nil {
		t.Fatal("expected error for empty public key")
	}

	_, err = NewPGPEncryptorFromKeys("public", "", "pass")
	if err == nil {
		t.Fatal("expected error for empty private key")
	}
}

func TestPGPEncryptor_Name(t *testing.T) {
	enc := &PGPEncryptor{}
	if enc.Name() != "PGP" {
		t.Errorf("expected name %q, got %q", "PGP", enc.Name())
	}
}

func TestPGPEncryptor_NonexistentKeyFile(t *testing.T) {
	_, err := NewPGPEncryptor("/nonexistent/public.key", "/nonexistent/private.key", "pass")
	if err == nil {
		t.Fatal("expected error for nonexistent key files")
	}
}
