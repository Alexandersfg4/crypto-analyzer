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

5. Get crypto data by passing coins to be analyzed
```sh
crypto-analyzer --coins=BTC,ETH
```
