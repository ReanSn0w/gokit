package base

// Validator — интерфейс для типов, которые умеют валидировать сами себя.
// Реализуется единственным методом Validate() error.
type Validator interface {
	Validate() error
}

// Validate проверяет, реализует ли data интерфейс Validator, и если да —
// вызывает Validate(). Проверяет как значение, так и указатель на него (&data),
// что позволяет использовать указательные receiver'ы.
// Возвращает nil, если тип не реализует Validator.
func Validate(data any) error {
	if data, ok := data.(Validator); ok {
		if err := data.Validate(); err != nil {
			return err
		}
	}

	if data, ok := any(&data).(Validator); ok {
		if err := data.Validate(); err != nil {
			return err
		}
	}

	return nil
}
