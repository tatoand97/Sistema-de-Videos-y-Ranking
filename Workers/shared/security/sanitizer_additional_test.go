package security

import (
	"strings"
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeLogInput_ControlCharacters(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "null byte",
			input:    "log\x00message",
			expected: "logmessage",
		},
		{
			name:     "bell character",
			input:    "log\x07message",
			expected: "logmessage",
		},
		{
			name:     "escape character",
			input:    "log\x1bmessage",
			expected: "logmessage",
		},
		{
			name:     "delete character",
			input:    "log\x7fmessage",
			expected: "logmessage",
		},
		{
			name:     "high control characters",
			input:    "log\x80\x9fmessage",
			expected: "log\x80\x9fmessage", // These are not in the regex range
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeLogInput(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSanitizeLogInput_InjectionAttempts(t *testing.T) {
	injectionTests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "CRLF injection",
			input:    "user input\r\nFAKE LOG ENTRY",
			expected: "user inputFAKE LOG ENTRY",
		},
		{
			name:     "log forging attempt",
			input:    "normal\nERROR: Fake error message",
			expected: "normalERROR: Fake error message",
		},
		{
			name:     "multiple newlines",
			input:    "test\n\n\nfake entry",
			expected: "testfake entry",
		},
		{
			name:     "mixed line endings",
			input:    "test\r\n\r\nfake\n\rentry",
			expected: "testfakeentry",
		},
	}
	
	for _, tt := range injectionTests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeLogInput(tt.input)
			assert.Equal(t, tt.expected, result)
			// Ensure no line breaks remain
			assert.NotContains(t, result, "\n")
			assert.NotContains(t, result, "\r")
		})
	}
}

func TestSanitizeLogInput_LengthLimiting(t *testing.T) {
	tests := []struct {
		name        string
		inputLength int
		expectTrunc bool
	}{
		{"short input", 50, false},
		{"exactly 100 chars", 100, false},
		{"101 chars", 101, true},
		{"very long input", 500, true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := strings.Repeat("a", tt.inputLength)
			result := SanitizeLogInput(input)
			
			if tt.expectTrunc {
				assert.Equal(t, 103, len(result)) // 100 + "..."
				assert.True(t, strings.HasSuffix(result, "..."))
			} else {
				assert.Equal(t, tt.inputLength, len(result))
				assert.NotContains(t, result, "...")
			}
		})
	}
}

func TestSanitizeLogInput_WhitespaceHandling(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "leading spaces",
			input:    "   message",
			expected: "message",
		},
		{
			name:     "trailing spaces",
			input:    "message   ",
			expected: "message",
		},
		{
			name:     "both ends",
			input:    "  message  ",
			expected: "message",
		},
		{
			name:     "internal spaces preserved",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "only spaces",
			input:    "   ",
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

func TestValidateFilename_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		expected bool
	}{
		{"single character", "a", true},
		{"single dot", ".", true},
		{"single dash", "-", true},
		{"single underscore", "_", true},
		{"numbers only", "123", true},
		{"mixed valid chars", "a1-b2_c3.txt", true},
		{"unicode characters", "файл.txt", false},
		{"emoji", "😀.txt", false},
		{"accented characters", "café.txt", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateFilename(tt.filename)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateFilename_SecurityCases(t *testing.T) {
	securityTests := []struct {
		name     string
		filename string
		expected bool
	}{
		{"path traversal unix", "../etc/passwd", false},
		{"path traversal windows", "..\\windows\\system32", false},
		{"absolute path unix", "/etc/passwd", false},
		{"absolute path windows", "C:\\Windows\\System32", false},
		{"null byte", "file\x00.txt", false},
		{"device names windows", "CON", true}, // Valid chars but dangerous
		{"device names windows ext", "PRN.txt", true}, // Valid chars
		{"hidden file", ".hidden", true}, // Valid chars
	}
	
	for _, tt := range securityTests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateFilename(tt.filename)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSanitizeFilename_UnicodeHandling(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "cyrillic characters",
			input:    "файл.txt",
			expected: "____.txt",
		},
		{
			name:     "chinese characters",
			input:    "文件.mp4",
			expected: "__.mp4",
		},
		{
			name:     "emoji",
			input:    "video😀.mp4",
			expected: "video_.mp4",
		},
		{
			name:     "accented characters",
			input:    "café.txt",
			expected: "caf_.txt",
		},
		{
			name:     "mixed unicode and ascii",
			input:    "test文件file.txt",
			expected: "test__file.txt",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeFilename(tt.input)
			assert.Equal(t, tt.expected, result)
			// Ensure result only contains safe characters
			assert.True(t, ValidateFilename(result) || result == "")
		})
	}
}

func TestSanitizeFilename_WindowsReservedNames(t *testing.T) {
	reservedNames := []string{
		"CON", "PRN", "AUX", "NUL",
		"COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9",
		"LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9",
	}
	
	for _, name := range reservedNames {
		t.Run(name, func(t *testing.T) {
			// Test both with and without extension
			result1 := SanitizeFilename(name)
			result2 := SanitizeFilename(name + ".txt")
			
			// Should not change the name (sanitizer doesn't handle reserved names)
			// but ValidateFilename should still work
			assert.Equal(t, name, result1)
			assert.Equal(t, name+".txt", result2)
			
			// Both should be valid according to character rules
			assert.True(t, ValidateFilename(result1))
			assert.True(t, ValidateFilename(result2))
		})
	}
}

func TestSanitizeFilename_ExtensionPreservation(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "preserve simple extension",
			input:    "file@name.txt",
			expected: "file_name.txt",
		},
		{
			name:     "preserve multiple dots",
			input:    "file@name.tar.gz",
			expected: "file_name.tar.gz",
		},
		{
			name:     "no extension",
			input:    "file@name",
			expected: "file_name",
		},
		{
			name:     "extension with unsafe chars",
			input:    "file.t@t",
			expected: "file.t_t",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeFilename(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSanitizeFilename_PerformanceWithLargeInput(t *testing.T) {
	// Test with very large input
	largeInput := strings.Repeat("a@b", 1000) + ".txt"
	
	result := SanitizeFilename(largeInput)
	
	// Should be truncated to 255 characters
	assert.LessOrEqual(t, len(result), 255)
	// Should not contain unsafe characters
	assert.NotContains(t, result, "@")
}

func TestRegexPatterns_Consistency(t *testing.T) {
	// Test that regex patterns work consistently
	testStrings := []string{
		"normal_file.txt",
		"file with spaces.txt",
		"file@#$.txt",
		"file\nwith\nnewlines.txt",
		"file\x00with\x1fnull.txt",
	}
	
	for _, str := range testStrings {
		t.Run(str, func(t *testing.T) {
			// Test log sanitization
			logResult := SanitizeLogInput(str)
			assert.NotContains(t, logResult, "\n")
			assert.NotContains(t, logResult, "\r")
			
			// Test filename validation and sanitization
			isValid := ValidateFilename(str)
			sanitized := SanitizeFilename(str)
			
			if isValid {
				// If original is valid, sanitized should be the same
				assert.Equal(t, str, sanitized)
			} else {
				// If original is invalid, sanitized should be valid or empty
				assert.True(t, ValidateFilename(sanitized) || sanitized == "")
			}
		})
	}
}

func TestSanitizeLogInput_ConcurrentSafety(t *testing.T) {
	// Test concurrent access to sanitization functions
	input := "test\nmessage\rwith\tcontrol\x00chars"
	expected := SanitizeLogInput(input)
	
	done := make(chan bool, 10)
	
	for i := 0; i < 10; i++ {
		go func() {
			defer func() { done <- true }()
			
			result := SanitizeLogInput(input)
			assert.Equal(t, expected, result)
		}()
	}
	
	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestValidateFilename_AllASCIICharacters(t *testing.T) {
	// Test all ASCII characters to ensure regex works correctly
	for i := 0; i < 128; i++ {
		char := string(rune(i))
		filename := "test" + char + "file.txt"
		
		isValid := ValidateFilename(filename)
		
		// Should be valid only for alphanumeric, dot, dash, underscore
		expectedValid := unicode.IsLetter(rune(i)) || unicode.IsDigit(rune(i)) || 
			char == "." || char == "-" || char == "_"
		
		if expectedValid {
			assert.True(t, isValid, "Character %d (%s) should be valid", i, char)
		} else {
			assert.False(t, isValid, "Character %d (%s) should be invalid", i, char)
		}
	}
}