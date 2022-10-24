package main

import (
	"bufio"
	"io"
	"net"
	"os"
	"strings"

	"github.com/ProtonMail/gopenpgp/v2/helper"
	log "github.com/sirupsen/logrus"
)

// Client model defines attributes of a client object.
type Client struct {
	id     string
	socket net.Conn
	data   chan []byte
}

var (
	passphrase []byte
	prvKey     string
)

func (client *Client) receive() {
	for {
		message := make([]byte, 4096)
		length, err := client.socket.Read(message)
		if err != nil {
			client.socket.Close()
			break
		}
		if length > 0 {
			messageParts := strings.Split(string(message), "$$$")
			clientID := messageParts[0]
			armor := messageParts[1]
			actualMessage, err := helper.DecryptMessageArmored(prvKey, passphrase, armor)
			if err == nil {
				log.Infof("New message from %s: %s", clientID, string(actualMessage))
			} else {
				log.Fatal(err)
			}
		}
	}
}

// StartClient method is used for starting a single client.
func StartClient() {
	publicKeyFile, ok := os.LookupEnv("PUBLIC_KEY_FILE")
	if !ok {
		log.Fatal("PUBLIC_KEY_FILE is not set")
	}

	privateKeyFile, ok := os.LookupEnv("PRIVATE_KEY_FILE")
	if !ok {
		log.Fatal("PRIVATE_KEY_FILE is not set")
	}

	pass, ok := os.LookupEnv("PASSPHRASE")
	if !ok {
		log.Fatal("PASSPHRASE is not set")
	} else {
		passphrase = []byte(pass)
	}

	publicKey, err := os.Open(publicKeyFile)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = publicKey.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	privateKey, err := os.Open(privateKeyFile)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = privateKey.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	public := new(strings.Builder)
	private := new(strings.Builder)

	io.Copy(public, publicKey)
	io.Copy(private, privateKey)

	prvKey = private.String()

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "12345"
	}

	hosts, ok := os.LookupEnv("HOSTS")
	if !ok {
		hosts = "localhost"
	}

	connection, err := net.Dial("tcp", hosts+":"+port)
	if err != nil {
		log.Error(err)
	}

	clientID := connection.LocalAddr().String()
	client := &Client{id: clientID, socket: connection}
	log.Infof("Starting client with ID : %s", clientID)

	go client.receive()
	for {
		reader := bufio.NewReader(os.Stdin)
		message, _ := reader.ReadString('\n')
		if strings.TrimRight(message, "\n") != "" {
			armor, err := helper.EncryptMessageArmored(public.String(), strings.TrimRight(message, "\n"))
			if err == nil {
				connection.Write([]byte(client.id + "$$$" + armor))
			} else {
				log.Fatal(err)
			}
		}
	}
}
