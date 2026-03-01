package models

import "time"

type ListingsLatestResponse struct {
	Data   []ListingsLatestData `json:"data"`
	Status struct {
		CreditCount  int       `json:"credit_count,omitzero"`
		Elapsed      int       `json:"elapsed,omitzero"`
		ErrorCode    int       `json:"error_code"`
		ErrorMessage any       `json:"error_message"`
		Notice       any       `json:"notice"`
		Timestamp    time.Time `json:"timestamp,omitzero"`
		TotalCount   int       `json:"total_count,omitzero"`
	} `json:"status"`
}

type ListingsLatestData struct {
	CirculatingSupply float64   `json:"circulating_supply,omitzero"`
	CmcRank           int       `json:"cmc_rank,omitzero"`
	DateAdded         time.Time `json:"date_added,omitzero"`
	ID                int       `json:"id,omitzero"`
	InfiniteSupply    bool      `json:"infinite_supply"`
	LastUpdated       time.Time `json:"last_updated,omitzero"`
	MaxSupply         *int      `json:"max_supply"`
	MintedMarketCap   float64   `json:"minted_market_cap,omitzero"`
	Name              string    `json:"name,omitzero"`
	NumMarketPairs    int       `json:"num_market_pairs,omitzero"`
	Platform          any       `json:"platform"`
	Quote             struct {
		Usd struct {
			FullyDilutedMarketCap float64   `json:"fully_diluted_market_cap,omitzero"`
			LastUpdated           time.Time `json:"last_updated,omitzero"`
			MarketCap             float64   `json:"market_cap,omitzero"`
			MarketCapDominance    float64   `json:"market_cap_dominance,omitzero"`
			PercentChange1h       float64   `json:"percent_change_1h,omitzero"`
			PercentChange24h      float64   `json:"percent_change_24h,omitzero"`
			PercentChange30d      float64   `json:"percent_change_30d,omitzero"`
			PercentChange60d      float64   `json:"percent_change_60d,omitzero"`
			PercentChange7d       float64   `json:"percent_change_7d,omitzero"`
			PercentChange90d      float64   `json:"percent_change_90d,omitzero"`
			Price                 float64   `json:"price,omitzero"`
			Tvl                   any       `json:"tvl"`
			Volume24h             float64   `json:"volume_24h,omitzero"`
			VolumeChange24h       float64   `json:"volume_change_24h,omitzero"`
		} `json:"USD"`
	} `json:"quote"`
	SelfReportedCirculatingSupply any      `json:"self_reported_circulating_supply"`
	SelfReportedMarketCap         any      `json:"self_reported_market_cap"`
	Slug                          string   `json:"slug,omitzero"`
	Symbol                        string   `json:"symbol,omitzero"`
	Tags                          []string `json:"tags"`
	TotalSupply                   float64  `json:"total_supply,omitzero"`
	TvlRatio                      any      `json:"tvl_ratio"`
}
