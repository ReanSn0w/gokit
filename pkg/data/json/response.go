package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ReanSn0w/gokit/pkg/base"
)

// errFailedToEncodeResponse is a sentinel error used inside [Response.Write]
// to prevent infinite recursion: if encoding the error response itself fails,
// no further encoding is attempted.
var (
	errFailedToEncodeResponse = fmt.Errorf("failed to encode response")
)

// NewResponse creates a *Response[T] whose fields depend on the dynamic type
// of data:
//   - [base.ErrorsMap]: Success = false, Errors = data
//   - error:            Success = false, Message = err.Error()
//   - anything else:    Success = true,  Data = data
func NewResponse[T any](data T) *Response[T] {
	switch t := any(data).(type) {
	case base.ErrorsMap:
		return &Response[T]{
			Success: false,
			Errors:  t,
		}
	case error:
		return &Response[T]{
			Success: false,
			Message: t.Error(),
		}
	default:
		return &Response[T]{
			Success: true,
			Data:    data,
		}
	}
}

// NewPlainResponse creates a successful [Response][any] that carries only a
// formatted text message and no payload data. The format string and args are
// processed by [fmt.Sprintf].
func NewPlainResponse(msg string, args ...any) *Response[any] {
	return &Response[any]{
		Success: true,
		Message: fmt.Sprintf(msg, args...),
	}
}

// Response is the universal JSON envelope for all API responses.
//
//   - Success indicates whether the request was handled successfully.
//   - Message is an optional human-readable description.
//   - Errors is an optional map of named validation or field errors.
//   - Data holds the response payload when Success is true.
type Response[T any] struct {
	// Success is true when the operation completed without errors.
	Success bool `json:"success"`

	// Message carries an optional human-readable description of the result
	// or the error that occurred.
	Message string `json:"message,omitempty"`

	// Errors holds a named collection of errors (e.g. validation failures).
	// Only populated when Success is false and the caller passed a
	// [base.ErrorsMap] to [NewResponse].
	Errors base.ErrorsMap `json:"error,omitempty"`

	// Data is the response payload. Omitted from JSON when it is the zero
	// value of T.
	Data T `json:"data,omitempty"`
}

// Write serialises the response to JSON and writes it to wr with the given
// HTTP status code. It sets the Content-Type header to "application/json".
//
// If JSON encoding fails, Write responds with 500 Internal Server Error and a
// plain error message. To avoid infinite recursion, only one level of
// fallback encoding is attempted.
func (r *Response[T]) Write(code int, wr http.ResponseWriter) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err := enc.Encode(r)

	if err != nil {
		if r.Message != errFailedToEncodeResponse.Error() {
			NewResponse(errFailedToEncodeResponse).
				Write(http.StatusInternalServerError, wr)
		}

		return
	}

	wr.Header().Set("Content-Type", "application/json")
	wr.WriteHeader(code)
	buf.WriteTo(wr)
}
