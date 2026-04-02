package report

import "github.com/Alexandersfg4/crypto-analyzer/internal/models"

type Data struct {
	News           models.GetNewsResponse
	FeatAndGreed   models.FearAndGreed
	MarketCap      models.MarketCap
	ListingsLatest []models.ListingsLatestData
	Protocols      models.GetProtocolsResponse
}
