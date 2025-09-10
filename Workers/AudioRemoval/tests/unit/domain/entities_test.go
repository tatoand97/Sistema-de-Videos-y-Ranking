package domain_test

import (
	"audioremoval/internal/domain"
	"testing"
	"time"

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

func TestProcessingResult_Creation(t *testing.T) {
	result := &domain.ProcessingResult{
		Success:      true,
		ErrorMessage: "",
		OutputPath:   "/path/to/output.mp4",
	}
	
	assert.True(t, result.Success)
	assert.Empty(t, result.ErrorMessage)
	assert.Equal(t, "/path/to/output.mp4", result.OutputPath)
}

func TestProcessingResult_WithError(t *testing.T) {
	result := &domain.ProcessingResult{
		Success:      false,
		ErrorMessage: "processing failed",
		OutputPath:   "",
	}
	
	assert.False(t, result.Success)
	assert.Equal(t, "processing failed", result.ErrorMessage)
	assert.Empty(t, result.OutputPath)
}