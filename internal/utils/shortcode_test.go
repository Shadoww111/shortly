package utils

import "testing"

func TestGenerateShortCode(t *testing.T) {
	code, err := GenerateShortCode(7)
	if err != nil {
		t.Fatal(err)
	}
	if len(code) != 7 {
		t.Errorf("expected length 7, got %d", len(code))
	}

	// uniqueness
	codes := make(map[string]bool)
	for i := 0; i < 1000; i++ {
		c, _ := GenerateShortCode(7)
		if codes[c] {
			t.Errorf("duplicate code: %s", c)
		}
		codes[c] = true
	}
}

func TestIsValidCustomCode(t *testing.T) {
	tests := []struct {
		code  string
		valid bool
	}{
		{"abc", true},
		{"my-link", true},
		{"test_123", true},
		{"ab", false},
		{"", false},
		{"has space", false},
		{"abcdefghijklmnopqrstu", false}, // 21 chars
	}
	for _, tt := range tests {
		if got := IsValidCustomCode(tt.code); got != tt.valid {
			t.Errorf("IsValidCustomCode(%q) = %v, want %v", tt.code, got, tt.valid)
		}
	}
}
