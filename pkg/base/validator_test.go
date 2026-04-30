package base

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── вспомогательные типы ─────────────────────────────────────────────────────

// nonValidator не реализует Validator.
type nonValidator struct{}

// valueValidator реализует Validator через value-receiver.
type valueValidator struct {
	err error
}

func (v valueValidator) Validate() error {
	return v.err
}

// ptrValidator реализует Validator через pointer-receiver.
type ptrValidator struct {
	err error
}

func (v *ptrValidator) Validate() error {
	return v.err
}

// ── тесты ────────────────────────────────────────────────────────────────────

func TestValidate_NonValidator(t *testing.T) {
	err := Validate(nonValidator{})
	assert.NoError(t, err)
}

func TestValidate_ValueReceiver_Valid(t *testing.T) {
	err := Validate(valueValidator{err: nil})
	assert.NoError(t, err)
}

func TestValidate_ValueReceiver_Invalid(t *testing.T) {
	expected := errors.New("invalid value")

	err := Validate(valueValidator{err: expected})

	require.Error(t, err)
	assert.Equal(t, expected, err)
}

func TestValidate_PointerReceiver_Valid(t *testing.T) {
	err := Validate(&ptrValidator{err: nil})
	assert.NoError(t, err)
}

func TestValidate_PointerReceiver_Invalid(t *testing.T) {
	expected := errors.New("invalid ptr")

	err := Validate(&ptrValidator{err: expected})

	require.Error(t, err)
	assert.Equal(t, expected, err)
}

func TestValidate_String(t *testing.T) {
	assert.NotPanics(t, func() {
		err := Validate("just a string")
		assert.NoError(t, err)
	})
}

func TestValidate_Nil(t *testing.T) {
	assert.NotPanics(t, func() {
		err := Validate(nil)
		assert.NoError(t, err)
	})
}
