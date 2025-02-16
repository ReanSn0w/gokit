package html

import (
	"context"
	"fmt"

	"github.com/ReanSn0w/gokit/pkg/composer"
)

func New(tag string, content ...composer.View) composer.Use {
	return composer.New(
		&node{
			tag:     tag,
			content: content,
		},
	)
}

func Inline(tag string) composer.Use {
	return composer.New(
		&node{
			inline: true,
			tag:    tag,
		},
	)
}

func Text(str string, args ...interface{}) composer.Use {
	return composer.External(fmt.Sprintf(str, args...))
}

type node struct {
	inline  bool
	tag     string
	content []composer.View
}

func (n *node) Body(ctx context.Context) composer.View {
	if n.inline {
		return composer.External("<" + n.tag + buildAttributesString(ctx) + "/>")
	}

	return composer.Group(
		composer.External("<"+n.tag+buildAttributesString(ctx)+">"),
		composer.Group(n.content...),
		composer.External("</"+n.tag+">"),
	)
}
