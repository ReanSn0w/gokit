package trainer

import "github.com/ollama/ollama/api"

// GenerateRequest - Добавляет метод
// для построения запроса к Ollama
type GenerateRequest interface {
	GenerateRequest(*Config) *api.GenerateRequest
}

// Match - добавляет к структуре метод
// для выполнения проверки
//
// Match служит для проверки результата
// Create cоздает экземпляр на основе *api.GenerateResponse
//
// Результат функции должен принимать
// значение от 0 до 1, где
// 1 - полностью подходящий результат
// 0 - результат не подходит
type Match[R any] interface {
	Match(c R) float64
	Create(*api.GenerateResponse) (Match[R], error)
}

// Case[R any, M Match[R]] - структура описывает тест кейс
type Case[R any] struct {
	Prompt GenerateRequest
	Want   Match[R]
}

// Result[R any, M Match[R]] - результат выполнения теста
type Result[R any] struct {
	Case[R]
	Have Match[R]

	Score float64
}

type Results[R any] []Result[R]

func (r Results[R]) BaddestCase() *Result[R] {
	var baddestCase *Result[R]

	for index := range r {
		if index == 0 {
			baddestCase = &r[index]
		}

		if baddestCase.Score > r[index].Score {
			baddestCase = &r[index]
		}
	}

	return baddestCase
}
