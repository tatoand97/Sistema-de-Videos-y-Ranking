package domain

import (
	"statesmachine/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVideo_Structure(t *testing.T) {
	video := &domain.Video{
		ID:           1,
		OriginalFile: "test.mp4",
		Status:       string(domain.StatusUploaded),
	}
	
	assert.NotNil(t, video)
	assert.Equal(t, "test.mp4", video.OriginalFile)
	assert.Equal(t, string(domain.StatusUploaded), video.Status)
	assert.Equal(t, uint(1), video.ID)
}

func TestVideoStatus_String(t *testing.T) {
	tests := []struct {
		status   domain.VideoStatus
		expected string
	}{
		{domain.StatusUploaded, "UPLOADED"},
		{domain.StatusTrimming, "TRIMMING"},
		{domain.StatusAdjustingRes, "ADJUSTING_RESOLUTION"},
		{domain.StatusRemovingAudio, "REMOVING_AUDIO"},
		{domain.StatusAddingWatermark, "ADDING_WATERMARK"},
		{domain.StatusAddingIntroOutro, "ADDING_INTRO_OUTRO"},
		{domain.StatusProcessed, "PROCESSED"},
		{domain.StatusFailed, "FAILED"},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.status))
		})
	}
}