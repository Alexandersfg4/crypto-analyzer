package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Alexandersfg4/crypto-analyzer/internal/coinstats"
	"github.com/Alexandersfg4/crypto-analyzer/internal/models"
)

var coins = flag.Bool("coins", false, "Observed tokens")

const (
	envAPIKey    = "COINSTATS_API_KEY"
	baseURL      = "https://openapiv1.coinstats.app"
	limitNews    = 20
	limitWorkers = 5
	timeoutWork  = time.Second * 10
)

func main() {
	apiKey, ok := os.LookupEnv(envAPIKey)
	if !ok {
		fmt.Println("API key not found")
		os.Exit(1)
	}

	ctx := context.Background()

	gotData, err := getData(ctx, apiKey)
	if err != nil {
		fmt.Println("getting data: ", err.Error())
		os.Exit(1)
	}

	showMarketCap(gotData.marketCap)
	showFearAndGreed(gotData.fearAndGreed)
	showNews(gotData.news)
}

type data struct {
	news         models.GetNewsResponse
	fearAndGreed models.FearAndGreed
	marketCap    models.MarketCap
}

func getData(ctx context.Context, apiKey string) (data, error) {
	var err error
	var result data
	var wg sync.WaitGroup
	var fearAndGreed models.FearAndGreed
	var marketCap models.MarketCap

	ctx, c := context.WithTimeout(ctx, timeoutWork)
	defer c()

	srv := coinstats.NewService(baseURL, apiKey)

	sem := make(chan struct{}, limitWorkers)
	newsCh := make(chan models.News)
	errorsCh := make(chan error)
	fearAndGreedCh := make(chan models.FearAndGreed)
	marketCapCh := make(chan models.MarketCap)
	mapNews := make(map[string]models.News)

	jobs := []func(){
		func() {
			defer func() {
				sem <- struct{}{}
			}()
			defer wg.Done()

			gotLatestNews, err := srv.GetNewsByType(ctx, models.NewsTypeLatest, limitNews)
			if err != nil {
				errorsCh <- fmt.Errorf("Error fetching latest news: %w", err)
				return
			}
			for _, n := range gotLatestNews {
				newsCh <- n
			}
		},
		func() {
			defer func() {
				sem <- struct{}{}
			}()
			defer wg.Done()

			gotTrendingNews, err := srv.GetNewsByType(ctx, models.NewsTypeTrending, limitNews)
			if err != nil {
				errorsCh <- fmt.Errorf("Error fetching tranding news: %w", err)
				return
			}
			for _, n := range gotTrendingNews {
				newsCh <- n
			}
		},
		func() {
			defer func() {
				sem <- struct{}{}
			}()
			defer wg.Done()

			gotFeatAndGreed, err := srv.GetFearAndGreed(ctx)
			if err != nil {
				errorsCh <- fmt.Errorf("Error getting fear and greed: %w", err)
				return
			}

			fearAndGreedCh <- gotFeatAndGreed
		},
		func() {
			defer func() {
				sem <- struct{}{}
			}()
			defer wg.Done()

			gotFeatAndGreed, err := srv.GetFearAndGreed(ctx)
			if err != nil {
				errorsCh <- fmt.Errorf("Error getting fear and greed: %w", err)
				return
			}

			fearAndGreedCh <- gotFeatAndGreed
		},
		func() {
			defer func() {
				sem <- struct{}{}
			}()
			defer wg.Done()

			gotMarketCap, err := srv.GetMarketCap(ctx)
			if err != nil {
				errorsCh <- fmt.Errorf("Error getting market cap: %w", err)
				return
			}

			marketCapCh <- gotMarketCap
		},
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case err = <-errorsCh:
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
			case err = <-errorsCh:
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
			case err = <-errorsCh:
				return
			case m, ok := <-marketCapCh:
				if !ok {
					return
				}

				marketCap = m
			}
		}
	}()

	for _, j := range jobs {
		wg.Add(1)

		sem <- struct{}{}

		go j()
	}

	wg.Wait()

	close(errorsCh)
	close(newsCh)
	close(fearAndGreedCh)

	if err != nil {
		return result, err
	}

	newsResult := make([]models.News, 0, len(mapNews))

	for _, v := range mapNews {
		newsResult = append(newsResult, v)
	}

	return data{
		news:         newsResult,
		fearAndGreed: fearAndGreed,
		marketCap:    marketCap,
	}, nil
}

func showNews(gotNews models.GetNewsResponse) {
	for _, news := range gotNews {
		fmt.Printf("Title: %s\n", news.Title)
		if news.Description != "" {
			fmt.Printf("Description: %s\n", news.Description)
		}
		fmt.Printf("Source: %s\n", news.Source)
		coins := make([]string, 0, len(news.Coins))
		for _, coin := range news.Coins {
			coins = append(coins, coin.CoinIDKeyWords)
		}
		if len(coins) > 0 {
			fmt.Println("Affected coins: ", coins)
		}
		fmt.Println()
	}
}

func showFearAndGreed(gotFearAndGreed models.FearAndGreed) {
	fmt.Println("Fear and Greed Index now")
	fmt.Printf("Value: %d\n", gotFearAndGreed.Now.Value)
	fmt.Printf("Classification: %s\n", gotFearAndGreed.Now.ValueClassification)
	fmt.Printf("Updated at: %s\n", gotFearAndGreed.Now.UpdateTime)
	fmt.Println("Fear and Greed Index yesterday")
	fmt.Printf("Value: %d\n", gotFearAndGreed.Yesterday.Value)
	fmt.Printf("Classification: %s\n", gotFearAndGreed.Yesterday.ValueClassification)
	fmt.Println("Fear and Greed Index last week")
	fmt.Printf("Value: %d\n", gotFearAndGreed.LastWeek.Value)
	fmt.Printf("Classification: %s\n", gotFearAndGreed.LastWeek.ValueClassification)
	fmt.Println()
}

func showMarketCap(gotMarketCap models.MarketCap) {
	fmt.Println("Market Cap")
	fmt.Printf("Total market capitalization of all cryptocurrencies : %d$\n", gotMarketCap.MarketCap)
	fmt.Printf("Total 24-hour trading volume across all cryptocurrencies: %d$\n", gotMarketCap.Volume)
	fmt.Printf("Bitcoin's percentage share of the total cryptocurrency market capitalization: %f%%\n", gotMarketCap.BtcDominance)
	fmt.Printf("24-hour change in total market capitalization: %f%%\n", gotMarketCap.MarketCapChange)
	fmt.Printf("24-hour change in total trading volume: %f%%\n", gotMarketCap.VolumeChange)
	fmt.Printf("24-hour change in Bitcoin dominance: %f%%\n", gotMarketCap.BtcDominanceChange)
	fmt.Println()
}
