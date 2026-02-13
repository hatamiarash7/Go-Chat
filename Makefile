BINARY   := go-chat
CMD_DIR  := ./cmd/go-chat
PKG      := github.com/hatamiarash7/go-chat/internal/version

VERSION  ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT   ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE     ?= $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
LDFLAGS  := -s -w -X $(PKG).version=$(VERSION) -X $(PKG).commit=$(COMMIT) -X $(PKG).buildDate=$(DATE)

-include .env

## Build & Run ---------------------------------------------------------------

.PHONY: build
build: clean ## Build the binary
	CGO_ENABLED=0 go build -trimpath -ldflags="$(LDFLAGS)" -o $(BINARY) $(CMD_DIR)

.PHONY: clean
clean: ## Remove build artifacts
	rm -f $(BINARY)

.PHONY: server
server: build ## Run the server
	START_MODE=server HOST=$(or $(HOST),localhost) PORT=$(or $(PORT),12345) ./$(BINARY)

.PHONY: client
client: build ## Run the client
	START_MODE=client \
		HOST=$(or $(HOST),localhost) \
		PORT=$(or $(PORT),12345) \
		ENCRYPTION=$(or $(ENCRYPTION),pgp) \
		PUBLIC_KEY_FILE=$(PUBLIC_KEY_FILE) \
		PRIVATE_KEY_FILE=$(PRIVATE_KEY_FILE) \
		PASSPHRASE="$(PASSPHRASE)" \
		./$(BINARY)

## Testing -------------------------------------------------------------------

.PHONY: test
test: ## Run all tests
	go test -v -race -count=1 ./...

.PHONY: test-short
test-short: ## Run tests (short mode)
	go test -short -race ./...

.PHONY: coverage
coverage: ## Run tests with coverage report
	go test -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

## Code Quality --------------------------------------------------------------

.PHONY: fmt
fmt: ## Format code
	gofmt -s -w .

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: lint
lint: ## Run golangci-lint (install: https://golangci-lint.run/welcome/install/)
	golangci-lint run ./...

.PHONY: check
check: fmt vet test ## Run all checks (format, vet, test)

## Docker --------------------------------------------------------------------

.PHONY: docker-build
docker-build: ## Build Docker image
	docker build -t go-chat-server:$(VERSION) .

.PHONY: docker-run
docker-run: ## Run Docker container
	docker run -it --rm -e PORT=$(or $(PORT),12345) -e HOST=0.0.0.0 -p $(or $(PORT),12345):$(or $(PORT),12345) go-chat-server:$(VERSION)

## Utilities -----------------------------------------------------------------

.PHONY: version
version: ## Show version info
	@echo "Version: $(VERSION)"
	@echo "Commit:  $(COMMIT)"
	@echo "Date:    $(DATE)"

.PHONY: deps
deps: ## Download and tidy dependencies
	go mod download
	go mod tidy

.PHONY: help
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help