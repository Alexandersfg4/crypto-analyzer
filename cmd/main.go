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
	fmt.Fprintln(w, "🔥 Top News")
	for _, news := range gotNews {
		fmt.Fprintf(w, "**>[%s](%s)\n", news.Title, news.Link)
		if news.Description != "" {
			fmt.Fprintf(w, ">%s\n", news.Description)
		}
		coins := make([]string, 0, len(news.Coins))
		for _, coin := range news.Coins {
			coins = append(coins, coin.CoinIDKeyWords)
		}
		if len(coins) > 0 {
			fmt.Fprintln(w, ">affected coins: ", coins)
		}
	}

	fmt.Fprintln(w)
}

func showFearAndGreed(w io.Writer, gotFearAndGreed models.FearAndGreed) {
	fmt.Fprintln(w, "😨 *Fear and Greed Index*")
	fmt.Fprintln(w, "_Fear and Greed Index today_")
	fmt.Fprintf(w, "Value: _%d_\n", gotFearAndGreed.Now.Value)
	fmt.Fprintf(w, "Classification: _%s_\n", gotFearAndGreed.Now.ValueClassification)
	fmt.Fprintf(w, "Updated at: _%s_\n", gotFearAndGreed.Now.UpdateTime)
	fmt.Fprintln(w, "_Fear and Greed Index yesterday_")
	fmt.Fprintf(w, "Value: _%d_\n", gotFearAndGreed.Yesterday.Value)
	fmt.Fprintf(w, "Classification: _%s_\n", gotFearAndGreed.Yesterday.ValueClassification)
	fmt.Fprintln(w, "_Fear and Greed Index last week_")
	fmt.Fprintf(w, "Value: _%d_\n", gotFearAndGreed.LastWeek.Value)
	fmt.Fprintf(w, "Classification: _%s_\n", gotFearAndGreed.LastWeek.ValueClassification)
	fmt.Fprintln(w)
}

func showMarketCap(w io.Writer, gotMarketCap models.MarketCap) {
	fmt.Fprintln(w, "📊 *Market Capitalization*")
	fmt.Fprintf(
		w,
		"Market Cap: _%d$_\n",
		gotMarketCap.MarketCap,
	)
	fmt.Fprintf(
		w,
		"Volume: _%d$_\n",
		gotMarketCap.Volume,
	)
	fmt.Fprintf(
		w,
		"BTC Dom: _%f%%_\n",
		gotMarketCap.BtcDominance,
	)
	fmt.Fprintf(
		w,
		"24-hour change in cap: _%.2f%%_\n",
		gotMarketCap.MarketCapChange,
	)
	fmt.Fprintf(w, "24-hour change in total trading volume: _%f%%_\n", gotMarketCap.VolumeChange)
	fmt.Fprintf(w, "24-hour change in Bitcoin dominance: _%f%%_\n", gotMarketCap.BtcDominanceChange)
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

	fmt.Fprintln(w, "₿ *Tokens*")
	if len(tokens) == 0 {
		for _, c := range gotCoins {
			showTokenInfo(w, c)
		}
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

	coindByChanges90d := slices.SortedFunc(gotCoins, func(a, b models.ListingsLatestData) int {
		if a.UsdQuote().PercentChange90d < b.UsdQuote().PercentChange90d {
			return -1
		}
		if a.UsdQuote().PercentChange90d > b.UsdQuote().PercentChange90d {
			return 1
		}

		return 0
	})

	fmt.Fprintln(w, "📈 *Gainers by 90d change*")
	for _, c := range coindByChanges90d[:5] {
		showTokenInfo(w, c)
	}
	fmt.Fprintln(w, "📉 *Losers by 90d change*")
	for _, c := range coindByChanges90d[len(coindByChanges90d)-5:] {
		showTokenInfo(w, c)
	}

	fmt.Fprintln(w)
}

func showTokenInfo(w io.Writer, c models.ListingsLatestData) {
	q := c.UsdQuote()
	fmt.Fprintf(w, "%s: _%.2f$_ (24h: %f, 7d: %f%%, 90d: %f%%)\n", c.Symbol, q.Price, q.Volume24h, q.PercentChange7d, q.PercentChange90d)
}

func showProtocols(w io.Writer, gotProtocols models.GetProtocolsResponse, protocols []string) {
	fmt.Fprintln(w, "🚀 *DEX*")
	if len(protocols) == 0 {
		for _, p := range gotProtocols {
			if p.Tvl > 0 {
				showProtocolData(w, p)
			}
		}
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
	fmt.Fprintln(w)
}

func showProtocolData(w io.Writer, p models.Data) {
	fmt.Fprintf(w, "*%s\n(%s)* - %s - %s:\n", p.Name, p.Symbol, p.Description, p.Category)
	fmt.Fprintf(w, "TVL: _%f$_(changed 1 hour: %f%%, 24h: %f%%, 7d: %f%%)\n", p.Tvl, p.Change1h, p.Change1d, p.Change7d)
}
