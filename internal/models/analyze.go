package models

type AnalyzeRequest struct {
	Protfolio      Protfolio
	MarketCap      MarketCap
	FeatAndGreed   FearAndGreed
	ListingsLatest []ListingsLatestData
	News           GetNewsResponse
	Protocols      GetProtocolsResponse
}

type Protfolio struct {
	Tokens    []string
	Protocols []string
}
