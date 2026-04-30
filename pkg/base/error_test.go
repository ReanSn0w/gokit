package base

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── ReadError ────────────────────────────────────────────────────────────────

func TestReadError_NormalCase(t *testing.T) {
	r := strings.NewReader("something went wrong")
	result, err := ReadError(r)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "readed error", result.Msg)
	assert.Equal(t, "something went wrong", result.Err.Error())
}

func TestReadError_EmptyBody(t *testing.T) {
	r := strings.NewReader("")
	result, err := ReadError(r)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "readed error", result.Msg)
	assert.Equal(t, "", result.Err.Error())
}

// ── Error ────────────────────────────────────────────────────────────────────

func TestNewError_FieldsAreSetCorrectly(t *testing.T) {
	cause := errors.New("db connection refused")
	err := NewError(cause, "user %d not found", 42)

	assert.Equal(t, "user 42 not found", err.Msg)
	assert.Equal(t, cause, err.Err)
}

func TestError_ErrorString(t *testing.T) {
	cause := errors.New("timeout")
	err := NewError(cause, "fetch failed")

	assert.Equal(t, "fetch failed: timeout", err.Error())
}

func TestError_Is_SameMsg(t *testing.T) {
	cause := errors.New("root")
	err := NewError(cause, "not found")
	sentinel := &Error{Msg: "not found"}

	// errors.Is должен вызывать наш метод Is и сравнивать по Msg
	assert.True(t, errors.Is(err, sentinel))
}

func TestError_Is_DifferentMsg(t *testing.T) {
	cause := errors.New("root")
	err := NewError(cause, "not found")
	sentinel := &Error{Msg: "unauthorized"}

	assert.False(t, errors.Is(err, sentinel))
}

func TestError_Is_NonErrorType(t *testing.T) {
	cause := errors.New("root")
	err := NewError(cause, "not found")

	// Сравнение с ошибкой другого типа должно вернуть false
	assert.False(t, errors.Is(err, errors.New("not found")))
}

func TestError_Unwrap_ReturnsCause(t *testing.T) {
	cause := errors.New("root cause")
	err := NewError(cause, "wrapper")

	assert.Equal(t, cause, err.Unwrap())
}

func TestError_Is_ThroughStdWrap(t *testing.T) {
	cause := errors.New("root")
	appErr := NewError(cause, "not found")

	// errors.Is должен проходить через цепочку fmt.Errorf %w
	wrapped := fmt.Errorf("service: %w", appErr)
	sentinel := &Error{Msg: "not found"}

	assert.True(t, errors.Is(wrapped, sentinel))
}

func TestError_As_ExtractsType(t *testing.T) {
	cause := errors.New("root")
	appErr := NewError(cause, "permission denied")
	wrapped := fmt.Errorf("handler: %w", appErr)

	var target *Error
	require.True(t, errors.As(wrapped, &target))
	assert.Equal(t, "permission denied", target.Msg)
}

func TestError_Unwrap_ChainThroughCause(t *testing.T) {
	root := errors.New("root")
	inner := NewError(root, "inner")
	outer := NewError(inner, "outer")

	// errors.Is должен дойти до root через два уровня Unwrap
	assert.True(t, errors.Is(outer, root))
}

// ── ErrorsMap ────────────────────────────────────────────────────────────────

func TestErrorsMap_Error_ContainsKeysAndMessages(t *testing.T) {
	em := ErrorsMap{
		"name":  errors.New("too short"),
		"email": errors.New("invalid format"),
	}

	msg := em.Error()

	assert.True(t, strings.Contains(msg, "name") && strings.Contains(msg, "too short"),
		"ожидался ключ 'name' и сообщение 'too short' в: %s", msg)
	assert.True(t, strings.Contains(msg, "email") && strings.Contains(msg, "invalid format"),
		"ожидался ключ 'email' и сообщение 'invalid format' в: %s", msg)
}

func TestErrorsMap_Is_MatchesContainedError(t *testing.T) {
	sentinel := errors.New("not found")

	em := ErrorsMap{
		"user": sentinel,
	}

	assert.True(t, errors.Is(em, sentinel))
}

func TestErrorsMap_Is_NoMatch(t *testing.T) {
	em := ErrorsMap{
		"user": errors.New("not found"),
	}

	assert.False(t, errors.Is(em, errors.New("unauthorized")))
}

func TestErrorsMap_Is_WithAppError(t *testing.T) {
	cause := errors.New("root")
	appErr := NewError(cause, "validation failed")
	sentinel := &Error{Msg: "validation failed"}

	em := ErrorsMap{
		"body": appErr,
	}

	// errors.Is на ErrorsMap должен дойти до AppError и сравнить по Msg
	assert.True(t, errors.Is(em, sentinel))
}

func TestErrorsMap_Unwrap_ReturnsAllErrors(t *testing.T) {
	e1 := errors.New("e1")
	e2 := errors.New("e2")

	em := ErrorsMap{
		"a": e1,
		"b": e2,
	}

	unwrapped := em.Unwrap()
	require.Len(t, unwrapped, 2)

	// Порядок не гарантирован — проверяем наличие
	assert.ElementsMatch(t, []error{e1, e2}, unwrapped)
}

func TestErrorsMap_Unwrap_EmptyMap(t *testing.T) {
	em := make(ErrorsMap)
	assert.Empty(t, em.Unwrap())
}
