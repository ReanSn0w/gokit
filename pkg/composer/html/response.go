package html

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/ReanSn0w/gokit/pkg/composer"
	"github.com/ReanSn0w/gokit/pkg/web"
)

func NewHTMLResponse(ctx context.Context, builder composer.Builder, code int, w http.ResponseWriter, content composer.View) {
	buffer := new(bytes.Buffer)
	err := Builder(ctx, builder, buffer, content)

	if err != nil {
		web.NewResponse(err).
			Write(http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(code)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buffer.WriteTo(buffer)
}

func Builder(ctx context.Context, builder composer.Builder, wr io.Writer, content composer.View) error {
	return builder(ctx, content, func(ctx context.Context, i interface{}) {
		bytes, ok := i.(string)
		if !ok {
			panic(fmt.Errorf("invalid external type: %v", i))
		}

		wr.Write([]byte(bytes))
	})
}
