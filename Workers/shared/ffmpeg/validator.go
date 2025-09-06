package ffmpeg

import (
	"fmt"
	"os/exec"
)

func ValidateFFmpeg() error {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return fmt.Errorf("ffmpeg not found: %w", err)
	}
	return nil
}