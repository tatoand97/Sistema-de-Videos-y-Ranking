package services

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type MP4VideoProcessingService struct{}

func NewMP4VideoProcessingService() *MP4VideoProcessingService {
	return &MP4VideoProcessingService{}
}

func (s *MP4VideoProcessingService) RemoveAudio(inputData []byte) ([]byte, error) {
	// Validate FFmpeg is available
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return nil, fmt.Errorf("ffmpeg not found: %w", err)
	}

	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "audio-removal-*")
	if err != nil {
		return nil, fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Write input file
	inputPath := filepath.Join(tmpDir, "input.mp4")
	if err := os.WriteFile(inputPath, inputData, 0600); err != nil {
		return nil, fmt.Errorf("write input file: %w", err)
	}

	// Output file
	outputPath := filepath.Join(tmpDir, "output.mp4")

	// Use FFmpeg to remove audio (copy video stream only)
	cmd := exec.Command("ffmpeg",
		"-i", inputPath,
		"-c:v", "copy", // Copy video stream without re-encoding
		"-an", // Remove audio stream
		"-y", // Overwrite output file
		outputPath,
	)

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("ffmpeg failed: %w", err)
	}

	// Read output file
	outputData, err := os.ReadFile(outputPath)
	if err != nil {
		return nil, fmt.Errorf("read output file: %w", err)
	}

	return outputData, nil
}
