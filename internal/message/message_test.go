package message

import (
	"net"
	"testing"
	"time"
)

func TestMessage_EncodeDecodeCycle(t *testing.T) {
	original := &Message{
		Sender: "192.168.1.1:12345",
		Data:   "encrypted-data-here",
	}

	data, err := original.Encode()
	if err != nil {
		t.Fatalf("encode failed: %v", err)
	}

	decoded, err := Decode(data)
	if err != nil {
		t.Fatalf("decode failed: %v", err)
	}

	if decoded.Sender != original.Sender {
		t.Errorf("sender: expected %q, got %q", original.Sender, decoded.Sender)
	}
	if decoded.Data != original.Data {
		t.Errorf("data: expected %q, got %q", original.Data, decoded.Data)
	}
}

func TestDecode_InvalidJSON(t *testing.T) {
	_, err := Decode([]byte("not json"))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestWriteRead_RoundTrip(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()
	defer clientConn.Close()

	original := &Message{
		Sender: "client-1",
		Data:   "Hello, encrypted world!",
	}

	errChan := make(chan error, 1)
	msgChan := make(chan *Message, 1)

	go func() {
		errChan <- Write(clientConn, original)
	}()

	go func() {
		msg, err := Read(serverConn)
		if err != nil {
			errChan <- err
			return
		}
		msgChan <- msg
	}()

	select {
	case err := <-errChan:
		if err != nil {
			t.Fatalf("write failed: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("write timed out")
	}

	select {
	case msg := <-msgChan:
		if msg.Sender != original.Sender {
			t.Errorf("sender: expected %q, got %q", original.Sender, msg.Sender)
		}
		if msg.Data != original.Data {
			t.Errorf("data: expected %q, got %q", original.Data, msg.Data)
		}
	case err := <-errChan:
		t.Fatalf("read failed: %v", err)
	case <-time.After(2 * time.Second):
		t.Fatal("read timed out")
	}
}

func TestWriteRead_MultipleMessages(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()
	defer clientConn.Close()

	messages := []*Message{
		{Sender: "a", Data: "first"},
		{Sender: "b", Data: "second"},
		{Sender: "c", Data: "third"},
	}

	go func() {
		for _, msg := range messages {
			if err := Write(clientConn, msg); err != nil {
				t.Errorf("write failed: %v", err)
				return
			}
		}
	}()

	for i, expected := range messages {
		msg, err := Read(serverConn)
		if err != nil {
			t.Fatalf("read %d failed: %v", i, err)
		}
		if msg.Sender != expected.Sender || msg.Data != expected.Data {
			t.Errorf("message %d: expected {%s, %s}, got {%s, %s}",
				i, expected.Sender, expected.Data, msg.Sender, msg.Data)
		}
	}
}

func TestWriteRead_LargeMessage(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()
	defer clientConn.Close()

	largeData := make([]byte, 100*1024)
	for i := range largeData {
		largeData[i] = byte('A' + (i % 26))
	}

	original := &Message{
		Sender: "client-large",
		Data:   string(largeData),
	}

	go func() {
		if err := Write(clientConn, original); err != nil {
			t.Errorf("write failed: %v", err)
		}
	}()

	msg, err := Read(serverConn)
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}

	if msg.Data != original.Data {
		t.Error("large message data mismatch")
	}
}

func TestRead_ConnectionClosed(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	clientConn.Close()

	_, err := Read(serverConn)
	if err == nil {
		t.Fatal("expected error when reading from closed connection")
	}
	serverConn.Close()
}
