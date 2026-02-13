// Go-Chat: A simple encrypted chat system using PUB/SUB model.
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"

	"github.com/hatamiarash7/go-chat/internal/client"
	"github.com/hatamiarash7/go-chat/internal/config"
	"github.com/hatamiarash7/go-chat/internal/encryption"
	"github.com/hatamiarash7/go-chat/internal/server"
	"github.com/hatamiarash7/go-chat/internal/version"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	// Load .env file if present (non-fatal if missing).
	_ = godotenv.Load()

	// Check for version flag.
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Println(version.Info())
		os.Exit(0)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	switch cfg.Mode {
	case config.ModeServer:
		runServer(cfg)
	case config.ModeClient:
		runClient(cfg)
	}
}

func runServer(cfg *config.Config) {
	srv := server.New(cfg.Address())

	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// Wait for interrupt signal for graceful shutdown.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	srv.Stop()
}

func runClient(cfg *config.Config) {
	enc, err := createEncryptor(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize encryption: %v", err)
	}

	c := client.New(cfg.Address(), enc)

	if err := c.Connect(); err != nil {
		log.Fatalf("Connection failed: %v", err)
	}

	// Handle interrupt signal.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		c.Stop()
	}()

	if err := c.Start(); err != nil {
		log.Fatalf("Client error: %v", err)
	}
}

func createEncryptor(cfg *config.Config) (encryption.Encryptor, error) {
	switch cfg.Encryption {
	case config.EncryptionAES:
		return encryption.NewAESEncryptor(cfg.Passphrase)
	case config.EncryptionPGP:
		return encryption.NewPGPEncryptor(cfg.PublicKeyFile, cfg.PrivateKeyFile, cfg.Passphrase)
	default:
		return nil, fmt.Errorf("unsupported encryption: %s", cfg.Encryption)
	}
}
