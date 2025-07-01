package web

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"time"
)

func NewDebugClient(cl HTTPClient, body bool) HTTPClient {
	return &DebugClient{
		body:   body,
		client: cl,
	}
}

type (
	DebugClient struct {
		body   bool
		client HTTPClient
	}

	HTTPClient interface {
		Do(*http.Request) (*http.Response, error)
	}
)

func (r *DebugClient) Do(req *http.Request) (*http.Response, error) {
	t := time.Now()

	reqData, err := httputil.DumpRequest(req, r.body)
	if err != nil {
		return nil, err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}

	respData, err := httputil.DumpResponse(resp, r.body)
	if err != nil {
		return nil, err
	}

	distance := time.Since(t)
	fmt.Printf(
		"[%s] %s (%s)\n\nRequest:\n%s\n\nResponse:\n%s\n\n",
		resp.Status, resp.Request.URL, distance, string(reqData), string(respData))

	return resp, nil
}
