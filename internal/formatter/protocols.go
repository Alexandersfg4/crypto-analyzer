package formatter

import (
	"fmt"
	"io"
	"strings"

	"github.com/Alexandersfg4/crypto-analyzer/internal/models"
	"github.com/dustin/go-humanize"
)

func Protocols(w io.Writer, gotProtocols models.GetProtocolsResponse, protocols []string) {
	fmt.Fprintln(w, "🚀 *DEX*")
	var n int
	if len(protocols) == 0 {
		for _, p := range gotProtocols {
			if p.Tvl > 0 {
				n++
				protocolData(w, n, p)
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
			n++
			protocolData(w, n, p)
		}
	}
	fmt.Fprintln(w)
}

func protocolData(w io.Writer, n int, p models.Data) {
	fmt.Fprintf(w, "%d.*%s | %s* | %s | %s:\n", n, p.Name, p.Symbol, p.Description, p.Category)
	fmt.Fprintf(w, "TVL: *%s$* - changed 1 hour: _%.2f%%_, 24h: _%.2f%%_, 7d: _%.2f%%_\n", humanize.Commaf(p.Tvl), p.Change1h, p.Change1d, p.Change7d)
}
