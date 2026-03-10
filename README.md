## Crypto Analyzer
Fetch real-time cryptocurrency market data, including global capitalization, Fear and Greed Index, token prices you request, and the latest industry news — absolutely **FREE**.

1. Get API key from [coinstats](https://openapi.coinstats.app/) and [coinmarketcap](https://pro.coinmarketcap.com/account)
2. Set the API keys in your environment
```sh
export COINSTATS_API_KEY=<YOUR_API_KEY>
export API_KEY_COINMARKETCAP=<YOUR_API_KEY>
```

3. Build the binary
```sh
make build
```

4. Install the binary (installs to `/usr/local/bin/`)
```sh
make install
```

5. Get crypto data
```sh
crypto-analyzer --protocols=AAVE,DRIFT --tokens=BTC,ETH,SUI
```

Result example:
```sh
<Market Capitalization>
Total market capitalization of all cryptocurrencies : 2415020674552$
Total 24-hour trading volume across all cryptocurrencies: 63418354819$
Bitcoin's percentage share of the total cryptocurrency market capitalization: 55.890000%
24-hour change in total market capitalization: -0.930000%
24-hour change in total trading volume: -29.480000%
24-hour change in Bitcoin dominance: 0.050000%
</Market Capitalization>

<Fear and Greed Index now>
Value: 19
Classification: Extreme fear
Updated at: 2026-03-08 11:38:10.035 +0000 UTC
Fear and Greed Index yesterday
Value: 20
Classification: Fear
Fear and Greed Index last week
Value: 18
Classification: Extreme fear
</Fear and Greed Index now>

<TOKENS>
Name: Bitcoin
Symbol: BTC
<Quotes>
<USD>
Price: 67538.555973
Volume for 24h: 25940025103.004490
Market Cap: 1350737350189.560059
Price changed 1 hour: -0.642788%
Price changed 24 hours: -0.830310%
Price changed 7 days: 1.514453%
Price changed 90 days: -26.578206%
</USD>
</Quotes>
Name: Ethereum
Symbol: ETH
<Quotes>
<USD>
Price: 1952.336917
Volume for 24h: 12106748090.971827
Market Cap: 235631551065.220001
Price changed 1 hour: -1.074365%
Price changed 24 hours: -1.762906%
Price changed 7 days: -1.761769%
Price changed 90 days: -37.985086%
</USD>
</Quotes>
</TOKENS>

<NEWS>
Title: 🚨 JPMORGAN CEO JAMIE DIMON ADMITS HE WAS WRONG ABOUT BITCOIN AND CRYPTO He now says “It’s real” and “It will be used by all of us.” From one of Wall Street’s biggest skeptics to openly embracing it! 🚀
Source: That Martini Guy Twitter
Link: https://x.com/MartiniGuyYT/status/2030604098173870570?utm_medium=referral&utm_source=coinstats
Title: Florida Senate Approves Groundbreaking Stablecoin Licensing Framework
Description: Florida’s senate passed a bill requiring stablecoin issuers to obtain a state license. The legislation aims to complement federal policies and strengthen consumer protections.
Source: Cointurk News EN
Link: https://en.coin-turk.com/florida-senate-approves-groundbreaking-stablecoin-licensing-framework/?utm_medium=referral&utm_source=coinstats
</NEWS>

<PROTOCOLS>
Name: Aave V3
Symbol: AAVE
Description: Earn interest, borrow assets, and build applications
Category: Lending
TVL: 26172718879.414017$
Price changed 1 hour: 0.389966%
Price changed 24 hours: 0.091805%
Price changed 7 days: -0.114796%
</PROTOCOLS>
```
