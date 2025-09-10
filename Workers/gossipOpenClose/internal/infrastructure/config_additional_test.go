package infrastructure

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_DefaultValues(t *testing.T) {
	// Clear environment variables
	envVars := []string{
		"POSTGRES_HOST", "POSTGRES_PORT", "POSTGRES_USER", "POSTGRES_PASSWORD", "POSTGRES_DB",
		"MINIO_ENDPOINT", "MINIO_ACCESS_KEY", "MINIO_SECRET_KEY", "MINIO_USE_SSL",
		"RABBITMQ_URL", "RAW_BUCKET", "PROCESSED_BUCKET", "STATE_QUEUE",
		"INTRO_SECONDS", "OUTRO_SECONDS", "TARGET_WIDTH", "TARGET_HEIGHT", "FPS",
		"LOGO_PATH",
	}
	
	for _, env := range envVars {
		os.Unsetenv(env)
	}
	
	config := NewConfig()
	
	// Test default values
	assert.Equal(t, "localhost", config.PostgresHost)
	assert.Equal(t, "5432", config.PostgresPort)
	assert.Equal(t, "postgres", config.PostgresUser)
	assert.Equal(t, "password", config.PostgresPassword)
	assert.Equal(t, "videorank", config.PostgresDB)
	
	assert.Equal(t, "localhost:9000", config.MinioEndpoint)
	assert.Equal(t, "minioadmin", config.MinioAccessKey)
	assert.Equal(t, "minioadmin", config.MinioSecretKey)
	assert.False(t, config.MinioUseSSL)
	
	assert.Equal(t, "amqp://guest:guest@localhost:5672/", config.RabbitMQURL)
	assert.Equal(t, "raw-videos", config.RawBucket)
	assert.Equal(t, "processed-videos", config.ProcessedBucket)
	assert.Equal(t, "state-machine", config.StateQueue)
	
	assert.Equal(t, 2.5, config.IntroSeconds)
	assert.Equal(t, 2.5, config.OutroSeconds)
	assert.Equal(t, 1920, config.TargetWidth)
	assert.Equal(t, 1080, config.TargetHeight)
	assert.Equal(t, 30, config.FPS)
	assert.Equal(t, "./assets/nba-logo-removebg-preview.png", config.LogoPath)
}

func TestConfig_EnvironmentOverrides(t *testing.T) {
	// Set environment variables
	testEnvVars := map[string]string{
		"POSTGRES_HOST":     "test-host",
		"POSTGRES_PORT":     "5433",
		"POSTGRES_USER":     "test-user",
		"POSTGRES_PASSWORD": "test-password",
		"POSTGRES_DB":       "test-db",
		"MINIO_ENDPOINT":    "test-minio:9001",
		"MINIO_ACCESS_KEY":  "test-access",
		"MINIO_SECRET_KEY":  "test-secret",
		"MINIO_USE_SSL":     "true",
		"RABBITMQ_URL":      "amqp://test:test@test-rabbit:5672/",
		"RAW_BUCKET":        "test-raw",
		"PROCESSED_BUCKET":  "test-processed",
		"STATE_QUEUE":       "test-state",
		"INTRO_SECONDS":     "3.0",
		"OUTRO_SECONDS":     "1.5",
		"TARGET_WIDTH":      "1280",
		"TARGET_HEIGHT":     "720",
		"FPS":               "25",
		"LOGO_PATH":         "/custom/logo.png",
	}
	
	for key, value := range testEnvVars {
		t.Setenv(key, value)
	}
	
	config := NewConfig()
	
	// Test overridden values
	assert.Equal(t, "test-host", config.PostgresHost)
	assert.Equal(t, "5433", config.PostgresPort)
	assert.Equal(t, "test-user", config.PostgresUser)
	assert.Equal(t, "test-password", config.PostgresPassword)
	assert.Equal(t, "test-db", config.PostgresDB)
	
	assert.Equal(t, "test-minio:9001", config.MinioEndpoint)
	assert.Equal(t, "test-access", config.MinioAccessKey)
	assert.Equal(t, "test-secret", config.MinioSecretKey)
	assert.True(t, config.MinioUseSSL)
	
	assert.Equal(t, "amqp://test:test@test-rabbit:5672/", config.RabbitMQURL)
	assert.Equal(t, "test-raw", config.RawBucket)
	assert.Equal(t, "test-processed", config.ProcessedBucket)
	assert.Equal(t, "test-state", config.StateQueue)
	
	assert.Equal(t, 3.0, config.IntroSeconds)
	assert.Equal(t, 1.5, config.OutroSeconds)
	assert.Equal(t, 1280, config.TargetWidth)
	assert.Equal(t, 720, config.TargetHeight)
	assert.Equal(t, 25, config.FPS)
	assert.Equal(t, "/custom/logo.png", config.LogoPath)
}

func TestConfig_InvalidEnvironmentValues(t *testing.T) {
	// Set invalid environment variables
	testEnvVars := map[string]string{
		"MINIO_USE_SSL":  "invalid-bool",
		"INTRO_SECONDS":  "invalid-float",
		"OUTRO_SECONDS":  "not-a-number",
		"TARGET_WIDTH":   "invalid-int",
		"TARGET_HEIGHT":  "not-an-int",
		"FPS":            "invalid-fps",
	}
	
	for key, value := range testEnvVars {
		t.Setenv(key, value)
	}
	
	config := NewConfig()
	
	// Should fall back to defaults for invalid values
	assert.False(t, config.MinioUseSSL) // Default false
	assert.Equal(t, 2.5, config.IntroSeconds) // Default 2.5
	assert.Equal(t, 2.5, config.OutroSeconds) // Default 2.5
	assert.Equal(t, 1920, config.TargetWidth) // Default 1920
	assert.Equal(t, 1080, config.TargetHeight) // Default 1080
	assert.Equal(t, 30, config.FPS) // Default 30
}

func TestConfig_EmptyEnvironmentValues(t *testing.T) {
	// Set empty environment variables
	testEnvVars := map[string]string{
		"POSTGRES_HOST":     "",
		"POSTGRES_PORT":     "",
		"MINIO_USE_SSL":     "",
		"INTRO_SECONDS":     "",
		"TARGET_WIDTH":      "",
	}
	
	for key, value := range testEnvVars {
		t.Setenv(key, value)
	}
	
	config := NewConfig()
	
	// Should use defaults for empty values
	assert.Equal(t, "localhost", config.PostgresHost)
	assert.Equal(t, "5432", config.PostgresPort)
	assert.False(t, config.MinioUseSSL)
	assert.Equal(t, 2.5, config.IntroSeconds)
	assert.Equal(t, 1920, config.TargetWidth)
}

func TestConfig_BooleanParsing(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{"true lowercase", "true", true},
		{"TRUE uppercase", "TRUE", true},
		{"True mixed case", "True", true},
		{"1 as true", "1", true},
		{"false lowercase", "false", false},
		{"FALSE uppercase", "FALSE", false},
		{"0 as false", "0", false},
		{"empty as false", "", false},
		{"invalid as false", "invalid", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("MINIO_USE_SSL", tt.value)
			config := NewConfig()
			assert.Equal(t, tt.expected, config.MinioUseSSL)
		})
	}
}

func TestConfig_NumericParsing(t *testing.T) {
	tests := []struct {
		name         string
		envVar       string
		value        string
		expectedInt  int
		expectedFloat float64
	}{
		{"valid integer", "TARGET_WIDTH", "1280", 1280, 0},
		{"valid float", "INTRO_SECONDS", "3.5", 0, 3.5},
		{"zero integer", "TARGET_WIDTH", "0", 0, 0},
		{"zero float", "INTRO_SECONDS", "0.0", 0, 0.0},
		{"negative integer", "TARGET_WIDTH", "-100", -100, 0},
		{"negative float", "INTRO_SECONDS", "-1.5", 0, -1.5},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(tt.envVar, tt.value)
			config := NewConfig()
			
			if tt.envVar == "TARGET_WIDTH" {
				assert.Equal(t, tt.expectedInt, config.TargetWidth)
			} else if tt.envVar == "INTRO_SECONDS" {
				assert.Equal(t, tt.expectedFloat, config.IntroSeconds)
			}
		})
	}
}

func TestConfig_DatabaseConnectionString(t *testing.T) {
	t.Setenv("POSTGRES_HOST", "db-host")
	t.Setenv("POSTGRES_PORT", "5433")
	t.Setenv("POSTGRES_USER", "dbuser")
	t.Setenv("POSTGRES_PASSWORD", "dbpass")
	t.Setenv("POSTGRES_DB", "testdb")
	
	config := NewConfig()
	
	// Test that all database config is properly set
	assert.Equal(t, "db-host", config.PostgresHost)
	assert.Equal(t, "5433", config.PostgresPort)
	assert.Equal(t, "dbuser", config.PostgresUser)
	assert.Equal(t, "dbpass", config.PostgresPassword)
	assert.Equal(t, "testdb", config.PostgresDB)
}

func TestConfig_MinioConfiguration(t *testing.T) {
	t.Setenv("MINIO_ENDPOINT", "minio.example.com:9000")
	t.Setenv("MINIO_ACCESS_KEY", "access123")
	t.Setenv("MINIO_SECRET_KEY", "secret456")
	t.Setenv("MINIO_USE_SSL", "true")
	
	config := NewConfig()
	
	assert.Equal(t, "minio.example.com:9000", config.MinioEndpoint)
	assert.Equal(t, "access123", config.MinioAccessKey)
	assert.Equal(t, "secret456", config.MinioSecretKey)
	assert.True(t, config.MinioUseSSL)
}

func TestConfig_ProcessingParameters(t *testing.T) {
	t.Setenv("INTRO_SECONDS", "1.5")
	t.Setenv("OUTRO_SECONDS", "3.0")
	t.Setenv("TARGET_WIDTH", "1280")
	t.Setenv("TARGET_HEIGHT", "720")
	t.Setenv("FPS", "25")
	t.Setenv("LOGO_PATH", "/custom/path/logo.png")
	
	config := NewConfig()
	
	assert.Equal(t, 1.5, config.IntroSeconds)
	assert.Equal(t, 3.0, config.OutroSeconds)
	assert.Equal(t, 1280, config.TargetWidth)
	assert.Equal(t, 720, config.TargetHeight)
	assert.Equal(t, 25, config.FPS)
	assert.Equal(t, "/custom/path/logo.png", config.LogoPath)
}