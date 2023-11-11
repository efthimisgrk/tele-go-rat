# tele-go-rat

**tele-go-rat** is a straightforward, cross-platform RAT written in Go. It leverages Telegram for communication.

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
git clone https://github.com/efthimisgrk/tele-go-rat.git
```

2. Create a Telegram bot and get your API token. You can follow the official Telegram documentation to create a bot: [Creating a new bot](https://core.telegram.org/bots#botfather).

3. Get your chatId by sending a message to yourself on the Telegram web application: `https://web.telegram.org/z/#<CHAT_ID>`

4. Add the API token and the chatId to your environment variables with the name "BOT_TOKEN" and "CHAT_ID" respectively.

5. Build the teleGoRAT binary:

```shell
cd tele-go-rat
go build
```

6. Run tele-go-rat:

```shell
./tele-go-rat
```

7. Start communicating with your bot via Telegram to execute commands and manage the remote system.

## Usage
- Send `/help` to get a list of available commands.
- Have fun :)

## License
This project is licensed under the MIT License - see the LICENSE file for details.

## Disclaimer
**tele-go-rat** is intended for educational and research purposes only. The author is not responsible for any misuse or damage caused by this software.
