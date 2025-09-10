package domain_test

import (
	"api/internal/domain/entities"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestVideo_TableName(t *testing.T) {
	video := entities.Video{}
	assert.Equal(t, "video", video.TableName())
}

func TestVideoStatus_Constants(t *testing.T) {
	tests := []struct {
		name     string
		status   entities.VideoStatus
		expected string
	}{
		{"StatusUploaded", entities.StatusUploaded, "UPLOADED"},
		{"StatusTrimming", entities.StatusTrimming, "TRIMMING"},
		{"StatusAdjustingRes", entities.StatusAdjustingRes, "ADJUSTING_RESOLUTION"},
		{"StatusAddingWatermark", entities.StatusAddingWatermark, "ADDING_WATERMARK"},
		{"StatusRemovingAudio", entities.StatusRemovingAudio, "REMOVING_AUDIO"},
		{"StatusAddingIntroOutro", entities.StatusAddingIntroOutro, "ADDING_INTRO_OUTRO"},
		{"StatusProcessed", entities.StatusProcessed, "PROCESSED"},
		{"StatusPublished", entities.StatusPublished, "PUBLISHED"},
		{"StatusFailed", entities.StatusFailed, "FAILED"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.status))
		})
	}
}

func TestAllVideoStatuses(t *testing.T) {
	statuses := entities.AllVideoStatuses()

	expected := []entities.VideoStatus{
		entities.StatusUploaded,
		entities.StatusTrimming,
		entities.StatusAdjustingRes,
		entities.StatusAddingWatermark,
		entities.StatusRemovingAudio,
		entities.StatusAddingIntroOutro,
		entities.StatusProcessed,
		entities.StatusPublished,
		entities.StatusFailed,
	}

	assert.Equal(t, expected, statuses)
	assert.Len(t, statuses, 9)
}

func TestVideo_Creation(t *testing.T) {
	now := time.Now()
	processedAt := now.Add(time.Hour)

	video := entities.Video{
		VideoID:       1,
		UserID:        123,
		Title:         "Test Video",
		OriginalFile:  "https://example.com/original.mp4",
		ProcessedFile: stringPtr("https://example.com/processed.mp4"),
		Status:        string(entities.StatusProcessed),
		UploadedAt:    now,
		ProcessedAt:   &processedAt,
	}

	assert.Equal(t, uint(1), video.VideoID)
	assert.Equal(t, uint(123), video.UserID)
	assert.Equal(t, "Test Video", video.Title)
	assert.Equal(t, "https://example.com/original.mp4", video.OriginalFile)
	assert.NotNil(t, video.ProcessedFile)
	assert.Equal(t, "https://example.com/processed.mp4", *video.ProcessedFile)
	assert.Equal(t, string(entities.StatusProcessed), video.Status)
	assert.Equal(t, now, video.UploadedAt)
	assert.NotNil(t, video.ProcessedAt)
	assert.Equal(t, processedAt, *video.ProcessedAt)
}

func stringPtr(s string) *string {
	return &s
}
