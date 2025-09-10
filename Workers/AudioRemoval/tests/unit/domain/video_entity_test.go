package domain_test

import (
	"audioremoval/internal/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestVideo_NewVideo(t *testing.T) {
	now := time.Now()
	video := &domain.Video{
		ID:        "video-123",
		Filename:  "test.mp4",
		Status:    domain.StatusPending,
		CreatedAt: now,
	}

	assert.Equal(t, "video-123", video.ID)
	assert.Equal(t, "test.mp4", video.Filename)
	assert.Equal(t, domain.StatusPending, video.Status)
	assert.Equal(t, now, video.CreatedAt)
	assert.Nil(t, video.ProcessedAt)
}

func TestVideo_WithProcessedAt(t *testing.T) {
	now := time.Now()
	processedAt := now.Add(time.Hour)
	
	video := &domain.Video{
		ID:          "video-123",
		Filename:    "test.mp4",
		Status:      domain.StatusCompleted,
		CreatedAt:   now,
		ProcessedAt: &processedAt,
	}

	assert.NotNil(t, video.ProcessedAt)
	assert.Equal(t, processedAt, *video.ProcessedAt)
}

func TestVideo_StatusTransitions(t *testing.T) {
	video := &domain.Video{
		ID:       "video-123",
		Filename: "test.mp4",
		Status:   domain.StatusPending,
	}

	// Test status transitions
	video.Status = domain.StatusProcessing
	assert.Equal(t, domain.StatusProcessing, video.Status)

	video.Status = domain.StatusCompleted
	assert.Equal(t, domain.StatusCompleted, video.Status)

	video.Status = domain.StatusFailed
	assert.Equal(t, domain.StatusFailed, video.Status)
}

func TestVideo_EmptyFields(t *testing.T) {
	video := &domain.Video{}

	assert.Empty(t, video.ID)
	assert.Empty(t, video.Filename)
	assert.Empty(t, video.Status)
	assert.True(t, video.CreatedAt.IsZero())
	assert.Nil(t, video.ProcessedAt)
}

func TestProcessingStatus_StringValues(t *testing.T) {
	tests := []struct {
		status   domain.ProcessingStatus
		expected string
	}{
		{domain.StatusPending, "pending"},
		{domain.StatusProcessing, "processing"},
		{domain.StatusCompleted, "completed"},
		{domain.StatusFailed, "failed"},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.status))
		})
	}
}

func TestProcessingResult_Success(t *testing.T) {
	result := &domain.ProcessingResult{
		Success:      true,
		ErrorMessage: "",
		OutputPath:   "/path/to/output.mp4",
	}

	assert.True(t, result.Success)
	assert.Empty(t, result.ErrorMessage)
	assert.Equal(t, "/path/to/output.mp4", result.OutputPath)
}

func TestProcessingResult_Failure(t *testing.T) {
	result := &domain.ProcessingResult{
		Success:      false,
		ErrorMessage: "processing failed",
		OutputPath:   "",
	}

	assert.False(t, result.Success)
	assert.Equal(t, "processing failed", result.ErrorMessage)
	assert.Empty(t, result.OutputPath)
}

func TestProcessingResult_Empty(t *testing.T) {
	result := &domain.ProcessingResult{}

	assert.False(t, result.Success)
	assert.Empty(t, result.ErrorMessage)
	assert.Empty(t, result.OutputPath)
}