package composer_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/ReanSn0w/gokit/pkg/composer"
)

func TestIfTrue(t *testing.T) {
	hello := Text("Hello")(
		composer.Hidden(false),
	)

	val := StringBuilder(hello) // Hello

	if val != "Hello" {
		t.Error("If(true) not working")
	}
}

func TestIfFalse(t *testing.T) {
	hello := Text("Hello")(
		composer.Hidden(true),
	)

	val := StringBuilder(hello) // ""

	if val != "" {
		t.Error("If(false) not working")
	}
}

func TestReplace(t *testing.T) {
	hello := Text("Hello")(
		composer.Replace(Text("World")),
	)

	val := StringBuilder(hello) // World

	if val != "World" {
		t.Error("Replace not working")
	}
}

func TestIfReplace(t *testing.T) {
	hello := Text("Hello")(
		composer.If(true, composer.Replace(Text("World"))),
	)

	val := StringBuilder(hello) // World

	if val != "World" {
		t.Error("IfReplace not working")
	}
}

func TestIfReplaceFalse(t *testing.T) {
	hello := Text("Hello")(
		composer.If(false, composer.Replace(Text("World"))),
	)

	val := StringBuilder(hello) // Hello

	if val != "Hello" {
		t.Error("IfReplaceFalse not working")
	}
}

func TestFor(t *testing.T) {
	hello := composer.For(3, func(index int) composer.View {
		return Text("Hello")
	})

	val := StringBuilder(hello) // HelloHelloHello

	if val != "HelloHelloHello" {
		t.Error("For not working")
	}
}

func TestForZero(t *testing.T) {
	hello := composer.For(0, func(index int) composer.View {
		return Text("Hello")
	})

	val := StringBuilder(hello) // ""

	if val != "" {
		t.Error("For(0) not working")
	}
}

func TestForNegative(t *testing.T) {
	hello := composer.For(-1, func(index int) composer.View {
		return Text("Hello")
	})

	val := StringBuilder(hello) // ""

	if val != "" {
		t.Error("For(-1) not working")
	}
}

func TestViewBuilderWithValues(t *testing.T) {
	hello := composer.Group(
		Text("Hello"),
		Text("World"),
	)

	val := StringBuilder(hello) // HelloWorld

	if val != "HelloWorld" {
		t.Error("ViewBuilder not working")
	}
}

func TestContext(t *testing.T) {
	res := false
	key := "k"

	hello := composer.Group(
		Text("Hello"),
		Text("World")(
			composer.Context(func(ctx context.Context) context.Context {
				return context.WithValue(ctx, &key, "value")
			}),
			composer.Context(func(ctx context.Context) context.Context {
				val, ok := ctx.Value(&key).(string)

				if val == "value" && ok {
					res = true
				}

				return ctx
			}),
		),
	)

	_ = StringBuilder(hello) // HelloWorld

	if !res {
		t.Error("Context not working")
	}
}

func TestNil(t *testing.T) {
	hello := composer.Group(
		Text("Hello"),
		nil,
		Text("World"),
	)

	val := StringBuilder(hello) // HelloWorld

	if val != "HelloWorld" {
		t.Error("nil not working")
	}
}

func TestNilModificator(t *testing.T) {
	hello := composer.Group(
		Text("Hello"),
		Text("World")(
			nil,
		),
	)

	val := StringBuilder(hello) // HelloWorld

	if val != "HelloWorld" {
		t.Error("nil modificator not working")
	}
}

// MARK: - TextBuilder

func Text(text string) composer.Use {
	return composer.External(text)
}

func StringBuilder(item composer.View) string {
	buffer := new(bytes.Buffer)

	composer.UnsafeBuilder(context.TODO(), item, func(ctx context.Context, i interface{}) {
		buffer.WriteString(i.(string))
	})

	return buffer.String()
}
