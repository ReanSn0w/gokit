package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"git.papkovda.ru/library/gokit/pkg/tool"
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

func NewResponseErrorFromReader(r io.Reader) (*Response[map[string]interface{}], error) {
	data := make(GenericJsonError)
	err := json.NewDecoder(r).Decode(&data)

	return &Response[map[string]interface{}]{
		Success: false,
		Error:   data,
	}, err
}

type GenericJsonError map[string]interface{}

func (g GenericJsonError) Error() string {
	return fmt.Errorf("generic response error: %v", g).Error()
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
