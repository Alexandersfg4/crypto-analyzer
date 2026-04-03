# Crypto Analyzer

Telegram bot for fetching real-time cryptocurrency market data, including global capitalization, Fear and Greed Index, token prices, DeFi protocols TVL, and the latest industry news.

<img src="./assets/report_example.png" alt="drawing" width="200"/>

## Prerequisites

1. Get API keys from [CoinStats](https://openapi.coinstats.app/) and [CoinMarketCap](https://pro.coinmarketcap.com/account)
2. Create a Telegram bot and get the API token from [@BotFather](https://t.me/BotFather)
3. Get your Telegram User ID (you can use [@userinfobot](https://t.me/userinfobot) to find it)

## Environment Variables

```sh
export COINSTATS_API_KEY=<YOUR_API_KEY>
export API_KEY_COINMARKETCAP=<YOUR_API_KEY>
export TELEGRAM_API_TOKEN=<YOUR_TELEGRAM_BOT_TOKEN>
export TELEGRAM_USER_ID=<YOUR_TELEGRAM_USER_ID>
```

## Build & Run

Build the binary:
```sh
make build
```

Install the binary (installs to `/usr/local/bin/`):
```sh
make install
```

Run the bot:
```sh
crypto-analyzer
```

## Bot Commands

- `/report` — Generate and send crypto market report
- `/tokens ETH,BTC,SUIA` — Set observed tokens (comma-separated)
- `/protocols AAVE,UNISWAP,DRIFT` — Set observed DeFi protocols (comma-separated)
- `/config` — Show current configuration
- `/cron 12:15` — Set daily report generation time (UTC timezone)

## Configuration

The bot stores your configuration in `crypto-analyzer.json` file in the working directory. This file contains:

- **tokens**: List of tokens to monitor
- **protocols**: List of DeFi protocols to track
- **cron_next_execution_time**: Scheduled report time

The configuration file is automatically created when you first set tokens or protocols using the bot commands.
