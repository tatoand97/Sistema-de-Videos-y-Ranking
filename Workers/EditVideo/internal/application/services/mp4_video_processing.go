package services

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// MP4VideoProcessingService cumple con domain.VideoProcessingService
// esperado por el worker heredado de TrimVideo, pero en lugar de recortar,
// normaliza a 16:9 @ 720p manteniendo proporciones (scale+pad) y SAR=1.
type MP4VideoProcessingService struct{}

func NewMP4VideoProcessingService() *MP4VideoProcessingService { return &MP4VideoProcessingService{} }

// TrimToMaxSeconds recibe el video como []byte y devuelve []byte.
// Mantenemos la firma para cumplir la interfaz, pero ignoramos el "trim" y
// normalizamos a 1280x720 (16:9) sin distorsión.
func (s *MP4VideoProcessingService) TrimToMaxSeconds(input []byte, maxSeconds int) ([]byte, error) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return nil, fmt.Errorf("ffmpeg not found: %w", err)
	}

	tmpDir, err := os.MkdirTemp("", "edit-*")
	if err != nil {
		return nil, fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	inputPath := filepath.Join(tmpDir, "input.mp4")
	outputPath := filepath.Join(tmpDir, "output.mp4")

	if err := os.WriteFile(inputPath, input, 0600); err != nil {
		return nil, fmt.Errorf("write input: %w", err)
	}

	// Pipeline de normalización a 16:9 720p
	args := []string{
		"-y", "-i", inputPath,
		"-vf", "scale=1280:720:force_original_aspect_ratio=decrease,pad=1280:720:(ow-iw)/2:(oh-ih)/2,setsar=1",
		"-c:v", "libx264", "-preset", "veryfast", "-crf", "20",
		"-c:a", "aac", "-b:a", "128k",
		outputPath,
	}

	cmd := exec.Command("ffmpeg", args...)
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("ffmpeg failed: %w", err)
	}

	return os.ReadFile(outputPath)
}

// (Opcional) Si en algún punto llamas a otro método genérico en tu código,
// puedes dejar un helper para pruebas locales con rutas:
// NO forma parte de la interfaz, pero a veces es útil.
/*
func (s *MP4VideoProcessingService) Normalize16x9_720File(inputPath string) (string, error) {
	outPath := filepath.Join(os.TempDir(), fmt.Sprintf("edit_out_%d.mp4", time.Now().UnixNano()))
	args := []string{
		"-y", "-i", inputPath,
		"-vf", "scale=1280:720:force_original_aspect_ratio=decrease,pad=1280:720:(ow-iw)/2:(oh-ih)/2,setsar=1",
		"-c:v", "libx264", "-preset", "veryfast", "-crf", "20",
		"-c:a", "aac", "-b:a", "128k",
		outPath,
	}
	cmd := exec.Command("ffmpeg", args...)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil { return "", fmt.Errorf("ffmpeg: %w", err) }
	return outPath, nil
}
*/