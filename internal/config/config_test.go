package config

import (
	"os"
	"testing"
)

func clearEnv() {
	os.Unsetenv("START_MODE")
	os.Unsetenv("HOST")
	os.Unsetenv("PORT")
	os.Unsetenv("ENCRYPTION")
	os.Unsetenv("PUBLIC_KEY_FILE")
	os.Unsetenv("PRIVATE_KEY_FILE")
	os.Unsetenv("PASSPHRASE")
}

func TestLoad_MissingStartMode(t *testing.T) {
	clearEnv()

	_, err := Load()
	if err == nil {
		t.Fatal("expected error when START_MODE is not set")
	}
}

func TestLoad_InvalidStartMode(t *testing.T) {
	clearEnv()
	os.Setenv("START_MODE", "invalid")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for invalid START_MODE")
	}
}

func TestLoad_ServerDefaults(t *testing.T) {
	clearEnv()
	os.Setenv("START_MODE", "server")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Mode != ModeServer {
		t.Errorf("expected mode %q, got %q", ModeServer, cfg.Mode)
	}
	if cfg.Host != "localhost" {
		t.Errorf("expected host %q, got %q", "localhost", cfg.Host)
	}
	if cfg.Port != "12345" {
		t.Errorf("expected port %q, got %q", "12345", cfg.Port)
	}
	if cfg.Encryption != EncryptionPGP {
		t.Errorf("expected encryption %q, got %q", EncryptionPGP, cfg.Encryption)
	}
}

func TestLoad_ServerCustom(t *testing.T) {
	clearEnv()
	os.Setenv("START_MODE", "SERVER")
	os.Setenv("HOST", "0.0.0.0")
	os.Setenv("PORT", "9999")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Mode != ModeServer {
		t.Errorf("expected mode %q, got %q", ModeServer, cfg.Mode)
	}
	if cfg.Host != "0.0.0.0" {
		t.Errorf("expected host %q, got %q", "0.0.0.0", cfg.Host)
	}
	if cfg.Port != "9999" {
		t.Errorf("expected port %q, got %q", "9999", cfg.Port)
	}
}

func TestLoad_ClientMissingPGPKeys(t *testing.T) {
	clearEnv()
	os.Setenv("START_MODE", "client")
	os.Setenv("PASSPHRASE", "test")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error when PUBLIC_KEY_FILE is missing for PGP")
	}
}

func TestLoad_ClientAES(t *testing.T) {
	clearEnv()
	os.Setenv("START_MODE", "client")
	os.Setenv("ENCRYPTION", "aes")
	os.Setenv("PASSPHRASE", "mysecretpassphrase")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Encryption != EncryptionAES {
		t.Errorf("expected encryption %q, got %q", EncryptionAES, cfg.Encryption)
	}
	if cfg.Passphrase != "mysecretpassphrase" {
		t.Errorf("expected passphrase %q, got %q", "mysecretpassphrase", cfg.Passphrase)
	}
}

func TestLoad_ClientAESMissingPassphrase(t *testing.T) {
	clearEnv()
	os.Setenv("START_MODE", "client")
	os.Setenv("ENCRYPTION", "aes")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error when PASSPHRASE is missing for AES")
	}
}

func TestLoad_InvalidEncryption(t *testing.T) {
	clearEnv()
	os.Setenv("START_MODE", "client")
	os.Setenv("ENCRYPTION", "rsa")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for invalid ENCRYPTION")
	}
}

func TestConfig_Address(t *testing.T) {
	cfg := &Config{Host: "127.0.0.1", Port: "8080"}
	expected := "127.0.0.1:8080"
	if addr := cfg.Address(); addr != expected {
		t.Errorf("expected %q, got %q", expected, addr)
	}
}
