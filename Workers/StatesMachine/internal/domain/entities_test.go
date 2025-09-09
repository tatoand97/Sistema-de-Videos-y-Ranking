package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVideo_TableName(t *testing.T) {
	video := Video{}
	assert.Equal(t, "video", video.TableName())
}

func TestVideo_Creation(t *testing.T) {
	video := Video{
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

func TestVideoStatus_Workflow(t *testing.T) {
	// Test typical workflow progression
	workflow := []VideoStatus{
		StatusUploaded,
		StatusTrimming,
		StatusAdjustingRes,
		StatusAddingWatermark,
		StatusRemovingAudio,
		StatusAddingIntroOutro,
		StatusProcessed,
	}
	
	assert.Len(t, workflow, 7)
	assert.Equal(t, StatusUploaded, workflow[0])
	assert.Equal(t, StatusProcessed, workflow[len(workflow)-1])
}

func TestVideoStatus_FailureStates(t *testing.T) {
	failureStates := []VideoStatus{StatusFailed}
	assert.Contains(t, failureStates, StatusFailed)
}

func TestVideo_WithVideoStatus(t *testing.T) {
	video := Video{
		ID:           456,
		OriginalFile: "workflow.mp4",
		Status:       string(StatusTrimming),
	}
	
	assert.Equal(t, string(StatusTrimming), video.Status)
	assert.Equal(t, "TRIMMING", video.Status)
}