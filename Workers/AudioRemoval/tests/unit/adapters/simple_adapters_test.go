package adapters_test

import (
	"audioremoval/internal/adapters"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewVideoRepository(t *testing.T) {
	repo := adapters.NewVideoRepository()
	assert.NotNil(t, repo)
}

func TestVideoRepository_FindByFilename(t *testing.T) {
	repo := adapters.NewVideoRepository()
	
	video, err := repo.FindByFilename("test.mp4")
	
	assert.NoError(t, err)
	assert.NotNil(t, video)
	assert.Equal(t, "test.mp4", video.Filename)
}

func TestVideoRepository_UpdateStatus(t *testing.T) {
	repo := adapters.NewVideoRepository()
	
	err := repo.UpdateStatus("123", "completed")
	
	assert.NoError(t, err)
}

func TestNewStorageRepository(t *testing.T) {
	repo := adapters.NewStorageRepository(nil)
	assert.NotNil(t, repo)
}