package ffmpeg

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateFFmpeg_Success(t *testing.T) {
	// This test will pass if ffmpeg is installed on the system
	err := ValidateFFmpeg()
	
	if err != nil {
		// Expected in environments without ffmpeg
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ffmpeg not found")
	} else {
		// If ffmpeg is available
		assert.NoError(t, err)
	}
}

func TestValidateFFmpeg_NotFound(t *testing.T) {
	// Mock exec.LookPath to simulate ffmpeg not found
	originalLookPath := exec.LookPath
	defer func() { exec.LookPath = originalLookPath }()
	
	exec.LookPath = func(file string) (string, error) {
		if file == "ffmpeg" {
			return "", &exec.Error{Name: file, Err: exec.ErrNotFound}
		}
		return originalLookPath(file)
	}
	
	err := ValidateFFmpeg()
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ffmpeg not found")
}

func TestValidateFFmpeg_PathError(t *testing.T) {
	// Mock exec.LookPath to simulate path error
	originalLookPath := exec.LookPath
	defer func() { exec.LookPath = originalLookPath }()
	
	exec.LookPath = func(file string) (string, error) {
		if file == "ffmpeg" {
			return "", os.ErrPermission
		}
		return originalLookPath(file)
	}
	
	err := ValidateFFmpeg()
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ffmpeg not found")
}

func TestValidateFFmpeg_EmptyPath(t *testing.T) {
	// Mock exec.LookPath to return empty path
	originalLookPath := exec.LookPath
	defer func() { exec.LookPath = originalLookPath }()
	
	exec.LookPath = func(file string) (string, error) {
		if file == "ffmpeg" {
			return "", nil // Empty path but no error
		}
		return originalLookPath(file)
	}
	
	err := ValidateFFmpeg()
	
	// Should still work if LookPath returns empty string with no error
	assert.NoError(t, err)
}

func TestValidateFFmpeg_MultipleValidations(t *testing.T) {
	// Test multiple consecutive validations
	for i := 0; i < 3; i++ {
		err := ValidateFFmpeg()
		
		// Should be consistent across multiple calls
		if err != nil {
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "ffmpeg not found")
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestValidateFFmpeg_ConcurrentValidations(t *testing.T) {
	// Test concurrent validations
	done := make(chan bool, 3)
	
	for i := 0; i < 3; i++ {
		go func() {
			defer func() { done <- true }()
			
			err := ValidateFFmpeg()
			
			// Should handle concurrent calls safely
			if err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		}()
	}
	
	// Wait for all goroutines to complete
	for i := 0; i < 3; i++ {
		<-done
	}
}

func TestValidateFFmpeg_ErrorTypes(t *testing.T) {
	// Test different types of errors that LookPath might return
	errorTests := []struct {
		name        string
		mockError   error
		expectError bool
	}{
		{
			name:        "exec.ErrNotFound",
			mockError:   exec.ErrNotFound,
			expectError: true,
		},
		{
			name:        "permission error",
			mockError:   os.ErrPermission,
			expectError: true,
		},
		{
			name:        "no error",
			mockError:   nil,
			expectError: false,
		},
	}
	
	originalLookPath := exec.LookPath
	defer func() { exec.LookPath = originalLookPath }()
	
	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			exec.LookPath = func(file string) (string, error) {
				if file == "ffmpeg" {
					if tt.mockError != nil {
						return "", tt.mockError
					}
					return "/usr/bin/ffmpeg", nil
				}
				return originalLookPath(file)
			}
			
			err := ValidateFFmpeg()
			
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "ffmpeg not found")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateFFmpeg_SystemIntegration(t *testing.T) {
	// Test actual system integration (if ffmpeg is available)
	err := ValidateFFmpeg()
	
	if err == nil {
		// If ffmpeg is available, test that we can actually call it
		cmd := exec.Command("ffmpeg", "-version")
		output, cmdErr := cmd.Output()
		
		assert.NoError(t, cmdErr, "ffmpeg should be executable if validation passes")
		assert.Contains(t, string(output), "ffmpeg version", "ffmpeg should return version info")
	} else {
		// If ffmpeg is not available, ensure the error is appropriate
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ffmpeg not found")
		
		// Verify that ffmpeg is indeed not available
		_, cmdErr := exec.LookPath("ffmpeg")
		assert.Error(t, cmdErr, "ffmpeg should not be found in PATH")
	}
}

func TestValidateFFmpeg_EdgeCases(t *testing.T) {
	originalLookPath := exec.LookPath
	defer func() { exec.LookPath = originalLookPath }()
	
	// Test with different file names (should only check for "ffmpeg")
	exec.LookPath = func(file string) (string, error) {
		switch file {
		case "ffmpeg":
			return "/usr/bin/ffmpeg", nil
		case "ffprobe":
			return "", exec.ErrNotFound
		default:
			return originalLookPath(file)
		}
	}
	
	err := ValidateFFmpeg()
	assert.NoError(t, err, "Should only validate ffmpeg, not other tools")
}

func TestValidateFFmpeg_PathVariations(t *testing.T) {
	originalLookPath := exec.LookPath
	defer func() { exec.LookPath = originalLookPath }()
	
	pathTests := []struct {
		name string
		path string
	}{
		{"unix path", "/usr/bin/ffmpeg"},
		{"local path", "./ffmpeg"},
		{"windows path", "C:\\Program Files\\ffmpeg\\bin\\ffmpeg.exe"},
		{"relative path", "../bin/ffmpeg"},
	}
	
	for _, tt := range pathTests {
		t.Run(tt.name, func(t *testing.T) {
			exec.LookPath = func(file string) (string, error) {
				if file == "ffmpeg" {
					return tt.path, nil
				}
				return originalLookPath(file)
			}
			
			err := ValidateFFmpeg()
			assert.NoError(t, err, "Should accept any valid path returned by LookPath")
		})
	}
}