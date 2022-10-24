package main

import (
	"net"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

// ClientManager model defines attributes of a manager object.
type ClientManager struct {
	clients    map[*Client]string
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func (manager *ClientManager) start() {
	for {
		select {
		case connection := <-manager.register:
			manager.clients[connection] = connection.id
			log.Println("[NEW CONNECTION] Address : %s", connection.id)
		case connection := <-manager.unregister:
			_, ok := manager.clients[connection]
			if ok {
				close(connection.data)
				delete(manager.clients, connection)
				log.Infof("[CONNECTION CLOSED] Address : %s", connection.id)
			}
		case message := <-manager.broadcast:
			messageParts := strings.Split(string(message), "$$$")
			messageClientID := messageParts[0]
			for connection := range manager.clients {
				if messageClientID == connection.id {
					continue
				}
				select {
				case connection.data <- message:
				default:
					close(connection.data)
					delete(manager.clients, connection)
				}
			}
		}
	}
}

func (manager *ClientManager) receive(client *Client) {
	for {
		message := make([]byte, 4096)
		length, err := client.socket.Read(message)
		if err != nil {
			manager.unregister <- client
			client.socket.Close()
			break
		}

		if length > 0 {
			messageParts := strings.Split(string(message), "$$$")
			clientID := messageParts[0]
			log.Infof("Message from : %s", clientID)
			manager.broadcast <- message
		}
	}
}

func (manager *ClientManager) send(client *Client) {
	defer client.socket.Close()
	for {
		select {
		case message, ok := <-client.data:
			if !ok {
				return
			}
			client.socket.Write(message)
		}
	}
}

// StartServer method is used for starting the server.
func StartServer() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "12345"
	}

	hosts, ok := os.LookupEnv("HOSTS")
	if !ok {
		hosts = "localhost"
	}

	log.Info("Starting server...")
	log.Infof("Accepting connections on %s:%s", hosts, port)
	listener, error := net.Listen("tcp", hosts+":"+port)
	if error != nil {
		log.Error(error)
	}

	manager := ClientManager{
		clients:    make(map[*Client]string),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}

	go manager.start()
	for {
		connection, _ := listener.Accept()
		if error != nil {
			log.Error(error)
		}

		client := &Client{id: connection.RemoteAddr().String(), socket: connection, data: make(chan []byte)}
		manager.register <- client
		go manager.receive(client)
		go manager.send(client)
	}
}
