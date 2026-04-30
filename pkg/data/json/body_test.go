package json_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	datajson "github.com/ReanSn0w/gokit/pkg/data/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── вспомогательные типы ────────────────────────────────────────────────────

// simpleBody — базовая структура для тестирования корректного декодирования.
type simpleBody struct {
	Name string `json:"name"`
}

// ptrValidated реализует base.Validator через pointer-receiver и всегда
// возвращает ошибку. Используется для тестирования валидационных ошибок.
type ptrValidated struct {
	Name string `json:"name"`
}

func (p *ptrValidated) Validate() error {
	return errors.New("always invalid")
}

// ── тесты Decoder ───────────────────────────────────────────────────────────

// TestDecoder_ValidJSON проверяет, что при корректном JSON-теле:
//   - следующий handler вызывается;
//   - контекст содержит правильно декодированное значение.
func TestDecoder_ValidJSON(t *testing.T) {
	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true

		body := datajson.GetBody[simpleBody](r.Context())
		require.NotNil(t, body)
		assert.Equal(t, "Alice", body.Name)

		w.WriteHeader(http.StatusOK)
	})

	handler := datajson.Decoder[simpleBody](next)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"Alice"}`))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.True(t, nextCalled, "next handler must be called on valid JSON")
	assert.Equal(t, http.StatusOK, rec.Code)
}

// TestDecoder_InvalidJSON проверяет, что при невалидном JSON-теле:
//   - следующий handler НЕ вызывается;
//   - возвращается 400;
//   - тело ответа содержит "success":false.
func TestDecoder_InvalidJSON(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("next handler must not be called on invalid JSON")
	})

	handler := datajson.Decoder[simpleBody](next)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("not json"))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), `"success":false`)
}

// TestDecoder_ValidationError проверяет, что при структуре, чей pointer-receiver
// Validate() возвращает ошибку, Decoder отвечает 400.
// Здесь T = *ptrValidated: JSON-декодер аллоцирует значение, а base.Validate
// вызывает (*ptrValidated).Validate(), которая всегда возвращает ошибку.
func TestDecoder_ValidationError(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("next handler must not be called when validation fails")
	})

	// T = *ptrValidated: декодер нальёт значение, base.Validate вызовет
	// (*ptrValidated).Validate() и получит ошибку → 400.
	handler := datajson.Decoder[*ptrValidated](next)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"test"}`))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), `"success":false`)
}

// ── тесты GetBody / SetBody ─────────────────────────────────────────────────

// TestGetBody_ReturnsValue проверяет, что GetBody возвращает то значение,
// которое было положено в контекст через SetBody.
func TestGetBody_ReturnsValue(t *testing.T) {
	want := simpleBody{Name: "test"}

	ctx, err := datajson.SetBody(context.Background(), want)
	require.NoError(t, err)

	got := datajson.GetBody[simpleBody](ctx)
	require.NotNil(t, got)
	assert.Equal(t, want, *got)
}

// TestSetBody_ValidData проверяет, что SetBody с валидными данными:
//   - возвращает nil-ошибку;
//   - кладёт данные в контекст так, что GetBody возвращает их корректно.
func TestSetBody_ValidData(t *testing.T) {
	want := simpleBody{Name: "hello"}

	ctx, err := datajson.SetBody(context.Background(), want)

	require.NoError(t, err)

	got := datajson.GetBody[simpleBody](ctx)
	require.NotNil(t, got)
	assert.Equal(t, want, *got)
}

// TestSetBody_ValidationError проверяет, что SetBody возвращает ошибку,
// если Validate() данных завершается неудачей, и не изменяет контекст.
func TestSetBody_ValidationError(t *testing.T) {
	origCtx := context.Background()

	// *ptrValidated всегда возвращает ошибку из Validate().
	newCtx, err := datajson.SetBody(origCtx, &ptrValidated{})

	require.Error(t, err)
	assert.Equal(t, origCtx, newCtx, "context must remain unchanged when validation fails")
}

// ── вспомогательная функция: декодирование тела ответа ─────────────────────

// decodeResponseBody удобный хелпер для декодирования JSON-тела ответа
// в map для последующих утверждений.
func decodeResponseBody(t *testing.T, body string) map[string]any {
	t.Helper()
	var m map[string]any
	require.NoError(t, json.Unmarshal([]byte(body), &m))
	return m
}
