package tool_test

import (
	"testing"
	"time"

	"github.com/ReanSn0w/gokit/pkg/tool"
)

func TestNewRandom(t *testing.T) {
	r := tool.NewRandom(true, false, false)
	if r == nil {
		t.Error("NewRandom() returned nil")
	}
}

func TestRandom_Generate(t *testing.T) {
	tests := []struct {
		name     string
		count    int
		alphabet bool
		numbers  bool
		symbols  bool
	}{
		{"All options", 10, true, true, true},
		{"Only alphabet", 5, true, false, false},
		{"Only numbers", 5, false, true, false},
		{"Only symbols", 5, false, false, true},
	}

	for _, tt := range tests {
		r := tool.NewRandom(tt.alphabet, tt.numbers, tt.symbols)

		t.Run(tt.name, func(t *testing.T) {
			result := r.Generate(tt.count)
			if len(result) != tt.count {
				t.Errorf("Generate() returned string of length %d, want %d", len(result), tt.count)
			}
			checkContent(t, result, tt.alphabet, tt.numbers, tt.symbols)
		})
	}
}

func TestRandom_Pseudo(t *testing.T) {
	seed := time.Now().Unix()
	tests := []struct {
		name     string
		count    int
		alphabet bool
		numbers  bool
		symbols  bool
	}{
		{"All options", 10, true, true, true},
		{"Only alphabet", 5, true, false, false},
		{"Only numbers", 5, false, true, false},
		{"Only symbols", 5, false, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tool.NewRandom(tt.alphabet, tt.numbers, tt.symbols)

			result1 := r.Pseudo(seed, tt.count)
			result2 := r.Pseudo(seed, tt.count)

			if result1 != result2 {
				t.Errorf("Pseudo() with same seed returned different results: %s and %s", result1, result2)
			}

			if len(result1) != tt.count {
				t.Errorf("Pseudo() returned string of length %d, want %d", len(result1), tt.count)
			}

			checkContent(t, result1, tt.alphabet, tt.numbers, tt.symbols)
		})
	}
}

func checkContent(t *testing.T, s string, alphabet, numbers, symbols bool) {
	hasAlphabet := false
	hasNumbers := false
	hasSymbols := false

	for _, ch := range s {
		switch {
		case ch >= 'A' && ch <= 'Z':
			hasAlphabet = true
		case ch >= '0' && ch <= '9':
			hasNumbers = true
		case ch == '!' || ch == '@' || ch == '#' || ch == '$' || ch == '%' || ch == '^' || ch == '&' || ch == '*' || ch == '<' || ch == '>' || ch == '?':
			hasSymbols = true
		}
	}

	if alphabet && !hasAlphabet {
		t.Error("String does not contain alphabet characters")
	}
	if numbers && !hasNumbers {
		t.Error("String does not contain numbers")
	}
	if symbols && !hasSymbols {
		t.Error("String does not contain symbols")
	}
	if !alphabet && hasAlphabet {
		t.Error("String contains alphabet characters when it shouldn't")
	}
	if !numbers && hasNumbers {
		t.Error("String contains numbers when it shouldn't")
	}
	if !symbols && hasSymbols {
		t.Error("String contains symbols when it shouldn't")
	}
}
