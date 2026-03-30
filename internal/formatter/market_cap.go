package formatter

import (
	"fmt"
	"io"

	"github.com/Alexandersfg4/crypto-analyzer/internal/models"
)

func MarketCap(w io.Writer, gotMarketCap models.MarketCap) {
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
		"BTC Dominance: _%f%%_\n",
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
