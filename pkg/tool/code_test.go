package tool

import (
	"testing"
	"time"
)

func TestCodeGenerator_Generate(t *testing.T) {
	cg := NewCodeGenerator(true, true, 5, 6)
	seed := "test_seed"

	code := cg.Generate(seed)
	if len(code) != 6 {
		t.Errorf("Generated code length expected to be 6, got %d", len(code))
	}
}

func TestCodeGenerator_Check(t *testing.T) {
	cg := NewCodeGenerator(true, true, 5, 6)
	seed := "test_seed"

	code := cg.Generate(seed)
	if !cg.Check(seed, code) {
		t.Errorf("Generated code should be valid")
	}

	invalidCode := "invalid"
	if cg.Check(seed, invalidCode) {
		t.Errorf("Invalid code should not be accepted")
	}
}

func TestCodeGenerator_GenerateDifferentCodes(t *testing.T) {
	cg := NewCodeGenerator(true, true, 5, 6)
	seed := "test_seed"

	code1 := cg.Generate(seed)
	time.Sleep(time.Minute) // Wait for a minute to ensure different time
	code2 := cg.Generate(seed)

	if code1 == code2 {
		t.Errorf("Generated codes should be different for different times")
	}
}

func TestCodeGenerator_CheckMultipleCodes(t *testing.T) {
	cg := NewCodeGenerator(true, true, 5, 6)
	seed := "test_seed"

	codes := make([]string, 5)
	for i := 0; i < 5; i++ {
		codes[i] = cg.Generate(seed)
		time.Sleep(time.Minute) // Wait for a minute to generate different codes
	}

	for _, code := range codes {
		if !cg.Check(seed, code) {
			t.Errorf("All generated codes within TTL should be valid")
		}
	}
}

func TestCodeGenerator_SeedConverter(t *testing.T) {
	cg := NewCodeGenerator(true, true, 5, 6)
	seed1 := "test_seed_1"
	seed2 := "test_seed_2"

	result1 := cg.seedConverter(seed1)
	result2 := cg.seedConverter(seed2)

	if result1 == result2 {
		t.Errorf("Different seeds should produce different results")
	}

	// Test consistency
	if result1 != cg.seedConverter(seed1) {
		t.Errorf("Seed converter should produce consistent results for the same input")
	}
}

func TestNewCodeGenerator(t *testing.T) {
	cg := NewCodeGenerator(true, false, 10, 8)

	if cg.ttl != 10 {
		t.Errorf("Expected TTL to be 10, got %d", cg.ttl)
	}

	if cg.len != 8 {
		t.Errorf("Expected length to be 8, got %d", cg.len)
	}

	if cg.generator == nil {
		t.Errorf("Generator should not be nil")
	}
}
