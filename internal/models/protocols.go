package models

type GetProtocolsResponse []Data

type DataChainTvls struct {
	Fantom    float64 `json:"Fantom"`
	Chiliz    float64 `json:"Chiliz"`
	Ronin     float64 `json:"Ronin"`
	Sonic     float64 `json:"Sonic"`
	Hedera    float64 `json:"Hedera"`
	Sui       float64 `json:"Sui"`
	OpBnb     float64 `json:"Op_Bnb"`
	Starknet  float64 `json:"Starknet"`
	Scroll    float64 `json:"Scroll"`
	Base      float64 `json:"Base"`
	Algorand  float64 `json:"Algorand"`
	Manta     float64 `json:"Manta"`
	ZkSyncEra float64 `json:"zkSync Era"`
	Binance   float64 `json:"Binance"`
	Optimism  float64 `json:"Optimism"`
	Plasma    float64 `json:"Plasma"`
	Stellar   float64 `json:"Stellar"`
	Litecoin  float64 `json:"Litecoin"`
	Near      float64 `json:"Near"`
	Avalanche float64 `json:"Avalanche"`
	Ethereum  float64 `json:"Ethereum"`
	Polkadot  float64 `json:"Polkadot"`
	Arbitrum  float64 `json:"Arbitrum"`
	Polygon   float64 `json:"Polygon"`
	Celo      float64 `json:"Celo"`
	Tron      float64 `json:"Tron"`
	Solana    float64 `json:"Solana"`
	TON       float64 `json:"TON"`
	Ripple    float64 `json:"Ripple"`
	Bitcoin   float64 `json:"Bitcoin"`
	Doge      float64 `json:"Doge"`
	Aptos     float64 `json:"Aptos"`
}

type DataTokenBreakdowns any

type Data struct {
	Id                   string              `json:"id"`
	Name                 string              `json:"name"`
	Address              any                 `json:"address"`
	Symbol               string              `json:"symbol"`
	Url                  string              `json:"url"`
	Description          string              `json:"description"`
	Chain                string              `json:"chain"`
	Logo                 string              `json:"logo"`
	Audits               string              `json:"audits"`
	GeckoId              any                 `json:"gecko_id"`
	CmcId                any                 `json:"cmcId"`
	Category             string              `json:"category"`
	Chains               []string            `json:"chains"`
	Module               string              `json:"module"`
	Twitter              string              `json:"twitter"`
	ListedAt             int                 `json:"listedAt"`
	Methodology          string              `json:"methodology"`
	MisrepresentedTokens bool                `json:"misrepresentedTokens"`
	TvlCodePath          string              `json:"tvlCodePath"`
	Slug                 string              `json:"slug"`
	Tvl                  float64             `json:"tvl"`
	ChainTvls            DataChainTvls       `json:"chainTvls"`
	Change1h             float64             `json:"change_1h"`
	Change1d             float64             `json:"change_1d"`
	Change7d             float64             `json:"change_7d"`
	TokenBreakdowns      DataTokenBreakdowns `json:"tokenBreakdowns"`
	Mcap                 any                 `json:"mcap"`
}
