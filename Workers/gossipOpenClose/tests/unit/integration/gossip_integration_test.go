package integration

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"gossipopenclose/internal/adapters"
	"gossipopenclose/internal/application/services"
	"gossipopenclose/internal/application/usecases"
	"gossipopenclose/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Integration test mocks
type IntegrationMockVideoRepository struct {
	mock.Mock
}

func (m *IntegrationMockVideoRepository) FindByFilename(filename string) (*domain.Video, error) {
	args := m.Called(filename)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Video), args.Error(1)
}

func (m *IntegrationMockVideoRepository) UpdateStatus(id string, status domain.ProcessingStatus) error {
	args := m.Called(id, status)
	return args.Error(0)
}

type IntegrationMockStorageRepository struct {
	mock.Mock
}

func (m *IntegrationMockStorageRepository) Download(bucket, filename string) ([]byte, error) {
	args := m.Called(bucket, filename)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *IntegrationMockStorageRepository) Upload(bucket, filename string, data []byte) error {
	args := m.Called(bucket, filename, data)
	return args.Error(0)
}

type IntegrationMockNotificationService struct {
	mock.Mock
}

func (m *IntegrationMockNotificationService) NotifyVideoProcessed(videoID, filename, bucketPath string) error {
	args := m.Called(videoID, filename, bucketPath)
	return args.Error(0)
}

type IntegrationMockProcessingService struct {
	mock.Mock
}

func (m *IntegrationMockProcessingService) Process(inputData []byte, logoPath string, introSeconds, outroSeconds float64, width, height, fps int) ([]byte, error) {
	args := m.Called(inputData, logoPath, introSeconds, outroSeconds, width, height, fps)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func TestGossipOpenClose_FullWorkflow_Success(t *testing.T) {
	// Setup mocks
	videoRepo := &IntegrationMockVideoRepository{}
	storageRepo := &IntegrationMockStorageRepository{}
	processingService := &IntegrationMockProcessingService{}
	notificationService := &IntegrationMockNotificationService{}
	
	// Create use case
	useCase := usecases.NewOpenCloseUseCase(
		videoRepo,
		storageRepo,
		processingService,
		notificationService,
		"raw-bucket",
		"processed-bucket",
		"logo.png",
		2.5,
		2.5,
		1920,
		1080,
		30,
	)
	
	// Create message handler
	messageHandler := adapters.NewMessageHandler(useCase)
	
	// Test data
	video := &domain.Video{
		ID:       "video-123",
		Filename: "test.mp4",
		Status:   domain.StatusPending,
		CreatedAt: time.Now(),
	}
	
	inputData := []byte("input video data")
	processedData := []byte("processed video data")
	
	message := adapters.VideoMessage{
		VideoID:  "video-123",
		Filename: "test.mp4",
	}
	
	messageBody, _ := json.Marshal(message)
	
	// Setup expectations
	videoRepo.On("FindByFilename", "test.mp4").Return(video, nil)
	videoRepo.On("UpdateStatus", "video-123", domain.StatusProcessing).Return(nil)
	storageRepo.On("Download", "raw-bucket", "test.mp4").Return(inputData, nil)
	processingService.On("Process", inputData, "logo.png", 2.5, 2.5, 1920, 1080, 30).Return(processedData, nil)
	storageRepo.On("Upload", "processed-bucket", "test.mp4", processedData).Return(nil)
	videoRepo.On("UpdateStatus", "video-123", domain.StatusCompleted).Return(nil)
	notificationService.On("NotifyVideoProcessed", "video-123", "test.mp4", "processed-bucket/test.mp4").Return(nil)
	
	// Execute full workflow
	err := messageHandler.HandleMessage(messageBody)
	
	// Verify
	assert.NoError(t, err)
	videoRepo.AssertExpectations(t)
	storageRepo.AssertExpectations(t)
	processingService.AssertExpectations(t)
	notificationService.AssertExpectations(t)
}

func TestGossipOpenClose_FullWorkflow_ProcessingFailure(t *testing.T) {
	// Setup mocks
	videoRepo := &IntegrationMockVideoRepository{}
	storageRepo := &IntegrationMockStorageRepository{}
	processingService := &IntegrationMockProcessingService{}
	notificationService := &IntegrationMockNotificationService{}
	
	// Create use case
	useCase := usecases.NewOpenCloseUseCase(
		videoRepo,
		storageRepo,
		processingService,
		notificationService,
		"raw-bucket",
		"processed-bucket",
		"logo.png",
		2.5,
		2.5,
		1920,
		1080,
		30,
	)
	
	// Create message handler
	messageHandler := adapters.NewMessageHandler(useCase)
	
	// Test data
	video := &domain.Video{
		ID:       "video-123",
		Filename: "test.mp4",
		Status:   domain.StatusPending,
		CreatedAt: time.Now(),
	}
	
	inputData := []byte("input video data")
	processingError := errors.New("processing failed")
	
	message := adapters.VideoMessage{
		VideoID:  "video-123",
		Filename: "test.mp4",
	}
	
	messageBody, _ := json.Marshal(message)
	
	// Setup expectations
	videoRepo.On("FindByFilename", "test.mp4").Return(video, nil)
	videoRepo.On("UpdateStatus", "video-123", domain.StatusProcessing).Return(nil)
	storageRepo.On("Download", "raw-bucket", "test.mp4").Return(inputData, nil)
	processingService.On("Process", inputData, "logo.png", 2.5, 2.5, 1920, 1080, 30).Return(nil, processingError)
	videoRepo.On("UpdateStatus", "video-123", domain.StatusFailed).Return(nil)
	
	// Execute workflow
	err := messageHandler.HandleMessage(messageBody)
	
	// Verify
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "process")
	videoRepo.AssertExpectations(t)
	storageRepo.AssertExpectations(t)
	processingService.AssertExpectations(t)
	notificationService.AssertNotCalled(t, "NotifyVideoProcessed")
}

func TestGossipOpenClose_MessageHandling_InvalidMessage(t *testing.T) {
	// Setup mocks
	videoRepo := &IntegrationMockVideoRepository{}
	storageRepo := &IntegrationMockStorageRepository{}
	processingService := &IntegrationMockProcessingService{}
	notificationService := &IntegrationMockNotificationService{}
	
	// Create use case
	useCase := usecases.NewOpenCloseUseCase(
		videoRepo,
		storageRepo,
		processingService,
		notificationService,
		"raw-bucket",
		"processed-bucket",
		"logo.png",
		2.5,
		2.5,
		1920,
		1080,
		30,
	)
	
	// Create message handler
	messageHandler := adapters.NewMessageHandler(useCase)
	
	// Invalid JSON message
	invalidMessage := []byte(`{"invalid": "json"`)
	
	// Execute
	err := messageHandler.HandleMessage(invalidMessage)
	
	// Verify
	assert.Error(t, err)
	videoRepo.AssertNotCalled(t, "FindByFilename")
	storageRepo.AssertNotCalled(t, "Download")
	processingService.AssertNotCalled(t, "Process")
	notificationService.AssertNotCalled(t, "NotifyVideoProcessed")
}

func TestGossipOpenClose_StatusTransitions(t *testing.T) {
	// Setup mocks
	videoRepo := &IntegrationMockVideoRepository{}
	storageRepo := &IntegrationMockStorageRepository{}
	processingService := &IntegrationMockProcessingService{}
	notificationService := &IntegrationMockNotificationService{}
	
	// Create use case
	useCase := usecases.NewOpenCloseUseCase(
		videoRepo,
		storageRepo,
		processingService,
		notificationService,
		"raw-bucket",
		"processed-bucket",
		"logo.png",
		2.5,
		2.5,
		1920,
		1080,
		30,
	)
	
	video := &domain.Video{
		ID:       "video-123",
		Filename: "test.mp4",
		Status:   domain.StatusPending,
		CreatedAt: time.Now(),
	}
	
	inputData := []byte("input video data")
	processedData := []byte("processed video data")
	
	// Test successful status transitions: pending -> processing -> completed
	videoRepo.On("FindByFilename", "test.mp4").Return(video, nil)
	
	// First call: pending -> processing
	videoRepo.On("UpdateStatus", "video-123", domain.StatusProcessing).Return(nil).Once()
	
	storageRepo.On("Download", "raw-bucket", "test.mp4").Return(inputData, nil)
	processingService.On("Process", inputData, "logo.png", 2.5, 2.5, 1920, 1080, 30).Return(processedData, nil)
	storageRepo.On("Upload", "processed-bucket", "test.mp4", processedData).Return(nil)
	
	// Second call: processing -> completed
	videoRepo.On("UpdateStatus", "video-123", domain.StatusCompleted).Return(nil).Once()
	
	notificationService.On("NotifyVideoProcessed", "video-123", "test.mp4", "processed-bucket/test.mp4").Return(nil)
	
	// Execute
	err := useCase.Execute("video-123", "test.mp4")
	
	// Verify
	assert.NoError(t, err)
	videoRepo.AssertExpectations(t)
	storageRepo.AssertExpectations(t)
	processingService.AssertExpectations(t)
	notificationService.AssertExpectations(t)
}

func TestGossipOpenClose_ErrorRecovery(t *testing.T) {
	// Setup mocks
	videoRepo := &IntegrationMockVideoRepository{}
	storageRepo := &IntegrationMockStorageRepository{}
	processingService := &IntegrationMockProcessingService{}
	notificationService := &IntegrationMockNotificationService{}
	
	// Create use case
	useCase := usecases.NewOpenCloseUseCase(
		videoRepo,
		storageRepo,
		processingService,
		notificationService,
		"raw-bucket",
		"processed-bucket",
		"logo.png",
		2.5,
		2.5,
		1920,
		1080,
		30,
	)
	
	video := &domain.Video{
		ID:       "video-123",
		Filename: "test.mp4",
		Status:   domain.StatusPending,
		CreatedAt: time.Now(),
	}
	
	downloadError := errors.New("download failed")
	
	// Test error recovery: when download fails, status should be set to failed
	videoRepo.On("FindByFilename", "test.mp4").Return(video, nil)
	videoRepo.On("UpdateStatus", "video-123", domain.StatusProcessing).Return(nil)
	storageRepo.On("Download", "raw-bucket", "test.mp4").Return(nil, downloadError)
	videoRepo.On("UpdateStatus", "video-123", domain.StatusFailed).Return(nil)
	
	// Execute
	err := useCase.Execute("video-123", "test.mp4")
	
	// Verify
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "download")
	videoRepo.AssertExpectations(t)
	storageRepo.AssertExpectations(t)
	processingService.AssertNotCalled(t, "Process")
	notificationService.AssertNotCalled(t, "NotifyVideoProcessed")
}