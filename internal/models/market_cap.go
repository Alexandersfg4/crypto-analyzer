package models

type MarketCap struct {
	BtcDominance       float64 `json:"btcDominance,omitzero"`
	BtcDominanceChange float64 `json:"btcDominanceChange,omitzero"`
	MarketCap          int     `json:"marketCap,omitzero"`
	MarketCapChange    float64 `json:"marketCapChange,omitzero"`
	Volume             int     `json:"volume,omitzero"`
	VolumeChange       float64 `json:"volumeChange,omitzero"`
}
