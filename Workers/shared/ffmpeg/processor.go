package ffmpeg

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

	// Validate and replace placeholders in args
	for i, arg := range args {
		if strings.Contains(arg, "../") || strings.Contains(arg, "..\\") {
			return nil, fmt.Errorf("invalid argument contains path traversal: %s", arg)
		}
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