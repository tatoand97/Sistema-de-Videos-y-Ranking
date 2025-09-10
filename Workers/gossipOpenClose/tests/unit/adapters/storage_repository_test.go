package adapters_test

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockStorageService struct {
	mock.Mock
}

func (m *MockStorageService) GetObject(bucket, filename string) (io.Reader, error) {
	args := m.Called(bucket, filename)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(io.Reader), args.Error(1)
}

func (m *MockStorageService) PutObject(bucket, filename string, reader io.Reader, size int64) error {
	args := m.Called(bucket, filename, reader, size)
	return args.Error(0)
}

func TestNewStorageRepository(t *testing.T) {
	storage := &MockStorageService{}
	repo := NewStorageRepository(storage)
	
	assert.NotNil(t, repo)
	assert.Equal(t, storage, repo.storage)
}

func TestStorageRepository_Download_Success(t *testing.T) {
	storage := &MockStorageService{}
	repo := NewStorageRepository(storage)
	
	bucket := "test-bucket"
	filename := "test.mp4"
	expectedData := []byte("test video data")
	
	reader := bytes.NewReader(expectedData)
	storage.On("GetObject", bucket, filename).Return(reader, nil)
	
	result, err := repo.Download(bucket, filename)
	
	assert.NoError(t, err)
	assert.Equal(t, expectedData, result)
	storage.AssertExpectations(t)
}

func TestStorageRepository_Download_Error(t *testing.T) {
	storage := &MockStorageService{}
	repo := NewStorageRepository(storage)
	
	bucket := "test-bucket"
	filename := "test.mp4"
	expectedError := errors.New("storage error")
	
	storage.On("GetObject", bucket, filename).Return(nil, expectedError)
	
	result, err := repo.Download(bucket, filename)
	
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, result)
	storage.AssertExpectations(t)
}

func TestStorageRepository_Upload_Success(t *testing.T) {
	storage := &MockStorageService{}
	repo := NewStorageRepository(storage)
	
	bucket := "test-bucket"
	filename := "test.mp4"
	data := []byte("test video data")
	
	storage.On("PutObject", bucket, filename, mock.AnythingOfType("*bytes.Reader"), int64(len(data))).Return(nil)
	
	err := repo.Upload(bucket, filename, data)
	
	assert.NoError(t, err)
	storage.AssertExpectations(t)
}

func TestStorageRepository_Upload_Error(t *testing.T) {
	storage := &MockStorageService{}
	repo := NewStorageRepository(storage)
	
	bucket := "test-bucket"
	filename := "test.mp4"
	data := []byte("test video data")
	expectedError := errors.New("upload error")
	
	storage.On("PutObject", bucket, filename, mock.AnythingOfType("*bytes.Reader"), int64(len(data))).Return(expectedError)
	
	err := repo.Upload(bucket, filename, data)
	
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	storage.AssertExpectations(t)
}