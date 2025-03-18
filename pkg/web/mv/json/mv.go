package json

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ReanSn0w/gokit/pkg/web"
)

var (
	jsonDecoderCtxKey = &jsonDecoderCtx{}
)

type jsonDecoderCtx struct{}

type Validate interface {
	Validate() error
}

func Decoder[T any](h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data T
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			web.NewResponse(err).Write(http.StatusBadRequest, w)
			return
		}

		if validate, ok := any(&data).(Validate); ok {
			if err := validate.Validate(); err != nil {
				web.NewResponse(err).Write(http.StatusBadRequest, w)
				return
			}
		}

		if validate, ok := any(data).(Validate); ok {
			if err := validate.Validate(); err != nil {
				web.NewResponse(err).Write(http.StatusBadRequest, w)
				return
			}
		}

		ctx := context.WithValue(r.Context(), jsonDecoderCtxKey, &data)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Get[T any](ctx context.Context) *T {
	return ctx.Value(jsonDecoderCtxKey).(*T)
}
