package utils

import "testing"

func TestIsValidURL(t *testing.T) {
	tests := []struct {
		url   string
		valid bool
	}{
		{"https://example.com", true},
		{"http://localhost:3000/path", true},
		{"https://sub.domain.io/a/b?q=1", true},
		{"ftp://files.com", false},
		{"not-a-url", false},
		{"", false},
		{"https://", false},
	}
	for _, tt := range tests {
		if got := IsValidURL(tt.url); got != tt.valid {
			t.Errorf("IsValidURL(%q) = %v, want %v", tt.url, got, tt.valid)
		}
	}
}

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		email string
		valid bool
	}{
		{"user@example.com", true},
		{"a@b.c", true},
		{"nope", false},
		{"@no.com", false},
	}
	for _, tt := range tests {
		if got := IsValidEmail(tt.email); got != tt.valid {
			t.Errorf("IsValidEmail(%q) = %v, want %v", tt.email, got, tt.valid)
		}
	}
}
