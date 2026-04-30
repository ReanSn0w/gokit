package json_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ReanSn0w/gokit/pkg/base"
	datajson "github.com/ReanSn0w/gokit/pkg/data/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewResponse_WithError проверяет, что NewResponse(error) создаёт
// Response с Success = false и заполненным Message.
func TestNewResponse_WithError(t *testing.T) {
	resp := datajson.NewResponse(errors.New("oops"))

	assert.False(t, resp.Success)
	assert.Equal(t, "oops", resp.Message)
	assert.Nil(t, resp.Errors)
}

// TestNewResponse_WithData проверяет, что NewResponse(data) создаёт
// Response с Success = true и заполненным Data.
func TestNewResponse_WithData(t *testing.T) {
	resp := datajson.NewResponse(42)

	assert.True(t, resp.Success)
	assert.Equal(t, 42, resp.Data)
	assert.Empty(t, resp.Message)
}

// TestNewResponse_WithErrorsMap проверяет, что NewResponse(base.ErrorsMap)
// создаёт Response с Success = false и непустым полем Errors.
func TestNewResponse_WithErrorsMap(t *testing.T) {
	em := base.ErrorsMap{"field": errors.New("required")}

	resp := datajson.NewResponse(em)

	assert.False(t, resp.Success)
	assert.NotNil(t, resp.Errors)
	assert.Empty(t, resp.Message)
}

// TestNewPlainResponse проверяет, что NewPlainResponse форматирует сообщение
// и создаёт Response с Success = true.
func TestNewPlainResponse(t *testing.T) {
	resp := datajson.NewPlainResponse("done %d", 3)

	assert.True(t, resp.Success)
	assert.Equal(t, "done 3", resp.Message)
}

// TestResponse_Write_StatusCode проверяет, что Write выставляет
// переданный HTTP-статус код.
func TestResponse_Write_StatusCode(t *testing.T) {
	rec := httptest.NewRecorder()

	datajson.NewResponse(42).Write(http.StatusCreated, rec)

	assert.Equal(t, http.StatusCreated, rec.Code)
}

// TestResponse_Write_ContentType проверяет, что Write выставляет заголовок
// Content-Type: application/json.
func TestResponse_Write_ContentType(t *testing.T) {
	rec := httptest.NewRecorder()

	datajson.NewResponse(42).Write(http.StatusOK, rec)

	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
}

// TestResponse_Write_Body проверяет, что записанное тело ответа корректно
// декодируется обратно в Response[int] с правильными полями.
func TestResponse_Write_Body(t *testing.T) {
	rec := httptest.NewRecorder()

	datajson.NewResponse(42).Write(http.StatusOK, rec)

	var got datajson.Response[int]
	err := json.NewDecoder(rec.Body).Decode(&got)
	require.NoError(t, err)

	assert.True(t, got.Success)
	assert.Equal(t, 42, got.Data)
	assert.Empty(t, got.Message)
}
