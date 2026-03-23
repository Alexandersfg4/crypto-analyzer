package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
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
	showMarketCap(output, gotData.marketCap)
	showFearAndGreed(output, gotData.fearAndGreed)
	showCoins(output, gotData.listingsLatest, opts.tokens)
	showNews(output, gotData.news)
	showProtocols(output, gotData.protocols, opts.protocols)

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

			gotCoins, err := srvCoinmarketcap.GetListingsLatest(ctx, 2, limitCoins)
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

func showNews(w io.Writer, gotNews models.GetNewsResponse) {
	fmt.Fprintln(w, "<NEWS>")
	for _, news := range gotNews {
		fmt.Fprintf(w, "Title: %s\n", news.Title)
		if news.Description != "" {
			fmt.Fprintf(w, "Description: %s\n", news.Description)
		}
		fmt.Fprintf(w, "Source: %s\n", news.Source)
		fmt.Fprintf(w, "Link: %s\n", news.Link)
		coins := make([]string, 0, len(news.Coins))
		for _, coin := range news.Coins {
			coins = append(coins, coin.CoinIDKeyWords)
		}
		if len(coins) > 0 {
			fmt.Fprintln(w, "Affected coins: ", coins)
		}
	}

	fmt.Fprintln(w, "</NEWS>")
	fmt.Fprintln(w)
}

func showFearAndGreed(w io.Writer, gotFearAndGreed models.FearAndGreed) {
	fmt.Fprintln(w, "<Fear and Greed Index now>")
	fmt.Fprintf(w, "Value: %d\n", gotFearAndGreed.Now.Value)
	fmt.Fprintf(w, "Classification: %s\n", gotFearAndGreed.Now.ValueClassification)
	fmt.Fprintf(w, "Updated at: %s\n", gotFearAndGreed.Now.UpdateTime)
	fmt.Fprintln(w, "Fear and Greed Index yesterday")
	fmt.Fprintf(w, "Value: %d\n", gotFearAndGreed.Yesterday.Value)
	fmt.Fprintf(w, "Classification: %s\n", gotFearAndGreed.Yesterday.ValueClassification)
	fmt.Fprintln(w, "Fear and Greed Index last week")
	fmt.Fprintf(w, "Value: %d\n", gotFearAndGreed.LastWeek.Value)
	fmt.Fprintf(w, "Classification: %s\n", gotFearAndGreed.LastWeek.ValueClassification)
	fmt.Fprintln(w, "</Fear and Greed Index now>")
	fmt.Fprintln(w)
}

func showMarketCap(w io.Writer, gotMarketCap models.MarketCap) {
	fmt.Fprintln(w, "<Market Capitalization>")
	fmt.Fprintf(
		w,
		"Total market capitalization of all cryptocurrencies : %d$\n",
		gotMarketCap.MarketCap,
	)
	fmt.Fprintf(
		w,
		"Total 24-hour trading volume across all cryptocurrencies: %d$\n",
		gotMarketCap.Volume,
	)
	fmt.Fprintf(
		w,
		"Bitcoin's percentage share of the total cryptocurrency market capitalization: %f%%\n",
		gotMarketCap.BtcDominance,
	)
	fmt.Fprintf(
		w,
		"24-hour change in total market capitalization: %f%%\n",
		gotMarketCap.MarketCapChange,
	)
	fmt.Fprintf(w, "24-hour change in total trading volume: %f%%\n", gotMarketCap.VolumeChange)
	fmt.Fprintf(w, "24-hour change in Bitcoin dominance: %f%%\n", gotMarketCap.BtcDominanceChange)
	fmt.Fprintln(w, "</Market Capitalization>")
	fmt.Fprintln(w)
}

func showCoins(w io.Writer, gotCoins []models.ListingsLatestData, tokens []string) {
	slices.SortStableFunc(gotCoins, func(a, b models.ListingsLatestData) int {
		if a.CmcRank < b.CmcRank {
			return -1
		}
		if a.CmcRank > b.CmcRank {
			return 1
		}

		return 0
	})

	fmt.Fprintln(w, "<TOKENS>")
	if len(tokens) == 0 {
		for _, c := range gotCoins {
			showTokenInfo(w, c)
		}
		fmt.Fprintln(w, "</TOKENS>")
		fmt.Fprintln(w)
		return
	}

	tokenSet := make(map[string]struct{}, len(tokens))
	for _, t := range tokens {
		tokenSet[strings.ToUpper(t)] = struct{}{}
	}

	for _, c := range gotCoins {
		if _, ok := tokenSet[c.Symbol]; !ok {
			continue
		}

		showTokenInfo(w, c)
	}
	fmt.Fprintln(w, "</TOKENS>")
	fmt.Fprintln(w)
}

func showTokenInfo(w io.Writer, c models.ListingsLatestData) {
	fmt.Fprintf(w, "Name: %s\n", c.Name)
	fmt.Fprintf(w, "Symbol: %s\n", c.Symbol)
	fmt.Fprintln(w, "<Quotes>")
	for _, q := range c.Quote {
		fmt.Fprintf(w, "<%s>\n", q.Symbol)
		fmt.Fprintf(w, "Price: %f\n", q.Price)
		fmt.Fprintf(w, "Volume for 24h: %f\n", q.Volume24h)
		fmt.Fprintf(w, "Market Cap: %f\n", q.MarketCap)
		fmt.Fprintf(w, "Price changed 1 hour: %f%%\n", q.PercentChange1h)
		fmt.Fprintf(w, "Price changed 24 hours: %f%%\n", q.PercentChange24h)
		fmt.Fprintf(w, "Price changed 7 days: %f%%\n", q.PercentChange7d)
		fmt.Fprintf(w, "Price changed 90 days: %f%%\n", q.PercentChange90d)
		fmt.Fprintf(w, "</%s>\n", q.Symbol)
	}
	fmt.Fprintln(w, "</Quotes>")
}

func showProtocols(w io.Writer, gotProtocols models.GetProtocolsResponse, protocols []string) {
	fmt.Fprintln(w, "<PROTOCOLS>")
	if len(protocols) == 0 {
		for _, p := range gotProtocols {
			if p.Tvl > 0 {
				showProtocolData(w, p)
			}
		}
		fmt.Fprintln(w, "</PROTOCOLS>")
		fmt.Fprintln(w)
		return
	}

	upperFilters := make(map[string]bool, len(protocols))
	for _, p := range protocols {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}

		upperFilters[strings.ToUpper(p)] = true
	}

	for _, p := range gotProtocols {
		if p.Tvl > 0 && upperFilters[p.Symbol] {
			showProtocolData(w, p)
		}
	}
	fmt.Fprintln(w, "</PROTOCOLS>")
	fmt.Fprintln(w)
}

func showProtocolData(w io.Writer, p models.Data) {
	fmt.Fprintf(w, "Name: %s\n", p.Name)
	fmt.Fprintf(w, "Symbol: %s\n", p.Symbol)
	fmt.Fprintf(w, "Description: %s\n", p.Description)
	fmt.Fprintf(w, "Category: %s\n", p.Category)
	fmt.Fprintf(w, "TVL: %f$\n", p.Tvl)
	fmt.Fprintf(w, "Price changed 1 hour: %f%%\n", p.Change1h)
	fmt.Fprintf(w, "Price changed 24 hours: %f%%\n", p.Change1d)
	fmt.Fprintf(w, "Price changed 7 days: %f%%\n", p.Change7d)
}
