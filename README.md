## Crypto Analyzer

1. Get API key from [coinstats](https://openapi.coinstats.app/)
2. Set the API key in the environment variable
```sh
export COINSTATS_API_KEY=<YOUR_API_KEY>
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
Title: KuCoin Teams Up with Zypto to Enable Direct Crypto Payments for Everyday Purchases
Description: KuCoin integrated its payment service with Zypto, allowing direct crypto spending in daily life. Users can now pay bills, top-up cards, and shop with over 50 cryptocurrencies in the app.
Source: Cointurk News EN
Link: https://en.coin-turk.com/kucoin-teams-up-with-zypto-to-enable-direct-crypto-payments-for-everyday-purchases/?utm_medium=referral&utm_source=coinstats
Affected coins:  [kucoin-shares]
Title: Massive Fund Inflows Highlight Growing Interest in Crypto ETFs
Description: On February 25, 2026, spot cryptocurrency exchange-traded funds (ETFs) in the United States experienced tremendous capital inflows, signaling a growing interest from institutional investors. A staggering $697.79 million found its way into these funds, reaffirming the sector&#8217;s evolution into a critical component of the financial landscape.
Source: BH NEWS
Link: https://en.bitcoinhaber.net/massive-fund-inflows-highlight-growing-interest-in-crypto-etfs?utm_medium=referral&utm_source=coinstats
Affected coins:  [vaulta]
Title: 🇨🇳 China just showcased AI humanoid Robots performing Kung Fu moves during a visit by German Chancellor WE ARE SOO COOKED
Source: Vivek⚡️ Twitter
Link: https://x.com/Vivek4real_/status/2027009708008345980?utm_medium=referral&utm_source=coinstats
Title: Gold Surges as Middle East Tensions Drive Safe-Haven Demand
Description: Rising geopolitical tensions in the Middle East are shaping how investors size risk, with safe-haven assets drawing attention as equities and crypto markets recalibrate. Fresh indicators show hedging behavior taking hold: oil flows from Iran are rising, while gold demand in key markets is climbing as traders seek ballast against potential disruption and macro volatility. [...]
Source: Crypto Breaking News
Link: https://www.cryptobreaking.com/gold-surges-as-middle-east/?utm_medium=referral&utm_source=coinstats
Title: Starknet’s Revolutionary strkBTC Launch: Transforming Bitcoin Privacy and DeFi Utility on Ethereum Layer 2
Description: BitcoinWorld
Source: Bitcoin World
Link: https://bitcoinworld.co.in/starknet-strkbtc-bitcoin-privacy-defi/?utm_medium=referral&utm_source=coinstats
Affected coins:  [bitcoin ethereum starknet]
Title: Best Crypto Presale 2026: DeepSnitch AI Surpasses $1.7M, Prints 175% and Pulls Fresh Capital as BMIC and Trd Network Give Back Gains
Description: Kraken has rolled out a new product called Flexline, introducing fixed-rate crypto-backed loans for users of its Kraken Pro platform. This development forces investors to ask: What is the best crypto presale to buy now to utilize this newfound liquidity?  While some are looking at speculative projects like BMIC or the pumping Trd Network, the
Source: CaptainAltcoin
Link: https://captainaltcoin.com/best-crypto-presale-2026-deepsnitch-ai-surpasses-1-7m-prints-175-and-pulls-fresh-capital-as-bmic-and-trd-network-give-back-gains/?utm_medium=referral&utm_source=coinstats
Title: Avalanche Price Prediction 2026-2030: Can AVAX Realistically Surge to $100?
Description: BitcoinWorld
Source: Bitcoin World
Link: https://bitcoinworld.co.in/avalanche-avax-price-prediction-2030-6/?utm_medium=referral&utm_source=coinstats
Affected coins:  [avalanche-2 3z2tRjNuQjoq6UDcw4zyEPD1Eb5KXMPYb4GWFzVT1DPg_solana surge-3]
Title: Miljarden PUMP tokens in beweging na distributie aankondiging
Description: Een opvallende transactie met PUMP tokens zorgt voor speculatie in de markt. Ongeveer 11,2 miljard PUMP, met een geschatte waarde van $21,22 miljoen, is een uur geleden overgemaakt naar crypto exchange Kraken. De tokens komen van een wallet die gelinkt is aan de Pump core treasury en vertegenwoordigen circa 1,12%...
Source: Blockchain Stories
Link: https://www.blockchainstories.com/2026/02/26/miljarden-pump-tokens-in-beweging-na-distributie-aankondiging/?utm_medium=referral&utm_source=coinstats
Affected coins:  [coredaoorg core-2 core]
Title: Morning Minute: Circle Earnings Highlight a Massive Day for Crypto
Description: It was a green day for stablecoin enthusiasts and crypto majors, and positive earnings from NVDA after hours signal the rally may have legs.
Source: Decrypt
Link: https://decrypt.co/359208/morning-minute-circle-earnings-highlight-a-massive-day-for-crypto?utm_medium=referral&utm_source=coinstats
Title: Dogecoin Price Eyes Breakout Amid Falling Wedge and Accumulation
Description: Dogecoin trades around $0.099, breaking a falling wedge and entering a third accumulation phase, signaling potential bullish continuation.
Source: Coinpaper
Link: https://coinpaper.com/14955/dogecoin-price-eyes-breakout-amid-falling-wedge-and-accumulation?utm_medium=referral&utm_source=coinstats
Affected coins:  [dogecoin 0x4206931337dc273a630d328da6441786bfad668f_ethereum]
Title: Hyperliquid Price Prediction: A Sober Analysis of HYPE’s Path to 2030
Description: BitcoinWorld
Source: Bitcoin World
Link: https://bitcoinworld.co.in/hyperliquid-hype-price-prediction-2030-6/?utm_medium=referral&utm_source=coinstats
Affected coins:  [hyperliquid vaulta]
Title: Playnance Pays $2M, Revenue Tops $5.3M Ahead of G-Token
Description: Playnance reports $2M in payouts and $5.3M revenue as it prepares to launch its G-Token across live on-chain platforms.
Source: Coinpaper
Link: https://coinpaper.com/14954/playnance-pays-2-m-revenue-tops-5-3-m-ahead-of-g-token?utm_medium=referral&utm_source=coinstats
Title: KuCoin Pay ve Zypto Ortaklığıyla Kripto Paralar Günlük Harcamalara Taşınıyor
Description: KuCoin Pay ile Zypto arasında entegrasyon gerçekleştirildi ve yeni ödeme imkânları sunuldu. Kullanıcılar, kripto bakiyeleriyle günlük alışveriş ve ödeme işlemleri yapabiliyor.
Source: Cointurk News TR
Link: https://coin-turk.com/kucoin-pay-ve-zypto-ortakligiyla-kripto-paralar-gunluk-harcamalara-tasiniyor?utm_source%3Drss%26%23038%3Butm_medium%3Drss%26%23038%3Butm_campaign%3Dkucoin-pay-ve-zypto-ortakligiyla-kripto-paralar-gunluk-harcamalara-tasiniyor&utm_medium=referral&utm_source=coinstats
Affected coins:  [kucoin-shares]
Title: MEXC, sınırlı bir süre için ürünlerinde kazanç artışı sağlıyor ve USDT esnek faiz oranını  ’ye kadar yükseltiyor
Description: Dünyanın en hızlı büyüyen dijital varlık borsası ve gerçek sıfır komisyonlu işlemlerin öncüsü MEXC , MEXC Earn tekliflerinde sınırlı süreli bir yükseltme duyurdu ve USDT Esnek Tasarruf hesaplarındaki Yıllık Faiz Oranını (APR)  &#8217;ye kadar çıkardı. Bu iyileştirme, yatırımcıların piyasa oynaklığıyla mücadele ederken istikrarlı ancak rekabetçi kazanç çözümlerine yönelik artan talebi karşılamakta ve MEXC&#8217;in çeşitli yatırımcı [&#8230;]
Source: Bitcoin Sistemi Turkish
Link: https://www.bitcoinsistemi.com/mexc-sinirli-bir-sure-icin-urunlerinde-kazanc-artisi-sagliyor-ve-usdt-esnek-faiz-oranini-%20ye-kadar-yukseltiyor/?utm_medium=referral&utm_source=coinstats
Affected coins:  [tether 10_hydration]
Title: Sygnum, 100 Milyar Dolarlık Kripto Hazinesi Pazarını Hedefleyen Yönetim Servisini Başlattı
Description: Sygnum, kripto varlık yönetiminde kurumsal düzeyde yeni bir hizmet başlattı.
Source: Cointurk News TR
Link: https://coin-turk.com/sygnum-100-milyar-dolarlik-kripto-hazinesi-pazarini-hedefleyen-yonetim-servisini-baslatti?utm_source%3Drss%26%23038%3Butm_medium%3Drss%26%23038%3Butm_campaign%3Dsygnum-100-milyar-dolarlik-kripto-hazinesi-pazarini-hedefleyen-yonetim-servisini-baslatti&utm_medium=referral&utm_source=coinstats
Title: Gold Price Forecast: Debasement Trade Signals Impending Surge to New Highs – TD Securities Analysis
Description: BitcoinWorld
Source: Bitcoin World
Link: https://bitcoinworld.co.in/gold-debasement-trade-new-highs/?utm_medium=referral&utm_source=coinstats
Affected coins:  [3z2tRjNuQjoq6UDcw4zyEPD1Eb5KXMPYb4GWFzVT1DPg_solana surge-3 surge-2]
Title: 21Shares inscrit STRC NA et propose un rendement adossé au bitcoin en Europe
Description: Les investisseurs européens disposent désormais d'un nouveau moyen d'accéder au rendement d'entreprise soutenu par Bitcoin. 21Shares a lancé un produit négocié en bourse (ETP) lié aux actions privilégiées émises par Strategy, la société dirigée par Michael Saylor et largement reconnue comme le plus grand détenteur public de Bitcoin.
Source: Cointribune FR
Link: https://www.cointribune.com/21shares-inscrit-strc-na-et-propose-un-rendement-adosse-au-bitcoin-en-europe/?utm_medium=referral&utm_source=coinstats
Affected coins:  [bitcoin 0xb0f70c0bd6fd87dbeb7c10dc692a2a6106817072_ethereum 0xb0f70c0bd6fd87dbeb7c10dc692a2a6106817072_megaeth]
Title: XRP Bulls Eye $1.65 as Golden Cross Sparks Reversal Hype
Description: XRP is forming a Golden Cross on the Ichimoku, signaling a potential bullish reversal toward $1.65.
Source: Coinpaper
Link: https://coinpaper.com/14956/xrp-bulls-eye-1-65-as-golden-cross-sparks-reversal-hype?utm_medium=referral&utm_source=coinstats
Affected coins:  [ripple cross-2 J1e5SjPnzrJAPhR2T1qgfyWyYxWRcsWVMzBxWScwNmd6_solana]
Title: Bitcoin’s ’10 AM dumps’ stop as Jane Street gets sued: ‘That’s all it took!’
Description: Is Bitcoin rallying on strength, or simply breathing without algorithmic pressure?
Source: AMBCrypto
Link: https://ambcrypto.com/bitcoins-10-am-dumps-stop-as-jane-street-gets-sued-thats-all-it-took/?utm_medium=referral&utm_source=coinstats
Affected coins:  [bitcoin 0xb0f70c0bd6fd87dbeb7c10dc692a2a6106817072_ethereum 0xb0f70c0bd6fd87dbeb7c10dc692a2a6106817072_megaeth]
Title: Here’s Why Stable (STABLE) Price Is Up Today
Description: Stable (STABLE) is ripping higher today. The token is up about 15% and trades near $0.03420, making it the top gainer in the market at the time of writing.&nbsp; Volume is the bigger headline, it’s up roughly 530%, which is exactly the kind of spike that turns a normal move into a scramble. This rally
Source: CaptainAltcoin
Link: https://captainaltcoin.com/heres-why-stable-stable-price-is-up-today/?utm_medium=referral&utm_source=coinstats
</NEWS>
```
