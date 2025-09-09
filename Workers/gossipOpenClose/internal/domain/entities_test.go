package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestProcessingStatus_Constants(t *testing.T) {
	tests := []struct {
		name     string
		status   ProcessingStatus
		expected string
	}{
		{"StatusPending", StatusPending, "pending"},
		{"StatusProcessing", StatusProcessing, "processing"},
		{"StatusCompleted", StatusCompleted, "completed"},
		{"StatusFailed", StatusFailed, "failed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.status))
		})
	}
}

func TestVideo_Creation(t *testing.T) {
	now := time.Now()
	processedAt := now.Add(time.Hour)
	
	video := &Video{
		ID:          "video-123",
		Filename:    "test.mp4",
		Status:      StatusPending,
		CreatedAt:   now,
		ProcessedAt: &processedAt,
	}
	
	assert.Equal(t, "video-123", video.ID)
	assert.Equal(t, "test.mp4", video.Filename)
	assert.Equal(t, StatusPending, video.Status)
	assert.Equal(t, now, video.CreatedAt)
	assert.Equal(t, processedAt, *video.ProcessedAt)
}

func TestVideo_WithoutProcessedAt(t *testing.T) {
	now := time.Now()
	
	video := &Video{
		ID:          "video-456",
		Filename:    "pending.mp4",
		Status:      StatusProcessing,
		CreatedAt:   now,
		ProcessedAt: nil,
	}
	
	assert.Equal(t, "video-456", video.ID)
	assert.Equal(t, "pending.mp4", video.Filename)
	assert.Equal(t, StatusProcessing, video.Status)
	assert.Equal(t, now, video.CreatedAt)
	assert.Nil(t, video.ProcessedAt)
}

func TestProcessingStatus_Workflow(t *testing.T) {
	workflow := []ProcessingStatus{
		StatusPending,
		StatusProcessing,
		StatusCompleted,
	}
	
	assert.Len(t, workflow, 3)
	assert.Equal(t, StatusPending, workflow[0])
	assert.Equal(t, StatusCompleted, workflow[len(workflow)-1])
}

func TestProcessingStatus_FailureState(t *testing.T) {
	assert.Equal(t, "failed", string(StatusFailed))
	assert.NotEqual(t, StatusFailed, StatusCompleted)
}