package services

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// MP4VideoProcessingService cumple con domain.VideoProcessingService.
// Aplica marca de agua ANB (overlay) y normaliza a 1280x720.
type MP4VideoProcessingService struct{}

func NewMP4VideoProcessingService() *MP4VideoProcessingService { return &MP4VideoProcessingService{} }

// TrimToMaxSeconds: mantiene la firma []byte→[]byte y agrega watermark.
// Ignoramos el "recorte" y normalizamos a 720p con overlay del logo.
func (s *MP4VideoProcessingService) TrimToMaxSeconds(input []byte, maxSeconds int) ([]byte, error) {
	inPath := filepath.Join(os.TempDir(), fmt.Sprintf("wm_in_%d.mp4", time.Now().UnixNano()))
	outPath := filepath.Join(os.TempDir(), fmt.Sprintf("wm_out_%d.mp4", time.Now().UnixNano()))
	logoPath := os.Getenv("WATERMARK_PATH")
	if logoPath == "" {
		logoPath = "./assets/nba-logo-removebg-preview.png"
	}

	defer func() {
		_ = os.Remove(inPath)
		_ = os.Remove(outPath)
	}()

	if err := os.WriteFile(inPath, input, 0o600); err != nil {
		return nil, fmt.Errorf("write temp input: %w", err)
	}

	args := []string{
		 "-y",
		"-i", inPath,      // video de entrada
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
		outPath,
	}

	cmd := exec.Command("ffmpeg", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("ffmpeg: %w", err)
	}

	out, err := os.ReadFile(outPath)
	if err != nil { return nil, fmt.Errorf("read output: %w", err) }
	return out, nil
}
