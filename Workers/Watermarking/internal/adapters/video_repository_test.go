package adapters

import (
	"testing"
	"watermarking/internal/domain"

	"github.com/stretchr/testify/assert"
)

func TestNewVideoRepository(t *testing.T) {
	repo := NewVideoRepository()
	assert.NotNil(t, repo)
}

func TestVideoRepository_FindByFilename(t *testing.T) {
	repo := NewVideoRepository()
	filename := "test-video.mp4"
	
	video, err := repo.FindByFilename(filename)
	
	assert.NoError(t, err)
	assert.NotNil(t, video)
	assert.Equal(t, filename, video.ID)
	assert.Equal(t, filename, video.Filename)
	assert.Equal(t, domain.StatusPending, video.Status)
	assert.False(t, video.CreatedAt.IsZero())
}

func TestVideoRepository_UpdateStatus(t *testing.T) {
	repo := NewVideoRepository()
	
	tests := []struct {
		name   string
		id     string
		status domain.ProcessingStatus
	}{
		{"update to processing", "video-1", domain.StatusProcessing},
		{"update to completed", "video-2", domain.StatusCompleted},
		{"update to failed", "video-3", domain.StatusFailed},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.UpdateStatus(tt.id, tt.status)
			assert.NoError(t, err)
		})
	}
}