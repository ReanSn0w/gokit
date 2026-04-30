package web

import (
	"errors"
	"net/http"
	"net/http/httputil"

	"github.com/ReanSn0w/gokit/pkg/data/json"
	"github.com/go-pkgz/lgr"
)

// Ping returns a middleware that intercepts GET /ping requests and responds
// with HTTP 200 and the plain-text body "pong". All other requests are passed
// unchanged to the next handler in the chain.
func Ping() func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/ping" {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("pong"))
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}

// DebugRequest returns a middleware that logs the full incoming HTTP request,
// including the request body for POST, PUT, and DELETE methods.
// When enabled is false the original handler is returned as-is, so there is
// no runtime overhead in production. Logging is performed via the provided
// lgr.L implementation.
func DebugRequest(enabled bool, log lgr.L) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		if !enabled {
			return h
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			dump, err := httputil.DumpRequest(
				r, r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodDelete,
			)

			if err == nil {
				log.Logf("[DEBUG] debug request:\n", string(dump))
			} else {
				log.Logf("[ERROR] debug request error: %v", err)
			}

			h.ServeHTTP(w, r)
		})
	}
}

// APIKey returns a middleware that authenticates requests by comparing the
// value of the HTTP request header named header against the expected key.
// If the header is missing or its value does not match, the middleware
// responds with HTTP 401 Unauthorized and a JSON error body.
// Matching requests are forwarded to the next handler.
func APIKey(header, key string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if requestKey := r.Header.Get(header); requestKey != key {
				json.NewResponse(errors.New("invalid api key")).
					Write(http.StatusUnauthorized, w)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
