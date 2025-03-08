package utils

import (
	"net/url"
	"strings"
)

func IsValidURL(raw string) bool {
	if raw == "" {
		return false
	}
	if !strings.HasPrefix(raw, "http://") && !strings.HasPrefix(raw, "https://") {
		return false
	}
	u, err := url.ParseRequestURI(raw)
	if err != nil {
		return false
	}
	return u.Host != ""
}

func IsValidEmail(email string) bool {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	return len(parts[0]) > 0 && strings.Contains(parts[1], ".")
}
