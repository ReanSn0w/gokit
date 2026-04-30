package web

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockLogger is a no-op implementation of lgr.L used in tests.
type mockLogger struct {
	lines []string
}

func (m *mockLogger) Logf(format string, args ...interface{}) {
	// intentionally a no-op; captured only when needed
}

// okHandler is a minimal next-handler that always replies 200 "ok".
var okHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
})

// ──────────────────────────────────────────────────────────────────────────────
// Ping
// ──────────────────────────────────────────────────────────────────────────────

func TestPing_PingPath(t *testing.T) {
	handler := Ping()(okHandler)

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "pong", rec.Body.String())
}

func TestPing_OtherPath(t *testing.T) {
	handler := Ping()(okHandler)

	req := httptest.NewRequest(http.MethodGet, "/other", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "ok", rec.Body.String())
}

// ──────────────────────────────────────────────────────────────────────────────
// DebugRequest
// ──────────────────────────────────────────────────────────────────────────────

func TestDebugRequest_Disabled(t *testing.T) {
	log := &mockLogger{}
	handler := DebugRequest(false, log)(okHandler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "ok", rec.Body.String())
}

func TestDebugRequest_Enabled(t *testing.T) {
	log := &mockLogger{}
	handler := DebugRequest(true, log)(okHandler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// middleware only logs — it must not block the request
	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "ok", rec.Body.String())
}

// ──────────────────────────────────────────────────────────────────────────────
// APIKey
// ──────────────────────────────────────────────────────────────────────────────

const (
	testAPIHeader = "X-API-Key"
	testAPIKey    = "secret-key"
)

func TestAPIKey_ValidKey(t *testing.T) {
	handler := APIKey(testAPIHeader, testAPIKey)(okHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(testAPIHeader, testAPIKey)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "ok", rec.Body.String())
}

func TestAPIKey_InvalidKey(t *testing.T) {
	handler := APIKey(testAPIHeader, testAPIKey)(okHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(testAPIHeader, "wrong-key")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid api key")
}

func TestAPIKey_MissingKey(t *testing.T) {
	handler := APIKey(testAPIHeader, testAPIKey)(okHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	// no header set
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid api key")
}
