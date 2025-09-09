package services

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type MP4VideoProcessingService struct{}

func NewMP4VideoProcessingService() *MP4VideoProcessingService { return &MP4VideoProcessingService{} }

// TrimToMaxSeconds recorta el video a maxSeconds usando ffmpeg (-c copy).
func (s *MP4VideoProcessingService) TrimToMaxSeconds(inputData []byte, maxSeconds int) ([]byte, error) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return nil, fmt.Errorf("ffmpeg not found: %w", err)
	}

	tmpDir, err := os.MkdirTemp("", "trim-*")
	if err != nil {
		return nil, fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	inputPath := filepath.Join(tmpDir, "input.mp4")
	outputPath := filepath.Join(tmpDir, "output.mp4")

	if err := os.WriteFile(inputPath, inputData, 0600); err != nil {
		return nil, fmt.Errorf("write input: %w", err)
	}

	args := []string{"-y", "-i", inputPath, "-t", fmt.Sprintf("%d", maxSeconds), "-c", "copy", outputPath}
	cmd := exec.Command("ffmpeg", args...)
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("ffmpeg failed: %w", err)
	}

	return os.ReadFile(outputPath)
}
