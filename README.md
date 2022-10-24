# Go Chat

It's a simple chat system using a **PUB/SUB** single channel and **GPG encryption**. Any message sent from a client ( publisher ) is routed to each other client ( subscriber ) on demand.

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

**Note:** You should escape any special characters for the `PASSPHRASE`.

## Usage

### Binaries

Coming soon

### Docker

Coming soon

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
