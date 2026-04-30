package view_test

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/ReanSn0w/gokit/pkg/view"
)

func TestIfTrue(t *testing.T) {
	hello := Text("Hello")(
		view.Hidden(false),
	)

	val := StringBuilder(hello) // Hello

	if val != "Hello" {
		t.Error("If(true) not working")
	}
}

func TestIfFalse(t *testing.T) {
	hello := Text("Hello")(
		view.Hidden(true),
	)

	val := StringBuilder(hello) // ""

	if val != "" {
		t.Error("If(false) not working")
	}
}

func TestReplace(t *testing.T) {
	hello := Text("Hello")(
		view.Replace(Text("World")),
	)

	val := StringBuilder(hello) // World

	if val != "World" {
		t.Error("Replace not working")
	}
}

func TestIfReplace(t *testing.T) {
	hello := Text("Hello")(
		view.If(true, view.Replace(Text("World"))),
	)

	val := StringBuilder(hello) // World

	if val != "World" {
		t.Error("IfReplace not working")
	}
}

func TestIfReplaceFalse(t *testing.T) {
	hello := Text("Hello")(
		view.If(false, view.Replace(Text("World"))),
	)

	val := StringBuilder(hello) // Hello

	if val != "Hello" {
		t.Error("IfReplaceFalse not working")
	}
}

func TestFor(t *testing.T) {
	hello := view.For(3, func(index int) view.View {
		return Text("Hello")
	})

	val := StringBuilder(hello) // HelloHelloHello

	if val != "HelloHelloHello" {
		t.Error("For not working")
	}
}

func TestForZero(t *testing.T) {
	hello := view.For(0, func(index int) view.View {
		return Text("Hello")
	})

	val := StringBuilder(hello) // ""

	if val != "" {
		t.Error("For(0) not working")
	}
}

func TestForNegative(t *testing.T) {
	hello := view.For(-1, func(index int) view.View {
		return Text("Hello")
	})

	val := StringBuilder(hello) // ""

	if val != "" {
		t.Error("For(-1) not working")
	}
}

func TestViewBuilderWithValues(t *testing.T) {
	hello := view.Group(
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

	hello := view.Group(
		Text("Hello"),
		Text("World")(
			view.Context(func(ctx context.Context) context.Context {
				return context.WithValue(ctx, &key, "value")
			}),
			view.Context(func(ctx context.Context) context.Context {
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
	hello := view.Group(
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
	hello := view.Group(
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

func Text(text string) view.Use {
	return view.External(text)
}

func StringBuilder(item view.View) string {
	buffer := new(bytes.Buffer)

	view.UnsafeBuilder(context.TODO(), item, func(ctx context.Context, i interface{}) {
		buffer.WriteString(i.(string))
	})

	return buffer.String()
}

// MARK: - panicView

// panicView — вспомогательный тип, реализующий view.View,
// который паникует с заданным значением при вызове Body.
type panicView struct{ val any }

func (p panicView) Body(_ context.Context) view.View { panic(p.val) }

// MARK: - SafeBuilder tests

func TestSafeBuilder_NormalFlow(t *testing.T) {
	tree := view.Group(
		Text("foo"),
		Text("bar"),
	)

	// Результат через SafeBuilder должен совпадать с UnsafeBuilder
	unsafeResult := StringBuilder(tree)

	buffer := new(bytes.Buffer)
	err := view.SafeBuilder(context.TODO(), tree, func(ctx context.Context, i interface{}) {
		buffer.WriteString(i.(string))
	})

	if err != nil {
		t.Fatalf("SafeBuilder вернул ошибку: %v", err)
	}
	if buffer.String() != unsafeResult {
		t.Errorf("ожидалось %q, получено %q", unsafeResult, buffer.String())
	}
}

func TestSafeBuilder_PanicWithError(t *testing.T) {
	origErr := errors.New("boom")
	v := panicView{val: origErr}

	err := view.SafeBuilder(context.TODO(), v, func(_ context.Context, _ interface{}) {})

	if err == nil {
		t.Fatal("SafeBuilder должен был вернуть ошибку, но вернул nil")
	}
	if err.Error() != "boom" {
		t.Errorf("ожидалось сообщение %q, получено %q", "boom", err.Error())
	}
}

func TestSafeBuilder_PanicWithNonError(t *testing.T) {
	v := panicView{val: "crash"}

	err := view.SafeBuilder(context.TODO(), v, func(_ context.Context, _ interface{}) {})

	if err == nil {
		t.Fatal("SafeBuilder должен был вернуть ошибку, но вернул nil")
	}
	if !strings.Contains(err.Error(), "crash") {
		t.Errorf("ожидалось сообщение содержащее %q, получено %q", "crash", err.Error())
	}
}

func TestSafeBuilder_NilView(t *testing.T) {
	err := view.SafeBuilder(context.TODO(), nil, func(_ context.Context, _ interface{}) {})

	if err != nil {
		t.Fatalf("SafeBuilder с nil-view должен вернуть nil, получено: %v", err)
	}
}

// MARK: - Closure & External tests

func TestClosure_CallsBuilder(t *testing.T) {
	v := view.Closure(func(ctx context.Context) view.View {
		return Text("from closure")
	})

	result := StringBuilder(v)

	if result != "from closure" {
		t.Errorf("ожидалось %q, получено %q", "from closure", result)
	}
}

func TestExternal_ValueIsPassedToCallback(t *testing.T) {
	v := view.External(42)

	var received interface{}
	view.UnsafeBuilder(context.TODO(), v, func(_ context.Context, i interface{}) {
		received = i
	})

	if received != 42 {
		t.Errorf("ожидалось 42, получено %v", received)
	}
}
