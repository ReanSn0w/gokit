package json

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/ReanSn0w/gokit/pkg/base"
)

// maxBufferSize is the maximum byte size of the bufio.Scanner buffer used by
// [Request.Stream]. Set to 512 KiB.
var (
	maxBufferSize = 1000 * 512 // 512kb
)

// NewRequest constructs a new [Request] with sensible defaults: HTTP GET
// method, empty headers, empty query parameters, and no body.
func NewRequest(cl HTTPClient, requestURL string) *Request {
	return &Request{
		client:  cl,
		method:  http.MethodGet,
		headers: make(map[string]string),
		query:   make(url.Values, 0),
		url:     requestURL,
		body:    nil,
	}
}

// HTTPClient is the interface for the underlying HTTP transport. The standard
// *http.Client satisfies this interface, and it can be replaced with a mock
// in tests.
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

// Request is a fluent builder for outgoing JSON HTTP requests. Build up the
// request with the Set* methods, then execute it via [Request.Do] or
// [Request.Stream].
type Request struct {
	// client is the HTTP transport used to execute the request.
	client HTTPClient

	// method is the HTTP verb (GET, POST, PUT, …).
	method string

	// headers holds the request headers to be sent.
	headers map[string]string

	// query holds the URL query parameters.
	query url.Values

	// url is the base request URL, without a query string.
	url string

	// body is the value to be JSON-encoded as the request body.
	// nil means no body.
	body any
}

// Optional conditionally applies a modifier function to the request. If
// enabled is true, fn is called and its result is returned; otherwise the
// unchanged receiver is returned.
//
// This is useful for conditional configuration in fluent chains:
//
//	req.Optional(authEnabled, func(r *Request) *Request {
//	    return r.SetHeader("Authorization", "Bearer "+token)
//	})
func (j *Request) Optional(enabled bool, fn func(j *Request) *Request) *Request {
	if enabled {
		return fn(j)
	}

	return j
}

// SetMethod overrides the HTTP method (default: GET).
func (j *Request) SetMethod(method string) *Request {
	j.method = method
	return j
}

// SetHeader adds or overwrites a single HTTP request header.
func (j *Request) SetHeader(name, val string) *Request {
	j.headers[name] = val
	return j
}

// SetQuery sets one or more values for the named URL query parameter.
// Multiple calls with the same name overwrite the previous values.
func (j *Request) SetQuery(name string, val ...string) *Request {
	j.query[name] = val
	return j
}

// SetBody sets the request body payload. The value is JSON-encoded when the
// request is executed via [Request.Do] or [Request.Stream].
func (j *Request) SetBody(body any) *Request {
	j.body = body
	return j
}

// Stream executes the request and returns a read-only channel of raw byte
// tokens. The tokens are produced by a [bufio.Scanner] using sf as the split
// function. Scanning runs in a separate goroutine; the channel is closed when
// scanning is complete or the response body is exhausted.
//
// If the response status code is >= 300, Stream reads the response body via
// [base.ReadError] and returns the resulting error without opening a channel.
func (j *Request) Stream(sf bufio.SplitFunc) (<-chan []byte, error) {
	request, err := j.makeRequest()
	if err != nil {
		return nil, err
	}

	resp, err := j.client.Do(request)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 300 {
		respError, err := base.ReadError(resp.Body)
		if err != nil {
			return nil, err
		} else {
			return nil, respError
		}
	}

	out := make(chan []byte)

	go func() {
		scanner := bufio.NewScanner(resp.Body)
		scanner.Split(sf)

		scanBuf := make([]byte, 0, maxBufferSize)
		scanner.Buffer(scanBuf, maxBufferSize)

		for scanner.Scan() {
			out <- scanner.Bytes()
		}

		close(out)
	}()

	return out, nil
}

// Do executes the request and JSON-decodes the successful response body into
// body. body must be a non-nil pointer to the target value.
//
// If the response status code is >= 300, Do reads the response body via
// [base.ReadError] and returns the resulting error without decoding.
func (j *Request) Do(body any) error {
	request, err := j.makeRequest()
	if err != nil {
		return err
	}

	resp, err := j.client.Do(request)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 300 {
		respError, err := base.ReadError(resp.Body)
		if err != nil {
			return err
		} else {
			return respError
		}
	}

	return json.NewDecoder(resp.Body).Decode(body)
}

// makeRequest builds a *http.Request from the current Request configuration.
// It appends the query string (if any), JSON-encodes the body (if any), and
// applies all stored headers.
func (j *Request) makeRequest() (*http.Request, error) {
	url := j.url
	if len(j.query) > 0 {
		url += "?" + j.query.Encode()
	}

	var bodyReader io.Reader
	if j.body != nil {
		buffer := new(bytes.Buffer)
		err := json.NewEncoder(buffer).Encode(j.body)
		if err != nil {
			return nil, err
		}

		bodyReader = buffer
	}

	request, err := http.NewRequest(j.method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	for key, value := range j.headers {
		request.Header.Set(key, value)
	}

	return request, nil
}
