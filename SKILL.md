---
name: crypto-analyzer
description: Fetches real-time cryptocurrency market data (market cap, Fear and Greed Index, top coin prices, and latest industry news). Use when users ask for a market overview, crypto sentiment, BTC/ETH prices, or recent crypto news.
compatibility: Requires internet access, the `crypto-analyzer` binary on PATH, and COINSTATS_API_KEY plus API_KEY_COINMARKETCAP environment variables.
metadata:
  homepage: https://github.com/Alexandersfg4/crypto-analyzer
  nanobot:
    emoji: "📈"
    requires:
      bins:
        - go
        - crypto-analyzer
    install:
      - id: download
        kind: shell
        command: git clone git@github.com:Alexandersfg4/crypto-analyzer.git && cd crypto-analyzer
        label: Download repository and go to the directory
      - id: build
        kind: shell
        command: make build && make install
        label: Build and install crypto-analyzer
---

# Crypto Analyzer

A Go-based CLI tool that aggregates Coinstats, CoinMarketCap and [Defillama](https://defillama.com) data to provide a comprehensive snapshot of the cryptocurrency market in four tagged sections.

## When to use (trigger phrases)

Use this skill when the user asks for:
- “How is the crypto market doing?”
- “Get a market overview.”
- “run crypto-analyzer”
- “What is the current Fear and Greed index?”
- “Show me the latest crypto news.”
- “Check Bitcoin/Ethereum prices.” -you should run with flags`crypto-analyzer --tokens=BTC,ETH`
- “How is my DEX protocol doing?(AAVE, UNI)” - you should run with flags `crypto-analyzer --protocols=AAVE,UNI`
- “How is my portfolio doing?(tokens: BTC, ETH, used DEX protocols: AAVE, UNI)” - you should run with flags `crypto-analyzer --tokens=BTC,ETH --protocols=AAVE,UNI`

## Quick start

Use this workflow to produce a single combined report of market stats, sentiment, prices, and news:
1. Ensure the `crypto-analyzer` binary is installed and available on your PATH.
2. Export `COINSTATS_API_KEY` and `API_KEY_COINMARKETCAP` in the current shell session.
3. Run `crypto-analyzer` and summarize the four output sections for the user.
4. Output contains TOKENS from TOP 200 coins by market cap and all available protocols with TVL > 0.

Optional Arguments(show info only for selected protocols/tokens):
- `--protocols` — comma-separated protocol symbols to filter (e.g., Aave, Drift, Kamino).
- `--tokens` — comma-separated token symbols to filter (e.g., BTC, ETH, SOL).

Example:
```sh
export COINSTATS_API_KEY=... API_KEY_COINMARKETCAP=...
crypto-analyzer --tokens=BTC,ETH,SUI --protocols=AAVE,Drift,Kamino
```

## Requirements & Setup
- Coinstats API key exported as `COINSTATS_API_KEY`.
- CoinMarketCap API key exported as `API_KEY_COINMARKETCAP`.
- Network access to the required APIs.
- Binary available on PATH (often `/usr/local/bin/` after install).

## Output data points
The command returns four distinct sections wrapped in XML-like tags for easy parsing. If any section is missing, report what you do have and note the missing portion.

1. Market Capitalization
Includes total market cap, 24h trading volume, and Bitcoin dominance percentage.
  - **Key metric**: Use the 24-hour change to describe the market trend (up/down/flat). If not present, say it is unavailable.
2. Fear and Greed Index
Provides market sentiment (0-100).
  - **Classifications**: Extreme Fear, Fear, Neutral, Greed, Extreme Greed.
  - Includes historical values (Yesterday and Last Week) to show sentiment trends.
3. Coins
Detailed stats for major assets (BTC, ETH, etc.):
  - Price: Current USD value.
  - Change: Percentage movement over 24 hours and 7 days.
4. News
A list of the latest headlines including:
  - Title & Description: Summary of the event.
  - Source & Link: Original URL for deep diving.
  - Affected Coins: Identifies which specific tokens are relevant to the news item.
5. Protocols data:
  - TVL  
  - Price changes

## Troubleshooting
If the command fails, ensure the binary is on your PATH (often `/usr/local/bin/`), network access is available, and `COINSTATS_API_KEY` and `API_KEY_COINMARKETCAP` are valid and exported in the current shell session.
