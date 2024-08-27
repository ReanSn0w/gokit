package web

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

var (
	maxBufferSize = 1000 * 512 // 512kb
)

func NewJsonRequest(cl *http.Client, requestURL string) *JsonRequest {
	return &JsonRequest{
		method:  http.MethodGet,
		headers: make(map[string]string),
		query:   make(url.Values, 0),
		url:     requestURL,
		body:    nil,
	}
}

type JsonRequest struct {
	client *http.Client

	method  string
	headers map[string]string
	query   url.Values
	url     string
	body    any
}

func (j *JsonRequest) Optional(enabled bool, fn func(j *JsonRequest) *JsonRequest) *JsonRequest {
	if enabled {
		return fn(j)
	}

	return j
}

func (j *JsonRequest) SetMethod(method string) *JsonRequest {
	j.method = method
	return j
}

func (j *JsonRequest) SetHeader(name, val string) *JsonRequest {
	j.headers[name] = val
	return j
}

func (j *JsonRequest) SetQuery(name string, val ...string) *JsonRequest {
	j.query[name] = val
	return j
}

func (j *JsonRequest) SetBody(body any) *JsonRequest {
	j.body = body
	return j
}

func (j *JsonRequest) Stream(sf bufio.SplitFunc) (<-chan []byte, error) {
	request, err := j.makeRequest()
	if err != nil {
		return nil, err
	}

	resp, err := j.client.Do(request)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 300 {
		respError, err := NewResponseErrorFromReader(resp.Body)
		if err != nil {
			return nil, err
		} else {
			return nil, respError.Error
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

func (j *JsonRequest) Do(body any) error {
	request, err := j.makeRequest()
	if err != nil {
		return err
	}

	resp, err := j.client.Do(request)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 300 {
		respError, err := NewResponseErrorFromReader(resp.Body)
		if err != nil {
			return err
		} else {
			return respError.Error
		}
	}

	return json.NewDecoder(resp.Body).Decode(body)
}

func (j *JsonRequest) makeRequest() (*http.Request, error) {
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
