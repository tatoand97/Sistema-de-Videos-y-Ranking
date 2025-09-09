package security

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeLogInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "normal input",
			input:    "normal log message",
			expected: "normal log message",
		},
		{
			name:     "input with newlines",
			input:    "log message\nwith newline",
			expected: "log messagewith newline",
		},
		{
			name:     "input with carriage return",
			input:    "log message\rwith carriage return",
			expected: "log messagewith carriage return",
		},
		{
			name:     "input with tab",
			input:    "log message\twith tab",
			expected: "log messagewith tab",
		},
		{
			name:     "input with control characters",
			input:    "log\x00message\x1fwith\x7fcontrol",
			expected: "logmessagewithcontrol",
		},
		{
			name:     "long input gets truncated",
			input:    strings.Repeat("a", 150),
			expected: strings.Repeat("a", 100) + "...",
		},
		{
			name:     "input with leading/trailing spaces",
			input:    "  spaced input  ",
			expected: "spaced input",
		},
		{
			name:     "empty input",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeLogInput(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateFilename(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		expected bool
	}{
		{
			name:     "valid simple filename",
			filename: "video.mp4",
			expected: true,
		},
		{
			name:     "valid filename with numbers",
			filename: "video123.mp4",
			expected: true,
		},
		{
			name:     "valid filename with underscore and dash",
			filename: "my_video-file.mp4",
			expected: true,
		},
		{
			name:     "invalid filename with spaces",
			filename: "my video.mp4",
			expected: false,
		},
		{
			name:     "invalid filename with slash",
			filename: "path/to/video.mp4",
			expected: false,
		},
		{
			name:     "invalid filename with backslash",
			filename: "path\\to\\video.mp4",
			expected: false,
		},
		{
			name:     "invalid empty filename",
			filename: "",
			expected: false,
		},
		{
			name:     "invalid too long filename",
			filename: strings.Repeat("a", 256),
			expected: false,
		},
		{
			name:     "valid max length filename",
			filename: strings.Repeat("a", 255),
			expected: true,
		},
		{
			name:     "invalid filename with special characters",
			filename: "video@#$.mp4",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateFilename(tt.filename)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		expected string
	}{
		{
			name:     "valid filename unchanged",
			filename: "video.mp4",
			expected: "video.mp4",
		},
		{
			name:     "filename with spaces replaced",
			filename: "my video.mp4",
			expected: "my_video.mp4",
		},
		{
			name:     "filename with path traversal",
			filename: "../../../etc/passwd",
			expected: "_.._.._etc_passwd",
		},
		{
			name:     "filename with special characters",
			filename: "video@#$%.mp4",
			expected: "video____.mp4",
		},
		{
			name:     "filename with leading dots removed",
			filename: "...hidden.mp4",
			expected: "hidden.mp4",
		},
		{
			name:     "filename with backslashes",
			filename: "path\\to\\video.mp4",
			expected: "path_to_video.mp4",
		},
		{
			name:     "long filename gets truncated",
			filename: strings.Repeat("a", 300) + ".mp4",
			expected: strings.Repeat("a", 255),
		},
		{
			name:     "filename with mixed unsafe characters",
			filename: "../my video@file.mp4",
			expected: "_my_video_file.mp4",
		},
		{
			name:     "empty filename",
			filename: "",
			expected: "",
		},
		{
			name:     "only dots filename",
			filename: "....",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeFilename(tt.filename)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSanitizeFilename_SecurityCases(t *testing.T) {
	// Test specific security-related cases
	securityTests := []struct {
		name     string
		filename string
		expected string
	}{
		{
			name:     "path traversal attack",
			filename: "../../../../etc/passwd",
			expected: "_.._.._.._etc_passwd",
		},
		{
			name:     "null byte injection",
			filename: "video.mp4\x00.txt",
			expected: "video.mp4_.txt",
		},
		{
			name:     "windows reserved names",
			filename: "CON.mp4",
			expected: "CON.mp4",
		},
		{
			name:     "unicode characters",
			filename: "vidéo.mp4",
			expected: "vid_o.mp4",
		},
	}

	for _, tt := range securityTests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeFilename(tt.filename)
			assert.Equal(t, tt.expected, result)
			// Ensure result is safe
			assert.True(t, ValidateFilename(result) || result == "")
		})
	}
}