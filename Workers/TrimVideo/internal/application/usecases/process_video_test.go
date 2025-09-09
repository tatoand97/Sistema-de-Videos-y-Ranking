package usecases

import (
	"errors"
	"testing"
	"time"
	"trimvideo/internal/domain"

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
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockStorageRepository) Upload(bucket, filename string, data []byte) error {
	args := m.Called(bucket, filename, data)
	return args.Error(0)
}

type MockVideoProcessingService struct {
	mock.Mock
}

func (m *MockVideoProcessingService) TrimToMaxSeconds(data []byte, maxSeconds int) ([]byte, error) {
	args := m.Called(data, maxSeconds)
	return args.Get(0).([]byte), args.Error(1)
}

type MockNotificationService struct {
	mock.Mock
}

func (m *MockNotificationService) NotifyVideoProcessed(videoID, filename, bucketPath string) error {
	args := m.Called(videoID, filename, bucketPath)
	return args.Error(0)
}

func TestNewProcessVideoUseCase(t *testing.T) {
	videoRepo := &MockVideoRepository{}
	storageRepo := &MockStorageRepository{}
	processingService := &MockVideoProcessingService{}
	notificationService := &MockNotificationService{}
	
	useCase := NewProcessVideoUseCase(
		videoRepo,
		storageRepo,
		processingService,
		notificationService,
		"raw-bucket",
		"processed-bucket",
		30,
	)
	
	assert.NotNil(t, useCase)
	assert.Equal(t, "raw-bucket", useCase.rawBucket)
	assert.Equal(t, "processed-bucket", useCase.processedBucket)
	assert.Equal(t, 30, useCase.maxSeconds)
}

func TestProcessVideoUseCase_Execute_Success(t *testing.T) {
	videoRepo := &MockVideoRepository{}
	storageRepo := &MockStorageRepository{}
	processingService := &MockVideoProcessingService{}
	notificationService := &MockNotificationService{}
	
	video := &domain.Video{
		ID:        "video-123",
		Filename:  "test.mp4",
		Status:    domain.StatusPending,
		CreatedAt: time.Now(),
	}
	
	inputData := []byte("input video data")
	processedData := []byte("trimmed video data")
	
	videoRepo.On("FindByFilename", "test.mp4").Return(video, nil)
	videoRepo.On("UpdateStatus", "video-123", domain.StatusProcessing).Return(nil)
	storageRepo.On("Download", "raw-bucket", "test.mp4").Return(inputData, nil)
	processingService.On("TrimToMaxSeconds", inputData, 30).Return(processedData, nil)
	storageRepo.On("Upload", "processed-bucket", "test.mp4", processedData).Return(nil)
	videoRepo.On("UpdateStatus", "video-123", domain.StatusCompleted).Return(nil)
	notificationService.On("NotifyVideoProcessed", "video-123", "test.mp4", "processed-bucket/test.mp4").Return(nil)
	
	useCase := NewProcessVideoUseCase(
		videoRepo,
		storageRepo,
		processingService,
		notificationService,
		"raw-bucket",
		"processed-bucket",
		30,
	)
	
	err := useCase.Execute("video-123", "test.mp4")
	
	assert.NoError(t, err)
	videoRepo.AssertExpectations(t)
	storageRepo.AssertExpectations(t)
	processingService.AssertExpectations(t)
	notificationService.AssertExpectations(t)
}

func TestProcessVideoUseCase_Execute_VideoNotFound(t *testing.T) {
	videoRepo := &MockVideoRepository{}
	storageRepo := &MockStorageRepository{}
	processingService := &MockVideoProcessingService{}
	notificationService := &MockNotificationService{}
	
	videoRepo.On("FindByFilename", "test.mp4").Return(nil, errors.New("video not found"))
	
	useCase := NewProcessVideoUseCase(
		videoRepo,
		storageRepo,
		processingService,
		notificationService,
		"raw-bucket",
		"processed-bucket",
		30,
	)
	
	err := useCase.Execute("video-123", "test.mp4")
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "find video")
	videoRepo.AssertExpectations(t)
}

func TestProcessVideoUseCase_Execute_ProcessingFails(t *testing.T) {
	videoRepo := &MockVideoRepository{}
	storageRepo := &MockStorageRepository{}
	processingService := &MockVideoProcessingService{}
	notificationService := &MockNotificationService{}
	
	video := &domain.Video{
		ID:       "video-123",
		Filename: "test.mp4",
		Status:   domain.StatusPending,
	}
	
	inputData := []byte("input video data")
	
	videoRepo.On("FindByFilename", "test.mp4").Return(video, nil)
	videoRepo.On("UpdateStatus", "video-123", domain.StatusProcessing).Return(nil)
	storageRepo.On("Download", "raw-bucket", "test.mp4").Return(inputData, nil)
	processingService.On("TrimToMaxSeconds", inputData, 30).Return([]byte{}, errors.New("processing failed"))
	videoRepo.On("UpdateStatus", "video-123", domain.StatusFailed).Return(nil)
	
	useCase := NewProcessVideoUseCase(
		videoRepo,
		storageRepo,
		processingService,
		notificationService,
		"raw-bucket",
		"processed-bucket",
		30,
	)
	
	err := useCase.Execute("video-123", "test.mp4")
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "processing")
	videoRepo.AssertExpectations(t)
	storageRepo.AssertExpectations(t)
	processingService.AssertExpectations(t)
}