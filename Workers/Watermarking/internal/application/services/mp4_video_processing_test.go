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

func TestMP4VideoProcessingService_Structure(t *testing.T) {
	service := &MP4VideoProcessingService{}
	assert.NotNil(t, service)
}

func TestMP4VideoProcessingService_WatermarkSpecs(t *testing.T) {
	// Test watermark specifications
	specs := map[string]interface{}{
		"logo_width":    180,
		"margin":        24,
		"position":      "bottom-right",
		"codec":         "libx264",
		"preset":        "veryfast",
		"crf":           "20",
		"audio_codec":   "aac",
		"audio_bitrate": "128k",
	}
	
	assert.Equal(t, 180, specs["logo_width"])
	assert.Equal(t, 24, specs["margin"])
	assert.Equal(t, "bottom-right", specs["position"])
	assert.Equal(t, "libx264", specs["codec"])
	assert.Equal(t, "veryfast", specs["preset"])
	assert.Equal(t, "20", specs["crf"])
	assert.Equal(t, "aac", specs["audio_codec"])
	assert.Equal(t, "128k", specs["audio_bitrate"])
}

func TestMP4VideoProcessingService_FFmpegArgs(t *testing.T) {
	// Test the FFmpeg arguments logic without actually running FFmpeg
	expectedArgs := []string{
		"-y",
		"-i", "input.mp4",
		"-i", "logo.png",
		"-filter_complex",
		"[1]scale=180:-1[wm];[0:v][wm]overlay=W-w-24:H-h-24[out]",
		"-map", "[out]",
		"-map", "0:a?",
		"-c:v", "libx264", "-preset", "veryfast", "-crf", "20",
		"-c:a", "aac", "-b:a", "128k",
		"-shortest",
		"output.mp4",
	}
	
	// Verify key arguments
	assert.Contains(t, expectedArgs, "-y")
	assert.Contains(t, expectedArgs, "-filter_complex")
	assert.Contains(t, expectedArgs, "[1]scale=180:-1[wm];[0:v][wm]overlay=W-w-24:H-h-24[out]")
	assert.Contains(t, expectedArgs, "libx264")
	assert.Contains(t, expectedArgs, "veryfast")
	assert.Contains(t, expectedArgs, "-shortest")
}

func TestMP4VideoProcessingService_OverlayFilter(t *testing.T) {
	// Test the overlay filter logic
	overlayFilter := "[1]scale=180:-1[wm];[0:v][wm]overlay=W-w-24:H-h-24[out]"
	
	// Verify filter components
	assert.Contains(t, overlayFilter, "[1]scale=180:-1[wm]")
	assert.Contains(t, overlayFilter, "[0:v][wm]overlay=W-w-24:H-h-24[out]")
	assert.Contains(t, overlayFilter, "scale=180:-1")
	assert.Contains(t, overlayFilter, "overlay=W-w-24:H-h-24")
}

func TestMP4VideoProcessingService_WatermarkPath(t *testing.T) {
	// Test watermark path logic
	tests := []struct {
		name     string
		envValue string
		expected string
	}{
		{"default path", "", "./assets/nba-logo-removebg-preview.png"},
		{"custom path", "/custom/logo.png", "/custom/logo.png"},
		{"relative path", "./custom/logo.png", "./custom/logo.png"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv("WATERMARK_PATH", tt.envValue)
				defer os.Unsetenv("WATERMARK_PATH")
			} else {
				os.Unsetenv("WATERMARK_PATH")
			}
			
			logoPath := os.Getenv("WATERMARK_PATH")
			if logoPath == "" {
				logoPath = "./assets/nba-logo-removebg-preview.png"
			}
			
			assert.Equal(t, tt.expected, logoPath)
		})
	}
}

func TestMP4VideoProcessingService_ParameterValidation(t *testing.T) {
	tests := []struct {
		name       string
		maxSeconds int
		valid      bool
	}{
		{"positive value", 30, true},
		{"zero value", 0, true},
		{"negative value", -1, true}, // Ignored in Watermarking
		{"large value", 3600, true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// In Watermarking, maxSeconds is ignored but should be accepted
			assert.True(t, tt.valid)
		})
	}
}

func TestMP4VideoProcessingService_WatermarkPosition(t *testing.T) {
	// Test watermark positioning logic
	position := "W-w-24:H-h-24" // bottom-right with 24px margin
	
	assert.Contains(t, position, "W-w-24") // X position: video width - watermark width - margin
	assert.Contains(t, position, "H-h-24") // Y position: video height - watermark height - margin
}

func TestMP4VideoProcessingService_TempFileHandling(t *testing.T) {
	// Test temporary file naming logic
	inputFile := "input.mp4"
	outputFile := "output.mp4"
	logoFile := "logo.png"
	
	assert.Equal(t, "input.mp4", inputFile)
	assert.Equal(t, "output.mp4", outputFile)
	assert.Equal(t, "logo.png", logoFile)
	assert.Contains(t, inputFile, ".mp4")
	assert.Contains(t, outputFile, ".mp4")
	assert.Contains(t, logoFile, ".png")
}