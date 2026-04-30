package query

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── тесты для Decoder middleware ─────────────────────────────────────────────

// TestDecoder_ValidQuery проверяет, что при корректном запросе middleware
// декодирует query-параметры и передаёт управление следующему обработчику.
func TestDecoder_ValidQuery(t *testing.T) {
	type nameQuery struct{ Name string }

	var capturedCtx context.Context

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedCtx = r.Context()
		w.WriteHeader(http.StatusOK)
	})

	handler := Decoder[nameQuery](next)
	req := httptest.NewRequest(http.MethodGet, "/?Name=Alice", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code, "ожидается 200 при корректном запросе")
	require.NotNil(t, capturedCtx, "контекст должен быть захвачен обработчиком")

	data := Get[nameQuery](capturedCtx)
	require.NotNil(t, data)
	assert.Equal(t, "Alice", data.Name)
}

// TestDecoder_InvalidQuery проверяет, что при невозможности декодировать
// query-параметры middleware отвечает 400 Bad Request.
func TestDecoder_InvalidQuery(t *testing.T) {
	type ageQuery struct{ Age int }

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Не должны сюда попасть.
		w.WriteHeader(http.StatusOK)
	})

	handler := Decoder[ageQuery](next)
	// "abc" не может быть распарсен в int → decode вернёт ошибку.
	req := httptest.NewRequest(http.MethodGet, "/?Age=abc", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ── тесты для Get ─────────────────────────────────────────────────────────────

// TestGet_ReturnsValue проверяет, что значение, положенное через Set,
// корректно извлекается через Get.
func TestGet_ReturnsValue(t *testing.T) {
	type payload struct{ Score int }

	ctx := context.Background()
	ctx, err := Set[payload](ctx, payload{Score: 99})
	require.NoError(t, err)

	got := Get[payload](ctx)
	require.NotNil(t, got)
	assert.Equal(t, 99, got.Score)
}

// TestGet_NotFound проверяет, что Get возвращает nil на пустом контексте.
func TestGet_NotFound(t *testing.T) {
	got := Get[struct{ Name string }](context.Background())
	assert.Nil(t, got)
}

// ── тесты для Set ─────────────────────────────────────────────────────────────

// TestSet_ValidData проверяет, что Set корректно сохраняет валидные данные
// в контекст и возвращает nil-ошибку.
func TestSet_ValidData(t *testing.T) {
	type item struct{ Label string }

	ctx := context.Background()
	newCtx, err := Set[item](ctx, item{Label: "hello"})
	require.NoError(t, err)

	got := Get[item](newCtx)
	require.NotNil(t, got)
	assert.Equal(t, "hello", got.Label)
}

// TestSet_ValidationError проверяет, что Set возвращает ошибку валидации
// и не изменяет исходный контекст.
func TestSet_ValidationError(t *testing.T) {
	ctx := context.Background()
	newCtx, err := Set[alwaysInvalidStruct](ctx, alwaysInvalidStruct{Name: "x"})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")
	// Set возвращает исходный контекст при ошибке — данные не должны попасть в контекст.
	assert.Nil(t, Get[alwaysInvalidStruct](newCtx), "контекст не должен содержать данные при ошибке валидации")
}
