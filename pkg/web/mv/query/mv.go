package query

import (
	"context"
	"net/http"

	"github.com/ReanSn0w/gokit/pkg/web"
)

var (
	queryDecoderCtxKey = &queryDecoderCtx{}
)

type queryDecoderCtx struct{}

// Получает данные из Query запроса, сохраняя их в структуру.
func Decoder[T any](h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data T
		err := Decode(r.URL.Query(), &data)

		if err != nil {
			web.NewResponse(err).Write(http.StatusBadRequest, w)
			return
		}

		ctx := context.WithValue(r.Context(), queryDecoderCtxKey, &data)
		r = r.WithContext(ctx)

		h.ServeHTTP(w, r)
	})
}

// Получает данные из Query запроса, сохраняя их в структуру.
func Get[T any](ctx context.Context) *T {
	if ctx == nil {
		return nil
	}

	data, ok := ctx.Value(queryDecoderCtxKey).(*T)
	if !ok {
		return nil
	}

	return data
}
