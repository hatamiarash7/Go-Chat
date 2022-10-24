include .env

build: ## Build project
	go build

server: ## Run the server
	START_MODE=server ./go-chat --mode server

client: ## Run the client
	START_MODE=client PUBLIC_KEY_FILE=$(PUBLIC_KEY_FILE) PRIVATE_KEY_FILE=$(PRIVATE_KEY_FILE) PASSPHRASE="$(PASSPHRASE)" ./go-chat --mode client

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: build server client help
.DEFAULT_GOAL := help