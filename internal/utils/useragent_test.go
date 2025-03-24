package utils

import "testing"

func TestParseUserAgent(t *testing.T) {
	tests := []struct {
		ua      string
		device  string
		browser string
		os      string
	}{
		{
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/120.0",
			"desktop", "Chrome", "Windows",
		},
		{
			"Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) Safari/605.1.15",
			"mobile", "Safari", "iOS",
		},
		{
			"Mozilla/5.0 (Linux; Android 14) Mobile Firefox/120.0",
			"mobile", "Firefox", "Android",
		},
	}
	for _, tt := range tests {
		d, b, o := ParseUserAgent(tt.ua)
		if d != tt.device || b != tt.browser || o != tt.os {
			t.Errorf("ParseUserAgent(%q) = (%s, %s, %s), want (%s, %s, %s)",
				tt.ua, d, b, o, tt.device, tt.browser, tt.os)
		}
	}
}
