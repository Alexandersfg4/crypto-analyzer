## Crypto Analyzer
Fetch real-time cryptocurrency market data, including global capitalization, Fear and Greed Index, top coin prices, and latest industry news absolutely **FREE**.

1. Get API key from [coinstats](https://openapi.coinstats.app/) and [coinmarketcap](https://pro.coinmarketcap.com/account)
2. Set the API key in the environment variable
```sh
export COINSTATS_API_KEY=<YOUR_API_KEY>
export API_KEY_COINMARKETCAP=<YOUR_API_KEY>
```

3. Build the binary
```sh
make build
```

4. Install the binary (placed it in /usr/local/bin/)
```sh
make install
```

5. Get crypto data
```sh
crypto-analyzer
```

Result example:
```sh
<Market Capitalization>
Total market capitalization of all cryptocurrencies : 2459306521633$
Total 24-hour trading volume across all cryptocurrencies: 144257676647$
Bitcoin's percentage share of the total cryptocurrency market capitalization: 55.320000%
24-hour change in total market capitalization: 2.520000%
24-hour change in total trading volume: 30.820000%
24-hour change in Bitcoin dominance: 0.070000%
</Market Capitalization>

<Fear and Greed Index now>
Value: 16
Classification: Extreme fear
Updated at: 2026-02-26 13:38:10.026 +0000 UTC
Fear and Greed Index yesterday
Value: 11
Classification: Extreme fear
Fear and Greed Index last week
Value: 24
Classification: Fear
</Fear and Greed Index now>

<COINS>
Name: Bitcoin
Symbol: BTC
Price: 67990.419873$
Volume: 43832461385.942734$
Market Cap: 1359473272674.195557$
Price changed 24 hours: 2.910000%
Price changed 7 days: 2.750000%
Name: Ethereum
Symbol: ETH
Price: 2066.917932$
Volume: 34325928331.640942$
Market Cap: 249461058468.436615$
Price changed 24 hours: 5.320000%
Price changed 7 days: 7.020000%
...
</COINS>

<NEWS>
Title: 🇨🇳 China just showcased AI humanoid Robots performing Kung Fu moves during a visit by German Chancellor WE ARE SOO COOKED
Source: Vivek⚡️ Twitter
Link: https://x.com/Vivek4real_/status/2027009708008345980?utm_medium=referral&utm_source=coinstats
Title: Gold Surges as Middle East Tensions Drive Safe-Haven Demand
Description: Rising geopolitical tensions in the Middle East are shaping how investors size risk, with safe-haven assets drawing attention as equities and crypto markets recalibrate. Fresh indicators show hedging behavior taking hold: oil flows from Iran are rising, while gold demand in key markets is climbing as traders seek ballast against potential disruption and macro volatility. [...]
Source: Crypto Breaking News
Link: https://www.cryptobreaking.com/gold-surges-as-middle-east/?utm_medium=referral&utm_source=coinstats
Title: KuCoin Teams Up with Zypto to Enable Direct Crypto Payments for Everyday Purchases
Description: KuCoin integrated its payment service with Zypto, allowing direct crypto spending in daily life. Users can now pay bills, top-up cards, and shop with over 50 cryptocurrencies in the app.
Source: Cointurk News EN
Link: https://en.coin-turk.com/kucoin-teams-up-with-zypto-to-enable-direct-crypto-payments-for-everyday-purchases/?utm_medium=referral&utm_source=coinstats
Affected coins:  [kucoin-shares]
...
</NEWS>
```
