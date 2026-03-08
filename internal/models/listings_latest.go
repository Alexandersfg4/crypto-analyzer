package models

import "time"

type ListingsLatestResponse struct {
	Data   []ListingsLatestData `json:"data"`
	Status struct {
		CreditCount  int       `json:"credit_count,omitzero"`
		Elapsed      int       `json:"elapsed,omitzero"`
		ErrorCode    string    `json:"error_code,omitzero"`
		ErrorMessage string    `json:"error_message"`
		Timestamp    time.Time `json:"timestamp,omitzero"`
	} `json:"status"`
}

type ListingsLatestData struct {
	CirculatingSupply float64   `json:"circulating_supply,omitzero"`
	CmcRank           int       `json:"cmc_rank,omitzero"`
	DateAdded         time.Time `json:"date_added,omitzero"`
	ID                int       `json:"id,omitzero"`
	InfiniteSupply    bool      `json:"infinite_supply"`
	LastUpdated       time.Time `json:"last_updated,omitzero"`
	MaxSupply         *float64  `json:"max_supply"`
	Name              string    `json:"name,omitzero"`
	NumMarketPairs    int       `json:"num_market_pairs,omitzero"`
	Platform          *struct {
		ID           int    `json:"id,omitzero"`
		Name         string `json:"name,omitzero"`
		Slug         string `json:"slug,omitzero"`
		Symbol       string `json:"symbol,omitzero"`
		TokenAddress string `json:"token_address,omitzero"`
	} `json:"platform"`
	Quote []struct {
		CexVolume24h          float64   `json:"cex_volume_24h,omitzero"`
		DexVolume24h          float64   `json:"dex_volume_24h,omitzero"`
		FullyDilutedMarketCap float64   `json:"fully_diluted_market_cap,omitzero"`
		ID                    int       `json:"id,omitzero"`
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
		Symbol                string    `json:"symbol,omitzero"`
		Tvl                   any       `json:"tvl"`
		Volume24h             float64   `json:"volume_24h,omitzero"`
		VolumeChange24h       float64   `json:"volume_change_24h,omitzero"`
	} `json:"quote"`
	SelfReportedCirculatingSupply any      `json:"self_reported_circulating_supply"`
	SelfReportedMarketCap         any      `json:"self_reported_market_cap"`
	Slug                          string   `json:"slug,omitzero"`
	Symbol                        string   `json:"symbol,omitzero"`
	Tags                          []string `json:"tags"`
	TotalSupply                   float64  `json:"total_supply,omitzero"`
	TvlRatio                      any      `json:"tvl_ratio"`
}
