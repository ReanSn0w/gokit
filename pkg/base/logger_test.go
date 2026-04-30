package base

import (
	"bytes"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// compile-time assertion: *StdLogger должен реализовывать Logger
var _ Logger = (*StdLogger)(nil)

func TestStdLogger_ImplementsLogger(t *testing.T) {
	// Если файл компилируется, var _ Logger = (*StdLogger)(nil) выше
	// уже гарантирует реализацию интерфейса. Здесь — дополнительная
	// runtime-проверка через присваивание.
	var l Logger = &StdLogger{}
	assert.NotNil(t, l)
}

func TestNewStdLogger_StoresLogger(t *testing.T) {
	buf := new(bytes.Buffer)
	l := log.New(buf, "", 0)

	sl := NewStdLogger(l)
	sl.Logf("ping")

	// Если переданный *log.Logger был сохранён и используется,
	// вывод окажется в buf, а не в os.Stderr.
	assert.Contains(t, buf.String(), "ping")
}

func TestStdLogger_Logf_WritesToLogger(t *testing.T) {
	buf := new(bytes.Buffer)
	l := log.New(buf, "", 0)
	sl := NewStdLogger(l)

	sl.Logf("hello %s", "world")

	assert.Contains(t, buf.String(), "hello world")
}
