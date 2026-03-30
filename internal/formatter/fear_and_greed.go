package formatter

import (
	"fmt"
	"io"

	"github.com/Alexandersfg4/crypto-analyzer/internal/models"
)

func FearAndGreed(w io.Writer, gotFearAndGreed models.FearAndGreed) {
	fmt.Fprintln(w, "😨 *Fear and Greed Index*")
	fmt.Fprintln(w, "_Fear and Greed Index today_")
	fmt.Fprintf(w, "Value: _%d_\n", gotFearAndGreed.Now.Value)
	fmt.Fprintf(w, "Classification: _%s_\n", gotFearAndGreed.Now.ValueClassification)
	fmt.Fprintf(w, "Updated at: _%s_\n", gotFearAndGreed.Now.UpdateTime)
	fmt.Fprintln(w, "_Fear and Greed Index yesterday_")
	fmt.Fprintf(w, "Value: _%d_\n", gotFearAndGreed.Yesterday.Value)
	fmt.Fprintf(w, "Classification: _%s_\n", gotFearAndGreed.Yesterday.ValueClassification)
	fmt.Fprintln(w, "_Fear and Greed Index last week_")
	fmt.Fprintf(w, "Value: _%d_\n", gotFearAndGreed.LastWeek.Value)
	fmt.Fprintf(w, "Classification: _%s_\n", gotFearAndGreed.LastWeek.ValueClassification)
	fmt.Fprintln(w)
}
