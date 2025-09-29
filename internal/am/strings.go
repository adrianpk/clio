package am

import "strings"

// Cap capitalizes the first letter of a string.
func Cap(s string) string {
	if s == "" {
		return s
	}

	return strings.ToUpper(s[:1]) + s[1:]
}
