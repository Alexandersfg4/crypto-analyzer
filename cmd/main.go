package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Alexandersfg4/crypto-analyzer/internal/coinmarketcap"
	"github.com/Alexandersfg4/crypto-analyzer/internal/coinstats"
	"github.com/Alexandersfg4/crypto-analyzer/internal/defillama"
	"github.com/Alexandersfg4/crypto-analyzer/internal/models"
)

const (
	envCoinstatsAPIKey     = "COINSTATS_API_KEY"
	envCoinmarketcapAPIKey = "API_KEY_COINMARKETCAP"
	limitNews              = 20
	limitWorkers           = 6
	limitCoins             = 500
	timeoutWork            = time.Second * 10
)

var protocols = flag.String("protocols", "", "comma-separated list of protocols")

func main() {
	flag.Parse()

	apiKeyCoinstats, ok := os.LookupEnv(envCoinstatsAPIKey)
	if !ok {
		fmt.Println("env COINSTATS_API_KEY not found")
		os.Exit(1)
	}

	apiKeyCoinmarketcap, ok := os.LookupEnv(envCoinmarketcapAPIKey)
	if !ok {
		fmt.Println("env API_KEY_COINMARKETCAP not found")
		os.Exit(1)
	}

	flag.Parse()

	ctx := context.Background()

	gotData, err := getData(ctx, apiKeyCoinstats, apiKeyCoinmarketcap)
	if err != nil {
		fmt.Println("getting data: ", err.Error())
		os.Exit(1)
	}

	showMarketCap(gotData.marketCap)
	showFearAndGreed(gotData.fearAndGreed)
	showCoins(gotData.listingsLatest)
	showNews(gotData.news)
	showProtocols(gotData.protocols)
}

type data struct {
	news           models.GetNewsResponse
	fearAndGreed   models.FearAndGreed
	marketCap      models.MarketCap
	listingsLatest models.ListingsLatestResponse
	protocols      models.GetProtocolsResponse
}

func getData(ctx context.Context, coinstatsApiKey, coinmarketcapApiKey string) (data, error) {
	var (
		result          data
		wg              sync.WaitGroup
		fearAndGreed    models.FearAndGreed
		marketCap       models.MarketCap
		coins           models.ListingsLatestResponse
		protocolsResult models.GetProtocolsResponse
	)

	ctx, c := context.WithTimeout(ctx, timeoutWork)
	defer c()

	srvCoinstats := coinstats.NewService(coinstatsApiKey)
	srvCoinmarketcap := coinmarketcap.NewService(coinmarketcapApiKey)
	srvDefillama := defillama.NewService()

	sem := make(chan struct{}, limitWorkers)
	newsCh := make(chan models.News)
	errorsCh := make(chan error, limitWorkers)
	fearAndGreedCh := make(chan models.FearAndGreed)
	marketCapCh := make(chan models.MarketCap)
	coinsCh := make(chan models.ListingsLatestResponse)
	protocolsCh := make(chan models.GetProtocolsResponse)
	mapNews := make(map[string]models.News)

	jobs := []func(){
		func() {
			defer func() {
				<-sem
			}()
			defer wg.Done()

			gotLatestNews, err := srvCoinstats.GetNewsByType(ctx, models.NewsTypeLatest, limitNews)
			if err != nil {
				errorsCh <- fmt.Errorf("error fetching latest news: %w", err)
				return
			}
			for _, n := range gotLatestNews {
				newsCh <- n
			}
		},
		func() {
			defer func() {
				<-sem
			}()
			defer wg.Done()

			gotTrendingNews, err := srvCoinstats.GetNewsByType(ctx, models.NewsTypeTrending, limitNews)
			if err != nil {
				errorsCh <- fmt.Errorf("error fetching trending news: %w", err)
				return
			}
			for _, n := range gotTrendingNews {
				newsCh <- n
			}
		},
		func() {
			defer func() {
				<-sem
			}()
			defer wg.Done()

			gotFearAndGreed, err := srvCoinstats.GetFearAndGreed(ctx)
			if err != nil {
				errorsCh <- fmt.Errorf("error getting fear and greed: %w", err)
				return
			}

			fearAndGreedCh <- gotFearAndGreed
		},
		func() {
			defer func() {
				<-sem
			}()
			defer wg.Done()

			gotMarketCap, err := srvCoinstats.GetMarketCap(ctx)
			if err != nil {
				errorsCh <- fmt.Errorf("error getting market cap: %w", err)
				return
			}

			marketCapCh <- gotMarketCap
		},
		func() {
			defer func() {
				<-sem
			}()
			defer wg.Done()

			gotCoins, err := srvCoinmarketcap.GetListingsLatest(ctx, limitCoins)
			if err != nil {
				errorsCh <- fmt.Errorf("error listings latests: %w", err)
				return
			}

			coinsCh <- gotCoins
		},
	}

	if *protocols != "" {
		jobs = append(jobs, func() {
			defer func() {
				<-sem
			}()
			defer wg.Done()

			gotProtocols, err := srvDefillama.GetProtocols(ctx)
			if err != nil {
				errorsCh <- fmt.Errorf("error getting protocols: %w", err)
				return
			}

			protocolsCh <- gotProtocols
		})
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case news, ok := <-newsCh:
				if !ok {
					return
				}

				mapNews[news.Title] = news
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case f, ok := <-fearAndGreedCh:
				if !ok {
					return
				}

				fearAndGreed = f
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case m, ok := <-marketCapCh:
				if !ok {
					return
				}

				marketCap = m
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case c, ok := <-coinsCh:
				if !ok {
					return
				}

				coins = c
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case r, ok := <-protocolsCh:
				if !ok {
					return
				}

				protocolsResult = r
			}
		}
	}()

	for _, j := range jobs {
		wg.Add(1)

		sem <- struct{}{}

		go j()
	}

	wg.Wait()

	close(newsCh)
	close(fearAndGreedCh)
	close(marketCapCh)
	close(coinsCh)
	close(errorsCh)
	close(protocolsCh)

	for e := range errorsCh {
		if e != nil {
			return result, e
		}
	}

	newsResult := make([]models.News, 0, len(mapNews))

	for _, v := range mapNews {
		newsResult = append(newsResult, v)
	}

	return data{
		news:           newsResult,
		fearAndGreed:   fearAndGreed,
		marketCap:      marketCap,
		listingsLatest: coins,
		protocols:      protocolsResult,
	}, nil
}

func showNews(gotNews models.GetNewsResponse) {
	fmt.Println("<NEWS>")
	for _, news := range gotNews {
		fmt.Printf("Title: %s\n", news.Title)
		if news.Description != "" {
			fmt.Printf("Description: %s\n", news.Description)
		}
		fmt.Printf("Source: %s\n", news.Source)
		fmt.Printf("Link: %s\n", news.Link)
		coins := make([]string, 0, len(news.Coins))
		for _, coin := range news.Coins {
			coins = append(coins, coin.CoinIDKeyWords)
		}
		if len(coins) > 0 {
			fmt.Println("Affected coins: ", coins)
		}
	}

	fmt.Println("</NEWS>")
	fmt.Println()
}

func showFearAndGreed(gotFearAndGreed models.FearAndGreed) {
	fmt.Println("<Fear and Greed Index now>")
	fmt.Printf("Value: %d\n", gotFearAndGreed.Now.Value)
	fmt.Printf("Classification: %s\n", gotFearAndGreed.Now.ValueClassification)
	fmt.Printf("Updated at: %s\n", gotFearAndGreed.Now.UpdateTime)
	fmt.Println("Fear and Greed Index yesterday")
	fmt.Printf("Value: %d\n", gotFearAndGreed.Yesterday.Value)
	fmt.Printf("Classification: %s\n", gotFearAndGreed.Yesterday.ValueClassification)
	fmt.Println("Fear and Greed Index last week")
	fmt.Printf("Value: %d\n", gotFearAndGreed.LastWeek.Value)
	fmt.Printf("Classification: %s\n", gotFearAndGreed.LastWeek.ValueClassification)
	fmt.Println("</Fear and Greed Index now>")
	fmt.Println()
}

func showMarketCap(gotMarketCap models.MarketCap) {
	fmt.Println("<Market Capitalization>")
	fmt.Printf(
		"Total market capitalization of all cryptocurrencies : %d$\n",
		gotMarketCap.MarketCap,
	)
	fmt.Printf(
		"Total 24-hour trading volume across all cryptocurrencies: %d$\n",
		gotMarketCap.Volume,
	)
	fmt.Printf(
		"Bitcoin's percentage share of the total cryptocurrency market capitalization: %f%%\n",
		gotMarketCap.BtcDominance,
	)
	fmt.Printf(
		"24-hour change in total market capitalization: %f%%\n",
		gotMarketCap.MarketCapChange,
	)
	fmt.Printf("24-hour change in total trading volume: %f%%\n", gotMarketCap.VolumeChange)
	fmt.Printf("24-hour change in Bitcoin dominance: %f%%\n", gotMarketCap.BtcDominanceChange)
	fmt.Println("</Market Capitalization>")
	fmt.Println()
}

func showCoins(gotCoins models.ListingsLatestResponse) {
	fmt.Println("<COINS>")
	for _, c := range gotCoins.Data {
		fmt.Printf("Name: %s\n", c.Name)
		fmt.Printf("Symbol: %s\n", c.Symbol)
		fmt.Printf("Price: %f$\n", c.Quote.Usd.Price)
		fmt.Printf("Volume: %f$\n", c.Quote.Usd.Volume24h)
		fmt.Printf("Market Cap: %f$\n", c.Quote.Usd.MarketCap)
		fmt.Printf("Price changed 1 hour: %f%%\n", c.Quote.Usd.PercentChange1h)
		fmt.Printf("Price changed 24 hours: %f%%\n", c.Quote.Usd.PercentChange24h)
		fmt.Printf("Price changed 7 days: %f%%\n", c.Quote.Usd.PercentChange7d)
		fmt.Printf("Price changed 90 days: %f%%\n", c.Quote.Usd.PercentChange90d)
	}
	fmt.Println("</COINS>")
	fmt.Println()
}

func showProtocols(gotProtocols models.GetProtocolsResponse) {
	fmt.Println("<PROTOCOLS>")
	lfp := strings.Split(*protocols, ",")
	lowerLfp := make([]string, len(lfp))
	for i, l := range lfp {
		lowerLfp[i] = strings.ToLower(l)
	}

	fmt.Println("GOT AMOUNT ", len(gotProtocols))
	for _, p := range gotProtocols {
		for _, l := range lowerLfp {
			if strings.Contains(strings.ToLower(p.Name), l) {
				fmt.Printf("Name: %s\n", p.Name)
				fmt.Printf("Symbol: %s\n", p.Symbol)
				fmt.Printf("Description: %f$\n", p.Description)
				fmt.Printf("Category: %f$\n", p.Category)
				fmt.Printf("TVL: %f$\n", p.Tvl)
				fmt.Printf("Price changed 1 hour: %f%%\n", p.Change1h)
				fmt.Printf("Price changed 24 hours: %f%%\n", p.Change1d)
				fmt.Printf("Price changed 7 days: %f%%\n", p.Change7d)
			} else {
				fmt.Println("NOT FOUND ", strings.ToLower(p.Name), " AND ", l)
			}
		}
	}
	fmt.Println("</PROTOCOLS>")
	fmt.Println()
}
