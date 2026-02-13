package server

import (
	"net"
	"testing"
	"time"

	"github.com/hatamiarash7/go-chat/internal/message"
)

func startTestServer(t *testing.T) *Server {
	t.Helper()
	srv := New("127.0.0.1:0")
	if err := srv.Start(); err != nil {
		t.Fatalf("failed to start server: %v", err)
	}
	return srv
}

func dialWithTimeout(addr string) (net.Conn, error) {
	return net.DialTimeout("tcp", addr, 2*time.Second)
}

func TestServer_StartStop(t *testing.T) {
	srv := startTestServer(t)
	defer srv.Stop()

	addr := srv.Addr()
	if addr == nil {
		t.Fatal("server address should not be nil")
	}
}

func TestServer_ClientConnect(t *testing.T) {
	srv := startTestServer(t)
	defer srv.Stop()

	addr := srv.Addr().String()

	conn, err := dialWithTimeout(addr)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	time.Sleep(50 * time.Millisecond)

	srv.mu.RLock()
	count := len(srv.clients)
	srv.mu.RUnlock()

	if count != 1 {
		t.Errorf("expected 1 client, got %d", count)
	}
}

func TestServer_MessageBroadcast(t *testing.T) {
	srv := startTestServer(t)
	defer srv.Stop()

	addr := srv.Addr().String()

	conn1, err := dialWithTimeout(addr)
	if err != nil {
		t.Fatalf("client 1 connect failed: %v", err)
	}
	defer conn1.Close()

	conn2, err := dialWithTimeout(addr)
	if err != nil {
		t.Fatalf("client 2 connect failed: %v", err)
	}
	defer conn2.Close()

	time.Sleep(50 * time.Millisecond)

	msg := &message.Message{
		Sender: conn1.LocalAddr().String(),
		Data:   "hello from client 1",
	}

	if err := message.Write(conn1, msg); err != nil {
		t.Fatalf("write failed: %v", err)
	}

	conn2.SetReadDeadline(time.Now().Add(2 * time.Second))
	received, err := message.Read(conn2)
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}

	if received.Data != msg.Data {
		t.Errorf("expected data %q, got %q", msg.Data, received.Data)
	}
	if received.Sender != msg.Sender {
		t.Errorf("expected sender %q, got %q", msg.Sender, received.Sender)
	}
}

func TestServer_NoEchoToSender(t *testing.T) {
	srv := startTestServer(t)
	defer srv.Stop()

	addr := srv.Addr().String()

	conn, err := dialWithTimeout(addr)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer conn.Close()

	time.Sleep(50 * time.Millisecond)

	msg := &message.Message{
		Sender: conn.LocalAddr().String(),
		Data:   "echo test",
	}

	if err := message.Write(conn, msg); err != nil {
		t.Fatalf("write failed: %v", err)
	}

	conn.SetReadDeadline(time.Now().Add(200 * time.Millisecond))
	_, err = message.Read(conn)
	if err == nil {
		t.Fatal("sender should not receive its own message")
	}
}

func TestServer_ClientDisconnect(t *testing.T) {
	srv := startTestServer(t)
	defer srv.Stop()

	addr := srv.Addr().String()

	conn, err := dialWithTimeout(addr)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}

	time.Sleep(50 * time.Millisecond)

	srv.mu.RLock()
	count := len(srv.clients)
	srv.mu.RUnlock()
	if count != 1 {
		t.Fatalf("expected 1 client, got %d", count)
	}

	conn.Close()
	time.Sleep(100 * time.Millisecond)

	srv.mu.RLock()
	count = len(srv.clients)
	srv.mu.RUnlock()
	if count != 0 {
		t.Errorf("expected 0 clients after disconnect, got %d", count)
	}
}

func TestServer_MultipleClients(t *testing.T) {
	srv := startTestServer(t)
	defer srv.Stop()

	addr := srv.Addr().String()

	conns := make([]net.Conn, 5)
	for i := range conns {
		conn, err := dialWithTimeout(addr)
		if err != nil {
			t.Fatalf("client %d connect failed: %v", i, err)
		}
		conns[i] = conn
	}
	defer func() {
		for _, c := range conns {
			c.Close()
		}
	}()

	time.Sleep(100 * time.Millisecond)

	srv.mu.RLock()
	count := len(srv.clients)
	srv.mu.RUnlock()

	if count != 5 {
		t.Errorf("expected 5 clients, got %d", count)
	}
}
