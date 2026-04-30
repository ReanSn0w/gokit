package base

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

// NewError создаёт новый Error с отформатированным сообщением и оборачивает
// переданную ошибку base как причину. Поддерживает форматирование в стиле fmt.Sprintf.
func NewError(base error, format string, args ...any) *Error {
	return &Error{
		Msg: fmt.Sprintf(format, args...),
		Err: base,
	}
}

// ReadError читает тело HTTP-ответа с ошибкой (как io.Reader) и оборачивает
// его содержимое в *Error с фиксированным Msg "readed error".
// Используется при обработке HTTP-ответов со статусом >= 300.
func ReadError(r io.Reader) (*Error, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r)
	if err != nil {
		return nil, err
	}

	return &Error{
		Msg: "readed error",
		Err: errors.New(buf.String()),
	}, nil
}

// Error — кастомная ошибка с текстовым сообщением и вложенной причиной.
// Поддерживает errors.Is (сравнение по Msg), errors.As и errors.Unwrap.
type Error struct {
	Msg string
	Err error
}

// Error реализует интерфейс error и возвращает строку вида "Msg: Err".
func (e *Error) Error() string {
	return e.Msg + ": " + e.Err.Error()
}

// Is позволяет errors.Is сравнивать ошибки по полю Msg, а не по указателю.
// Это даёт возможность использовать sentinel-ошибки вида &Error{Msg: "..."}
// для проверки типа через errors.Is, даже если это разные экземпляры.
func (e *Error) Is(target error) bool {
	var t *Error
	if errors.As(target, &t) {
		return e.Msg == t.Msg
	}
	return false
}

// Unwrap возвращает вложенную ошибку-причину, позволяя errors.Is и errors.As
// проходить по всей цепочке обёрток.
func (e *Error) Unwrap() error {
	return e.Err
}

// ErrorsMap — коллекция именованных ошибок, удобная для агрегации
// нескольких ошибок (например, результатов валидации).
// Реализует интерфейс error, а также поддерживает errors.Is и errors.Unwrap
// для многоуровневого обхода через все вложенные ошибки (Go 1.20+).
type ErrorsMap map[string]error

// Error реализует интерфейс error. Возвращает строку с количеством ошибок
// и перечислением каждой пары «ключ: сообщение».
func (e ErrorsMap) Error() string {
	buf := new(bytes.Buffer)

	fmt.Fprintf(buf, "has %d errors", len(e))

	for key, err := range e {
		fmt.Fprintf(buf, "\n%s: %v", key, err.Error())
	}

	return buf.String()
}

// Is возвращает true, если хотя бы одна из ошибок в карте удовлетворяет
// условию errors.Is(err, target). Позволяет использовать errors.Is напрямую
// с ErrorsMap как с целевой ошибкой.
func (e ErrorsMap) Is(target error) bool {
	for _, err := range e {
		if errors.Is(err, target) {
			return true
		}
	}

	return false
}

// Unwrap возвращает все ошибки из карты в виде среза. Поддерживает
// интерфейс многоуровневого разворачивания ошибок (Go 1.20+), что позволяет
// errors.Is и errors.As обходить все вложенные ошибки.
func (e ErrorsMap) Unwrap() []error {
	slice := make([]error, 0, len(e))
	for _, err := range e {
		slice = append(slice, err)
	}
	return slice
}
