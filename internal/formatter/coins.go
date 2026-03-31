package formatter

import (
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/Alexandersfg4/crypto-analyzer/internal/models"
)

func Coins(w io.Writer, gotCoins []models.ListingsLatestData, tokens []string) {
	slices.SortStableFunc(gotCoins, func(a, b models.ListingsLatestData) int {
		if a.CmcRank < b.CmcRank {
			return -1
		}
		if a.CmcRank > b.CmcRank {
			return 1
		}

		return 0
	})

	fmt.Fprintln(w, "₿ *Listed tokens*")
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

	coindByChanges90d := slices.Clone(gotCoins)
	slices.SortStableFunc(coindByChanges90d, func(a, b models.ListingsLatestData) int {
		if a.UsdQuote().PercentChange90d < b.UsdQuote().PercentChange90d {
			return 1
		}
		if a.UsdQuote().PercentChange90d > b.UsdQuote().PercentChange90d {
			return -1
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
	fmt.Fprintf(w, "`%s`: *%.2f$* | 1h: _%.2f_ | 24h: _%.2f_ | 7d: _%.2f%%_ | 90d: _%.2f%%_\n", c.Symbol, q.Price, q.PercentChange1h, q.PercentChange24h, q.PercentChange7d, q.PercentChange90d)
}
