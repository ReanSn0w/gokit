package web_test

import (
	"net/url"
	"reflect"
	"testing"

	"git.papkovda.ru/library/gokit/pkg/web"
)

// Define a struct to test the decoding into.
type TestStruct struct {
	Name  string `query:"name"`
	Age   int    `query:"age"`
	Email string `query:"email"`
}

func Test_DecodeQuery(t *testing.T) {
	tests := []struct {
		name     string
		values   url.Values
		expected TestStruct
		wantErr  bool
	}{
		{
			name: "valid values",
			values: url.Values{
				"name":  {"John Doe"},
				"age":   {"30"},
				"email": {"john.doe@example.com"},
			},
			expected: TestStruct{
				Name:  "John Doe",
				Age:   30,
				Email: "john.doe@example.com",
			},
			wantErr: false,
		},
		{
			name: "missing values",
			values: url.Values{
				"name": {"John Doe"},
			},
			expected: TestStruct{
				Name: "John Doe",
				Age:  0, // Assuming default int value
			},
			wantErr: false,
		},
		{
			name: "invalid age",
			values: url.Values{
				"name": {"John Doe"},
				"age":  {"invalid"},
			},
			expected: TestStruct{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result TestStruct

			err := web.DecodeQuery(tt.values, &result)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Decode() got = %v, expected %v", result, tt.expected)
			}
		})
	}
}
