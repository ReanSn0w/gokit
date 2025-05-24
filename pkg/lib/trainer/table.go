package trainer

type ScoreTable[R any] []Score[R]

type Score[R any] struct {
	Score       float64
	Config      *Config
	FailedCases Results[R]
}

func (st ScoreTable[R]) BestResult() *Score[R] {
	var bestResult = &Score[R]{Score: 0, Config: nil, FailedCases: make(Results[R], 0)}

	for index, score := range st {
		if score.Score > bestResult.Score {
			bestResult = &st[index]
		}
	}

	return bestResult
}
