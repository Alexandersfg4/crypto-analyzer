package report

import (
	"github.com/Alexandersfg4/crypto-analyzer/internal/providers/coinmarketcap"
	"github.com/Alexandersfg4/crypto-analyzer/internal/providers/coinstats"
	"github.com/Alexandersfg4/crypto-analyzer/internal/providers/defillama"
	"github.com/Alexandersfg4/crypto-analyzer/internal/providers/openrouter"
)

type Report struct {
	coinmarketcapSrv *coinmarketcap.Service
	coinstatsSrv     *coinstats.Service
	defillamaSrv     *defillama.Service
	openRouterSrv    *openrouter.Service
}

func New(
	coinmarketcapSrv *coinmarketcap.Service,
	coinstatsSrv *coinstats.Service,
	defillamaSrv *defillama.Service,
	openRouterSrv *openrouter.Service,
) *Report {
	return &Report{
		coinmarketcapSrv: coinmarketcapSrv,
		coinstatsSrv:     coinstatsSrv,
		defillamaSrv:     defillamaSrv,
	}
}
