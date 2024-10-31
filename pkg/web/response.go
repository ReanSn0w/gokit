package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"git.papkovda.ru/library/gokit/pkg/tool"
	"github.com/go-pkgz/lgr"
)

var (
	errFailedToEncodeResponse = errors.New("failed to encode response")
)

// NewResponse creates a new response
func NewResponse[T any](data T) *Response[T] {
	switch t := any(data).(type) {
	case tool.ErrorsMap:
		return &Response[T]{
			Success: false,
			Error:   t,
		}
	case error:
		return &Response[T]{
			Success: false,
			Error: tool.ErrorsMap{
				"error": t,
			},
		}
	default:
		return &Response[T]{
			Success: true,
			Data:    data,
		}
	}
}

// NewPlainResponse creates Response[any]
// without any data and send it to http.ResponseWriter
func NewPlainResponse(code int, wr http.ResponseWriter) {
	resp := Response[any]{
		Success: true,
	}

	resp.Write(code, wr)
}

// NewStreamResponse makes stream with Response[T]
// func returns chan <- Response[T] that sends values to response
// response end after you close channel
func NewStreamResponse[T any](code int, separator []byte, wr http.ResponseWriter) chan<- Response[T] {
	con := http.NewResponseController(wr)
	{
		wr.Header().Add("Content-Type", "application/octet-stream")
		wr.WriteHeader(http.StatusOK)
	}

	ch := make(chan Response[T])

	go func() {
		for msg := range ch {
			msg.write(wr)

			wr.Write(separator)

			err := con.Flush()
			if err != nil {
				lgr.Default().Logf("[ERROR] flush content error: %v", err)
			}
		}
	}()

	return ch
}

// NewResponseErrorFromReader
func NewResponseErrorFromReader(r io.Reader) (*Response[GenericJsonError], error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return &Response[GenericJsonError]{
		Success: false,
		Error:   GenericJsonError(data),
	}, err
}

type GenericJsonError []byte

func (g GenericJsonError) Error() string {
	return fmt.Sprintf("generic response error: %s", string(g))
}

// Response is a generic response structure
type (
	Response[T any] struct {
		Success bool  `json:"success"`
		Error   error `json:"message,omitempty"`
		Data    T     `json:"data,omitempty"`
	}
)

// Write writes the response to the writer
func (r *Response[T]) Write(code int, wr http.ResponseWriter) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err := enc.Encode(r)
	if err != nil {
		if r.Error != errFailedToEncodeResponse {
			NewResponse(errFailedToEncodeResponse).
				Write(http.StatusInternalServerError, wr)
		}

		return
	}

	wr.Header().Set("Content-Type", "application/json")
	wr.WriteHeader(code)
	buf.WriteTo(wr)
}

func (r *Response[T]) write(wr http.ResponseWriter) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err := enc.Encode(r)
	if err != nil {
		if r.Error != errFailedToEncodeResponse {
			NewResponse(errFailedToEncodeResponse).
				write(wr)
		}

		return
	}

	buf.WriteTo(wr)
}
