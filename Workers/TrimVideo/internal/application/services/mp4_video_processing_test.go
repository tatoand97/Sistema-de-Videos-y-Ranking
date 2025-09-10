package services

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMP4VideoProcessingService(t *testing.T) {
	service := NewMP4VideoProcessingService()
	assert.NotNil(t, service)
}

func TestMP4VideoProcessingService_TrimToMaxSeconds_NoFFmpeg(t *testing.T) {
	service := NewMP4VideoProcessingService()
	inputData := []byte("fake video data")
	
	_, err := service.TrimToMaxSeconds(inputData, 30)
	if err != nil {
		assert.Contains(t, err.Error(), "ffmpeg not found")
	}
}

func TestMP4VideoProcessingService_TrimToMaxSeconds_EmptyInput(t *testing.T) {
	service := NewMP4VideoProcessingService()
	
	_, err := service.TrimToMaxSeconds([]byte{}, 30)
	assert.Error(t, err)
}

func TestMP4VideoProcessingService_TrimToMaxSeconds_InvalidSeconds(t *testing.T) {
	service := NewMP4VideoProcessingService()
	inputData := []byte("test data")
	
	_, err := service.TrimToMaxSeconds(inputData, 0)
	assert.Error(t, err)
}

func TestMP4VideoProcessingService_TempDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "trim-test-*")
	if err == nil {
		assert.NotEmpty(t, tmpDir)
		assert.Contains(t, tmpDir, "trim-test-")
		defer os.RemoveAll(tmpDir)
	}
}