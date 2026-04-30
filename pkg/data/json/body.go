// Package json provides HTTP middleware and utilities for JSON-based API
// communication. It includes JSON request body decoding with validation,
// structured response encoding, and a fluent HTTP client for JSON APIs.
package json

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ReanSn0w/gokit/pkg/base"
)

// jsonBodyDecoderKey is the context key under which the decoded request body
// is stored by the [Decoder] middleware and retrieved by [GetBody].
const (
	jsonBodyDecoderKey = "json_body_decoder"
)

// Decoder is an HTTP middleware that decodes the JSON request body into a
// value of type T and stores a pointer to it in the request context.
//
// If decoding fails (malformed JSON), Decoder writes a 400 Bad Request
// response with a JSON error body and stops the handler chain.
// If the decoded value fails validation via [base.Validate], the same 400
// response is returned.
//
// On success, the decoded *T is placed in the context under
// [jsonBodyDecoderKey] and the next handler is called. The value can be
// retrieved downstream with [GetBody].
func Decoder[T any](h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data T
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			NewResponse(err).Write(http.StatusBadRequest, w)
			return
		}

		err = base.Validate(data)
		if err != nil {
			NewResponse(err).Write(http.StatusBadRequest, w)
			return
		}

		ctx := context.WithValue(r.Context(), jsonBodyDecoderKey, &data)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetBody retrieves the decoded request body of type *T from ctx.
//
// It panics if no body was stored in ctx (e.g. the [Decoder] middleware was
// not applied to the handler chain). Use [SetBody] to populate the context
// in tests or outside of an HTTP handler.
func GetBody[T any](ctx context.Context) *T {
	return ctx.Value(jsonBodyDecoderKey).(*T)
}

// SetBody validates val via [base.Validate] and, on success, stores a pointer
// to val in ctx under [jsonBodyDecoderKey]. It is the programmatic equivalent
// of what [Decoder] does for an incoming HTTP request — useful in tests or
// when constructing a context without going through HTTP decoding.
//
// Returns an error if validation fails; in that case ctx is returned unchanged.
func SetBody[T any](ctx context.Context, val T) (context.Context, error) {
	err := base.Validate(val)
	if err != nil {
		return ctx, err
	}

	ctx = context.WithValue(ctx, jsonBodyDecoderKey, &val)
	return ctx, nil
}
