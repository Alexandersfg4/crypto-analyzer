package models

type NewsType string

const (
	NewsTypeHandpicked = "handpicked"
	NewsTypeTrending   = "trending"
	NewsTypeLatest     = "latest"
	NewsTypeBullish    = "bullish"
	NewsTypeBearish    = "bearish"
)

type (
	GetNewsResponse []News
	News            struct {
		BigImg bool `json:"bigImg"`
		Coins  []struct {
			CoinIDKeyWords    string  `json:"coinIdKeyWords,omitzero"`
			CoinKeyWords      string  `json:"coinKeyWords,omitzero"`
			CoinNameKeyWords  string  `json:"coinNameKeyWords,omitzero"`
			CoinPercent       float64 `json:"coinPercent,omitzero"`
			CoinTitleKeyWords string  `json:"coinTitleKeyWords,omitzero"`
		} `json:"coins"`
		Content        bool     `json:"content,omitzero"`
		Description    string   `json:"description,omitzero"`
		FeedDate       int      `json:"feedDate,omitzero"`
		ID             string   `json:"id,omitzero"`
		ImgURL         string   `json:"imgUrl,omitzero"`
		IsFeatured     bool     `json:"isFeatured"`
		Link           string   `json:"link,omitzero"`
		Reactions      []any    `json:"reactions"`
		ReactionsCount any      `json:"reactionsCount"`
		RelatedCoins   []string `json:"relatedCoins"`
		SearchKeyWords []string `json:"searchKeyWords"`
		ShareURL       string   `json:"shareURL,omitzero"`
		Source         string   `json:"source,omitzero"`
		SourceLink     string   `json:"sourceLink,omitzero"`
		Title          string   `json:"title,omitzero"`
	}
)
