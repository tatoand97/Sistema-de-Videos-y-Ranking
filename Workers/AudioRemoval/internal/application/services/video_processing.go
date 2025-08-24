package services

import (
	"os"
	"os/exec"
)

type FFmpegVideoProcessingService struct{}

func NewFFmpegVideoProcessingService() *FFmpegVideoProcessingService {
	return &FFmpegVideoProcessingService{}
}

func (s *FFmpegVideoProcessingService) RemoveAudio(inputData []byte) ([]byte, error) {
	tmpInput := "/tmp/input_video.mp4"
	tmpOutput := "/tmp/output_video.mp4"

	if err := os.WriteFile(tmpInput, inputData, 0644); err != nil {
		return nil, err
	}
	defer os.Remove(tmpInput)

	cmd := exec.Command("ffmpeg", "-i", tmpInput, "-an", "-c:v", "copy", tmpOutput, "-y")
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	defer os.Remove(tmpOutput)

	return os.ReadFile(tmpOutput)
}