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
	limitCoins             = 100
	timeoutWork            = time.Second * 10
)

var (
	protocols = flag.String("protocols", "", "comma-separated list of protocols")
	tokens    = flag.String("tokens", "", "comma-separated list of tokens")
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

	showMarketCap(gotData.marketCap)
	showFearAndGreed(gotData.fearAndGreed)
	showCoins(gotData.listingsLatest, opts.tokens)
	showNews(gotData.news)
	showProtocols(gotData.protocols, opts.protocols)
}

type data struct {
	news           models.GetNewsResponse
	fearAndGreed   models.FearAndGreed
	marketCap      models.MarketCap
	listingsLatest models.ListingsLatestResponse
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
		result  data
		wg      sync.WaitGroup
		mu      sync.Mutex
		newsMap = make(map[string]models.News)
		errCh   = make(chan error, 8)
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
	}

	if opts.hasProtocols() {
		jobs = append(jobs, func() {
			gotProtocols, err := srvDefillama.GetProtocols(ctx)
			if err != nil {
				errCh <- fmt.Errorf("error getting protocols: %w", err)
				return
			}
			mu.Lock()
			result.protocols = gotProtocols
			mu.Unlock()
		})
	}

	if opts.hasTokens() {
		jobs = append(jobs, func() {
			gotCoins, err := srvCoinmarketcap.GetListingsLatest(ctx, limitCoins)
			if err != nil {
				errCh <- fmt.Errorf("error listings latests: %w", err)
				return
			}
			mu.Lock()
			result.listingsLatest = gotCoins
			mu.Unlock()
		})
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

func showCoins(gotCoins models.ListingsLatestResponse, tokens []string) {
	fmt.Println("<TOKENS>")
	if len(tokens) == 0 {
		fmt.Println("</TOKENS>")
		fmt.Println()
		return
	}

	tokenSet := make(map[string]struct{}, len(tokens))
	for _, t := range tokens {
		tokenSet[strings.ToUpper(t)] = struct{}{}
	}

	for _, c := range gotCoins.Data {
		if _, ok := tokenSet[c.Symbol]; !ok {
			continue
		}

		fmt.Printf("Name: %s\n", c.Name)
		fmt.Printf("Symbol: %s\n", c.Symbol)
		fmt.Println("<Quotes>")
		for _, q := range c.Quote {
			fmt.Printf("<%s>\n", q.Symbol)
			fmt.Printf("Price: %f\n", q.Price)
			fmt.Printf("Volume for 24h: %f\n", q.Volume24h)
			fmt.Printf("Market Cap: %f\n", q.MarketCap)
			fmt.Printf("Price changed 1 hour: %f%%\n", q.PercentChange1h)
			fmt.Printf("Price changed 24 hours: %f%%\n", q.PercentChange24h)
			fmt.Printf("Price changed 7 days: %f%%\n", q.PercentChange7d)
			fmt.Printf("Price changed 90 days: %f%%\n", q.PercentChange90d)
			fmt.Printf("</%s>\n", q.Symbol)
		}
		fmt.Println("</Quotes>")
	}
	fmt.Println("</TOKENS>")
	fmt.Println()
}

func showProtocols(gotProtocols models.GetProtocolsResponse, protocols []string) {
	fmt.Println("<PROTOCOLS>")
	if len(protocols) == 0 {
		fmt.Println("</PROTOCOLS>")
		fmt.Println()
		return
	}

	lowerFilters := make([]string, 0, len(protocols))
	for _, p := range protocols {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		lowerFilters = append(lowerFilters, strings.ToLower(p))
	}

	for _, p := range gotProtocols {
		for _, l := range lowerFilters {
			if strings.Contains(strings.ToLower(p.Name), l) && p.Tvl > 0 {
				fmt.Printf("Name: %s\n", p.Name)
				fmt.Printf("Symbol: %s\n", p.Symbol)
				fmt.Printf("Description: %s\n", p.Description)
				fmt.Printf("Category: %s\n", p.Category)
				fmt.Printf("TVL: %f$\n", p.Tvl)
				fmt.Printf("Price changed 1 hour: %f%%\n", p.Change1h)
				fmt.Printf("Price changed 24 hours: %f%%\n", p.Change1d)
				fmt.Printf("Price changed 7 days: %f%%\n", p.Change7d)
			}
		}
	}
	fmt.Println("</PROTOCOLS>")
	fmt.Println()
}
