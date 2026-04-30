package json_test

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	datajson "github.com/ReanSn0w/gokit/pkg/data/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── mock HTTP-клиент ─────────────────────────────────────────────────────────

// mockHTTPClient реализует datajson.HTTPClient для тестирования. Сохраняет
// последний входящий запрос в поле req, чтобы тесты могли его проверить.
type mockHTTPClient struct {
	resp *http.Response
	err  error
	req  *http.Request // захваченный запрос из последнего вызова Do
}

func (m *mockHTTPClient) Do(r *http.Request) (*http.Response, error) {
	m.req = r
	return m.resp, m.err
}

// okResponse создаёт 200-ответ с переданным телом.
func okResponse(body string) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

// ── вспомогательные типы ────────────────────────────────────────────────────

// nameResult используется для декодирования тестовых JSON-ответов.
type nameResult struct {
	Name string `json:"name"`
}

// ── тесты NewRequest ─────────────────────────────────────────────────────────

// TestNewRequest_Defaults проверяет, что NewRequest создаёт запрос
// с методом GET и корректным URL.
func TestNewRequest_Defaults(t *testing.T) {
	cl := &mockHTTPClient{resp: okResponse("null")}

	req := datajson.NewRequest(cl, "http://example.com")

	var out any
	err := req.Do(&out)
	require.NoError(t, err)

	assert.Equal(t, http.MethodGet, cl.req.Method)
	assert.Equal(t, "http://example.com", cl.req.URL.String())
}

// ── тесты сеттеров ───────────────────────────────────────────────────────────

// TestRequest_SetMethod проверяет, что SetMethod изменяет HTTP-метод запроса.
func TestRequest_SetMethod(t *testing.T) {
	cl := &mockHTTPClient{resp: okResponse("null")}

	req := datajson.NewRequest(cl, "http://example.com").SetMethod(http.MethodPost)

	var out any
	err := req.Do(&out)
	require.NoError(t, err)

	assert.Equal(t, http.MethodPost, cl.req.Method)
}

// TestRequest_SetHeader проверяет, что SetHeader добавляет заголовок,
// который виден в итоговом *http.Request.
func TestRequest_SetHeader(t *testing.T) {
	cl := &mockHTTPClient{resp: okResponse("null")}

	req := datajson.NewRequest(cl, "http://example.com").
		SetHeader("X-Custom", "value")

	var out any
	err := req.Do(&out)
	require.NoError(t, err)

	assert.Equal(t, "value", cl.req.Header.Get("X-Custom"))
}

// TestRequest_SetQuery проверяет, что SetQuery добавляет query-параметр,
// который виден в URL итогового запроса.
func TestRequest_SetQuery(t *testing.T) {
	cl := &mockHTTPClient{resp: okResponse("null")}

	req := datajson.NewRequest(cl, "http://example.com").
		SetQuery("page", "1")

	var out any
	err := req.Do(&out)
	require.NoError(t, err)

	assert.Equal(t, "1", cl.req.URL.Query().Get("page"))
}

// ── тесты Optional ───────────────────────────────────────────────────────────

// TestRequest_Optional_Enabled проверяет, что при enabled = true
// функция-модификатор применяется к запросу.
func TestRequest_Optional_Enabled(t *testing.T) {
	cl := &mockHTTPClient{resp: okResponse("null")}

	req := datajson.NewRequest(cl, "http://example.com").
		Optional(true, func(r *datajson.Request) *datajson.Request {
			return r.SetHeader("X-Applied", "yes")
		})

	var out any
	err := req.Do(&out)
	require.NoError(t, err)

	assert.Equal(t, "yes", cl.req.Header.Get("X-Applied"))
}

// TestRequest_Optional_Disabled проверяет, что при enabled = false
// функция-модификатор НЕ применяется.
func TestRequest_Optional_Disabled(t *testing.T) {
	cl := &mockHTTPClient{resp: okResponse("null")}

	req := datajson.NewRequest(cl, "http://example.com").
		Optional(false, func(r *datajson.Request) *datajson.Request {
			return r.SetHeader("X-Applied", "yes")
		})

	var out any
	err := req.Do(&out)
	require.NoError(t, err)

	assert.Empty(t, cl.req.Header.Get("X-Applied"))
}

// ── тесты Do ─────────────────────────────────────────────────────────────────

// TestRequest_Do_Success проверяет, что при статусе 200 и корректном JSON
// Do декодирует тело ответа в переданный указатель.
func TestRequest_Do_Success(t *testing.T) {
	cl := &mockHTTPClient{
		resp: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"name":"Bob"}`)),
		},
	}

	req := datajson.NewRequest(cl, "http://example.com")

	var result nameResult
	err := req.Do(&result)

	require.NoError(t, err)
	assert.Equal(t, "Bob", result.Name)
}

// TestRequest_Do_ErrorStatus проверяет, что при статусе >= 300
// Do возвращает ошибку.
func TestRequest_Do_ErrorStatus(t *testing.T) {
	cl := &mockHTTPClient{
		resp: &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       io.NopCloser(strings.NewReader("bad input")),
		},
	}

	req := datajson.NewRequest(cl, "http://example.com")

	var result any
	err := req.Do(&result)

	require.Error(t, err)
}

// TestRequest_Do_ClientError проверяет, что если HTTPClient.Do возвращает
// ошибку транспортного уровня, Do пробрасывает её вызывающей стороне.
func TestRequest_Do_ClientError(t *testing.T) {
	cl := &mockHTTPClient{
		err: errors.New("timeout"),
	}

	req := datajson.NewRequest(cl, "http://example.com")

	var result any
	err := req.Do(&result)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "timeout")
}
