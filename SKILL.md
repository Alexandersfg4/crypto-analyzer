---
name: crypto-analyzer
description: Fetch real-time cryptocurrency market data, including global capitalization, Fear and Greed Index, top coin prices, and latest industry news.
homepage: https://github.com/Alexandersfg4/crypto-analyzer
metadata: 
{
  "nanobot": {
    "emoji": "📈",
    "requires": { "bins": ["crypto-analyzer"] },
    "install": [
      {
        "id": "download",
        "kind": "shell",
        "command": "git clone git@github.com:Alexandersfg4/crypto-analyzer.git && cd crypto-analyzer",
        "label": "Download repository and go to the directory"
      },
      {
        "id": "build",
        "kind": "shell",
        "command": "make build && make install",
        "label": "Build and install crypto-analyzer"
      }
    ]
  }
}
---

# Crypto Analyzer

A Go-based CLI tool that aggregates data from the Coinstats API to provide a comprehensive snapshot of the cryptocurrency market.

## When to use (trigger phrases)

Use this skill when the user asks for:
- “How is the crypto market doing?”
- “What is the current Fear and Greed index?”
- “Show me the latest crypto news.”
- “Check Bitcoin/Ethereum prices.”
- “Get a market overview.”
- “run crypto-analyzer”

## Quick start

The tool provides a combined report of market stats, sentiment, prices, and news in one go:

```bash
crypto-analyzer
```

## Requirements & Setup
You must have a Coinstats API Key set in your environment:
```bash
export COINSTATS_API_KEY=<YOUR_API_KEY>
```

## Output Data Points
The command returns four distinct sections wrapped in XML-like tags for easy parsing:
1. Market Capitalization
Includes total market cap, 24h trading volume, and Bitcoin dominance percentage.
  - **Key Metric**: Look at "24-hour change" to determine if the market is trending up or down.
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
  
## Troubleshooting
If the command fails, ensure the binary is in your path (usually /usr/local/bin/) and the COINSTATS_API_KEY is valid and exported in the current shell session.
