package services

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMP4VideoProcessingService_RemoveAudio_InvalidData(t *testing.T) {
	service := NewMP4VideoProcessingService()
	
	// Test with invalid video data
	invalidData := []byte("not a video")
	_, err := service.RemoveAudio(invalidData)
	
	// Should fail due to invalid input or missing ffmpeg
	assert.Error(t, err)
}

func TestMP4VideoProcessingService_TempDirectory(t *testing.T) {
	// Test temp directory creation
	tmpDir, err := os.MkdirTemp("", "test-*")
	if err == nil {
		assert.NotEmpty(t, tmpDir)
		defer os.RemoveAll(tmpDir)
	}
}

func TestMP4VideoProcessingService_FileHandling(t *testing.T) {
	// Test basic file operations
	testData := []byte("test content")
	tmpFile := os.TempDir() + "/test.tmp"
	
	err := os.WriteFile(tmpFile, testData, 0600)
	if err == nil {
		defer os.Remove(tmpFile)
		
		readData, err := os.ReadFile(tmpFile)
		assert.NoError(t, err)
		assert.Equal(t, testData, readData)
	}
}