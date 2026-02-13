// Package client implements the Go-Chat client for sending and receiving
// encrypted messages through the relay server.
package client

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/hatamiarash7/go-chat/internal/encryption"
	"github.com/hatamiarash7/go-chat/internal/message"
	log "github.com/sirupsen/logrus"
)

// Client represents a chat client that connects to the relay server.
type Client struct {
	id        string
	conn      net.Conn
	encryptor encryption.Encryptor
	address   string
	ctx       context.Context
	cancel    context.CancelFunc
}

// New creates a new Client instance.
func New(address string, encryptor encryption.Encryptor) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		encryptor: encryptor,
		address:   address,
		ctx:       ctx,
		cancel:    cancel,
	}
}

// Connect establishes a TCP connection to the server.
func (c *Client) Connect() error {
	conn, err := net.Dial("tcp", c.address)
	if err != nil {
		return fmt.Errorf("failed to connect to server at %s: %w", c.address, err)
	}

	c.conn = conn
	c.id = conn.LocalAddr().String()
	log.Infof("Connected to server as %s (encryption: %s)", c.id, c.encryptor.Name())

	return nil
}

// Start begins receiving messages and reading user input.
// This method blocks until the context is cancelled or an error occurs.
func (c *Client) Start() error {
	if c.conn == nil {
		return fmt.Errorf("not connected to server")
	}

	go c.receiveMessages()

	return c.readInput()
}

// Stop gracefully disconnects the client.
func (c *Client) Stop() {
	log.Info("Disconnecting...")
	c.cancel()
	if c.conn != nil {
		c.conn.Close()
	}
}

// receiveMessages listens for incoming messages from the server.
func (c *Client) receiveMessages() {
	defer c.cancel()

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			msg, err := message.Read(c.conn)
			if err != nil {
				select {
				case <-c.ctx.Done():
					return
				default:
					log.Errorf("Connection lost: %v", err)
					return
				}
			}

			plaintext, err := c.encryptor.Decrypt(msg.Data)
			if err != nil {
				log.Warnf("Failed to decrypt message from %s: %v", msg.Sender, err)
				continue
			}

			fmt.Printf("\r[%s]: %s\n> ", msg.Sender, plaintext)
		}
	}
}

// readInput reads user input from stdin and sends encrypted messages.
func (c *Client) readInput() error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("> ")

	for {
		select {
		case <-c.ctx.Done():
			return nil
		default:
			text, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read input: %w", err)
			}

			text = strings.TrimSpace(text)
			if text == "" {
				fmt.Print("> ")
				continue
			}

			// Handle quit command.
			if text == "/quit" || text == "/exit" {
				c.Stop()
				return nil
			}

			encrypted, err := c.encryptor.Encrypt(text)
			if err != nil {
				log.Errorf("Failed to encrypt message: %v", err)
				fmt.Print("> ")
				continue
			}

			msg := &message.Message{
				Sender: c.id,
				Data:   encrypted,
			}

			if err := message.Write(c.conn, msg); err != nil {
				return fmt.Errorf("failed to send message: %w", err)
			}

			fmt.Print("> ")
		}
	}
}
