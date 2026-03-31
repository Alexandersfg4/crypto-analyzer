package formatter

import (
	"fmt"
	"io"

	"github.com/Alexandersfg4/crypto-analyzer/internal/models"
)

func FearAndGreed(w io.Writer, gotFearAndGreed models.FearAndGreed) {
	fmt.Fprintln(w, "😨 *Fear and Greed Index*")
	fmt.Fprintf(w, "Fear and Greed Index today: *%d*(%s)\n", gotFearAndGreed.Now.Value, gotFearAndGreed.Now.ValueClassification)
	fmt.Fprintf(w, "Updated at: _%s_\n", gotFearAndGreed.Now.UpdateTime)
	fmt.Fprintf(w, "Fear and Greed Index yesterday: *%d*(%s)\n", gotFearAndGreed.Yesterday.Value, gotFearAndGreed.Yesterday.ValueClassification)
	fmt.Fprintf(w, "Fear and Greed Index last week: *%d*(%s)\n", gotFearAndGreed.LastWeek.Value, gotFearAndGreed.LastWeek.ValueClassification)
	fmt.Fprintln(w)
}
