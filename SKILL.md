---
name: crypto-analyzer
description: Analyze cryptocurrency market data using the CoinStats API. Use when you need to fetch and summarize market cap, fear-and-greed index, news, or top coins, or when wiring CLI usage, configuration, and output formatting for this project.
---

# Crypto Analyzer Skill

## Goals
- Fetch market cap, fear-and-greed index, news, and optionally coins.
- Keep output concise and readable.
- Respect environment-based configuration.

## Required inputs
- Read `COINSTATS_API_KEY` from environment.

## Workflow
1. Parse flags early.
2. Validate `COINSTATS_API_KEY` and fail fast if missing.
3. Create a context with a bounded timeout.
4. Fetch data concurrently with a small worker limit.
5. Deduplicate news by title.
6. Print market cap, fear-and-greed, and news; print coins only when `-coins` is set.

## Implementation notes
- Keep `Service` usage isolated in a single orchestration function.
- Avoid duplicate API calls for the same dataset.
- Ensure goroutines cannot leak on early error or context cancellation.
- Use buffered error channel and aggregate the first error.

## Output expectations
- Use stable section headings: `Market Cap`, `Fear and Greed Index`, `NEWS`, `COINS`.
- Print timestamps and numeric values with clear units.
