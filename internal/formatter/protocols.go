package formatter

import (
	"fmt"
	"io"
	"strings"

	"github.com/Alexandersfg4/crypto-analyzer/internal/models"
)

func Protocols(w io.Writer, gotProtocols models.GetProtocolsResponse, protocols []string) {
	fmt.Fprintln(w, "🚀 *DEX*")
	if len(protocols) == 0 {
		for _, p := range gotProtocols {
			if p.Tvl > 0 {
				protocolData(w, p)
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
			protocolData(w, p)
		}
	}
	fmt.Fprintln(w)
}

func protocolData(w io.Writer, p models.Data) {
	fmt.Fprintf(w, "*%s - %s* - %s - %s:\n", p.Name, p.Symbol, p.Description, p.Category)
	fmt.Fprintf(w, "TVL: _%f$_ - changed 1 hour: %f%%, 24h: %f%%, 7d: %f%%\n", p.Tvl, p.Change1h, p.Change1d, p.Change7d)
}
