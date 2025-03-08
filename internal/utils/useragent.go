package utils

import "strings"

// ParseUserAgent extracts device, browser, and OS from a user agent string.
func ParseUserAgent(ua string) (device, browser, os string) {
	ua = strings.ToLower(ua)

	// device
	switch {
	case strings.Contains(ua, "mobile") || strings.Contains(ua, "android") && !strings.Contains(ua, "tablet"):
		device = "mobile"
	case strings.Contains(ua, "tablet") || strings.Contains(ua, "ipad"):
		device = "tablet"
	default:
		device = "desktop"
	}

	// browser
	switch {
	case strings.Contains(ua, "firefox"):
		browser = "Firefox"
	case strings.Contains(ua, "edg"):
		browser = "Edge"
	case strings.Contains(ua, "opr") || strings.Contains(ua, "opera"):
		browser = "Opera"
	case strings.Contains(ua, "chrome") || strings.Contains(ua, "chromium"):
		browser = "Chrome"
	case strings.Contains(ua, "safari"):
		browser = "Safari"
	default:
		browser = "Other"
	}

	// os
	switch {
	case strings.Contains(ua, "windows"):
		os = "Windows"
	case strings.Contains(ua, "mac os") || strings.Contains(ua, "macintosh"):
		os = "macOS"
	case strings.Contains(ua, "linux"):
		os = "Linux"
	case strings.Contains(ua, "android"):
		os = "Android"
	case strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad"):
		os = "iOS"
	default:
		os = "Other"
	}

	return
}
