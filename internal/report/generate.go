package report

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Alexandersfg4/crypto-analyzer/internal/formatter"
	"github.com/Alexandersfg4/crypto-analyzer/internal/models"
	"golang.org/x/sync/errgroup"
)

const (
	timeoutWork = time.Minute * 3
	limitNews   = 20
	limitCoins  = 100
)

func (r *Report) Generate(ctx context.Context, cfg models.Config) (models.Report, error) {
	ctx, cancel := context.WithTimeout(ctx, timeoutWork)
	defer cancel()

	var (
		fearAndGreedData models.FearAndGreed
		marketCapData    models.MarketCap
		protocolsData    models.GetProtocolsResponse
		newsMap          = make(map[string]models.News)
		listingsPage1    []models.ListingsLatestData
		listingsPage2    []models.ListingsLatestData
		mu               sync.Mutex
	)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		news, err := r.coinstatsSrv.GetNewsByType(ctx, models.NewsTypeLatest, limitNews)
		if err != nil {
			return fmt.Errorf("error fetching latest news: %w", err)
		}
		mu.Lock()
		for _, n := range news {
			newsMap[n.Title] = n
		}
		mu.Unlock()
		return nil
	})

	g.Go(func() error {
		news, err := r.coinstatsSrv.GetNewsByType(ctx, models.NewsTypeTrending, limitNews)
		if err != nil {
			return fmt.Errorf("error fetching trending news: %w", err)
		}
		mu.Lock()
		for _, n := range news {
			newsMap[n.Title] = n
		}
		mu.Unlock()
		return nil
	})

	g.Go(func() error {
		data, err := r.coinstatsSrv.GetFearAndGreed(ctx)
		if err != nil {
			return fmt.Errorf("error getting fear and greed: %w", err)
		}
		mu.Lock()
		fearAndGreedData = data
		mu.Unlock()
		return nil
	})

	g.Go(func() error {
		data, err := r.coinstatsSrv.GetMarketCap(ctx)
		if err != nil {
			return fmt.Errorf("error getting market cap: %w", err)
		}
		mu.Lock()
		marketCapData = data
		mu.Unlock()
		return nil
	})

	g.Go(func() error {
		data, err := r.defillamaSrv.GetProtocols(ctx)
		if err != nil {
			return fmt.Errorf("error getting protocols: %w", err)
		}
		mu.Lock()
		protocolsData = data
		mu.Unlock()
		return nil
	})

	g.Go(func() error {
		data, err := r.coinmarketcapSrv.GetListingsLatest(ctx, 1, limitCoins)
		if err != nil {
			return fmt.Errorf("error getting listings page 1: %w", err)
		}
		mu.Lock()
		listingsPage1 = data.Data
		mu.Unlock()
		return nil
	})

	g.Go(func() error {
		data, err := r.coinmarketcapSrv.GetListingsLatest(ctx, 101, limitCoins)
		if err != nil {
			return fmt.Errorf("error getting listings page 2: %w", err)
		}
		mu.Lock()
		listingsPage2 = data.Data
		mu.Unlock()
		return nil
	})

	if err := g.Wait(); err != nil {
		return models.Report{}, err
	}

	listingsLatestData := append(listingsPage1, listingsPage2...)

	newsResult := make([]models.News, 0, len(newsMap))
	for _, v := range newsMap {
		newsResult = append(newsResult, v)
	}

	cap := &strings.Builder{}
	formatter.MarketCap(cap, marketCapData)
	formatter.FearAndGreed(cap, fearAndGreedData)

	tokensAll := &strings.Builder{}
	formatter.CoinsAll(tokensAll, listingsLatestData)

	tokensInPortfolio := &strings.Builder{}
	formatter.CoinsInPortfolio(tokensInPortfolio, listingsLatestData, cfg.Tokens)

	tokensGainersAndLosers := &strings.Builder{}
	formatter.CoinsGainersAndLosers(tokensGainersAndLosers, listingsLatestData)

	proto := &strings.Builder{}
	formatter.Protocols(proto, protocolsData, cfg.Protocols)

	news := &strings.Builder{}
	formatter.News(news, newsResult)

	analyzeResult, err := r.openRouterSrv.Analyze(ctx, cfg.OpenrouterModel,
		strings.Join([]string{
			cap.String(),
			fmt.Sprintf("tokens in portfolio: %s", strings.Join(cfg.Tokens, ", ")),
			proto.String(),
			news.String(),
			tokensAll.String(),
		}, "\n"))
	if err != nil {
		return models.Report{}, fmt.Errorf("error analyzing: %w", err)
	}

	return models.Report{
		MarketCap: cap.String(),
		Tokens: models.Tokens{
			InPortfolio:       tokensInPortfolio.String(),
			GainersAndLoosers: tokensGainersAndLosers.String(),
		},
		AISummary: analyzeResult,
	}, nil
}
