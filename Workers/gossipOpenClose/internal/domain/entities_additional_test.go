package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestVideo_Creation(t *testing.T) {
	now := time.Now()
	video := Video{
		ID:        "test-id",
		Filename:  "test.mp4",
		Status:    StatusPending,
		CreatedAt: now,
	}
	
	assert.Equal(t, "test-id", video.ID)
	assert.Equal(t, "test.mp4", video.Filename)
	assert.Equal(t, StatusPending, video.Status)
	assert.Equal(t, now, video.CreatedAt)
	assert.Nil(t, video.ProcessedAt)
}

func TestVideo_WithProcessedTime(t *testing.T) {
	now := time.Now()
	processedTime := now.Add(5 * time.Minute)
	
	video := Video{
		ID:          "test-id",
		Filename:    "test.mp4",
		Status:      StatusCompleted,
		CreatedAt:   now,
		ProcessedAt: &processedTime,
	}
	
	assert.Equal(t, StatusCompleted, video.Status)
	assert.NotNil(t, video.ProcessedAt)
	assert.Equal(t, processedTime, *video.ProcessedAt)
}

func TestProcessingStatus_Constants(t *testing.T) {
	assert.Equal(t, ProcessingStatus("pending"), StatusPending)
	assert.Equal(t, ProcessingStatus("processing"), StatusProcessing)
	assert.Equal(t, ProcessingStatus("completed"), StatusCompleted)
	assert.Equal(t, ProcessingStatus("failed"), StatusFailed)
}

func TestProcessingStatus_StringConversion(t *testing.T) {
	tests := []struct {
		status   ProcessingStatus
		expected string
	}{
		{StatusPending, "pending"},
		{StatusProcessing, "processing"},
		{StatusCompleted, "completed"},
		{StatusFailed, "failed"},
	}
	
	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.status))
		})
	}
}

func TestVideo_StatusTransitions(t *testing.T) {
	video := Video{
		ID:       "test-id",
		Filename: "test.mp4",
		Status:   StatusPending,
	}
	
	// Test valid status transitions
	validTransitions := []ProcessingStatus{
		StatusProcessing,
		StatusCompleted,
		StatusFailed,
	}
	
	for _, status := range validTransitions {
		video.Status = status
		assert.Equal(t, status, video.Status)
	}
}

func TestVideo_EmptyValues(t *testing.T) {
	video := Video{}
	
	assert.Empty(t, video.ID)
	assert.Empty(t, video.Filename)
	assert.Empty(t, video.Status)
	assert.True(t, video.CreatedAt.IsZero())
	assert.Nil(t, video.ProcessedAt)
}

func TestVideo_LongFilename(t *testing.T) {
	longFilename := "very_long_filename_that_might_cause_issues_in_some_systems_" +
		"with_many_characters_and_special_symbols_123456789.mp4"
	
	video := Video{
		ID:       "test-id",
		Filename: longFilename,
		Status:   StatusPending,
	}
	
	assert.Equal(t, longFilename, video.Filename)
}

func TestVideo_SpecialCharactersInFilename(t *testing.T) {
	specialFilenames := []string{
		"file with spaces.mp4",
		"file-with-dashes.mp4",
		"file_with_underscores.mp4",
		"file.with.dots.mp4",
		"file(with)parentheses.mp4",
		"file[with]brackets.mp4",
		"file{with}braces.mp4",
	}
	
	for _, filename := range specialFilenames {
		t.Run(filename, func(t *testing.T) {
			video := Video{
				ID:       "test-id",
				Filename: filename,
				Status:   StatusPending,
			}
			
			assert.Equal(t, filename, video.Filename)
		})
	}
}

func TestVideo_TimeComparison(t *testing.T) {
	now := time.Now()
	later := now.Add(1 * time.Hour)
	
	video := Video{
		ID:          "test-id",
		Filename:    "test.mp4",
		Status:      StatusCompleted,
		CreatedAt:   now,
		ProcessedAt: &later,
	}
	
	assert.True(t, video.ProcessedAt.After(video.CreatedAt))
	
	duration := video.ProcessedAt.Sub(video.CreatedAt)
	assert.Equal(t, 1*time.Hour, duration)
}

func TestProcessingStatus_Comparison(t *testing.T) {
	statuses := []ProcessingStatus{
		StatusPending,
		StatusProcessing,
		StatusCompleted,
		StatusFailed,
	}
	
	// Test that each status is equal to itself
	for _, status := range statuses {
		assert.Equal(t, status, status)
	}
	
	// Test that different statuses are not equal
	for i, status1 := range statuses {
		for j, status2 := range statuses {
			if i != j {
				assert.NotEqual(t, status1, status2)
			}
		}
	}
}

func TestVideo_CopySemantics(t *testing.T) {
	original := Video{
		ID:       "original-id",
		Filename: "original.mp4",
		Status:   StatusPending,
		CreatedAt: time.Now(),
	}
	
	// Test that copying works as expected
	copy := original
	copy.ID = "copy-id"
	copy.Status = StatusProcessing
	
	// Original should remain unchanged
	assert.Equal(t, "original-id", original.ID)
	assert.Equal(t, StatusPending, original.Status)
	
	// Copy should have new values
	assert.Equal(t, "copy-id", copy.ID)
	assert.Equal(t, StatusProcessing, copy.Status)
	
	// But filename and created time should be the same
	assert.Equal(t, original.Filename, copy.Filename)
	assert.Equal(t, original.CreatedAt, copy.CreatedAt)
}