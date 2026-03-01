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
		AvailableSupply   float64 `json:"availableSupply,omitzero"`
		Color             string  `json:"color,omitempty,omitzero"`
		ContractAddress   string  `json:"contractAddress,omitempty,omitzero"`
		ContractAddresses []struct {
			Blockchain      string `json:"blockchain,omitzero"`
			ContractAddress string `json:"contractAddress,omitzero"`
		} `json:"contractAddresses"`
		Decimals              int      `json:"decimals,omitempty,omitzero"`
		Explorers             []string `json:"explorers"`
		FullyDilutedValuation float64  `json:"fullyDilutedValuation,omitzero"`
		Icon                  string   `json:"icon,omitzero"`
		ID                    string   `json:"id,omitzero"`
		IsStable              bool     `json:"isStable,omitempty,omitzero"`
		MarketCap             float64  `json:"marketCap,omitzero"`
		Name                  string   `json:"name,omitzero"`
		Price                 float64  `json:"price,omitzero"`
		PriceBtc              float64  `json:"priceBtc,omitzero"`
		PriceChange1D         float64  `json:"priceChange1d"`
		PriceChange1H         float64  `json:"priceChange1h"`
		PriceChange1W         float64  `json:"priceChange1w"`
		Rank                  int      `json:"rank,omitzero"`
		RedditURL             string   `json:"redditUrl,omitempty,omitzero"`
		Slug                  string   `json:"slug,omitzero"`
		Symbol                string   `json:"symbol,omitzero"`
		TotalSupply           float64  `json:"totalSupply,omitzero"`
		TwitterURL            string   `json:"twitterUrl,omitempty,omitzero"`
		Volume                float64  `json:"volume"`
		WebsiteURL            string   `json:"websiteUrl,omitempty,omitzero"`
	} `json:"result"`
}
