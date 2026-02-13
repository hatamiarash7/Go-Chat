// Package message provides message encoding, decoding, and length-prefixed
// framing for reliable TCP transport.
package message

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
)

// MaxMessageSize is the maximum allowed message size (1 MB).
const MaxMessageSize = 1 << 20

// Message represents a chat message exchanged between clients via the server.
type Message struct {
	Sender string `json:"sender"`
	Data   string `json:"data"`
}

// Encode serializes a Message to JSON bytes.
func (m *Message) Encode() ([]byte, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("failed to encode message: %w", err)
	}
	return data, nil
}

// Decode deserializes JSON bytes into a Message.
func Decode(data []byte) (*Message, error) {
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, fmt.Errorf("failed to decode message: %w", err)
	}
	return &msg, nil
}

// Write sends a length-prefixed message over a TCP connection.
// Format: [4 bytes big-endian length][message payload]
func Write(conn net.Conn, msg *Message) error {
	data, err := msg.Encode()
	if err != nil {
		return err
	}

	if len(data) > MaxMessageSize {
		return fmt.Errorf("message too large: %d bytes (max %d)", len(data), MaxMessageSize)
	}

	// Write 4-byte length header.
	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, uint32(len(data)))
	if _, err := conn.Write(header); err != nil {
		return fmt.Errorf("failed to write message header: %w", err)
	}

	// Write message payload.
	if _, err := conn.Write(data); err != nil {
		return fmt.Errorf("failed to write message payload: %w", err)
	}

	return nil
}

// Read reads a length-prefixed message from a TCP connection.
// Returns the decoded Message or an error (including io.EOF on connection close).
func Read(conn net.Conn) (*Message, error) {
	// Read 4-byte length header.
	header := make([]byte, 4)
	if _, err := io.ReadFull(conn, header); err != nil {
		return nil, err
	}

	length := binary.BigEndian.Uint32(header)
	if length == 0 {
		return nil, fmt.Errorf("invalid message: zero length")
	}
	if length > MaxMessageSize {
		return nil, fmt.Errorf("message too large: %d bytes (max %d)", length, MaxMessageSize)
	}

	// Read message payload.
	payload := make([]byte, length)
	if _, err := io.ReadFull(conn, payload); err != nil {
		return nil, fmt.Errorf("failed to read message payload: %w", err)
	}

	return Decode(payload)
}
