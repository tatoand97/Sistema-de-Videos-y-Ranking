package adapters_test

import (
	"audioremoval/internal/adapters"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStorageRepository(t *testing.T) {
	repo := adapters.NewStorageRepository(nil)
	assert.NotNil(t, repo)
}

func TestStorageRepository_Download(t *testing.T) {
	repo := adapters.NewStorageRepository(nil)
	
	data, err := repo.Download("test-bucket", "test.mp4")
	
	assert.NoError(t, err)
	assert.NotNil(t, data)
}

func TestStorageRepository_Upload(t *testing.T) {
	repo := adapters.NewStorageRepository(nil)
	testData := []byte("test video data")
	
	err := repo.Upload("test-bucket", "test.mp4", testData)
	
	assert.NoError(t, err)
}

func TestStorageRepository_DownloadEmptyBucket(t *testing.T) {
	repo := adapters.NewStorageRepository(nil)
	
	data, err := repo.Download("", "test.mp4")
	
	assert.NoError(t, err)
	assert.NotNil(t, data)
}

func TestStorageRepository_UploadEmptyFilename(t *testing.T) {
	repo := adapters.NewStorageRepository(nil)
	testData := []byte("test video data")
	
	err := repo.Upload("test-bucket", "", testData)
	
	assert.NoError(t, err)
}

func TestStorageRepository_UploadNilData(t *testing.T) {
	repo := adapters.NewStorageRepository(nil)
	
	err := repo.Upload("test-bucket", "test.mp4", nil)
	
	assert.NoError(t, err)
}