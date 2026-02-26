package models

type Coins struct {
	Meta struct {
		HasNextPage     bool `json:"hasNextPage,omitzero"`
		HasPreviousPage bool `json:"hasPreviousPage"`
		ItemCount       int  `json:"itemCount,omitzero"`
		Limit           int  `json:"limit,omitzero"`
		Page            int  `json:"page,omitzero"`
		PageCount       int  `json:"pageCount,omitzero"`
	} `json:"meta"`
	Result []struct {
		AvailableSupply       int      `json:"availableSupply,omitzero"`
		AvgChange             float64  `json:"avgChange,omitzero"`
		ContractAddress       any      `json:"contractAddress"`
		ContractAddresses     []any    `json:"contractAddresses"`
		Decimals              int      `json:"decimals,omitzero"`
		Explorers             []string `json:"explorers"`
		FullyDilutedValuation int      `json:"fullyDilutedValuation,omitzero"`
		Icon                  string   `json:"icon,omitzero"`
		ID                    string   `json:"id,omitzero"`
		LiquidityScore        float64  `json:"liquidityScore,omitzero"`
		MarketCap             int      `json:"marketCap,omitzero"`
		MarketCapScore        float64  `json:"marketCapScore,omitzero"`
		Name                  string   `json:"name,omitzero"`
		Price                 float64  `json:"price,omitzero"`
		PriceBtc              float64  `json:"priceBtc,omitzero"`
		PriceChange1D         float64  `json:"priceChange1d,omitzero"`
		PriceChange1H         float64  `json:"priceChange1h,omitzero"`
		PriceChange1W         float64  `json:"priceChange1w,omitzero"`
		Rank                  int      `json:"rank,omitzero"`
		RedditURL             string   `json:"redditUrl,omitzero"`
		RiskScore             float64  `json:"riskScore,omitzero"`
		Symbol                string   `json:"symbol,omitzero"`
		TotalSupply           int      `json:"totalSupply,omitzero"`
		TwitterURL            string   `json:"twitterUrl,omitzero"`
		VolatilityScore       float64  `json:"volatilityScore,omitzero"`
		Volume                int      `json:"volume,omitzero"`
		WebsiteURL            string   `json:"websiteUrl,omitzero"`
	} `json:"result"`
}
