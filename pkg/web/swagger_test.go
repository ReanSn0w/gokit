package web

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewScalarHandler_StatusCode(t *testing.T) {
	handler := NewScalarHandler("My API", "/openapi.json")

	req := httptest.NewRequest(http.MethodGet, "/docs", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
}

func TestNewScalarHandler_ContainsTitle(t *testing.T) {
	title := "Awesome API"
	handler := NewScalarHandler(title, "/openapi.json")

	req := httptest.NewRequest(http.MethodGet, "/docs", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Contains(t, rec.Body.String(), title)
}

func TestNewScalarHandler_ContainsDocURL(t *testing.T) {
	docURL := "/api/v1/openapi.json"
	handler := NewScalarHandler("My API", docURL)

	req := httptest.NewRequest(http.MethodGet, "/docs", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Contains(t, rec.Body.String(), docURL)
}

func TestNewScalarHandler_ContentType(t *testing.T) {
	handler := NewScalarHandler("My API", "/openapi.json")

	req := httptest.NewRequest(http.MethodGet, "/docs", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	body := rec.Body.String()
	assert.NotEmpty(t, body)
	assert.Contains(t, body, "<!DOCTYPE html>")
}
