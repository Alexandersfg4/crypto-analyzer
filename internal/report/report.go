package report

import (
	"context"

	"github.com/Alexandersfg4/crypto-analyzer/internal/models"
)

type CoinstatsProvider interface {
	GetMarketCap(ctx context.Context) (models.MarketCap, error)
	GetFearAndGreed(ctx context.Context) (models.FearAndGreed, error)
	GetNewsByType(ctx context.Context, newsType models.NewsType, limit int) (models.GetNewsResponse, error)
	GetCoins(ctx context.Context, limit int) (models.Coins, error)
}

type CoinmarketcapProvider interface {
	GetListingsLatest(ctx context.Context, start int, limit int) (models.ListingsLatestResponse, error)
}

type DefillamaProvider interface {
	GetProtocols(ctx context.Context) (models.GetProtocolsResponse, error)
}

type OpenRouterProvider interface {
	Analyze(ctx context.Context, model, data string) (string, error)
}

func New(
	coinmarketcapSrv CoinmarketcapProvider,
	coinstatsSrv CoinstatsProvider,
	defillamaSrv DefillamaProvider,
	openRouterSrv OpenRouterProvider,
) *Report {
	return &Report{
		coinmarketcapSrv: coinmarketcapSrv,
		coinstatsSrv:     coinstatsSrv,
		defillamaSrv:     defillamaSrv,
		openRouterSrv:    openRouterSrv,
	}
}

type Report struct {
	coinmarketcapSrv CoinmarketcapProvider
	coinstatsSrv     CoinstatsProvider
	defillamaSrv     DefillamaProvider
	openRouterSrv    OpenRouterProvider
}
