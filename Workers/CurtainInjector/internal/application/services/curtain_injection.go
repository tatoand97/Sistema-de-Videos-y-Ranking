package services

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

type CurtainInjectionService struct{}

func NewCurtainInjectionService() *CurtainInjectionService {
	return &CurtainInjectionService{}
}

func (s *CurtainInjectionService) InjectCurtains(inputData []byte, curtainInPath, curtainOutPath string) ([]byte, error) {
	tempDir, err := ioutil.TempDir("", "curtain_injection")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	inputFile := filepath.Join(tempDir, "input.mp4")
	outputFile := filepath.Join(tempDir, "output.mp4")

	if err := ioutil.WriteFile(inputFile, inputData, 0644); err != nil {
		return nil, fmt.Errorf("failed to write input file: %w", err)
	}

	// Debug: log the command being executed
	fmt.Printf("FFmpeg command: ffmpeg -i %s -i %s -i %s -filter_complex '[0:v][1:v][2:v]concat=n=3:v=1[outv]' -map '[outv]' -c:v libx264 -an -y %s\n", 
		curtainInPath, inputFile, curtainOutPath, outputFile)

	cmd := exec.Command("ffmpeg", "-i", curtainInPath, "-i", inputFile, "-i", curtainOutPath, 
		"-filter_complex", "[0:v][1:v][2:v]concat=n=3:v=1[outv]", 
		"-map", "[outv]", "-c:v", "libx264", "-an", "-y", outputFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("ffmpeg failed: %s, output: %s", err, string(output))
	}

	result, err := ioutil.ReadFile(outputFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read output file: %w", err)
	}

	return result, nil
}