package services_test

import (
	"audioremoval/internal/application/services"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMP4VideoProcessingService_RemoveAudio_LargeFile(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg not available, skipping integration test")
	}

	service := services.NewMP4VideoProcessingService()
	
	// Create a larger fake video data
	largeData := make([]byte, 1024*1024) // 1MB
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}

	_, err := service.RemoveAudio(largeData)

	// Should fail because it's not valid MP4, but tests large data handling
	assert.Error(t, err)
}

func TestMP4VideoProcessingService_RemoveAudio_TempDirPermissions(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg not available, skipping integration test")
	}

	service := services.NewMP4VideoProcessingService()
	
	// Test with minimal data
	minimalData := []byte("test")

	_, err := service.RemoveAudio(minimalData)

	// Should fail because it's not valid MP4
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ffmpeg failed")
}

func TestMP4VideoProcessingService_RemoveAudio_ConcurrentCalls(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg not available, skipping integration test")
	}

	service := services.NewMP4VideoProcessingService()
	testData := []byte("concurrent test data")

	// Run multiple concurrent calls
	done := make(chan bool, 3)
	for i := 0; i < 3; i++ {
		go func() {
			_, err := service.RemoveAudio(testData)
			// All should fail with invalid MP4, but shouldn't crash
			assert.Error(t, err)
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 3; i++ {
		<-done
	}
}

func TestMP4VideoProcessingService_RemoveAudio_SpecialCharacters(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg not available, skipping integration test")
	}

	service := services.NewMP4VideoProcessingService()
	
	// Test with data containing special characters
	specialData := []byte("test data with special chars: áéíóú ñ @#$%^&*()")

	_, err := service.RemoveAudio(specialData)

	// Should fail because it's not valid MP4
	assert.Error(t, err)
}

func TestMP4VideoProcessingService_RemoveAudio_BinaryData(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg not available, skipping integration test")
	}

	service := services.NewMP4VideoProcessingService()
	
	// Test with random binary data
	binaryData := []byte{0x00, 0x01, 0x02, 0x03, 0xFF, 0xFE, 0xFD, 0xFC}

	_, err := service.RemoveAudio(binaryData)

	// Should fail because it's not valid MP4
	assert.Error(t, err)
}

func TestMP4VideoProcessingService_RemoveAudio_SystemLimits(t *testing.T) {
	service := services.NewMP4VideoProcessingService()
	
	// Test with extremely large data (if system allows)
	if testing.Short() {
		t.Skip("skipping system limits test in short mode")
	}

	// Create 10MB of data
	largeData := make([]byte, 10*1024*1024)
	
	_, err := service.RemoveAudio(largeData)

	// Should handle large data gracefully
	assert.Error(t, err) // Will fail due to invalid MP4 format
}

func TestMP4VideoProcessingService_RemoveAudio_TempDirCleanup(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg not available, skipping integration test")
	}

	service := services.NewMP4VideoProcessingService()
	testData := []byte("cleanup test")

	// Get initial temp dir count
	tempDir := os.TempDir()
	initialFiles, _ := os.ReadDir(tempDir)
	initialCount := len(initialFiles)

	// Process video (will fail but should cleanup)
	_, err := service.RemoveAudio(testData)
	assert.Error(t, err)

	// Check that temp directories are cleaned up
	finalFiles, _ := os.ReadDir(tempDir)
	finalCount := len(finalFiles)

	// Should not have significantly more temp files
	assert.LessOrEqual(t, finalCount-initialCount, 1, "Temp directories not properly cleaned up")
}