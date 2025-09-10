package domain

import (
	"testing"
	"time"
	"watermarking/internal/domain"

	"github.com/stretchr/testify/assert"
)

func TestProcessingStatus_Constants(t *testing.T) {
	tests := []struct {
		name     string
		status   domain.ProcessingStatus
		expected string
	}{
		{"StatusPending", domain.StatusPending, "pending"},
		{"StatusProcessing", domain.StatusProcessing, "processing"},
		{"StatusCompleted", domain.StatusCompleted, "completed"},
		{"StatusFailed", domain.StatusFailed, "failed"},
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
	
	video := &domain.Video{
		ID:          "video-123",
		Filename:    "test.mp4",
		Status:      domain.StatusPending,
		CreatedAt:   now,
		ProcessedAt: &processedAt,
	}
	
	assert.Equal(t, "video-123", video.ID)
	assert.Equal(t, "test.mp4", video.Filename)
	assert.Equal(t, domain.StatusPending, video.Status)
	assert.Equal(t, now, video.CreatedAt)
	assert.Equal(t, processedAt, *video.ProcessedAt)
}

func TestVideo_WithoutProcessedAt(t *testing.T) {
	now := time.Now()
	
	video := &domain.Video{
		ID:          "video-456",
		Filename:    "pending.mp4",
		Status:      domain.StatusProcessing,
		CreatedAt:   now,
		ProcessedAt: nil,
	}
	
	assert.Equal(t, "video-456", video.ID)
	assert.Equal(t, "pending.mp4", video.Filename)
	assert.Equal(t, domain.StatusProcessing, video.Status)
	assert.Equal(t, now, video.CreatedAt)
	assert.Nil(t, video.ProcessedAt)
}