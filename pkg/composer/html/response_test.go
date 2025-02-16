package html_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/ReanSn0w/gokit/pkg/composer"
	"github.com/ReanSn0w/gokit/pkg/composer/html"
	"github.com/ReanSn0w/gokit/pkg/composer/html/attr"
	"github.com/ReanSn0w/gokit/pkg/composer/html/tag"
	"github.com/stretchr/testify/assert"
)

func TestResponse(t *testing.T) {
	cases := []struct {
		Name     string
		Content  composer.View
		Output   string
		SafeMode bool
		HasError bool
	}{
		{
			Name:    "Empty",
			Content: nil,
			Output:  "",
		},
		{
			Name:    "Plain Text",
			Content: tag.Text("Hello, world!"),
			Output:  "Hello, world!",
		},
		{
			Name:    "One Element",
			Content: tag.P(tag.Text("Hello, world!")),
			Output:  "<p>Hello, world!</p>",
		},
		{
			Name: "Three",
			Content: tag.Html(
				tag.Body(
					tag.H1(
						tag.Text("Hello, world!"),
					),
				),
			),
			Output: "<html><body><h1>Hello, world!</h1></body></html>",
		},
		{
			Name: "Attributed Element",
			Content: tag.P(tag.Text("hello, world!"))(
				attr.Class.Set("one", "two"),
			),
			Output: "<p class=\"one two\">hello, world!</p>",
		},
		{
			Name:     "Invalid External Type",
			Content:  composer.External(0),
			SafeMode: true,
			HasError: true,
		},
	}

	for _, c := range cases {
		var err error
		buffer := new(bytes.Buffer)
		if c.SafeMode {
			err = html.Builder(context.TODO(), composer.SafeBuilder, buffer, c.Content)
		} else {
			err = html.Builder(context.TODO(), composer.UnsafeBuilder, buffer, c.Content)
		}
		if c.HasError {
			assert.Error(t, err, "case $s", c.Name)
			continue
		}

		if err != nil {
			assert.NoError(t, err, "case %s", c.Name)
			continue
		}

		assert.Equal(t, c.Output, buffer.String(), "case %s", c.Name)
	}
}

func BenchmarkResponse(b *testing.B) {
	cases := []struct {
		Name    string
		Content composer.View
	}{
		{
			Name:    "Empty",
			Content: nil,
		},
		{
			Name:    "Plain Text",
			Content: tag.Text("Hello, world!"),
		},
		{
			Name:    "One Element",
			Content: tag.P(tag.Text("Hello, world!")),
		},
		{
			Name: "Three",
			Content: tag.Html(
				tag.Body(
					tag.H1(
						tag.Text("Hello, world!"),
					),
				),
			),
		},
		{
			Name: "Attributed Element",
			Content: tag.P(tag.Text("hello, world!"))(
				attr.Class.Set("one", "two"),
			),
		},
		{
			Name: "Large Page",
			Content: composer.Group(
				tag.Doctype(),
				tag.Html(
					tag.Head(
						tag.Title(tag.Text("Тестовая страница")),
						tag.Meta()(attr.Charset.Set("utf-8")),
						tag.Link()(attr.Src.Set("https://somesite.com/styles/style.css")),
					),
					tag.Body(
						tag.Header(
							tag.P(tag.Text("Название сайта"))(attr.Class.Set("title")),
							tag.P(tag.Text("Слоган")),
						),
						tag.Section(
							composer.For(10, func(i int) composer.View {
								return tag.Article(
									tag.H3(tag.Text("Название статьи %v", i)),
									tag.P(tag.Text("Краткое содержание статьи")),
								)
							}),
						)(attr.Class.Set("articles")),
						tag.Footer(
							tag.P(tag.Text("Название сайта"))(attr.Class.Set("title")),
							tag.P(tag.Text("Слоган")),
						),
					),
				),
			),
		},
		{
			Name: "Multiple Pages",
			Content: composer.Group(
				tag.Doctype(),
				tag.Html(
					tag.Head(
						tag.Title(tag.Text("Тестовая страница")),
						tag.Meta()(attr.Charset.Set("utf-8")),
						tag.Link()(attr.Src.Set("https://somesite.com/styles/style.css")),
					),
					tag.Body(
						tag.Header(
							tag.P(tag.Text("Название сайта"))(attr.Class.Set("title")),
							tag.P(tag.Text("Слоган")),
						),
						tag.Section(
							composer.For(10, func(i int) composer.View {
								return tag.Article(
									tag.H3(tag.Text("Название статьи %v", i)),
									tag.P(tag.Text("Краткое содержание статьи")),
								)
							}),
						)(attr.Class.Set("articles")),
						tag.Footer(
							tag.P(tag.Text("Название сайта"))(attr.Class.Set("title")),
							tag.P(tag.Text("Слоган")),
						),
					),
				),
				tag.Html(
					tag.Head(
						tag.Title(tag.Text("Тестовая страница")),
						tag.Meta()(attr.Charset.Set("utf-8")),
						tag.Link()(attr.Src.Set("https://somesite.com/styles/style.css")),
					),
					tag.Body(
						tag.Header(
							tag.P(tag.Text("Название сайта"))(attr.Class.Set("title")),
							tag.P(tag.Text("Слоган")),
						),
						tag.Section(
							composer.For(10, func(i int) composer.View {
								return tag.Article(
									tag.H3(tag.Text("Название статьи %v", i)),
									tag.P(tag.Text("Краткое содержание статьи")),
								)
							}),
						)(attr.Class.Set("articles")),
						tag.Footer(
							tag.P(tag.Text("Название сайта"))(attr.Class.Set("title")),
							tag.P(tag.Text("Слоган")),
						),
					),
				),
				tag.Html(
					tag.Head(
						tag.Title(tag.Text("Тестовая страница")),
						tag.Meta()(attr.Charset.Set("utf-8")),
						tag.Link()(attr.Src.Set("https://somesite.com/styles/style.css")),
					),
					tag.Body(
						tag.Header(
							tag.P(tag.Text("Название сайта"))(attr.Class.Set("title")),
							tag.P(tag.Text("Слоган")),
						),
						tag.Section(
							composer.For(10, func(i int) composer.View {
								return tag.Article(
									tag.H3(tag.Text("Название статьи %v", i)),
									tag.P(tag.Text("Краткое содержание статьи")),
								)
							}),
						)(attr.Class.Set("articles")),
						tag.Footer(
							tag.P(tag.Text("Название сайта"))(attr.Class.Set("title")),
							tag.P(tag.Text("Слоган")),
						),
					),
				),
				tag.Html(
					tag.Head(
						tag.Title(tag.Text("Тестовая страница")),
						tag.Meta()(attr.Charset.Set("utf-8")),
						tag.Link()(attr.Src.Set("https://somesite.com/styles/style.css")),
					),
					tag.Body(
						tag.Header(
							tag.P(tag.Text("Название сайта"))(attr.Class.Set("title")),
							tag.P(tag.Text("Слоган")),
						),
						tag.Section(
							composer.For(10, func(i int) composer.View {
								return tag.Article(
									tag.H3(tag.Text("Название статьи %v", i)),
									tag.P(tag.Text("Краткое содержание статьи")),
								)
							}),
						)(attr.Class.Set("articles")),
						tag.Footer(
							tag.P(tag.Text("Название сайта"))(attr.Class.Set("title")),
							tag.P(tag.Text("Слоган")),
						),
					),
				),
				tag.Html(
					tag.Head(
						tag.Title(tag.Text("Тестовая страница")),
						tag.Meta()(attr.Charset.Set("utf-8")),
						tag.Link()(attr.Src.Set("https://somesite.com/styles/style.css")),
					),
					tag.Body(
						tag.Header(
							tag.P(tag.Text("Название сайта"))(attr.Class.Set("title")),
							tag.P(tag.Text("Слоган")),
						),
						tag.Section(
							composer.For(10, func(i int) composer.View {
								return tag.Article(
									tag.H3(tag.Text("Название статьи %v", i)),
									tag.P(tag.Text("Краткое содержание статьи")),
								)
							}),
						)(attr.Class.Set("articles")),
						tag.Footer(
							tag.P(tag.Text("Название сайта"))(attr.Class.Set("title")),
							tag.P(tag.Text("Слоган")),
						),
					),
				),
				tag.Html(
					tag.Head(
						tag.Title(tag.Text("Тестовая страница")),
						tag.Meta()(attr.Charset.Set("utf-8")),
						tag.Link()(attr.Src.Set("https://somesite.com/styles/style.css")),
					),
					tag.Body(
						tag.Header(
							tag.P(tag.Text("Название сайта"))(attr.Class.Set("title")),
							tag.P(tag.Text("Слоган")),
						),
						tag.Section(
							composer.For(10, func(i int) composer.View {
								return tag.Article(
									tag.H3(tag.Text("Название статьи %v", i)),
									tag.P(tag.Text("Краткое содержание статьи")),
								)
							}),
						)(attr.Class.Set("articles")),
						tag.Footer(
							tag.P(tag.Text("Название сайта"))(attr.Class.Set("title")),
							tag.P(tag.Text("Слоган")),
						),
					),
				),
			),
		},
	}

	for _, c := range cases {
		b.Run(c.Name+"_SafeMode_Off", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				buffer := new(bytes.Buffer)
				err := html.Builder(context.TODO(), composer.UnsafeBuilder, buffer, c.Content)
				if err != nil {
					b.Fatalf("unexpected error: %v", err)
				}
			}
		})

		b.Run(c.Name+"_SafeMode_On", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				buffer := new(bytes.Buffer)
				err := html.Builder(context.TODO(), composer.SafeBuilder, buffer, c.Content)
				if err != nil {
					b.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}
