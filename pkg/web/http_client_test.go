package web_test

import (
	"io"
	"net/http"
	"strconv"
	"testing"

	"github.com/ReanSn0w/gokit/pkg/web"
	"github.com/stretchr/testify/assert"
)

func TestHTTPClient_Do(t *testing.T) {
	cl := web.NewDebugClient(http.DefaultClient, true)

	req, err := http.NewRequest("GET", "https://httpbin.org/get", nil)
	assert.NoError(t, err, "failed to create request")

	resp, err := cl.Do(req)
	assert.NoError(t, err, "Do returned error")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "expected status code 200")

	defer resp.Body.Close()

	contentLength := resp.Header.Get("Content-Length")
	if contentLength != "" {
		expectedLength, err := strconv.Atoi(contentLength)
		assert.NoError(t, err, "failed to parse Content-Length")
		bodyBytes, err := io.ReadAll(resp.Body)
		assert.NoError(t, err, "failed to read response body")
		assert.Equal(t, expectedLength, len(bodyBytes), "response body length does not match Content-Length header")
	}
}

func TestJSONRequest(t *testing.T) {
	cl := web.NewDebugClient(http.DefaultClient, false)

	data := map[string]any{}
	err := web.NewJsonRequest(cl, "https://httpbin.org/get").Do(&data)
	assert.NoError(t, err, "failed to create JSON request")

	assert.NotNil(t, data, "response data should not be nil")
	assert.NotEmpty(t, data, "response data should not be empty")
	assert.NotEmpty(t, data["url"], "response data must have url field")
}
