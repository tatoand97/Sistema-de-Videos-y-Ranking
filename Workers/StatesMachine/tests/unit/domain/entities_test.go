package domain

import (
	"statesmachine/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVideo_TableName(t *testing.T) {
	video := domain.Video{}
	assert.Equal(t, "video", video.TableName())
}

func TestVideo_Creation(t *testing.T) {
	video := domain.Video{
		ID:           123,
		OriginalFile: "test.mp4",
		Status:       "UPLOADED",
	}
	
	assert.Equal(t, uint(123), video.ID)
	assert.Equal(t, "test.mp4", video.OriginalFile)
	assert.Equal(t, "UPLOADED", video.Status)
}

func TestVideoStatus_Constants(t *testing.T) {
	tests := []struct {
		name     string
		status   domain.VideoStatus
		expected string
	}{
		{"StatusUploaded", domain.StatusUploaded, "UPLOADED"},
		{"StatusTrimming", domain.StatusTrimming, "TRIMMING"},
		{"StatusAdjustingRes", domain.StatusAdjustingRes, "ADJUSTING_RESOLUTION"},
		{"StatusAddingWatermark", domain.StatusAddingWatermark, "ADDING_WATERMARK"},
		{"StatusRemovingAudio", domain.StatusRemovingAudio, "REMOVING_AUDIO"},
		{"StatusAddingIntroOutro", domain.StatusAddingIntroOutro, "ADDING_INTRO_OUTRO"},
		{"StatusProcessed", domain.StatusProcessed, "PROCESSED"},
		{"StatusFailed", domain.StatusFailed, "FAILED"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.status))
		})
	}
}

func TestVideoStatus_Workflow(t *testing.T) {
	// Test typical workflow progression
	workflow := []domain.VideoStatus{
		domain.StatusUploaded,
		domain.StatusTrimming,
		domain.StatusAdjustingRes,
		domain.StatusAddingWatermark,
		domain.StatusRemovingAudio,
		domain.StatusAddingIntroOutro,
		domain.StatusProcessed,
	}
	
	assert.Len(t, workflow, 7)
	assert.Equal(t, domain.StatusUploaded, workflow[0])
	assert.Equal(t, domain.StatusProcessed, workflow[len(workflow)-1])
}