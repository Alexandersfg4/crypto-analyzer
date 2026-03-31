package formatter

import (
	"fmt"
	"io"

	"github.com/Alexandersfg4/crypto-analyzer/internal/models"
)

func News(w io.Writer, gotNews models.GetNewsResponse) {
	fmt.Fprintln(w, "🔥 *Top News*")
	for i, news := range gotNews {
		fmt.Fprintf(w, "%d. `%s`\n", i+1, news.Title)
		coins := make([]string, 0, len(news.Coins))
		for _, coin := range news.Coins {
			coins = append(coins, fmt.Sprintf("`%s`", coin.CoinNameKeyWords))
		}
		if len(coins) > 0 {
			fmt.Fprintln(w, "	Affected coins: ", coins)
		}
	}

	fmt.Fprintln(w)
}
