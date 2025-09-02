package services

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

type MP4VideoProcessingService struct{}

func NewMP4VideoProcessingService() *MP4VideoProcessingService { return &MP4VideoProcessingService{} }

// TrimToMaxSeconds recorta el video a maxSeconds usando ffmpeg (-c copy).
func (s *MP4VideoProcessingService) TrimToMaxSeconds(inputData []byte, maxSeconds int) ([]byte, error) {
	inFile, err := os.CreateTemp("/tmp", "in-*.mp4")
	if err != nil { return nil, err }
	defer os.Remove(inFile.Name())
	if _, err := inFile.Write(inputData); err != nil { return nil, err }
	inFile.Close()

	outFile, err := os.CreateTemp("/tmp", "out-*.mp4")
	if err != nil { return nil, err }
	outName := outFile.Name()
	outFile.Close()
	os.Remove(outName) // ffmpeg lo creará

	args := []string{"-y", "-i", inFile.Name(), "-t", fmt.Sprintf("%d", maxSeconds), "-c", "copy", outName}
	cmd := exec.Command("ffmpeg", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("ffmpeg trim error: %v; %s", err, stderr.String())
	}

	out, err := os.ReadFile(outName)
	if err != nil { return nil, err }
	defer os.Remove(outName)
	return out, nil
}
