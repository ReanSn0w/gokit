package trainer

import (
	"context"

	"github.com/go-pkgz/lgr"
	"github.com/ollama/ollama/api"
)

var (
	falseFlag = false
)

func New[M any, R Match[M]](
	log lgr.L,
	ollama *api.Client,
	model string,
	rules []string,
	cases []Case[R],
	generateConfigGenerator func(*Trainer[R], *Config, *Result[R]) *api.GenerateRequest,
	opts ...func(*Trainer[R]),
) *Trainer[R] {
	trainer := &Trainer[R]{
		log:                     log,
		ollama:                  ollama,
		model:                   model,
		Rules:                   rules,
		cases:                   cases,
		generateConfigGenerator: generateConfigGenerator,
	}

	for _, f := range opts {
		f(trainer)
	}

	return trainer
}

type Trainer[R any] struct {
	log    lgr.L
	ollama *api.Client

	cases []Case[R]

	generateConfigGenerator func(*Trainer[R], *Config, *Result[R]) *api.GenerateRequest

	model string
	Rules []string
}

type Config struct {
	SystemPrompt string
}

func (t *Trainer[R]) GenerateConfig(ctx context.Context, iterations int) (*Config, error) {
	var (
		resultsTable         = ScoreTable[R]{}
		maxPoints    float64 = float64(len(t.cases))
	)

	for i := range iterations {
		t.log.Logf("[INFO] ----- Iteration %v started", i+1)

		currentBest := resultsTable.BestResult()
		failedCase := currentBest.FailedCases.BaddestCase()

		newConfig, err := t.generateNewConfig(ctx, currentBest.Config, failedCase)
		if err != nil {
			return nil, err
		}

		t.log.Logf("[INFO] Test Configuration:\n%v", ObjectToString(newConfig))

		result := Score[R]{
			Score:       0,
			Config:      newConfig,
			FailedCases: make(Results[R], 0),
		}

		for index, tc := range t.cases {
			testResult, err := t.makeTest(ctx, newConfig, &tc)
			if err != nil {
				return nil, err
			}

			t.log.Logf("[INFO] Test: %v, Score: %v", index+1, testResult.Score)
			result.Score += testResult.Score
			if testResult.Score < 0.5 {
				result.FailedCases = append(result.FailedCases, *testResult)
			}
		}

		t.log.Logf("[INFO] Configuration Points: %v / %v", result.Score, maxPoints)
		resultsTable = append(resultsTable, result)
		if result.Score == maxPoints {
			break
		}
	}

	bestCase := resultsTable.BestResult()
	return bestCase.Config, nil
}

func (t *Trainer[R]) makeTest(ctx context.Context, config *Config, tc *Case[R]) (*Result[R], error) {
	req := tc.Prompt.GenerateRequest(config)
	req.Model = t.model
	req.Stream = &falseFlag

	apiResult, err := t.generateRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	var result Result[R]
	result.Case = *tc

	result.Have, err = result.Want.Create(apiResult)
	if err != nil {
		return nil, err
	}

	result.Score = result.Have.Match(result.Case.Want.(R))
	return &result, nil
}

func (t *Trainer[R]) generateNewConfig(ctx context.Context, c *Config, tr *Result[R]) (*Config, error) {
	req := t.generateConfigGenerator(t, c, tr)
	req.Model = t.model
	req.Stream = &falseFlag

	result, err := t.generateRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	conf := &Config{SystemPrompt: result.Response}
	return conf, nil
}

func (t *Trainer[R]) generateRequest(ctx context.Context, req *api.GenerateRequest) (*api.GenerateResponse, error) {
	result := &api.GenerateResponse{}
	err := t.ollama.Generate(ctx, req, func(gr api.GenerateResponse) error {
		result = &gr
		return nil
	})

	return result, err
}
