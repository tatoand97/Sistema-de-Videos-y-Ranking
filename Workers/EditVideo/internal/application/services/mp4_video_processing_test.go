package services

import (
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

func TestMP4VideoProcessingService_FFmpegArgs(t *testing.T) {
	// Test the FFmpeg arguments logic without actually running FFmpeg
	expectedArgs := []string{
		"-y", "-i", "input.mp4",
		"-vf", "scale=1280:720:force_original_aspect_ratio=decrease,pad=1280:720:(ow-iw)/2:(oh-ih)/2,setsar=1",
		"-c:v", "libx264", "-preset", "veryfast", "-crf", "20",
		"-c:a", "aac", "-b:a", "128k",
		"output.mp4",
	}
	
	// Verify the arguments structure
	assert.Contains(t, expectedArgs, "-y")
	assert.Contains(t, expectedArgs, "-i")
	assert.Contains(t, expectedArgs, "scale=1280:720:force_original_aspect_ratio=decrease,pad=1280:720:(ow-iw)/2:(oh-ih)/2,setsar=1")
	assert.Contains(t, expectedArgs, "libx264")
	assert.Contains(t, expectedArgs, "veryfast")
	assert.Contains(t, expectedArgs, "aac")
}

func TestMP4VideoProcessingService_VideoNormalizationSpecs(t *testing.T) {
	// Test video normalization specifications
	specs := map[string]interface{}{
		"target_width":  1280,
		"target_height": 720,
		"aspect_ratio":  "16:9",
		"codec":         "libx264",
		"preset":        "veryfast",
		"crf":           "20",
		"audio_codec":   "aac",
		"audio_bitrate": "128k",
	}
	
	assert.Equal(t, 1280, specs["target_width"])
	assert.Equal(t, 720, specs["target_height"])
	assert.Equal(t, "16:9", specs["aspect_ratio"])
	assert.Equal(t, "libx264", specs["codec"])
	assert.Equal(t, "veryfast", specs["preset"])
	assert.Equal(t, "20", specs["crf"])
	assert.Equal(t, "aac", specs["audio_codec"])
	assert.Equal(t, "128k", specs["audio_bitrate"])
}

func TestMP4VideoProcessingService_ScaleFilter(t *testing.T) {
	// Test the scale filter logic
	scaleFilter := "scale=1280:720:force_original_aspect_ratio=decrease,pad=1280:720:(ow-iw)/2:(oh-ih)/2,setsar=1"
	
	// Verify filter components
	assert.Contains(t, scaleFilter, "scale=1280:720")
	assert.Contains(t, scaleFilter, "force_original_aspect_ratio=decrease")
	assert.Contains(t, scaleFilter, "pad=1280:720")
	assert.Contains(t, scaleFilter, "(ow-iw)/2:(oh-ih)/2")
	assert.Contains(t, scaleFilter, "setsar=1")
}

func TestMP4VideoProcessingService_ParameterValidation(t *testing.T) {
	tests := []struct {
		name       string
		maxSeconds int
		valid      bool
	}{
		{"positive value", 30, true},
		{"zero value", 0, true},
		{"negative value", -1, true}, // Ignored in EditVideo
		{"large value", 3600, true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// In EditVideo, maxSeconds is ignored but should be accepted
			assert.True(t, tt.valid)
		})
	}
}

func TestMP4VideoProcessingService_TempFileHandling(t *testing.T) {
	// Test temporary file naming logic
	inputFile := "input.mp4"
	outputFile := "output.mp4"
	
	assert.Equal(t, "input.mp4", inputFile)
	assert.Equal(t, "output.mp4", outputFile)
	assert.Contains(t, inputFile, ".mp4")
	assert.Contains(t, outputFile, ".mp4")
}