package formatter

import (
	"fmt"
	"io"

	"github.com/Alexandersfg4/crypto-analyzer/internal/models"
	"github.com/samber/lo"
)

func News(w io.Writer, gotNews models.GetNewsResponse) {
	fmt.Fprintln(w, "🔥 *Top News*")
	for i, news := range gotNews {
		fmt.Fprintf(w, "%d. [%s](%s)\n", i+1, news.Title, news.Link)
		coins := make(map[string]bool, len(news.Coins))
		for _, coin := range news.Coins {
			coins[fmt.Sprintf("`%s`", coin.CoinNameKeyWords)] = true
		}
		if len(coins) > 0 {
			fmt.Fprintln(w, "Affected coins: ", lo.Keys(coins))
		}
	}

	fmt.Fprintln(w)
}
