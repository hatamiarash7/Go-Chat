package main

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	mode, ok := os.LookupEnv("START_MODE")

	if !ok {
		log.Fatal("START_MODE is not set")
		os.Exit(0)
	}

	if strings.ToLower(mode) == "server" {
		StartServer()
	} else {
		StartClient()
	}
}
