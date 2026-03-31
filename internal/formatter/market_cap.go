package formatter

import (
	"fmt"
	"io"

	"github.com/Alexandersfg4/crypto-analyzer/internal/models"
	"github.com/dustin/go-humanize"
)

func MarketCap(w io.Writer, gotMarketCap models.MarketCap) {
	fmt.Fprintln(w, "📊 *Market Capitalization*")
	fmt.Fprintf(
		w,
		"Market Cap: _%s$_\n",
		humanize.Comma(gotMarketCap.MarketCap),
	)
	fmt.Fprintf(
		w,
		"Volume: _%s$_\n",
		humanize.Comma(gotMarketCap.Volume),
	)
	fmt.Fprintf(
		w,
		"BTC Dominance: _%.2f%%_\n",
		gotMarketCap.BtcDominance,
	)
	fmt.Fprintf(
		w,
		"24-hour change in cap: _%.2f%%_\n",
		gotMarketCap.MarketCapChange,
	)
	fmt.Fprintf(w, "24-hour change in total trading volume: _%.2f%%_\n", gotMarketCap.VolumeChange)
	fmt.Fprintf(w, "24-hour change in Bitcoin dominance: _%.2f%%_\n", gotMarketCap.BtcDominanceChange)
	fmt.Fprintln(w)
}
