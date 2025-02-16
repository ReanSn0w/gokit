package html

import (
	"context"
	"strings"

	"github.com/ReanSn0w/gokit/pkg/composer"
)

var (
	attributesCtxKey = &attributesCtx{}
)

type attributesCtx struct{}

type attributeMap map[string][]string

func buildAttributesString(ctx context.Context) string {
	attrs, ok := ctx.Value(attributesCtxKey).(attributeMap)
	if !ok {
		return ""
	}

	result := ""
	for key, values := range attrs {
		if values == nil {
			result += key + " "
		} else {
			result += " " + key + "=\"" + strings.Join(values, " ") + "\""
		}
	}

	return result
}

func PrepareAttrubuteMap(prepare func(attributeMap)) composer.With {
	return composer.Context(func(ctx context.Context) context.Context {
		attrs, ok := ctx.Value(attributesCtxKey).(attributeMap)
		if !ok {
			attrs = make(attributeMap)
		}

		prepare(attrs)
		ctx = context.WithValue(ctx, attributesCtxKey, attrs)
		return ctx
	})
}

func AddAttribute(key string, values ...string) composer.With {
	return PrepareAttrubuteMap(func(am attributeMap) {
		vals := am[key]
		vals = append(vals, values...)
		am[key] = vals
	})
}

func SetAttribute(key string, values ...string) composer.With {
	return PrepareAttrubuteMap(func(am attributeMap) {
		am[key] = values
	})
}

func UnsetAttribute(key string, values ...string) composer.With {
	return PrepareAttrubuteMap(func(am attributeMap) {
		for _, value := range values {
			vals := am[key]
			for index, val := range vals {
				if val == value {
					vals = append(vals[:index], vals[index+1:]...)
				}
			}
			am[key] = vals
		}
	})
}

func DropAttribute(key string) composer.With {
	return PrepareAttrubuteMap(func(am attributeMap) {
		delete(am, key)
	})
}
