# teleGoRAT

**teleGoRAT** is a straightforward, cross-platform RAT written in Go. It leverages Telegram for communication.

## Features

- Cross-platform: Works on Windows, Linux, and MacOS.
- Telegram-based communication: Utilizes the Telegram API for commands and data exchange.
- Simple setup: Easy to configure and get started.
- Extensible: You can expand its functionality by adding your commands.

## Getting Started

### Prerequisites

- Go (Golang) installed on your system.
- A valid Telegram account.

### Installation

1. Clone this repository:

```shell
git clone https://github.com/efthimisgrk/teleGoRAT.git
```

2. Create a Telegram bot and get your API token. You can follow the official Telegram documentation to create a bot: [Creating a new bot](https://core.telegram.org/bots#botfather).

3. Add the API token to your environment variables with the name "BOT_TOKEN".

4. Build the teleGoRAT binary:

```shell
cd teleGoRAT
go build
```

5. Run teleGoRAT:

```shell
./teleGoRAT
```

6. Start communicating with your bot via Telegram to execute commands and manage the remote system.

## Usage
- Send `/help` to get a list of available commands.
- Have fun :)

## License
This project is licensed under the MIT License - see the LICENSE file for details.

## Disclaimer
**teleGoRAT** is intended for educational and research purposes only. The author is not responsible for any misuse or damage caused by this software.
