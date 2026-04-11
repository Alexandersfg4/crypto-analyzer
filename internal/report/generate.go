package report

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Alexandersfg4/crypto-analyzer/internal/formatter"
	"github.com/Alexandersfg4/crypto-analyzer/internal/models"
)

const (
	timeoutWork = time.Second * 30
	limitNews   = 20
	limitCoins  = 100
)

func (r *Report) Generate(ctx context.Context, cfg models.Config) (models.Report, error) {
	data, err := r.getData(ctx, cfg)
	if err != nil {
		return models.Report{}, fmt.Errorf("error getting data: %w", err)
	}

	return data, nil
}

func (r *Report) getData(ctx context.Context, cfg models.Config) (models.Report, error) {
	var (
		fearAndGreedData           models.FearAndGreed
		marketCapData              models.MarketCap
		listingsLatestData         = make([]models.ListingsLatestData, 0, limitCoins*2)
		protocolsData              models.GetProtocolsResponse
		wg                         sync.WaitGroup
		mu                         sync.Mutex
		newsMap                    = make(map[string]models.News)
		errCh                      = make(chan error, 8)
		listingsLatestDataCh       = make(chan []models.ListingsLatestData)
		secondListingsLatestDoneCh = make(chan struct{})
	)

	ctx, cancel := context.WithTimeout(ctx, timeoutWork)
	defer cancel()

	jobs := []func(){
		func() {
			gotLatestNews, err := r.coinstatsSrv.GetNewsByType(ctx, models.NewsTypeLatest, limitNews)
			if err != nil {
				errCh <- fmt.Errorf("error fetching latest news: %w", err)
				return
			}
			mu.Lock()
			for _, n := range gotLatestNews {
				newsMap[n.Title] = n
			}
			mu.Unlock()
		},
		func() {
			gotTrendingNews, err := r.coinstatsSrv.GetNewsByType(ctx, models.NewsTypeTrending, limitNews)
			if err != nil {
				errCh <- fmt.Errorf("error fetching trending news: %w", err)
				return
			}
			mu.Lock()
			for _, n := range gotTrendingNews {
				newsMap[n.Title] = n
			}
			mu.Unlock()
		},
		func() {
			gotFearAndGreed, err := r.coinstatsSrv.GetFearAndGreed(ctx)
			if err != nil {
				errCh <- fmt.Errorf("error getting fear and greed: %w", err)
				return
			}
			mu.Lock()
			fearAndGreedData = gotFearAndGreed
			mu.Unlock()
		},
		func() {
			gotMarketCap, err := r.coinstatsSrv.GetMarketCap(ctx)
			if err != nil {
				errCh <- fmt.Errorf("error getting market cap: %w", err)
				return
			}
			mu.Lock()
			marketCapData = gotMarketCap
			mu.Unlock()
		},
		func() {
			gotProtocols, err := r.defillamaSrv.GetProtocols(ctx)
			if err != nil {
				errCh <- fmt.Errorf("error getting protocols: %w", err)
				return
			}
			mu.Lock()
			protocolsData = gotProtocols
			mu.Unlock()
		},
		func() {
			for data := range listingsLatestDataCh {
				mu.Lock()
				listingsLatestData = append(listingsLatestData, data...)
				mu.Unlock()
			}
		},
		func() {
			defer func() {
				<-secondListingsLatestDoneCh
				close(listingsLatestDataCh)
			}()

			gotCoins, err := r.coinmarketcapSrv.GetListingsLatest(ctx, 1, limitCoins)
			if err != nil {
				errCh <- fmt.Errorf("error listings latests: %w", err)
				return
			}
			listingsLatestDataCh <- gotCoins.Data
		},
		func() {
			defer func() {
				secondListingsLatestDoneCh <- struct{}{}
			}()

			gotCoins, err := r.coinmarketcapSrv.GetListingsLatest(ctx, 101, limitCoins)
			if err != nil {
				errCh <- fmt.Errorf("error listings latests: %w", err)
				return
			}
			listingsLatestDataCh <- gotCoins.Data
		},
	}

	for _, j := range jobs {
		wg.Add(1)
		go func(job func()) {
			defer wg.Done()
			job()
		}(j)
	}

	wg.Wait()
	close(errCh)

	for e := range errCh {
		if e != nil {
			return models.Report{}, e
		}
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

	analyzeResult, err := r.openRouterSrv.Analyze(ctx, cfg.OpenrouterModel,
		strings.Join([]string{
			cap.String(),
			tokensInPortfolio.String(),
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
