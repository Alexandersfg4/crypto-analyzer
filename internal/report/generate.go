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
	timeoutFetch   = time.Minute * 2
	timeoutAnalyze = time.Minute * 3
	limitNews      = 20
	limitCoins     = 100
)

func (r *Report) Generate(ctx context.Context, cfg models.Config) (models.Report, error) {
	var (
		fearAndGreedData models.FearAndGreed
		marketCapData    models.MarketCap
		protocolsData    models.GetProtocolsResponse
		newsMap          = make(map[string]models.News)
		listingsCh       = make(chan []models.ListingsLatestData, 3)
		listingsWg       sync.WaitGroup
		newsMU           sync.Mutex
	)

	listingsWg.Add(3)

	go func() {
		listingsWg.Wait()
		close(listingsCh)
	}()

	fetchCtx, fetchCancel := context.WithTimeout(ctx, timeoutFetch)
	defer fetchCancel()

	g, fetchCtx := errgroup.WithContext(fetchCtx)

	g.Go(func() error {
		news, err := r.coinstatsSrv.GetNewsByType(fetchCtx, models.NewsTypeLatest, limitNews)
		if err != nil {
			return fmt.Errorf("error fetching latest news: %w", err)
		}
		newsMU.Lock()
		for _, n := range news {
			newsMap[n.Title] = n
		}
		newsMU.Unlock()
		return nil
	})

	g.Go(func() error {
		news, err := r.coinstatsSrv.GetNewsByType(fetchCtx, models.NewsTypeTrending, limitNews)
		if err != nil {
			return fmt.Errorf("error fetching trending news: %w", err)
		}
		newsMU.Lock()
		for _, n := range news {
			newsMap[n.Title] = n
		}
		newsMU.Unlock()
		return nil
	})

	g.Go(func() error {
		data, err := r.coinstatsSrv.GetFearAndGreed(fetchCtx)
		if err != nil {
			return fmt.Errorf("error getting fear and greed: %w", err)
		}
		fearAndGreedData = data
		return nil
	})

	g.Go(func() error {
		data, err := r.coinstatsSrv.GetMarketCap(fetchCtx)
		if err != nil {
			return fmt.Errorf("error getting market cap: %w", err)
		}
		marketCapData = data
		return nil
	})

	g.Go(func() error {
		data, err := r.defillamaSrv.GetProtocols(fetchCtx)
		if err != nil {
			return fmt.Errorf("error getting protocols: %w", err)
		}
		protocolsData = data
		return nil
	})

	g.Go(func() error {
		defer listingsWg.Done()

		data, err := r.coinmarketcapSrv.GetListingsLatest(fetchCtx, 1, limitCoins)
		if err != nil {
			return fmt.Errorf("error getting listings page 1: %w", err)
		}
		listingsCh <- data.Data
		return nil
	})

	g.Go(func() error {
		defer listingsWg.Done()

		data, err := r.coinmarketcapSrv.GetListingsLatest(fetchCtx, 101, limitCoins)
		if err != nil {
			return fmt.Errorf("error getting listings page 2: %w", err)
		}
		listingsCh <- data.Data
		return nil
	})

	g.Go(func() error {
		defer listingsWg.Done()

		data, err := r.coinmarketcapSrv.GetListingsLatest(fetchCtx, 201, limitCoins)
		if err != nil {
			return fmt.Errorf("error getting listings page 2: %w", err)
		}
		listingsCh <- data.Data
		return nil
	})

	if err := g.Wait(); err != nil {
		return models.Report{}, fmt.Errorf("error getting listings: %w", err)
	}

	listingsLatestData := make([]models.ListingsLatestData, 0, 300)
	for data := range listingsCh {
		listingsLatestData = append(listingsLatestData, data...)
	}

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

	analyzeCtx, analyzeCancel := context.WithTimeout(ctx, timeoutAnalyze)
	defer analyzeCancel()

	analyzeResult, err := r.openRouterSrv.Analyze(analyzeCtx, cfg.OpenrouterModel,
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
