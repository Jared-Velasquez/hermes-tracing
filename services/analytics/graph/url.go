package graph

import (
	"strings"
)

// extractHostFromURL parses a URL string and returns the host (without port)
func extractHostFromURL(urlStr string) string {
	// Simple parsing without importing net/url to keep it lightweight
	// Handles: "https://api.stripe.com:443/v1/charges" -> "api.stripe.com"

	// Remove scheme
	if idx := strings.Index(urlStr, "://"); idx != -1 {
		urlStr = urlStr[idx+3:]
	}

	// Remove path
	if idx := strings.Index(urlStr, "/"); idx != -1 {
		urlStr = urlStr[:idx]
	}

	// Remove userinfo if present
	if idx := strings.Index(urlStr, "@"); idx != -1 {
		urlStr = urlStr[idx+1:]
	}

	// Remove port
	if idx := strings.LastIndex(urlStr, ":"); idx != -1 {
		// Make sure it's a port, not part of IPv6
		if !strings.Contains(urlStr[idx:], "]") {
			urlStr = urlStr[:idx]
		}
	}

	return urlStr
}