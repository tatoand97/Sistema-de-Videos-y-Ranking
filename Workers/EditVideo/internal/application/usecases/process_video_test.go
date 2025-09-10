package usecases

import (
	"editvideo/internal/domain"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockVideoRepository struct {
	mock.Mock
}

func (m *MockVideoRepository) FindByFilename(filename string) (*domain.Video, error) {
	args := m.Called(filename)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Video), args.Error(1)
}

func (m *MockVideoRepository) UpdateStatus(id string, status domain.ProcessingStatus) error {
	args := m.Called(id, status)
	return args.Error(0)
}

type MockStorageRepository struct {
	mock.Mock
}

func (m *MockStorageRepository) Download(bucket, filename string) ([]byte, error) {
	args := m.Called(bucket, filename)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockStorageRepository) Upload(bucket, filename string, data []byte) error {
	args := m.Called(bucket, filename, data)
	return args.Error(0)
}

type MockProcessingService struct {
	mock.Mock
}

func (m *MockProcessingService) TrimToMaxSeconds(data []byte, maxSeconds int) ([]byte, error) {
	args := m.Called(data, maxSeconds)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

type MockNotificationService struct {
	mock.Mock
}

func (m *MockNotificationService) NotifyVideoProcessed(videoID, filename, bucketPath string) error {
	args := m.Called(videoID, filename, bucketPath)
	return args.Error(0)
}

func TestNewEditVideoUseCase(t *testing.T) {
	mockVideoRepo := &MockVideoRepository{}
	mockStorageRepo := &MockStorageRepository{}
	mockProcessingService := &MockProcessingService{}
	mockNotificationService := &MockNotificationService{}

	useCase := NewEditVideoUseCase(
		mockVideoRepo,
		mockStorageRepo,
		mockProcessingService,
		mockNotificationService,
		"raw",
		"processed",
		30,
	)

	assert.NotNil(t, useCase)
	assert.Equal(t, "raw", useCase.rawBucket)
	assert.Equal(t, "processed", useCase.processedBucket)
	assert.Equal(t, 30, useCase.maxSeconds)
}

func TestEditVideoUseCase_Execute_Success(t *testing.T) {
	mockVideoRepo := &MockVideoRepository{}
	mockStorageRepo := &MockStorageRepository{}
	mockProcessingService := &MockProcessingService{}
	mockNotificationService := &MockNotificationService{}

	video := &domain.Video{ID: "123", Filename: "test.mp4"}
	inputData := []byte("video data")
	processedData := []byte("processed data")

	mockVideoRepo.On("FindByFilename", "test.mp4").Return(video, nil)
	mockVideoRepo.On("UpdateStatus", "123", domain.StatusProcessing).Return(nil)
	mockStorageRepo.On("Download", "raw", "test.mp4").Return(inputData, nil)
	mockProcessingService.On("TrimToMaxSeconds", inputData, 30).Return(processedData, nil)
	mockStorageRepo.On("Upload", "processed", "test.mp4", processedData).Return(nil)
	mockVideoRepo.On("UpdateStatus", "123", domain.StatusCompleted).Return(nil)
	mockNotificationService.On("NotifyVideoProcessed", "123", "test.mp4", "processed/test.mp4").Return(nil)

	useCase := NewEditVideoUseCase(
		mockVideoRepo,
		mockStorageRepo,
		mockProcessingService,
		mockNotificationService,
		"raw",
		"processed",
		30,
	)

	err := useCase.Execute("123", "test.mp4")

	assert.NoError(t, err)
	mockVideoRepo.AssertExpectations(t)
	mockStorageRepo.AssertExpectations(t)
	mockProcessingService.AssertExpectations(t)
	mockNotificationService.AssertExpectations(t)
}

func TestEditVideoUseCase_Execute_VideoNotFound(t *testing.T) {
	mockVideoRepo := &MockVideoRepository{}
	mockStorageRepo := &MockStorageRepository{}
	mockProcessingService := &MockProcessingService{}
	mockNotificationService := &MockNotificationService{}

	mockVideoRepo.On("FindByFilename", "missing.mp4").Return(nil, errors.New("video not found"))

	useCase := NewEditVideoUseCase(
		mockVideoRepo,
		mockStorageRepo,
		mockProcessingService,
		mockNotificationService,
		"raw",
		"processed",
		30,
	)

	err := useCase.Execute("123", "missing.mp4")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "find video")
	mockVideoRepo.AssertExpectations(t)
}

func TestEditVideoUseCase_Execute_DownloadError(t *testing.T) {
	mockVideoRepo := &MockVideoRepository{}
	mockStorageRepo := &MockStorageRepository{}
	mockProcessingService := &MockProcessingService{}
	mockNotificationService := &MockNotificationService{}

	video := &domain.Video{ID: "123", Filename: "test.mp4"}

	mockVideoRepo.On("FindByFilename", "test.mp4").Return(video, nil)
	mockVideoRepo.On("UpdateStatus", "123", domain.StatusProcessing).Return(nil)
	mockStorageRepo.On("Download", "raw", "test.mp4").Return(nil, errors.New("download failed"))
	mockVideoRepo.On("UpdateStatus", "123", domain.StatusFailed).Return(nil)

	useCase := NewEditVideoUseCase(
		mockVideoRepo,
		mockStorageRepo,
		mockProcessingService,
		mockNotificationService,
		"raw",
		"processed",
		30,
	)

	err := useCase.Execute("123", "test.mp4")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "download")
	mockVideoRepo.AssertExpectations(t)
	mockStorageRepo.AssertExpectations(t)
}

func TestEditVideoUseCase_Execute_ProcessingError(t *testing.T) {
	mockVideoRepo := &MockVideoRepository{}
	mockStorageRepo := &MockStorageRepository{}
	mockProcessingService := &MockProcessingService{}
	mockNotificationService := &MockNotificationService{}

	video := &domain.Video{ID: "123", Filename: "test.mp4"}
	inputData := []byte("video data")

	mockVideoRepo.On("FindByFilename", "test.mp4").Return(video, nil)
	mockVideoRepo.On("UpdateStatus", "123", domain.StatusProcessing).Return(nil)
	mockStorageRepo.On("Download", "raw", "test.mp4").Return(inputData, nil)
	mockProcessingService.On("TrimToMaxSeconds", inputData, 30).Return(nil, errors.New("processing failed"))
	mockVideoRepo.On("UpdateStatus", "123", domain.StatusFailed).Return(nil)

	useCase := NewEditVideoUseCase(
		mockVideoRepo,
		mockStorageRepo,
		mockProcessingService,
		mockNotificationService,
		"raw",
		"processed",
		30,
	)

	err := useCase.Execute("123", "test.mp4")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "processing")
	mockVideoRepo.AssertExpectations(t)
	mockStorageRepo.AssertExpectations(t)
	mockProcessingService.AssertExpectations(t)
}