package web

import (
	"net/http"
	"net/http/httputil"

	"github.com/go-pkgz/lgr"
)

// Ping writes pong on /ping request
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

// DebugLogger prints debug log to func
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
