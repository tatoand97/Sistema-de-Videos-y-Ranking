package security

import (
	"regexp"
	"strings"
)

var (
	logSanitizeRegex = regexp.MustCompile(`[\r\n\t\x00-\x1f\x7f-\x9f]`)
	filenameSafeRegex = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
	filenameCleanRegex = regexp.MustCompile(`[^a-zA-Z0-9._-]`)
)

// SanitizeLogInput removes potentially dangerous characters from log input
// to prevent log injection attacks (CWE-117)
func SanitizeLogInput(input string) string {
	// Remove newlines, carriage returns, and control characters
	sanitized := logSanitizeRegex.ReplaceAllString(input, "")
	
	// Limit length to prevent log flooding
	if len(sanitized) > 100 {
		sanitized = sanitized[:100] + "..."
	}
	
	return strings.TrimSpace(sanitized)
}

// ValidateFilename checks if filename contains only safe characters
func ValidateFilename(filename string) bool {
	// Allow alphanumeric, dots, hyphens, underscores
	return filenameSafeRegex.MatchString(filename) && len(filename) > 0 && len(filename) <= 255
}

// SanitizeFilename removes unsafe characters from filename
func SanitizeFilename(filename string) string {
	// Remove path traversal attempts and unsafe characters
	sanitized := filenameCleanRegex.ReplaceAllString(filename, "_")
	
	// Remove leading dots to prevent hidden files
	sanitized = strings.TrimLeft(sanitized, ".")
	
	// Limit length
	if len(sanitized) > 255 {
		sanitized = sanitized[:255]
	}
	
	return sanitized
}