package formatter

import (
	"fmt"
	"io"

	"github.com/Alexandersfg4/crypto-analyzer/internal/models"
)

func News(w io.Writer, gotNews models.GetNewsResponse) {
	fmt.Fprintln(w, "🔥 Top News")
	for _, news := range gotNews {
		fmt.Fprintf(w, "[%s](%s)\n", news.Title, news.Link)
		if news.Description != "" {
			fmt.Fprintf(w, "%s\n", news.Description)
		}
		coins := make([]string, 0, len(news.Coins))
		for _, coin := range news.Coins {
			coins = append(coins, coin.CoinIDKeyWords)
		}
		if len(coins) > 0 {
			fmt.Fprintln(w, "affected coins: ", coins)
		}
	}

	fmt.Fprintln(w)
}
