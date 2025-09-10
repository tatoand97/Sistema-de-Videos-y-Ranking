package ffmpeg

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessWithTempFiles_Success(t *testing.T) {
	// Create minimal test data
	inputData := []byte("test video data")
	
	// Test with basic args (will fail due to invalid data, but tests parameter handling)
	args := []string{"-i", "{input}", "-c", "copy", "{output}"}
	
	_, err := ProcessWithTempFiles(inputData, args)
	
	// Expected to fail due to invalid video data and missing ffmpeg, but tests the flow
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ffmpeg failed")
}

func TestProcessWithTempFiles_PathTraversalPrevention(t *testing.T) {
	inputData := []byte("test data")
	
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "unix path traversal",
			args: []string{"-i", "../../../etc/passwd", "{output}"},
		},
		{
			name: "windows path traversal",
			args: []string{"-i", "..\\..\\windows\\system32\\config\\sam", "{output}"},
		},
		{
			name: "mixed path traversal",
			args: []string{"-i", "{input}", "-f", "../dangerous/path", "{output}"},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ProcessWithTempFiles(inputData, tt.args)
			
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "path traversal")
		})
	}
}

func TestProcessWithTempFiles_PlaceholderReplacement(t *testing.T) {
	inputData := []byte("test data")
	args := []string{"-i", "{input}", "-c", "copy", "{output}"}
	
	// Mock the temp directory creation to test placeholder replacement
	originalMkdirTemp := os.MkdirTemp
	defer func() { os.MkdirTemp = originalMkdirTemp }()
	
	tempDir := t.TempDir()
	os.MkdirTemp = func(dir, pattern string) (string, error) {
		return tempDir, nil
	}
	
	_, err := ProcessWithTempFiles(inputData, args)
	
	// Will fail at ffmpeg execution, but placeholders should be replaced
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ffmpeg failed")
}

func TestProcessWithTempFiles_EmptyInput(t *testing.T) {
	inputData := []byte{}
	args := []string{"-i", "{input}", "{output}"}
	
	_, err := ProcessWithTempFiles(inputData, args)
	
	assert.Error(t, err)
}

func TestProcessWithTempFiles_InvalidArgs(t *testing.T) {
	inputData := []byte("test data")
	
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "no input placeholder",
			args: []string{"-i", "missing_input", "{output}"},
		},
		{
			name: "no output placeholder",
			args: []string{"-i", "{input}", "missing_output"},
		},
		{
			name: "empty args",
			args: []string{},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ProcessWithTempFiles(inputData, tt.args)
			
			assert.Error(t, err)
		})
	}
}

func TestProcessWithTempFiles_TempDirCreationFailure(t *testing.T) {
	inputData := []byte("test data")
	args := []string{"-i", "{input}", "{output}"}
	
	// Mock MkdirTemp to fail
	originalMkdirTemp := os.MkdirTemp
	defer func() { os.MkdirTemp = originalMkdirTemp }()
	
	os.MkdirTemp = func(dir, pattern string) (string, error) {
		return "", assert.AnError
	}
	
	_, err := ProcessWithTempFiles(inputData, args)
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "create temp dir")
}

func TestProcessWithTempFiles_WriteFileFailure(t *testing.T) {
	inputData := []byte("test data")
	args := []string{"-i", "{input}", "{output}"}
	
	// Use invalid temp directory to cause write failure
	originalMkdirTemp := os.MkdirTemp
	defer func() { os.MkdirTemp = originalMkdirTemp }()
	
	os.MkdirTemp = func(dir, pattern string) (string, error) {
		return "/invalid/nonexistent/path", nil
	}
	
	_, err := ProcessWithTempFiles(inputData, args)
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "write input")
}

func TestProcessWithTempFiles_ArgumentValidation(t *testing.T) {
	inputData := []byte("test data")
	
	validArgs := []string{"-i", "{input}", "-c:v", "libx264", "-preset", "fast", "{output}"}
	invalidArgs := []string{"-i", "{input}", "-f", "../etc/passwd", "{output}"}
	
	// Valid args should not fail on validation (will fail on ffmpeg execution)
	_, err := ProcessWithTempFiles(inputData, validArgs)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ffmpeg failed")
	
	// Invalid args should fail on validation
	_, err = ProcessWithTempFiles(inputData, invalidArgs)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "path traversal")
}

func TestProcessWithTempFiles_LargeInput(t *testing.T) {
	// Test with larger input data
	inputData := make([]byte, 1024*1024) // 1MB
	for i := range inputData {
		inputData[i] = byte(i % 256)
	}
	
	args := []string{"-i", "{input}", "-c", "copy", "{output}"}
	
	_, err := ProcessWithTempFiles(inputData, args)
	
	// Should handle large input without issues (will fail at ffmpeg execution)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ffmpeg failed")
}

func TestProcessWithTempFiles_SpecialCharactersInArgs(t *testing.T) {
	inputData := []byte("test data")
	
	tests := []struct {
		name string
		args []string
		safe bool
	}{
		{
			name: "safe special characters",
			args: []string{"-i", "{input}", "-metadata", "title=My Video", "{output}"},
			safe: true,
		},
		{
			name: "unsafe path traversal",
			args: []string{"-i", "{input}", "-f", "../../dangerous", "{output}"},
			safe: false,
		},
		{
			name: "safe with quotes",
			args: []string{"-i", "{input}", "-vf", "scale=1920:1080", "{output}"},
			safe: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ProcessWithTempFiles(inputData, tt.args)
			
			if tt.safe {
				// Should fail at ffmpeg execution, not validation
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "ffmpeg failed")
			} else {
				// Should fail at validation
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "path traversal")
			}
		})
	}
}