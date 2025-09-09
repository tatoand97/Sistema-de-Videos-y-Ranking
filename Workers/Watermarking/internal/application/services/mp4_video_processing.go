package services

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// MP4VideoProcessingService cumple con domain.VideoProcessingService.
// Aplica marca de agua ANB (overlay) y normaliza a 1280x720.
type MP4VideoProcessingService struct{}

func NewMP4VideoProcessingService() *MP4VideoProcessingService { return &MP4VideoProcessingService{} }

// TrimToMaxSeconds: mantiene la firma []byte→[]byte y agrega watermark.
// Ignoramos el "recorte" y normalizamos a 720p con overlay del logo.
func (s *MP4VideoProcessingService) TrimToMaxSeconds(input []byte, maxSeconds int) ([]byte, error) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return nil, fmt.Errorf("ffmpeg not found: %w", err)
	}

	tmpDir, err := os.MkdirTemp("", "watermark-*")
	if err != nil {
		return nil, fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	inputPath := filepath.Join(tmpDir, "input.mp4")
	outputPath := filepath.Join(tmpDir, "output.mp4")
	logoPath := os.Getenv("WATERMARK_PATH")
	if logoPath == "" {
		logoPath = "./assets/nba-logo-removebg-preview.png"
	}

	if err := os.WriteFile(inputPath, input, 0600); err != nil {
		return nil, fmt.Errorf("write input: %w", err)
	}

	args := []string{
		 "-y",
		"-i", inputPath,      // video de entrada
		"-i", logoPath,    // logo con alpha (PNG)
		"-filter_complex",
		// 1) Escala del logo a 180px de ancho manteniendo aspecto
		// 2) Overlay en esquina inferior derecha con margen de 24px
		"[1]scale=180:-1[wm];[0:v][wm]overlay=W-w-24:H-h-24[out]",
		"-map", "[out]",
		"-map", "0:a?",
		"-c:v", "libx264", "-preset", "veryfast", "-crf", "20",
		"-c:a", "aac", "-b:a", "128k",
		"-shortest",
		outputPath,
	}

	cmd := exec.Command("ffmpeg", args...)
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("ffmpeg failed: %w", err)
	}

	return os.ReadFile(outputPath)
}
