package formatter

import (
	"fmt"
	"io"
	"strings"

	"github.com/Alexandersfg4/crypto-analyzer/internal/models"
	"github.com/samber/lo"
)

func News(w io.Writer, gotNews models.GetNewsResponse) {
	fmt.Fprintln(w, "🔥 *Top News*")
	for i, news := range gotNews {
		coins := make(map[string]bool, len(news.Coins))
		for _, coin := range news.Coins {
			coins[fmt.Sprintf("`%s`", coin.CoinNameKeyWords)] = true
		}

		var affectedCoins string
		if len(coins) > 0 {
			affectedCoins = strings.Join(lo.Keys(coins), ", ")
		}

		fmt.Fprintf(w, "%d. [%s](%s) %s\n", i+1, news.Title, news.Link, affectedCoins)

	}

	fmt.Fprintln(w)
}
