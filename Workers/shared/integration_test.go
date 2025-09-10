package shared

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	
	"shared/ffmpeg"
	"shared/security"
)

func TestSharedLibraries_Integration(t *testing.T) {
	// Test integration between different shared components
	
	// Test filename sanitization with FFmpeg processing
	unsafeFilename := "../../../etc/passwd"
	safeFilename := security.SanitizeFilename(unsafeFilename)
	
	assert.NotEqual(t, unsafeFilename, safeFilename)
	assert.True(t, security.ValidateFilename(safeFilename) || safeFilename == "")
}

func TestFFmpegSecurity_Integration(t *testing.T) {
	// Test that FFmpeg processor uses security validation
	
	inputData := []byte("test video data")
	
	// Test with sanitized arguments
	unsafeArg := "../../../etc/passwd"
	args := []string{"-i", "{input}", "-f", unsafeArg, "{output}"}
	
	_, err := ffmpeg.ProcessWithTempFiles(inputData, args)
	
	// Should fail due to path traversal detection
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "path traversal")
}

func TestConfigurationSecurity_Integration(t *testing.T) {
	// Test configuration validation with environment variables
	
	// Set up test environment
	t.Setenv("TEST_RABBITMQ_URL", "amqp://user:pass@localhost:5672/")
	t.Setenv("TEST_MINIO_ENDPOINT", "localhost:9000")
	t.Setenv("TEST_MINIO_ACCESS_KEY", "testaccess")
	t.Setenv("TEST_MINIO_SECRET_KEY", "testsecret123")
	t.Setenv("TEST_QUEUE_NAME", "test-queue")
	
	// Validate configuration
	rabbitURL, err := security.GetEnvWithValidation("TEST_RABBITMQ_URL", "", true)
	assert.NoError(t, err)
	
	err = security.ValidateRabbitMQURL(rabbitURL)
	assert.Error(t, err) // Should fail due to weak credentials
	
	// Test MinIO config
	endpoint, _ := security.GetEnvWithValidation("TEST_MINIO_ENDPOINT", "", true)
	accessKey, _ := security.GetEnvWithValidation("TEST_MINIO_ACCESS_KEY", "", true)
	secretKey, _ := security.GetEnvWithValidation("TEST_MINIO_SECRET_KEY", "", true)
	
	err = security.ValidateMinIOConfig(endpoint, accessKey, secretKey)
	assert.NoError(t, err)
	
	// Test queue name
	queueName, _ := security.GetEnvWithValidation("TEST_QUEUE_NAME", "", true)
	err = security.ValidateQueueName(queueName)
	assert.NoError(t, err)
}

func TestLoggingSecurity_Integration(t *testing.T) {
	// Test that log sanitization works with various inputs
	
	maliciousInputs := []string{
		"user input\nFAKE: Unauthorized access",
		"normal input\r\nERROR: System compromised",
		"test\x00\x1f\x7fwith control chars",
	}
	
	for _, input := range maliciousInputs {
		sanitized := security.SanitizeLogInput(input)
		
		// Ensure no injection characters remain
		assert.NotContains(t, sanitized, "\n")
		assert.NotContains(t, sanitized, "\r")
		assert.NotContains(t, sanitized, "\x00")
		
		// Ensure length is controlled
		assert.LessOrEqual(t, len(sanitized), 103) // 100 + "..."
	}
}

func TestFileProcessingSecurity_Integration(t *testing.T) {
	// Test file processing with security validation
	
	testCases := []struct {
		name     string
		filename string
		safe     bool
	}{
		{"safe filename", "video.mp4", true},
		{"unsafe path traversal", "../../../etc/passwd", false},
		{"unsafe special chars", "file@#$.mp4", false},
		{"unsafe spaces", "my video.mp4", false},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Validate filename
			isValid := security.ValidateFilename(tc.filename)
			assert.Equal(t, tc.safe, isValid)
			
			// Sanitize filename
			sanitized := security.SanitizeFilename(tc.filename)
			
			if tc.safe {
				// Safe filenames should remain unchanged
				assert.Equal(t, tc.filename, sanitized)
			} else {
				// Unsafe filenames should be sanitized
				assert.NotEqual(t, tc.filename, sanitized)
				assert.True(t, security.ValidateFilename(sanitized) || sanitized == "")
			}
		})
	}
}

func TestEnvironmentValidation_Integration(t *testing.T) {
	// Test complete environment validation workflow
	
	// Set up complete test environment
	testEnv := map[string]string{
		"RABBITMQ_URL":      "amqp://secure:password123@rabbitmq.example.com:5672/",
		"MINIO_ENDPOINT":    "minio.example.com:9000",
		"MINIO_ACCESS_KEY":  "secure_access_key",
		"MINIO_SECRET_KEY":  "secure_secret_key_123",
		"RAW_BUCKET":        "raw-videos",
		"PROCESSED_BUCKET":  "processed-videos",
		"STATE_QUEUE":       "state-machine",
		"MAX_WORKERS":       "5",
		"TIMEOUT_SECONDS":   "30",
	}
	
	for key, value := range testEnv {
		t.Setenv(key, value)
	}
	
	// Validate all configuration
	errors := []error{}
	
	// RabbitMQ validation
	if rabbitURL, err := security.GetEnvWithValidation("RABBITMQ_URL", "", true); err != nil {
		errors = append(errors, err)
	} else if err := security.ValidateRabbitMQURL(rabbitURL); err != nil {
		errors = append(errors, err)
	}
	
	// MinIO validation
	endpoint, err1 := security.GetEnvWithValidation("MINIO_ENDPOINT", "", true)
	accessKey, err2 := security.GetEnvWithValidation("MINIO_ACCESS_KEY", "", true)
	secretKey, err3 := security.GetEnvWithValidation("MINIO_SECRET_KEY", "", true)
	
	if err1 != nil || err2 != nil || err3 != nil {
		errors = append(errors, err1, err2, err3)
	} else if err := security.ValidateMinIOConfig(endpoint, accessKey, secretKey); err != nil {
		errors = append(errors, err)
	}
	
	// Queue validation
	for _, queueEnv := range []string{"RAW_BUCKET", "PROCESSED_BUCKET", "STATE_QUEUE"} {
		if queueName, err := security.GetEnvWithValidation(queueEnv, "", true); err != nil {
			errors = append(errors, err)
		} else if err := security.ValidateQueueName(queueName); err != nil {
			errors = append(errors, err)
		}
	}
	
	// Integer validation
	if _, err := security.GetIntEnvWithValidation("MAX_WORKERS", 1, 1, 100); err != nil {
		errors = append(errors, err)
	}
	
	if _, err := security.GetIntEnvWithValidation("TIMEOUT_SECONDS", 10, 1, 300); err != nil {
		errors = append(errors, err)
	}
	
	// Filter out nil errors
	var actualErrors []error
	for _, err := range errors {
		if err != nil {
			actualErrors = append(actualErrors, err)
		}
	}
	
	// All validations should pass
	assert.Empty(t, actualErrors, "Configuration validation should pass with valid environment")
}

func TestSecurityDefenseInDepth_Integration(t *testing.T) {
	// Test multiple layers of security validation
	
	maliciousInput := "../../../etc/passwd\nFAKE LOG ENTRY\x00\x1f"
	
	// Layer 1: Log sanitization
	sanitizedLog := security.SanitizeLogInput(maliciousInput)
	assert.NotContains(t, sanitizedLog, "\n")
	assert.NotContains(t, sanitizedLog, "\x00")
	
	// Layer 2: Filename sanitization
	sanitizedFilename := security.SanitizeFilename(maliciousInput)
	assert.True(t, security.ValidateFilename(sanitizedFilename) || sanitizedFilename == "")
	
	// Layer 3: FFmpeg argument validation
	args := []string{"-i", "{input}", "-f", maliciousInput, "{output}"}
	_, err := ffmpeg.ProcessWithTempFiles([]byte("test"), args)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "path traversal")
}

func TestPerformance_Integration(t *testing.T) {
	// Test performance of security functions with various input sizes
	
	inputSizes := []int{10, 100, 1000, 10000}
	
	for _, size := range inputSizes {
		t.Run(fmt.Sprintf("size_%d", size), func(t *testing.T) {
			// Create test input
			input := strings.Repeat("a@b\n", size/4)
			
			// Test log sanitization performance
			start := time.Now()
			result := security.SanitizeLogInput(input)
			logDuration := time.Since(start)
			
			assert.NotEmpty(t, result)
			assert.Less(t, logDuration, 100*time.Millisecond, "Log sanitization should be fast")
			
			// Test filename sanitization performance
			start = time.Now()
			result = security.SanitizeFilename(input)
			filenameDuration := time.Since(start)
			
			assert.NotEmpty(t, result)
			assert.Less(t, filenameDuration, 100*time.Millisecond, "Filename sanitization should be fast")
		})
	}
}

func TestErrorHandling_Integration(t *testing.T) {
	// Test error handling across different components
	
	// Test FFmpeg with invalid input
	_, err := ffmpeg.ProcessWithTempFiles(nil, []string{})
	assert.Error(t, err)
	
	// Test validation with invalid config
	err = security.ValidateRabbitMQURL("invalid-url")
	assert.Error(t, err)
	
	err = security.ValidateMinIOConfig("", "", "")
	assert.Error(t, err)
	
	// Test environment validation with missing required vars
	_, err = security.GetEnvWithValidation("NONEXISTENT_REQUIRED_VAR", "", true)
	assert.Error(t, err)
	
	// Test integer validation with invalid values
	_, err = security.GetIntEnvWithValidation("NONEXISTENT_INT", 10, 1, 100)
	assert.NoError(t, err) // Should return default
	
	t.Setenv("INVALID_INT", "not_a_number")
	_, err = security.GetIntEnvWithValidation("INVALID_INT", 10, 1, 100)
	assert.Error(t, err)
}

// Helper imports for the test
import (
	"fmt"
	"strings"
	"time"
)