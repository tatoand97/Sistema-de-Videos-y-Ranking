package adapters_test

import (
	"audioremoval/internal/adapters"
	"audioremoval/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewVideoRepository_Constructor(t *testing.T) {
	repo := adapters.NewVideoRepository()
	assert.NotNil(t, repo)
}

func TestVideoRepository_FindByFilename_Success(t *testing.T) {
	repo := adapters.NewVideoRepository()
	
	video, err := repo.FindByFilename("test.mp4")
	
	assert.NoError(t, err)
	assert.NotNil(t, video)
	assert.Equal(t, "test.mp4", video.Filename)
	assert.Equal(t, domain.StatusPending, video.Status)
}

func TestVideoRepository_FindByFilename_EmptyFilename(t *testing.T) {
	repo := adapters.NewVideoRepository()
	
	video, err := repo.FindByFilename("")
	
	assert.NoError(t, err)
	assert.NotNil(t, video)
	assert.Equal(t, "", video.Filename)
}

func TestVideoRepository_UpdateStatus_Success(t *testing.T) {
	repo := adapters.NewVideoRepository()
	
	err := repo.UpdateStatus("video-123", domain.StatusProcessing)
	
	assert.NoError(t, err)
}

func TestVideoRepository_UpdateStatus_AllStatuses(t *testing.T) {
	repo := adapters.NewVideoRepository()
	statuses := []domain.ProcessingStatus{
		domain.StatusPending,
		domain.StatusProcessing,
		domain.StatusCompleted,
		domain.StatusFailed,
	}
	
	for _, status := range statuses {
		err := repo.UpdateStatus("video-123", status)
		assert.NoError(t, err)
	}
}

func TestVideoRepository_UpdateStatus_EmptyID(t *testing.T) {
	repo := adapters.NewVideoRepository()
	
	err := repo.UpdateStatus("", domain.StatusCompleted)
	
	assert.NoError(t, err)
}