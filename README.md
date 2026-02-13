# Go Chat

[![Release version][badge_release_version]][link_releases]
![Project language][badge_language]
[![Release][release_badge]][release_link]
[![License][badge_license]][link_license]
[![Image size][badge_size_latest]][link_docker_hub]

A simple encrypted chat system using a **PUB/SUB** single channel. Supports **PGP** and **AES-256-GCM** end-to-end encryption. Any message sent from a client (publisher) is encrypted and routed to every other client (subscriber) on the channel.

- [Go Chat](#go-chat)
  - [Architecture](#architecture)
  - [Encryption Modes](#encryption-modes)
    - [PGP Mode (Default)](#pgp-mode-default)
    - [AES-256-GCM Mode](#aes-256-gcm-mode)
  - [Configure](#configure)
  - [Key Generation](#key-generation)
    - [PGP Keys](#pgp-keys)
    - [AES Mode (No Key Files Needed)](#aes-mode-no-key-files-needed)
  - [Usage](#usage)
    - [Binaries](#binaries)
    - [Docker](#docker)
    - [Build from source](#build-from-source)
  - [Development](#development)
  - [Project Structure](#project-structure)
  - [Support](#support)
  - [Contributing](#contributing)
  - [Issues](#issues)

## Architecture

All messages are encrypted client-side before being sent to the relay server. The server never sees plaintext — it only routes opaque encrypted payloads between clients.

```text
             ┌────────┐                  ┌────────┐
Message ─────► Client │                  │ Client ├──Decrypt──► Message
             └───┬────┘                  └────▲───┘
                 │                            │
                 │           ┌────────┐       │
                 └─Encrypt───► Server ├───────┘
                             └────────┘
```

Messages use **length-prefixed binary framing** (4-byte big-endian header + JSON payload) for reliable delivery over TCP, replacing the previous delimiter-based protocol.

## Encryption Modes

| Mode  | Algorithm              | Type       | Key Material               | Best For                      |
| ----- | ---------------------- | ---------- | -------------------------- | ----------------------------- |
| `pgp` | GPG/PGP                | Asymmetric | Public + Private key files | High-security, identity-based |
| `aes` | AES-256-GCM + Argon2id | Symmetric  | Shared passphrase          | Simple setup, fast encryption |

### PGP Mode (Default)

Uses public-key cryptography. Each client encrypts with the shared public key and decrypts with the private key + passphrase. Ideal when you need strong identity-based security.

### AES-256-GCM Mode

Uses authenticated symmetric encryption with:

- **Argon2id** key derivation (memory-hard, 64MB, 3 iterations) from the shared passphrase
- **Random 12-byte nonce** per message (no nonce reuse)
- **Authenticated encryption** (confidentiality + integrity + tamper detection)

This is the recommended mode for simplicity and performance.

## Configure

Configuration is done through environment variables. You can create a `.env` file in the project root:

```env
START_MODE=server
HOST=localhost
PORT=12345
ENCRYPTION=pgp
PUBLIC_KEY_FILE=./keys/public.asc
PRIVATE_KEY_FILE=./keys/private.asc
PASSPHRASE=your-secure-passphrase
```

| Variable           | Description                             | Default     | Required                   |
| ------------------ | --------------------------------------- | ----------- | -------------------------- |
| `START_MODE`       | Application mode                        | —           | Yes (`server` or `client`) |
| `PORT`             | TCP port                                | `12345`     | No                         |
| `HOST`             | Bind/connect address                    | `localhost` | No                         |
| `ENCRYPTION`       | Encryption algorithm                    | `pgp`       | No (`pgp` or `aes`)        |
| `PUBLIC_KEY_FILE`  | Path to PGP public key                  | —           | PGP mode only              |
| `PRIVATE_KEY_FILE` | Path to PGP private key                 | —           | PGP mode only              |
| `PASSPHRASE`       | PGP key passphrase or AES shared secret | —           | Client mode                |

> **Note**: Escape any special characters in `PASSPHRASE` when using shell commands.

## Key Generation

### PGP Keys

Generate a PGP key pair using GPG:

```bash
# Generate a new key pair (follow the prompts)
gpg --full-generate-key

# List your keys to find the key ID
gpg --list-keys

# Export public key (replace KEY_ID with your key ID or email)
gpg --armor --export KEY_ID > keys/public.asc

# Export private key
gpg --armor --export-secret-keys KEY_ID > keys/private.asc
```

**Quick generation** (non-interactive, for testing):

```bash
mkdir -p keys

# Generate key with predefined settings
gpg --batch --gen-key <<EOF
Key-Type: RSA
Key-Length: 4096
Subkey-Type: RSA
Subkey-Length: 4096
Name-Real: Go Chat User
Name-Email: chat@example.com
Expire-Date: 0
Passphrase: your-secure-passphrase
%commit
EOF

# Export keys
gpg --armor --export chat@example.com > keys/public.asc
gpg --armor --export-secret-keys chat@example.com > keys/private.asc
```

> [!WARNING]
> You can't use PGP keys on smart-cards like YubiKey because the `gopenpgp` requires raw private key material in memory.

### AES Mode (No Key Files Needed)

For AES-256-GCM mode, only a shared passphrase is needed:

```bash
# All clients must use the same passphrase
export ENCRYPTION=aes
export PASSPHRASE="a-strong-shared-secret-at-least-16-chars"
```

Generate a strong random passphrase:

```bash
# Using openssl
openssl rand -base64 32

# Using /dev/urandom
head -c 32 /dev/urandom | base64
```

## Usage

### Binaries

Download the latest release from the [releases page](https://github.com/hatamiarash7/Go-Chat/releases/latest).

**Server:**

```bash
START_MODE=server PORT=12345 ./go-chat
```

**Client (PGP):**

```bash
START_MODE=client ENCRYPTION=pgp \
  PUBLIC_KEY_FILE=./keys/public.asc \
  PRIVATE_KEY_FILE=./keys/private.asc \
  PASSPHRASE="your-passphrase" \
  HOST=server-address PORT=12345 \
  ./go-chat
```

**Client (AES):**

```bash
START_MODE=client ENCRYPTION=aes \
  PASSPHRASE="shared-secret" \
  HOST=server-address PORT=12345 \
  ./go-chat
```

**Client Commands:**

- Type a message and press Enter to send
- `/quit` or `/exit` to disconnect
- `Ctrl+C` for graceful shutdown

> **Note**: For macOS, you may need to allow the binary: `System Preferences > Security & Privacy > General > Open Anyway`.

### Docker

Run the server in a Docker container:

```bash
docker run -it -p 12345:12345 hatamiarash7/go-chat-server
```

With custom settings:

```bash
docker run -it -p 9999:9999 \
  -e PORT=9999 \
  -e HOST=0.0.0.0 \
  hatamiarash7/go-chat-server
```

> The Docker image defaults: `HOST=0.0.0.0`, `PORT=12345`, `START_MODE=server`.

### Build from source

```bash
# Build
make build

# Run server
make server

# Run client (set env vars first or create .env)
make client

# Show version
./go-chat --version
```

## Development

```bash
# Run all tests with race detection
make test

# Run tests with coverage report
make coverage

# Format code
make fmt

# Run go vet
make vet

# Run all checks (format + vet + test)
make check

# Build Docker image
make docker-build

# Tidy dependencies
make deps
```

## Project Structure

```text
.
├── cmd/
│   └── go-chat/
│       └── main.go              # Application entrypoint
├── internal/
│   ├── client/
│   │   └── client.go            # Chat client implementation
│   ├── config/
│   │   ├── config.go            # Configuration management
│   │   └── config_test.go       # Config tests
│   ├── encryption/
│   │   ├── encryption.go        # Encryptor interface
│   │   ├── aes.go               # AES-256-GCM implementation
│   │   ├── pgp.go               # PGP implementation
│   │   └── encryption_test.go   # Encryption tests
│   ├── message/
│   │   ├── message.go           # Message framing protocol
│   │   └── message_test.go      # Message tests
│   ├── server/
│   │   ├── server.go            # Relay server implementation
│   │   └── server_test.go       # Server tests
│   └── version/
│       └── version.go           # Build-time version info
├── .dockerignore
├── .env.example
├── .gitignore
├── Dockerfile
├── go.mod
├── go.sum
├── LICENSE
├── Makefile
└── README.md
```

---

## Support

[![Donate with Bitcoin](https://img.shields.io/badge/Bitcoin-bc1qmmh6vt366yzjt3grjxjjqynrrxs3frun8gnxrz-orange)](https://donatebadges.ir/donate/Bitcoin/bc1qmmh6vt366yzjt3grjxjjqynrrxs3frun8gnxrz) [![Donate with Ethereum](https://img.shields.io/badge/Ethereum-0x0831bD72Ea8904B38Be9D6185Da2f930d6078094-blueviolet)](https://donatebadges.ir/donate/Ethereum/0x0831bD72Ea8904B38Be9D6185Da2f930d6078094)

[![ko-fi](https://www.ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/D1D1WGU9)

<div><a href="https://payping.ir/@hatamiarash7"><img src="https://cdn.payping.ir/statics/Payping-logo/Trust/blue.svg" height="128" width="128"></a></div>

## Contributing

Contributions are welcome! Here's how you can help:

1. Fork the repository
2. Create your feature branch: `git checkout -b feature/amazing-feature`
3. Commit your changes: `git commit -m 'Add amazing feature'`
4. Push to the branch: `git push origin feature/amazing-feature`
5. Open a Pull Request

## Issues

Each project may have many problems. Contributing to the better development of this project by reporting them. 👍

[release_badge]: https://github.com/hatamiarash7/Go-Chat/actions/workflows/release.yaml/badge.svg
[release_link]: https://github.com/hatamiarash7/Go-Chat/actions/workflows/release.yaml
[link_license]: https://github.com/hatamiarash7/go-chat/blob/master/LICENSE
[badge_license]: https://img.shields.io/github/license/hatamiarash7/go-chat.svg?longCache=true
[badge_size_latest]: https://img.shields.io/docker/image-size/hatamiarash7/go-chat-server/latest?maxAge=30
[link_docker_hub]: https://hub.docker.com/r/hatamiarash7/go-chat-server/
[badge_release_version]: https://img.shields.io/github/release/hatamiarash7/go-chat.svg?maxAge=30&label=Release
[link_releases]: https://github.com/hatamiarash7/go-chat/releases
[badge_language]: https://img.shields.io/github/go-mod/go-version/hatamiarash7/go-chat?longCache=true
