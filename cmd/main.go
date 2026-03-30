package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Alexandersfg4/crypto-analyzer/internal/coinmarketcap"
	"github.com/Alexandersfg4/crypto-analyzer/internal/coinstats"
	"github.com/Alexandersfg4/crypto-analyzer/internal/defillama"
	"github.com/Alexandersfg4/crypto-analyzer/internal/formatter"
	"github.com/Alexandersfg4/crypto-analyzer/internal/models"
)

const (
	envCoinstatsAPIKey     = "COINSTATS_API_KEY"
	envCoinmarketcapAPIKey = "API_KEY_COINMARKETCAP"
	limitNews              = 20
	limitCoins             = 100
	timeoutWork            = time.Second * 20
)

var (
	protocols = flag.String("protocols", "", "comma-separated list of protocols")
	tokens    = flag.String("tokens", "", "comma-separated list of tokens")
	store     = flag.String("store", "", "path to store output")
)

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

	ctx := context.Background()

	opts := requestOptions{
		tokens:    parseCSV(*tokens),
		protocols: parseCSV(*protocols),
	}

	gotData, err := getData(ctx, apiKeyCoinstats, apiKeyCoinmarketcap, opts)
	if err != nil {
		fmt.Println("getting data: ", err.Error())
		os.Exit(1)
	}

	output := &bytes.Buffer{}
	formatter.MarketCap(output, gotData.marketCap)
	formatter.FearAndGreed(output, gotData.fearAndGreed)
	formatter.Coins(output, gotData.listingsLatest, opts.tokens)
	formatter.News(output, gotData.news)
	formatter.Protocols(output, gotData.protocols, opts.protocols)

	if *store != "" {
		now := time.Now()
		filename := fmt.Sprintf("crypto-analyzer-%d-%d-%d.txt", now.Year(), int(now.Month()), now.Day())
		dirPath := filepath.Join(*store, filename)

		if err := os.MkdirAll(*store, 0755); err != nil {
			fmt.Println("creating directory: ", err.Error())
			os.Exit(1)
		}

		if err := os.WriteFile(dirPath, output.Bytes(), 0644); err != nil {
			fmt.Println("writing file: ", err.Error())
			os.Exit(1)
		}
		fmt.Println("Output saved to:", dirPath)
	} else {
		io.Copy(os.Stdout, output)
	}
}

type data struct {
	news           models.GetNewsResponse
	fearAndGreed   models.FearAndGreed
	marketCap      models.MarketCap
	listingsLatest []models.ListingsLatestData
	protocols      models.GetProtocolsResponse
}

type requestOptions struct {
	tokens    []string
	protocols []string
}

func (o requestOptions) hasTokens() bool {
	return len(o.tokens) > 0
}

func (o requestOptions) hasProtocols() bool {
	return len(o.protocols) > 0
}

func parseCSV(value string) []string {
	if value == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		result = append(result, p)
	}
	return result
}

func getData(ctx context.Context, coinstatsApiKey, coinmarketcapApiKey string, opts requestOptions) (data, error) {
	var (
		result                     data
		wg                         sync.WaitGroup
		mu                         sync.Mutex
		newsMap                    = make(map[string]models.News)
		errCh                      = make(chan error, 8)
		listingsLatestDataCh       = make(chan []models.ListingsLatestData)
		secondListingsLatestDoneCh = make(chan struct{})
	)

	ctx, cancel := context.WithTimeout(ctx, timeoutWork)
	defer cancel()

	srvCoinstats := coinstats.NewService(coinstatsApiKey)
	srvCoinmarketcap := coinmarketcap.NewService(coinmarketcapApiKey)
	srvDefillama := defillama.NewService()

	jobs := []func(){
		func() {
			gotLatestNews, err := srvCoinstats.GetNewsByType(ctx, models.NewsTypeLatest, limitNews)
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
			gotTrendingNews, err := srvCoinstats.GetNewsByType(ctx, models.NewsTypeTrending, limitNews)
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
			gotFearAndGreed, err := srvCoinstats.GetFearAndGreed(ctx)
			if err != nil {
				errCh <- fmt.Errorf("error getting fear and greed: %w", err)
				return
			}
			mu.Lock()
			result.fearAndGreed = gotFearAndGreed
			mu.Unlock()
		},
		func() {
			gotMarketCap, err := srvCoinstats.GetMarketCap(ctx)
			if err != nil {
				errCh <- fmt.Errorf("error getting market cap: %w", err)
				return
			}
			mu.Lock()
			result.marketCap = gotMarketCap
			mu.Unlock()
		},
		func() {
			gotProtocols, err := srvDefillama.GetProtocols(ctx)
			if err != nil {
				errCh <- fmt.Errorf("error getting protocols: %w", err)
				return
			}
			mu.Lock()
			result.protocols = gotProtocols
			mu.Unlock()
		},
		func() {
			for data := range listingsLatestDataCh {
				mu.Lock()
				result.listingsLatest = append(result.listingsLatest, data...)
				mu.Unlock()
			}
		},
		func() {
			defer func() {
				<-secondListingsLatestDoneCh
				close(listingsLatestDataCh)
			}()

			gotCoins, err := srvCoinmarketcap.GetListingsLatest(ctx, 1, limitCoins)
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

			gotCoins, err := srvCoinmarketcap.GetListingsLatest(ctx, 101, limitCoins)
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
			return result, e
		}
	}

	newsResult := make([]models.News, 0, len(newsMap))
	for _, v := range newsMap {
		newsResult = append(newsResult, v)
	}
	result.news = newsResult

	return result, nil
}
