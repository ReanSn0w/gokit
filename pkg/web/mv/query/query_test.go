package query_test

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/ReanSn0w/gokit/pkg/web/mv/query"
)

// Define a struct to test the decoding into.
type TestStruct struct {
	Name    string   `query:"name"`
	Age     int      `query:"age"`
	Email   string   `query:"email"`
	Hobbies []string `query:"hobbies"`
}

func Test_Decode(t *testing.T) {
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
		{
			name: "with hobbies",
			values: url.Values{
				"name":    {"John Doe"},
				"age":     {"30"},
				"hobbies": {"reading", "coding"},
			},
			expected: TestStruct{
				Name:    "John Doe",
				Age:     30,
				Hobbies: []string{"reading", "coding"},
			},
			wantErr: false,
		},
		{
			name: "with hobbies and email",
			values: url.Values{
				"name":    {"John Doe"},
				"age":     {"30"},
				"hobbies": {"reading", "coding"},
				"email":   {"john.doe@example.com"},
			},
			expected: TestStruct{
				Name:    "John Doe",
				Age:     30,
				Hobbies: []string{"reading", "coding"},
				Email:   "john.doe@example.com",
			},
			wantErr: false,
		},
		{
			name: "missing hobbies",
			values: url.Values{
				"name": {"John Doe"},
				"age":  {"30"},
			},
			expected: TestStruct{
				Name:  "John Doe",
				Age:   30,
				Email: "",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result TestStruct

			err := query.Decode(tt.values, &result)
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
