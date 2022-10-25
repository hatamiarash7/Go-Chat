# Go Chat

[![Release version][badge_release_version]][link_releases]
![Project language][badge_language]
[![Release][release_badge]][release_link]
[![License][badge_license]][link_license]
[![Image size][badge_size_latest]][link_docker_hub]

It's a simple chat system using a **PUB/SUB** single channel and **GPG encryption**. Any message sent from a client ( publisher ) is routed to each other client ( subscriber ) on demand.

- [Go Chat](#go-chat)
  - [How-to](#how-to)
  - [Configure](#configure)
  - [Usage](#usage)
    - [Binaries](#binaries)
    - [Docker](#docker)
    - [Build from source](#build-from-source)
  - [Support 💛](#support-)
  - [Contributing 🤝](#contributing-)
  - [Issues](#issues)

## How-to

All messages are encrypted using GPG and the public key of the recipient. The encrypted message is then sent to the server and routed to the recipient. The recipient then decrypts the message using their private key.

```text
             ┌────────┐                  ┌────────┐
Message ─────► Client │                  │ Client ├──Decrypt──►
             └───┬────┘                  └────▲───┘
                 │                            │
                 │           ┌────────┐       │
                 └─Encrypt───► Server ├───────┘
                             └────────┘
```

## Configure

You need to have some environment variables set up. You can do this by creating a `.env` file in the root of the project. See the example:

```env
PORT=
HOST=
PUBLIC_KEY_FILE=
PRIVATE_KEY_FILE=
PASSPHRASE=
```

| Variable         | Description                         | Default     |
| ---------------- | ----------------------------------- | ----------- |
| START_MODE       | The mode to start the application   | `server`    |
| PORT             | The port the server will listen on. | `12345`     |
| HOST             | The host the server will listen on. | `localhost` |
| PUBLIC_KEY_FILE  | The path to the public key file.    |             |
| PRIVATE_KEY_FILE | The path to the private key file.   |             |
| PASSPHRASE       | The passphrase for the private key. |             |

To separate the server from the client, you can use the `START_MODE` variable. This can be set to either `server` or `client`. If it is set to `server`, the server will start. If it is set to `client`, the client will start.

> **Note**: You should escape any special characters for the `PASSPHRASE`.

## Usage

### Binaries

Download the latest release from the [releases page](https://github.com/hatamiarash7/Go-Chat/releases/latest) based on your operating system and architecture.

**Server:**

```bash
START_MODE=server ./go-chat-linux-amd64
```

**Client:**

Configure your `.env` file and then run the following command:

```bash
START_MODE=client ./go-chat-linux-amd64
```

> **Note**: For MacOS, you should allow the application to run. You can do this by going to `System Preferences > Security & Privacy > General` and then click `Open Anyway`.

### Docker

I think we don't need a Docker image for the **Client**. But if you want to run the **Server** in a Docker container, you can use the published image.

```bash
docker run -it hatamiarash7/go-chat-server
```

Use `PORT` and `HOST` environment variables to configure the server.

```bash
docker run -it -e PORT=1234 -e HOST=0.0.0.0 hatamiarash7/go-chat-server
```

> **Note**: Default `PORT` is `12345` and default `HOST` is `0.0.0.0` for Docker.

### Build from source

To build the application from source, run the `build` target:

```bash
make build
```

First, you should build you need to start the server. You can do this by running the following command:

```bash
make server
```

Then, you can run any number of clients. You can do this by running the following command:

```bash
make client
```

Note that you should set required environment variables before running the client.

---

## Support 💛

[![Donate with Bitcoin](https://en.cryptobadges.io/badge/micro/bc1qmmh6vt366yzjt3grjxjjqynrrxs3frun8gnxrz)](https://en.cryptobadges.io/donate/bc1qmmh6vt366yzjt3grjxjjqynrrxs3frun8gnxrz) [![Donate with Ethereum](https://en.cryptobadges.io/badge/micro/0x0831bD72Ea8904B38Be9D6185Da2f930d6078094)](https://en.cryptobadges.io/donate/0x0831bD72Ea8904B38Be9D6185Da2f930d6078094)

[![ko-fi](https://www.ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/D1D1WGU9)

<div><a href="https://payping.ir/@hatamiarash7"><img src="https://cdn.payping.ir/statics/Payping-logo/Trust/blue.svg" height="128" width="128"></a></div>

## Contributing 🤝

Don't be shy and reach out to us if you want to contribute 😉

1. Fork it!
2. Create your feature branch: `git checkout -b my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin my-new-feature`
5. Submit a pull request

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
