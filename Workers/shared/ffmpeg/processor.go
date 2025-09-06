package ffmpeg

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func ProcessWithTempFiles(inputData []byte, args []string) ([]byte, error) {
	tmpDir, err := os.MkdirTemp("", "ffmpeg-*")
	if err != nil {
		return nil, fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	inputPath := filepath.Join(tmpDir, "input.mp4")
	outputPath := filepath.Join(tmpDir, "output.mp4")

	if err := os.WriteFile(inputPath, inputData, 0600); err != nil {
		return nil, fmt.Errorf("write input: %w", err)
	}

	// Replace placeholders in args
	for i, arg := range args {
		if arg == "{input}" {
			args[i] = inputPath
		} else if arg == "{output}" {
			args[i] = outputPath
		}
	}

	cmd := exec.Command("ffmpeg", args...)
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("ffmpeg failed: %w", err)
	}

	return os.ReadFile(outputPath)
}