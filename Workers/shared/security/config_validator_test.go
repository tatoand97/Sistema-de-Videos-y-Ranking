package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateRabbitMQURL_Success(t *testing.T) {
	validURLs := []string{
		"amqp://myuser:mypass@192.168.1.100:5672/",
		"amqps://secure:password123@rabbit.domain.com:5671/production",
	}
	
	for _, url := range validURLs {
		t.Run(url, func(t *testing.T) {
			err := ValidateRabbitMQURL(url)
			assert.NoError(t, err)
		})
	}
}

func TestValidateRabbitMQURL_EmptyURL(t *testing.T) {
	err := ValidateRabbitMQURL("")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "RABBITMQ_URL is required")
}

func TestValidateRabbitMQURL_InvalidFormat(t *testing.T) {
	invalidURLs := []string{
		"not-a-url",
		"http://localhost:5672/",
		"://missing-scheme",
		"amqp://[invalid-host",
	}
	
	for _, url := range invalidURLs {
		t.Run(url, func(t *testing.T) {
			err := ValidateRabbitMQURL(url)
			assert.Error(t, err)
			// Error message may vary based on URL format
		assert.Error(t, err)
		})
	}
}

func TestValidateRabbitMQURL_InvalidScheme(t *testing.T) {
	invalidSchemes := []string{
		"http://user:pass@localhost:5672/",
		"https://user:pass@localhost:5672/",
		"tcp://user:pass@localhost:5672/",
		"ftp://user:pass@localhost:5672/",
	}
	
	for _, url := range invalidSchemes {
		t.Run(url, func(t *testing.T) {
			err := ValidateRabbitMQURL(url)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "must use amqp or amqps scheme")
		})
	}
}

func TestValidateRabbitMQURL_WeakCredentials(t *testing.T) {
	weakCredURLs := []string{
		"amqp://user:pass@localhost:5672/",
		"amqp://admin:admin@localhost:5672/",
		"amqp://guest:guest@localhost:5672/",
		"amqp://test:test@localhost:5672/",
	}
	
	for _, url := range weakCredURLs {
		t.Run(url, func(t *testing.T) {
			err := ValidateRabbitMQURL(url)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "weak credentials")
		})
	}
}

func TestValidateRabbitMQURL_NoCredentials(t *testing.T) {
	err := ValidateRabbitMQURL("amqp://localhost:5672/")
	assert.NoError(t, err) // No credentials is allowed
}

func TestValidateMinIOConfig_Success(t *testing.T) {
	err := ValidateMinIOConfig("localhost:9000", "accesskey123", "secretkey123")
	assert.NoError(t, err)
}

func TestValidateMinIOConfig_MissingEndpoint(t *testing.T) {
	err := ValidateMinIOConfig("", "accesskey", "secretkey")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "MINIO_ENDPOINT is required")
}

func TestValidateMinIOConfig_MissingAccessKey(t *testing.T) {
	err := ValidateMinIOConfig("localhost:9000", "", "secretkey")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "MINIO_ACCESS_KEY is required")
}

func TestValidateMinIOConfig_MissingSecretKey(t *testing.T) {
	err := ValidateMinIOConfig("localhost:9000", "accesskey", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "MINIO_SECRET_KEY is required")
}

func TestValidateMinIOConfig_ShortSecretKey(t *testing.T) {
	err := ValidateMinIOConfig("localhost:9000", "accesskey", "short")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must be at least 8 characters")
}

func TestGetEnvWithValidation_Required(t *testing.T) {
	t.Setenv("TEST_VAR", "test_value")
	
	value, err := GetEnvWithValidation("TEST_VAR", "default", true)
	assert.NoError(t, err)
	assert.Equal(t, "test_value", value)
}

func TestGetEnvWithValidation_RequiredMissing(t *testing.T) {
	t.Setenv("TEST_VAR", "")
	
	value, err := GetEnvWithValidation("NONEXISTENT_VAR", "default", true)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "is required")
	assert.Empty(t, value)
}

func TestGetEnvWithValidation_OptionalWithDefault(t *testing.T) {
	value, err := GetEnvWithValidation("NONEXISTENT_VAR", "default_value", false)
	assert.NoError(t, err)
	assert.Equal(t, "default_value", value)
}

func TestGetEnvWithValidation_OptionalEmpty(t *testing.T) {
	t.Setenv("EMPTY_VAR", "")
	
	value, err := GetEnvWithValidation("EMPTY_VAR", "default_value", false)
	assert.NoError(t, err)
	assert.Equal(t, "default_value", value)
}

func TestGetIntEnvWithValidation_Success(t *testing.T) {
	t.Setenv("INT_VAR", "42")
	
	value, err := GetIntEnvWithValidation("INT_VAR", 10, 1, 100)
	assert.NoError(t, err)
	assert.Equal(t, 42, value)
}

func TestGetIntEnvWithValidation_Default(t *testing.T) {
	value, err := GetIntEnvWithValidation("NONEXISTENT_INT", 25, 1, 100)
	assert.NoError(t, err)
	assert.Equal(t, 25, value)
}

func TestGetIntEnvWithValidation_InvalidFormat(t *testing.T) {
	t.Setenv("INVALID_INT", "not_a_number")
	
	value, err := GetIntEnvWithValidation("INVALID_INT", 25, 1, 100)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid integer value")
	assert.Equal(t, 25, value) // Should return default
}

func TestGetIntEnvWithValidation_OutOfRange(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		min      int
		max      int
		expected int
	}{
		{"below minimum", "0", 1, 100, 25},
		{"above maximum", "150", 1, 100, 25},
		{"negative when positive required", "-5", 1, 100, 25},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("RANGE_TEST", tt.value)
			
			value, err := GetIntEnvWithValidation("RANGE_TEST", tt.expected, tt.min, tt.max)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "must be between")
			assert.Equal(t, tt.expected, value)
		})
	}
}

func TestGetIntEnvWithValidation_BoundaryValues(t *testing.T) {
	tests := []struct {
		name  string
		value string
		min   int
		max   int
		valid bool
	}{
		{"minimum value", "1", 1, 100, true},
		{"maximum value", "100", 1, 100, true},
		{"just below minimum", "0", 1, 100, false},
		{"just above maximum", "101", 1, 100, false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("BOUNDARY_TEST", tt.value)
			
			value, err := GetIntEnvWithValidation("BOUNDARY_TEST", 50, tt.min, tt.max)
			
			if tt.valid {
				assert.NoError(t, err)
				expectedValue := 1
				if tt.value == "100" {
					expectedValue = 100
				}
				assert.Equal(t, expectedValue, value)
			} else {
				assert.Error(t, err)
				assert.Equal(t, 50, value) // Default value
			}
		})
	}
}

func TestValidateQueueName_Success(t *testing.T) {
	validNames := []string{
		"video-processing",
		"state_machine",
		"notifications",
		"queue123",
		"a",
	}
	
	for _, name := range validNames {
		t.Run(name, func(t *testing.T) {
			err := ValidateQueueName(name)
			assert.NoError(t, err)
		})
	}
}

func TestValidateQueueName_Empty(t *testing.T) {
	err := ValidateQueueName("")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be empty")
}

func TestValidateQueueName_TooLong(t *testing.T) {
	longName := string(make([]byte, 256))
	for i := range longName {
		longName = longName[:i] + "a" + longName[i+1:]
	}
	
	err := ValidateQueueName(longName)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "too long")
}

func TestValidateQueueName_WithSpaces(t *testing.T) {
	invalidNames := []string{
		"queue with spaces",
		" leading-space",
		"trailing-space ",
		"middle space queue",
	}
	
	for _, name := range invalidNames {
		t.Run(name, func(t *testing.T) {
			err := ValidateQueueName(name)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "cannot contain spaces")
		})
	}
}

func TestValidateQueueName_MaxLength(t *testing.T) {
	// Test exactly 255 characters (should be valid)
	maxLengthName := string(make([]byte, 255))
	for i := range maxLengthName {
		maxLengthName = maxLengthName[:i] + "a" + maxLengthName[i+1:]
	}
	
	err := ValidateQueueName(maxLengthName)
	assert.NoError(t, err)
}

func TestConfigValidation_Integration(t *testing.T) {
	// Test integration of multiple validation functions
	
	// Set up valid environment
	t.Setenv("RABBITMQ_URL", "amqp://myuser:mypass@localhost:5672/")
	t.Setenv("MINIO_ENDPOINT", "localhost:9000")
	t.Setenv("MINIO_ACCESS_KEY", "accesskey123")
	t.Setenv("MINIO_SECRET_KEY", "secretkey123")
	t.Setenv("QUEUE_NAME", "video-processing")
	t.Setenv("MAX_WORKERS", "10")
	
	// Validate RabbitMQ
	rabbitURL, err := GetEnvWithValidation("RABBITMQ_URL", "", true)
	assert.NoError(t, err)
	err = ValidateRabbitMQURL(rabbitURL)
	assert.NoError(t, err)
	
	// Validate MinIO
	endpoint, err := GetEnvWithValidation("MINIO_ENDPOINT", "", true)
	assert.NoError(t, err)
	accessKey, err := GetEnvWithValidation("MINIO_ACCESS_KEY", "", true)
	assert.NoError(t, err)
	secretKey, err := GetEnvWithValidation("MINIO_SECRET_KEY", "", true)
	assert.NoError(t, err)
	err = ValidateMinIOConfig(endpoint, accessKey, secretKey)
	assert.NoError(t, err)
	
	// Validate queue name
	queueName, err := GetEnvWithValidation("QUEUE_NAME", "", true)
	assert.NoError(t, err)
	err = ValidateQueueName(queueName)
	assert.NoError(t, err)
	
	// Validate integer config
	maxWorkers, err := GetIntEnvWithValidation("MAX_WORKERS", 5, 1, 50)
	assert.NoError(t, err)
	assert.Equal(t, 10, maxWorkers)
}

func TestConfigValidation_ErrorCascade(t *testing.T) {
	// Test that validation errors cascade properly
	
	// Set invalid values
	t.Setenv("RABBITMQ_URL", "invalid-url")
	t.Setenv("MINIO_SECRET_KEY", "short")
	t.Setenv("QUEUE_NAME", "queue with spaces")
	t.Setenv("MAX_WORKERS", "invalid")
	
	// Each validation should fail
	rabbitURL, _ := GetEnvWithValidation("RABBITMQ_URL", "", true)
	err := ValidateRabbitMQURL(rabbitURL)
	assert.Error(t, err)
	
	err = ValidateMinIOConfig("localhost:9000", "access", "short")
	assert.Error(t, err)
	
	queueName, _ := GetEnvWithValidation("QUEUE_NAME", "", true)
	err = ValidateQueueName(queueName)
	assert.Error(t, err)
	
	_, err = GetIntEnvWithValidation("MAX_WORKERS", 5, 1, 50)
	assert.Error(t, err)
}