# Digiflazz Bot

An unofficial Telegram bot for [Digiflazz](https://digiflazz.com), built with [Go](https://go.dev). Very lightweight — uses less than 10 MB of memory.

## Features

- Check account balance
- Make prepaid transactions
- Browse available products
- More features coming soon

## Installation

### Download Binary

Grab the latest release from the [Releases page](https://github.com/fidrasofyan/digiflazz-bot/releases).

### Build from Source

```sh
git clone https://github.com/fidrasofyan/digiflazz-bot.git
cd digiflazz-bot
make build
```

Make sure you have [Go](https://go.dev) installed. The compiled binary will be available in the `bin/` directory.

## Usage

1. Create a new Telegram bot with [@BotFather](https://t.me/botfather).
2. Configure your bot by editing the `.env` file.
3. Start the bot: `./bin/digiflazz-bot start`
4. Configure your domain for webhook.

You can also use Docker. See the [Dockerfile](https://github.com/fidrasofyan/digiflazz-bot/blob/main/Dockerfile) and [compose.example.yaml](https://github.com/fidrasofyan/digiflazz-bot/blob/main/compose.example.yaml) for details.

## Screenshots

<img width="300" height="378" alt="product_list" src="https://github.com/user-attachments/assets/b3f75cf7-695b-424d-9e22-18ac9cad23cb" />
<img width="300" height="378" alt="transaction" src="https://github.com/user-attachments/assets/65c3eb71-b54f-416f-9f3a-a32fbd40eb5a" />

## Notes

- This project is unofficial and not affiliated with Digiflazz.
- Always keep your `.env` credentials secure.

## License

MIT License
