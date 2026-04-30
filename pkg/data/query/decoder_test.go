package query

import (
	"errors"
	"net/url"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── вспомогательные типы ─────────────────────────────────────────────────────

// alwaysInvalidStruct реализует base.Validator через value-receiver и всегда
// возвращает ошибку. Используется для проверки того, что decode вызывает
// base.Validate и пробрасывает ошибку валидации.
type alwaysInvalidStruct struct {
	Name string
}

func (s alwaysInvalidStruct) Validate() error {
	return errors.New("validation failed")
}

// ── тесты для decode ─────────────────────────────────────────────────────────

func TestDecode_StringField(t *testing.T) {
	var s struct{ Name string }
	err := decode(url.Values{"Name": {"Alice"}}, &s)
	require.NoError(t, err)
	assert.Equal(t, "Alice", s.Name)
}

func TestDecode_TagOverride(t *testing.T) {
	var s struct {
		Name string `query:"name"`
	}
	err := decode(url.Values{"name": {"Bob"}}, &s)
	require.NoError(t, err)
	assert.Equal(t, "Bob", s.Name)
}

func TestDecode_IntField(t *testing.T) {
	var s struct{ Age int }
	err := decode(url.Values{"Age": {"42"}}, &s)
	require.NoError(t, err)
	assert.Equal(t, 42, s.Age)
}

func TestDecode_UintField(t *testing.T) {
	var s struct{ Count uint }
	err := decode(url.Values{"Count": {"7"}}, &s)
	require.NoError(t, err)
	assert.Equal(t, uint(7), s.Count)
}

func TestDecode_FloatField(t *testing.T) {
	var s struct{ Value float64 }
	err := decode(url.Values{"Value": {"3.14"}}, &s)
	require.NoError(t, err)
	assert.InDelta(t, 3.14, s.Value, 1e-9)
}

func TestDecode_BoolField(t *testing.T) {
	var s struct{ Active bool }
	err := decode(url.Values{"Active": {"true"}}, &s)
	require.NoError(t, err)
	assert.True(t, s.Active)
}

func TestDecode_SliceStringField(t *testing.T) {
	var s struct{ Tags []string }
	err := decode(url.Values{"Tags": {"a", "b", "c"}}, &s)
	require.NoError(t, err)
	assert.Equal(t, []string{"a", "b", "c"}, s.Tags)
}

func TestDecode_SliceIntField(t *testing.T) {
	var s struct{ IDs []int }
	err := decode(url.Values{"IDs": {"1", "2", "3"}}, &s)
	require.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, s.IDs)
}

func TestDecode_MissingField(t *testing.T) {
	var s struct{ Name string }
	err := decode(url.Values{}, &s)
	require.NoError(t, err)
	assert.Equal(t, "", s.Name, "отсутствующее поле должно оставаться нулевым")
}

func TestDecode_NonPointer(t *testing.T) {
	s := struct{ Name string }{}
	err := decode(url.Values{}, s)
	require.Error(t, err)
}

func TestDecode_NilPointer(t *testing.T) {
	var s *struct{ Name string }
	err := decode(url.Values{}, s)
	require.Error(t, err)
}

func TestDecode_NonStruct(t *testing.T) {
	n := 42
	err := decode(url.Values{}, &n)
	require.Error(t, err)
}

func TestDecode_ValidationError(t *testing.T) {
	s := alwaysInvalidStruct{}
	err := decode(url.Values{}, &s)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")
}

func TestDecode_UnsupportedType(t *testing.T) {
	var s struct {
		Ch chan int
	}
	// Передаём значение для поля Chan — setFieldValue будет вызван и вернёт ошибку.
	err := decode(url.Values{"Ch": {"something"}}, &s)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "неподдерживаемый")
}

// ── тесты для setFieldValue ──────────────────────────────────────────────────

func TestSetFieldValue_InvalidInt(t *testing.T) {
	var s struct{ Age int }
	field := reflect.ValueOf(&s).Elem().Field(0)
	err := setFieldValue(field, []string{"abc"})
	require.Error(t, err)
}

func TestSetFieldValue_InvalidBool(t *testing.T) {
	var s struct{ Active bool }
	field := reflect.ValueOf(&s).Elem().Field(0)
	err := setFieldValue(field, []string{"maybe"})
	require.Error(t, err)
}

func TestSetFieldValue_UnsupportedSliceElement(t *testing.T) {
	// []chan int — срез из неподдерживаемого типа элемента.
	var s struct{ Chs []chan int }
	field := reflect.ValueOf(&s).Elem().Field(0)
	err := setFieldValue(field, []string{"something"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "неподдерживаемый тип элемента среза")
}
