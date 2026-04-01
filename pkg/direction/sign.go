package direction

const (
	goUpSign   = "🚀"
	goDownSign = "🔻"
)

type Number interface {
	~int | ~int64 | ~float64
}

func Sign[T Number](percents ...T) string {
	var directionSign string

	m := make(map[Type]int)
	for _, p := range percents {
		if p < 0 {
			m[GoDownType]++
		} else if p > 0 {
			m[GoUpType]++
		}
	}

	switch {
	case m[GoUpType] == len(percents):
		directionSign = goUpSign
	case m[GoDownType] == len(percents):
		directionSign = goDownSign
	}

	return directionSign
}
