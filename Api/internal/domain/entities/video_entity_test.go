package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVideo_TableName(t *testing.T) {
	video := &Video{}
	assert.Equal(t, "video", video.TableName())
}

func TestVideoStatus_Constants(t *testing.T) {
	tests := []struct {
		name     string
		status   VideoStatus
		expected string
	}{
		{"StatusUploaded", StatusUploaded, "UPLOADED"},
		{"StatusTrimming", StatusTrimming, "TRIMMING"},
		{"StatusAdjustingRes", StatusAdjustingRes, "ADJUSTING_RESOLUTION"},
		{"StatusAddingWatermark", StatusAddingWatermark, "ADDING_WATERMARK"},
		{"StatusRemovingAudio", StatusRemovingAudio, "REMOVING_AUDIO"},
		{"StatusAddingIntroOutro", StatusAddingIntroOutro, "ADDING_INTRO_OUTRO"},
		{"StatusProcessed", StatusProcessed, "PROCESSED"},
		{"StatusFailed", StatusFailed, "FAILED"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.status))
		})
	}
}

func TestAllVideoStatuses(t *testing.T) {
	statuses := AllVideoStatuses()
	
	expectedStatuses := []VideoStatus{
		StatusUploaded,
		StatusTrimming,
		StatusAdjustingRes,
		StatusAddingWatermark,
		StatusRemovingAudio,
		StatusAddingIntroOutro,
		StatusProcessed,
		StatusFailed,
	}
	
	assert.Equal(t, expectedStatuses, statuses)
	assert.Len(t, statuses, 8)
}

func TestVideo_Creation(t *testing.T) {
	processedFile := "processed.mp4"
	video := &Video{
		VideoID:       1,
		UserID:        123,
		Title:         "Test Video",
		OriginalFile:  "original.mp4",
		ProcessedFile: &processedFile,
		Status:        string(StatusUploaded),
	}
	
	assert.Equal(t, uint(1), video.VideoID)
	assert.Equal(t, uint(123), video.UserID)
	assert.Equal(t, "Test Video", video.Title)
	assert.Equal(t, "original.mp4", video.OriginalFile)
	assert.Equal(t, "processed.mp4", *video.ProcessedFile)
	assert.Equal(t, "UPLOADED", video.Status)
}