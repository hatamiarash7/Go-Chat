// Package server implements the Go-Chat relay server.
//
// The server accepts TCP connections, manages connected clients, and broadcasts
// messages from any client to all other connected clients (PUB/SUB model).
package server

import (
	"context"
	"net"
	"sync"

	"github.com/hatamiarash7/go-chat/internal/message"
	log "github.com/sirupsen/logrus"
)

// client represents a connected chat client.
type client struct {
	id     string
	conn   net.Conn
	send   chan *message.Message
	server *Server
}

// Server manages client connections and message broadcasting.
type Server struct {
	address    string
	listener   net.Listener
	clients    map[*client]struct{}
	register   chan *client
	unregister chan *client
	broadcast  chan *message.Message
	mu         sync.RWMutex
	ctx        context.Context
	cancel     context.CancelFunc
}

// New creates a new Server instance bound to the given address.
func New(address string) *Server {
	ctx, cancel := context.WithCancel(context.Background())
	return &Server{
		address:    address,
		clients:    make(map[*client]struct{}),
		register:   make(chan *client),
		unregister: make(chan *client),
		broadcast:  make(chan *message.Message, 256),
		ctx:        ctx,
		cancel:     cancel,
	}
}

// Start begins listening for connections and processing messages.
func (s *Server) Start() error {
	var err error
	s.listener, err = net.Listen("tcp", s.address)
	if err != nil {
		return err
	}

	log.Infof("Server listening on %s", s.address)

	go s.processEvents()
	go s.acceptConnections()

	return nil
}

// Stop gracefully shuts down the server.
func (s *Server) Stop() {
	log.Info("Shutting down server...")
	s.cancel()

	if s.listener != nil {
		s.listener.Close()
	}

	s.mu.Lock()
	for c := range s.clients {
		close(c.send)
		c.conn.Close()
		delete(s.clients, c)
	}
	s.mu.Unlock()

	log.Info("Server stopped")
}

// Addr returns the listener's address, useful for tests with port 0.
func (s *Server) Addr() net.Addr {
	if s.listener != nil {
		return s.listener.Addr()
	}
	return nil
}

// acceptConnections listens for new TCP connections.
func (s *Server) acceptConnections() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.ctx.Done():
				return
			default:
				log.Errorf("Failed to accept connection: %v", err)
				continue
			}
		}

		c := &client{
			id:     conn.RemoteAddr().String(),
			conn:   conn,
			send:   make(chan *message.Message, 64),
			server: s,
		}

		s.register <- c
	}
}

// processEvents handles client registration, unregistration, and broadcasts.
func (s *Server) processEvents() {
	for {
		select {
		case <-s.ctx.Done():
			return

		case c := <-s.register:
			s.mu.Lock()
			s.clients[c] = struct{}{}
			s.mu.Unlock()
			log.Infof("Client connected: %s (total: %d)", c.id, len(s.clients))
			go s.readFromClient(c)
			go s.writeToClient(c)

		case c := <-s.unregister:
			s.mu.Lock()
			if _, ok := s.clients[c]; ok {
				delete(s.clients, c)
				close(c.send)
				c.conn.Close()
				log.Infof("Client disconnected: %s (total: %d)", c.id, len(s.clients))
			}
			s.mu.Unlock()

		case msg := <-s.broadcast:
			s.mu.RLock()
			for c := range s.clients {
				// Don't echo back to sender.
				if c.id == msg.Sender {
					continue
				}
				select {
				case c.send <- msg:
				default:
					// Client's send buffer is full; disconnect.
					log.Warnf("Client %s send buffer full, disconnecting", c.id)
					go func(cl *client) { s.unregister <- cl }(c)
				}
			}
			s.mu.RUnlock()
		}
	}
}

// readFromClient reads messages from a client and broadcasts them.
func (s *Server) readFromClient(c *client) {
	defer func() {
		s.unregister <- c
	}()

	for {
		msg, err := message.Read(c.conn)
		if err != nil {
			select {
			case <-s.ctx.Done():
			default:
				log.Debugf("Read error from %s: %v", c.id, err)
			}
			return
		}

		log.Infof("Message from %s", msg.Sender)
		s.broadcast <- msg
	}
}

// writeToClient sends queued messages to a client.
func (s *Server) writeToClient(c *client) {
	for msg := range c.send {
		if err := message.Write(c.conn, msg); err != nil {
			log.Debugf("Write error to %s: %v", c.id, err)
			s.unregister <- c
			return
		}
	}
}
