package trainer_test

import (
	"testing"

	"github.com/ReanSn0w/gokit/pkg/lib/trainer"
	"github.com/stretchr/testify/assert"
)

func Test_ObjectToString(t *testing.T) {
	testCases := []struct {
		Name   string
		Input  any
		Output string
	}{
		{
			Name: "Plain struct",
			Input: struct {
				Field string
			}{
				Field: "Hello world",
			},
			Output: "# Field\nHello world",
		},
		{
			Name: "Many Fields",
			Input: struct {
				Firstname string
				Lastname  string
			}{
				Firstname: "Олег",
				Lastname:  "Олегов",
			},
			Output: "# Firstname\nОлег\n\n# Lastname\nОлегов",
		},
		{
			Name: "Included Structs",
			Input: struct {
				FIO struct {
					Firstname string
					Lastname  string
				}
			}{
				FIO: struct {
					Firstname string
					Lastname  string
				}{
					Firstname: "Олег",
					Lastname:  "Олегов",
				},
			},
			Output: "# FIO\n\n## Firstname\nОлег\n\n## Lastname\nОлегов",
		},
		{
			Name: "Included Structs With Pointer",
			Input: struct {
				FIO *struct {
					Firstname string
					Lastname  string
				}
			}{
				FIO: &struct {
					Firstname string
					Lastname  string
				}{
					Firstname: "Олег",
					Lastname:  "Олегов",
				},
			},
			Output: "# FIO\n\n## Firstname\nОлег\n\n## Lastname\nОлегов",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			val := trainer.ObjectToString(tc.Input)
			assert.Equal(t, tc.Output, val)
		})
	}
}

func Test_MakeSchema(t *testing.T) {
	testCases := []struct {
		Name   string
		Input  any
		Output []byte
	}{
		{
			Name:   "empty struct",
			Input:  struct{}{},
			Output: []byte(`{"type":"object","properties":{}}`),
		},
		{
			Name: "simple struct",
			Input: struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			}{},
			Output: []byte(`{"type":"object","properties":{"id":{"type":"integer"},"name":{"type":"string"}},"required": ["id", "name"]}`),
		},
		{
			Name: "struct with slice and nested struct",
			Input: struct {
				List []string `json:"list"`
				Sub  struct {
					Value bool `json:"value"`
				} `json:"sub"`
			}{},
			Output: []byte(`{"type":"object","properties":{"list":{"type":"array","items":{"type":"string"}},"sub":{"type":"object","properties":{"value":{"type":"boolean"}},"required":["value"]}},"required":["list","sub"]}`),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			val := trainer.MakeSchema(tc.Input)
			assert.JSONEq(t, string(tc.Output), string(val))
		})
	}
}
