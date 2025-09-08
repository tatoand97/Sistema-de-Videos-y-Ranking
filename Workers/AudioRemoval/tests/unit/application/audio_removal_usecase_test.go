package application_test

import (
	"context"
	"errors"
	"testing"
	"github.com/stretchr/testify/assert"
)

type AudioRemovalUseCase struct {
	processor AudioProcessor
	storage   StorageService
}

type AudioProcessor interface {
	RemoveAudio(ctx context.Context, inputPath, outputPath string) error
}

type StorageService interface {
	Upload(ctx context.Context, filePath string) (string, error)
	Download(ctx context.Context, url string) (string, error)
}

func NewAudioRemovalUseCase(processor AudioProcessor, storage StorageService) *AudioRemovalUseCase {
	return &AudioRemovalUseCase{
		processor: processor,
		storage:   storage,
	}
}

func (uc *AudioRemovalUseCase) ProcessVideo(ctx context.Context, videoID, inputURL string) (string, error) {
	inputPath, err := uc.storage.Download(ctx, inputURL)
	if err != nil {
		return "", err
	}

	outputPath := inputPath + "_no_audio.mp4"
	if err := uc.processor.RemoveAudio(ctx, inputPath, outputPath); err != nil {
		return "", err
	}

	return uc.storage.Upload(ctx, outputPath)
}

type MockAudioProcessor struct {
	calls []string
	err   error
}

func (m *MockAudioProcessor) RemoveAudio(ctx context.Context, inputPath, outputPath string) error {
	m.calls = append(m.calls, "RemoveAudio")
	return m.err
}

type MockStorageService struct {
	calls      []string
	downloadResult string
	uploadResult   string
	err            error
}

func (m *MockStorageService) Upload(ctx context.Context, filePath string) (string, error) {
	m.calls = append(m.calls, "Upload")
	return m.uploadResult, m.err
}

func (m *MockStorageService) Download(ctx context.Context, url string) (string, error) {
	m.calls = append(m.calls, "Download")
	return m.downloadResult, m.err
}

func TestAudioRemovalUseCase_ProcessVideo(t *testing.T) {
	tests := []struct {
		name     string
		videoID  string
		inputURL string
		mockStorage *MockStorageService
		mockProcessor *MockAudioProcessor
		wantErr  bool
		want     string
	}{
		{
			name:     "successful processing",
			videoID:  "video-123",
			inputURL: "https://example.com/video.mp4",
			mockStorage: &MockStorageService{
				downloadResult: "/tmp/video.mp4",
				uploadResult: "https://example.com/processed.mp4",
			},
			mockProcessor: &MockAudioProcessor{},
			wantErr: false,
			want:    "https://example.com/processed.mp4",
		},
		{
			name:     "download fails",
			videoID:  "video-123",
			inputURL: "https://example.com/video.mp4",
			mockStorage: &MockStorageService{
				err: errors.New("download failed"),
			},
			mockProcessor: &MockAudioProcessor{},
			wantErr: true,
		},
		{
			name:     "processing fails",
			videoID:  "video-123",
			inputURL: "https://example.com/video.mp4",
			mockStorage: &MockStorageService{
				downloadResult: "/tmp/video.mp4",
			},
			mockProcessor: &MockAudioProcessor{
				err: errors.New("processing failed"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usecase := NewAudioRemovalUseCase(tt.mockProcessor, tt.mockStorage)
			
			result, err := usecase.ProcessVideo(context.Background(), tt.videoID, tt.inputURL)
			
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, result)
			}
		})
	}
}