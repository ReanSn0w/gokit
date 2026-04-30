package query

import (
	"context"
	"net/http"

	"github.com/ReanSn0w/gokit/pkg/base"
	"github.com/ReanSn0w/gokit/pkg/data/json"
)

const (
	// queryDecoderKey — ключ контекста для хранения декодированных
	// query-параметров, которые укладываются middleware [Decoder].
	queryDecoderKey = "query_decoder"
)

// Decoder — HTTP middleware, декодирующий query-параметры запроса в структуру
// типа T с помощью [decode].
//
// При ошибке декодирования или валидации отвечает клиенту кодом
// 400 Bad Request и JSON-телом с описанием ошибки, после чего прерывает
// цепочку обработчиков.
//
// При успешном декодировании *T помещается в контекст запроса под ключом
// [queryDecoderKey] и запрос передаётся следующему обработчику h.
// Получить значение из контекста можно с помощью [Get].
func Decoder[T any](h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data T
		err := decode(r.URL.Query(), &data)

		if err != nil {
			json.NewResponse(err).
				Write(http.StatusBadRequest, w)
			return
		}

		ctx := context.WithValue(r.Context(), queryDecoderKey, &data)
		r = r.WithContext(ctx)

		h.ServeHTTP(w, r)
	})
}

// Get извлекает ранее декодированные query-параметры *T из контекста ctx.
// Возвращает nil, если данные не найдены — например, если middleware [Decoder]
// не был применён к цепочке обработчиков или был использован другой тип T.
func Get[T any](ctx context.Context) *T {
	data, ok := ctx.Value(queryDecoderKey).(*T)
	if !ok {
		return nil
	}

	return data
}

// Set кладёт значение val типа T в контекст ctx напрямую, минуя
// HTTP-декодирование. Полезно в тестах или при ручном формировании контекста.
//
// Перед сохранением выполняет валидацию через [base.Validate]. Если валидация
// не прошла — возвращает ошибку, контекст остаётся неизменным.
// При успехе возвращает новый контекст с сохранённым *T.
func Set[T any](ctx context.Context, val T) (context.Context, error) {
	err := base.Validate(val)
	if err != nil {
		return ctx, err
	}

	ctx = context.WithValue(ctx, queryDecoderKey, &val)
	return ctx, nil
}
