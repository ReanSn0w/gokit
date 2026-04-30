package view

import (
	"context"
	"fmt"
	"sync"
)

type (
	// With это функция, которая принимает View
	// и возвращает View. Данный тип служит для модификации
	// View в своем теле.
	With func(View) View

	// Use это функция, которая применяет модификации ко View.
	Use func(args ...With) View

	// View является еденицей в системе построения gew
	//
	// В стандартной поставке в вашем распоряжжении находится
	// несколько таких элементов:
	// - Group для группировки нескольких view в один
	// - Closure выполняющий действие над view
	// - External для обертки внешних значений в View
	View interface {
		Body(context.Context) View
	}

	// Builder — тип функции для построения View-дерева.
	// Получает контекст, корневой View и коллбэк ext, который вызывается
	// для каждого терминального *external-узла. Используется как тип
	// для кастомных реализаций обхода дерева.
	Builder func(ctx context.Context, view View, ext func(context.Context, interface{})) error
)

// Group представляет из себя набор компонентов
// для пайплайна построения
func Group(elements ...View) Use {
	return New(group(elements))
}

// External обертка для внешних типов
// в пайплайне построения view
func External(content interface{}) Use {
	return New(&external{
		Content: content,
	})
}

// Closure возвращает View, функции переданной в нее
// в качестве аргумента.
//
// Через данную функцию реализованы такие функции как:
// For
func Closure(builder func(context.Context) View) Use {
	return New(&closure{
		builder: builder,
	})
}

// MARK: - Модификаци
//
// Код написанный ниже отвечает за модификации View

// New Возвращает функцию, которая
// соответствует интерфейсу View.
//
// Она предназначена для удоного применения модификаторов к View.
func New(view View) Use {
	return func(args ...With) View {
		for i := len(args) - 1; i >= 0; i-- {
			if args[i] == nil {
				continue
			}

			view = args[i](view)
		}

		return view
	}
}

// Body добавляет соответствие интерфейсу View для Applyer
func (a Use) Body(ctx context.Context) View {
	return a()
}

// Context - позволяет изменить контекст построения View.
func Context(prepare func(ctx context.Context) context.Context) With {
	return func(view View) View {
		return &contexted{
			content: view,
			prepare: prepare,
		}
	}
}

// If возвращает модификатор, которые выполнит переданные в него
// модификаторы в случае, если condition == true
func If(condition bool, modificators ...With) With {
	if condition {
		return func(v View) View {
			for _, modificator := range modificators {
				v = modificator(v)
			}

			return v
		}
	}

	return func(v View) View {
		return v
	}
}

// For возвращает View, который содержит count элементов,
// созданных функцией builder.
func For(count int, builder func(int) View) Use {
	// Исключение для нулевого количества элементов.
	// Так же оно сработает в случае если count < 0.
	if count < 1 {
		return Group()
	}

	return Closure(func(ctx context.Context) View {
		items := make([]View, 0, count)

		for index := 0; index < int(count); index++ {
			items = append(items, builder(index))
		}

		return Group(items...)
	})
}

// Hidden - модификатор позволяет скрыть элемент View из построения.
func Hidden(condition bool) With {
	return func(view View) View {
		if !condition {
			return view
		}

		return nil
	}
}

// Replace - заменяет элемент View на другой.
func Replace(view View) With {
	return func(_ View) View {
		return view
	}
}

// MARK: Встроенные типы
// Код описанный ниже служит для работы базовых
// функций построения пайплайна и хранения компонентов
// его модификации

type (
	// group — внутренний тип, представляющий срез View-элементов.
	// При построении UnsafeBuilder обходит элементы группы последовательно.
	group []View

	// contexted — внутренний тип, оборачивающий View и подменяющий контекст
	// перед вызовом дочернего Body. Функция prepare вызывается в UnsafeBuilder
	// до передачи контекста вниз по дереву.
	contexted struct {
		prepare func(context.Context) context.Context
		content View
	}

	// external — внутренний терминальный узел дерева, хранящий внешнее значение.
	// При обходе UnsafeBuilder вызывает ext-коллбэк с этим значением.
	external struct {
		Content interface{}
	}

	// closure — внутренний тип, чья реализация Body делегирована
	// переданной функции-строителю. Используется функцией Closure и For.
	closure struct {
		builder func(context.Context) View
	}
)

// Body возвращает дочерний content без изменений.
// Применение prepare к контексту выполняется в UnsafeBuilder
// до вызова этого метода, поэтому здесь оно не повторяется.
func (cv *contexted) Body(context context.Context) View {
	return cv.content
}

// Body возвращает nil; реальный последовательный обход элементов
// группы выполняет UnsafeBuilder через type-switch.
func (g group) Body(ctx context.Context) View {
	return nil
}

// Body возвращает nil; реальный вызов ext-коллбэка
// со значением Content выполняет UnsafeBuilder через type-switch.
func (e *external) Body(ctx context.Context) View {
	return nil
}

// Body вызывает хранимую функцию builder с переданным контекстом
// и возвращает полученный View.
func (cv *closure) Body(context context.Context) View {
	return cv.builder(context)
}

// UnsafeBuilder рекурсивно обходит дерево View и для каждого
// *external-узла вызывает коллбэк ext(ctx, value).
//
// Обрабатывает следующие случаи:
//   - nil — немедленный возврат без действий;
//   - group — последовательный рекурсивный обход всех элементов;
//   - *contexted — применяет prepare к контексту, затем рекурсивно обходит content;
//   - *external — вызывает ext(ctx, Content);
//   - любой другой View — рекурсивно обходит результат Body(ctx).
//
// Не защищён от паник; для безопасного обхода используйте SafeBuilder.
func UnsafeBuilder(ctx context.Context, view View, ext func(context.Context, interface{})) error {
	switch v := view.(type) {
	case nil:
		return nil
	case group:
		for i := range v {
			UnsafeBuilder(ctx, v[i], ext)
		}
	case *contexted:
		newCtx := v.prepare(ctx)
		UnsafeBuilder(newCtx, v.Body(newCtx), ext)
	case *external:
		ext(ctx, v.Content)
	default:
		UnsafeBuilder(ctx, v.Body(ctx), ext)
	}

	return nil
}

// SafeBuilder выполняет то же самое что UnsafeBuilder, но запускает обход
// в отдельной горутине с defer/recover.
//
// Если в процессе обхода произошла паника:
//   - со значением типа error — возвращает его напрямую;
//   - с любым другим значением — оборачивает в "unknown error: <значение>".
//
// Блокируется до завершения обхода и возвращает nil при успешном выполнении.
func SafeBuilder(ctx context.Context, view View, ext func(ctx context.Context, i interface{})) error {
	var externalErr error
	wg := new(sync.WaitGroup)
	wg.Add(1)

	go func() {
		defer func() {
			if data := recover(); data != nil {
				if err, ok := data.(error); ok {
					externalErr = err
				} else {
					externalErr = fmt.Errorf("unknown error: %v", data)
				}
			}

			wg.Done()
		}()

		_ = UnsafeBuilder(ctx, view, ext)
	}()

	wg.Wait()
	return externalErr
}
