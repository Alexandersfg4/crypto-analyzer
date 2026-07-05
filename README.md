# Crypto Analyzer

Service for generating a cryptocurrency market report and sending it to Telegram. The binary is designed for one-shot execution, which makes it suitable for Linux `cron`.

![Report example](assets/report_example.png)

## Prerequisites

1. Get API keys from [CoinStats](https://openapi.coinstats.app/), [CoinMarketCap](https://pro.coinmarketcap.com/account), and [OpenRouter](https://openrouter.ai/settings/keys).
2. Create a Telegram bot and get the API token from [@BotFather](https://t.me/BotFather).
3. Get the Telegram chat ID where the report should be delivered.

## Environment Variables

```sh
export COINSTATS_API_KEY=<YOUR_API_KEY>
export API_KEY_COINMARKETCAP=<YOUR_API_KEY>
export OPENROUTER_API_KEY=<YOUR_API_KEY>
export TELEGRAM_API_TOKEN=<YOUR_TELEGRAM_BOT_TOKEN>
```

## Build

```sh
make build
```

## Run Once

```sh
./bin/crypto-analyzer \
  -telegram-chat-id 123456789 \
  -openrouter-model google/gemini-2.5-pro \
  -tokens BTC,ETH,SUI \
  -protocols AAVE,UNISWAP,DRIFT
```

Arguments:

- `-telegram-chat-id` — Telegram chat ID for delivery
- `-openrouter-model` — OpenRouter model used for AI summary
- `-tokens` — comma-separated list of portfolio tokens
- `-protocols` — comma-separated list of DeFi protocols

## Cron Example

```cron
0 9 * * * /usr/local/bin/crypto-analyzer \
  -telegram-chat-id 123456789 \
  -openrouter-model google/gemini-2.5-pro \
  -tokens BTC,ETH,SUI \
  -protocols AAVE,UNISWAP,DRIFT >> /var/log/crypto-analyzer.log 2>&1
```

## What the Report Includes

- global market capitalization
- Fear and Greed Index
- selected portfolio tokens
- top gainers and losers
- selected DeFi protocols TVL
- AI-generated market summary
