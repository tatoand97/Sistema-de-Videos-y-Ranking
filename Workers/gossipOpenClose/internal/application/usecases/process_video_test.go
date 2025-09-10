package usecases

import (
	"errors"
	"testing"
	"time"

	"gossipopenclose/internal/domain"
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

func (m *MockProcessingService) Process(inputData []byte, logoPath string, introSeconds, outroSeconds float64, width, height, fps int) ([]byte, error) {
	args := m.Called(inputData, logoPath, introSeconds, outroSeconds, width, height, fps)
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

func TestNewOpenCloseUseCase(t *testing.T) {
	videoRepo := &MockVideoRepository{}
	storageRepo := &MockStorageRepository{}
	processingService := &MockProcessingService{}
	notificationService := &MockNotificationService{}
	
	uc := NewOpenCloseUseCase(
		videoRepo,
		storageRepo,
		nil, // Will be set manually
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
	
	assert.NotNil(t, uc)
	assert.Equal(t, "raw-bucket", uc.rawBucket)
	assert.Equal(t, "processed-bucket", uc.processedBucket)
	assert.Equal(t, "logo.png", uc.logoPath)
	assert.Equal(t, 2.5, uc.introSeconds)
	assert.Equal(t, 2.5, uc.outroSeconds)
	assert.Equal(t, 1920, uc.targetW)
	assert.Equal(t, 1080, uc.targetH)
	assert.Equal(t, 30, uc.fps)
}

func TestOpenCloseUseCase_Execute_Success(t *testing.T) {
	videoRepo := &MockVideoRepository{}
	storageRepo := &MockStorageRepository{}
	processingService := &MockProcessingService{}
	notificationService := &MockNotificationService{}
	
	uc := &OpenCloseUseCase{
		videoRepo:           videoRepo,
		storageRepo:         storageRepo,
		processingService:   processingService,
		notificationService: notificationService,
		rawBucket:          "raw-bucket",
		processedBucket:    "processed-bucket",
		logoPath:           "logo.png",
		introSeconds:       2.5,
		outroSeconds:       2.5,
		targetW:            1920,
		targetH:            1080,
		fps:                30,
	}
	
	video := &domain.Video{
		ID:       "video-123",
		Filename: "test.mp4",
		Status:   domain.StatusPending,
		CreatedAt: time.Now(),
	}
	
	inputData := []byte("input video data")
	processedData := []byte("processed video data")
	
	videoRepo.On("FindByFilename", "test.mp4").Return(video, nil)
	videoRepo.On("UpdateStatus", "video-123", domain.StatusProcessing).Return(nil)
	storageRepo.On("Download", "raw-bucket", "test.mp4").Return(inputData, nil)
	processingService.On("Process", inputData, "logo.png", 2.5, 2.5, 1920, 1080, 30).Return(processedData, nil)
	storageRepo.On("Upload", "processed-bucket", "test.mp4", processedData).Return(nil)
	videoRepo.On("UpdateStatus", "video-123", domain.StatusCompleted).Return(nil)
	notificationService.On("NotifyVideoProcessed", "video-123", "test.mp4", "processed-bucket/test.mp4").Return(nil)
	
	err := uc.Execute("video-123", "test.mp4")
	
	assert.NoError(t, err)
	videoRepo.AssertExpectations(t)
	storageRepo.AssertExpectations(t)
	processingService.AssertExpectations(t)
	notificationService.AssertExpectations(t)
}

func TestOpenCloseUseCase_Execute_FindVideoError(t *testing.T) {
	videoRepo := &MockVideoRepository{}
	storageRepo := &MockStorageRepository{}
	processingService := &MockProcessingService{}
	notificationService := &MockNotificationService{}
	
	uc := &OpenCloseUseCase{
		videoRepo:           videoRepo,
		storageRepo:         storageRepo,
		processingService:   processingService,
		notificationService: notificationService,
		rawBucket:          "raw-bucket",
		processedBucket:    "processed-bucket",
	}
	
	expectedError := errors.New("video not found")
	videoRepo.On("FindByFilename", "test.mp4").Return(nil, expectedError)
	
	err := uc.Execute("video-123", "test.mp4")
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "find video")
	videoRepo.AssertExpectations(t)
}

func TestOpenCloseUseCase_Execute_UpdateStatusError(t *testing.T) {
	videoRepo := &MockVideoRepository{}
	storageRepo := &MockStorageRepository{}
	processingService := &MockProcessingService{}
	notificationService := &MockNotificationService{}
	
	uc := &OpenCloseUseCase{
		videoRepo:           videoRepo,
		storageRepo:         storageRepo,
		processingService:   processingService,
		notificationService: notificationService,
		rawBucket:          "raw-bucket",
		processedBucket:    "processed-bucket",
	}
	
	video := &domain.Video{
		ID:       "video-123",
		Filename: "test.mp4",
		Status:   domain.StatusPending,
	}
	
	expectedError := errors.New("update status failed")
	videoRepo.On("FindByFilename", "test.mp4").Return(video, nil)
	videoRepo.On("UpdateStatus", "video-123", domain.StatusProcessing).Return(expectedError)
	
	err := uc.Execute("video-123", "test.mp4")
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "update status")
	videoRepo.AssertExpectations(t)
}

func TestOpenCloseUseCase_Execute_DownloadError(t *testing.T) {
	videoRepo := &MockVideoRepository{}
	storageRepo := &MockStorageRepository{}
	processingService := &MockProcessingService{}
	notificationService := &MockNotificationService{}
	
	uc := &OpenCloseUseCase{
		videoRepo:           videoRepo,
		storageRepo:         storageRepo,
		processingService:   processingService,
		notificationService: notificationService,
		rawBucket:          "raw-bucket",
		processedBucket:    "processed-bucket",
	}
	
	video := &domain.Video{
		ID:       "video-123",
		Filename: "test.mp4",
		Status:   domain.StatusPending,
	}
	
	expectedError := errors.New("download failed")
	videoRepo.On("FindByFilename", "test.mp4").Return(video, nil)
	videoRepo.On("UpdateStatus", "video-123", domain.StatusProcessing).Return(nil)
	storageRepo.On("Download", "raw-bucket", "test.mp4").Return(nil, expectedError)
	videoRepo.On("UpdateStatus", "video-123", domain.StatusFailed).Return(nil)
	
	err := uc.Execute("video-123", "test.mp4")
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "download")
	videoRepo.AssertExpectations(t)
	storageRepo.AssertExpectations(t)
}

func TestOpenCloseUseCase_Execute_ProcessingError(t *testing.T) {
	videoRepo := &MockVideoRepository{}
	storageRepo := &MockStorageRepository{}
	processingService := &MockProcessingService{}
	notificationService := &MockNotificationService{}
	
	uc := &OpenCloseUseCase{
		videoRepo:           videoRepo,
		storageRepo:         storageRepo,
		processingService:   processingService,
		notificationService: notificationService,
		rawBucket:          "raw-bucket",
		processedBucket:    "processed-bucket",
		logoPath:           "logo.png",
		introSeconds:       2.5,
		outroSeconds:       2.5,
		targetW:            1920,
		targetH:            1080,
		fps:                30,
	}
	
	video := &domain.Video{
		ID:       "video-123",
		Filename: "test.mp4",
		Status:   domain.StatusPending,
	}
	
	inputData := []byte("input video data")
	expectedError := errors.New("processing failed")
	
	videoRepo.On("FindByFilename", "test.mp4").Return(video, nil)
	videoRepo.On("UpdateStatus", "video-123", domain.StatusProcessing).Return(nil)
	storageRepo.On("Download", "raw-bucket", "test.mp4").Return(inputData, nil)
	processingService.On("Process", inputData, "logo.png", 2.5, 2.5, 1920, 1080, 30).Return(nil, expectedError)
	videoRepo.On("UpdateStatus", "video-123", domain.StatusFailed).Return(nil)
	
	err := uc.Execute("video-123", "test.mp4")
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "process")
	videoRepo.AssertExpectations(t)
	storageRepo.AssertExpectations(t)
	processingService.AssertExpectations(t)
}

func TestOpenCloseUseCase_Execute_UploadError(t *testing.T) {
	videoRepo := &MockVideoRepository{}
	storageRepo := &MockStorageRepository{}
	processingService := &MockProcessingService{}
	notificationService := &MockNotificationService{}
	
	uc := &OpenCloseUseCase{
		videoRepo:           videoRepo,
		storageRepo:         storageRepo,
		processingService:   processingService,
		notificationService: notificationService,
		rawBucket:          "raw-bucket",
		processedBucket:    "processed-bucket",
		logoPath:           "logo.png",
		introSeconds:       2.5,
		outroSeconds:       2.5,
		targetW:            1920,
		targetH:            1080,
		fps:                30,
	}
	
	video := &domain.Video{
		ID:       "video-123",
		Filename: "test.mp4",
		Status:   domain.StatusPending,
	}
	
	inputData := []byte("input video data")
	processedData := []byte("processed video data")
	expectedError := errors.New("upload failed")
	
	videoRepo.On("FindByFilename", "test.mp4").Return(video, nil)
	videoRepo.On("UpdateStatus", "video-123", domain.StatusProcessing).Return(nil)
	storageRepo.On("Download", "raw-bucket", "test.mp4").Return(inputData, nil)
	processingService.On("Process", inputData, "logo.png", 2.5, 2.5, 1920, 1080, 30).Return(processedData, nil)
	storageRepo.On("Upload", "processed-bucket", "test.mp4", processedData).Return(expectedError)
	videoRepo.On("UpdateStatus", "video-123", domain.StatusFailed).Return(nil)
	
	err := uc.Execute("video-123", "test.mp4")
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "upload")
	videoRepo.AssertExpectations(t)
	storageRepo.AssertExpectations(t)
	processingService.AssertExpectations(t)
}

func TestOpenCloseUseCase_Execute_FinalStatusUpdateError(t *testing.T) {
	videoRepo := &MockVideoRepository{}
	storageRepo := &MockStorageRepository{}
	processingService := &MockProcessingService{}
	notificationService := &MockNotificationService{}
	
	uc := &OpenCloseUseCase{
		videoRepo:           videoRepo,
		storageRepo:         storageRepo,
		processingService:   processingService,
		notificationService: notificationService,
		rawBucket:          "raw-bucket",
		processedBucket:    "processed-bucket",
		logoPath:           "logo.png",
		introSeconds:       2.5,
		outroSeconds:       2.5,
		targetW:            1920,
		targetH:            1080,
		fps:                30,
	}
	
	video := &domain.Video{
		ID:       "video-123",
		Filename: "test.mp4",
		Status:   domain.StatusPending,
	}
	
	inputData := []byte("input video data")
	processedData := []byte("processed video data")
	expectedError := errors.New("final status update failed")
	
	videoRepo.On("FindByFilename", "test.mp4").Return(video, nil)
	videoRepo.On("UpdateStatus", "video-123", domain.StatusProcessing).Return(nil)
	storageRepo.On("Download", "raw-bucket", "test.mp4").Return(inputData, nil)
	processingService.On("Process", inputData, "logo.png", 2.5, 2.5, 1920, 1080, 30).Return(processedData, nil)
	storageRepo.On("Upload", "processed-bucket", "test.mp4", processedData).Return(nil)
	videoRepo.On("UpdateStatus", "video-123", domain.StatusCompleted).Return(expectedError)
	
	err := uc.Execute("video-123", "test.mp4")
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "update final status")
	videoRepo.AssertExpectations(t)
	storageRepo.AssertExpectations(t)
	processingService.AssertExpectations(t)
}

func TestOpenCloseUseCase_Execute_NotificationError(t *testing.T) {
	videoRepo := &MockVideoRepository{}
	storageRepo := &MockStorageRepository{}
	processingService := &MockProcessingService{}
	notificationService := &MockNotificationService{}
	
	uc := &OpenCloseUseCase{
		videoRepo:           videoRepo,
		storageRepo:         storageRepo,
		processingService:   processingService,
		notificationService: notificationService,
		rawBucket:          "raw-bucket",
		processedBucket:    "processed-bucket",
		logoPath:           "logo.png",
		introSeconds:       2.5,
		outroSeconds:       2.5,
		targetW:            1920,
		targetH:            1080,
		fps:                30,
	}
	
	video := &domain.Video{
		ID:       "video-123",
		Filename: "test.mp4",
		Status:   domain.StatusPending,
	}
	
	inputData := []byte("input video data")
	processedData := []byte("processed video data")
	notificationError := errors.New("notification failed")
	
	videoRepo.On("FindByFilename", "test.mp4").Return(video, nil)
	videoRepo.On("UpdateStatus", "video-123", domain.StatusProcessing).Return(nil)
	storageRepo.On("Download", "raw-bucket", "test.mp4").Return(inputData, nil)
	processingService.On("Process", inputData, "logo.png", 2.5, 2.5, 1920, 1080, 30).Return(processedData, nil)
	storageRepo.On("Upload", "processed-bucket", "test.mp4", processedData).Return(nil)
	videoRepo.On("UpdateStatus", "video-123", domain.StatusCompleted).Return(nil)
	notificationService.On("NotifyVideoProcessed", "video-123", "test.mp4", "processed-bucket/test.mp4").Return(notificationError)
	
	// Notification error should not fail the entire process
	err := uc.Execute("video-123", "test.mp4")
	
	assert.NoError(t, err) // Should still succeed despite notification error
	videoRepo.AssertExpectations(t)
	storageRepo.AssertExpectations(t)
	processingService.AssertExpectations(t)
	notificationService.AssertExpectations(t)
}

func TestEnvFloat(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		envValue string
		defValue float64
		expected float64
	}{
		{"valid float", "TEST_FLOAT", "2.5", 1.0, 2.5},
		{"invalid float", "TEST_FLOAT", "invalid", 1.0, 1.0},
		{"empty env", "NONEXISTENT", "", 1.0, 1.0},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				t.Setenv(tt.key, tt.envValue)
			}
			
			result := envFloat(tt.key, tt.defValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEnvInt(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		envValue string
		defValue int
		expected int
	}{
		{"valid int", "TEST_INT", "30", 25, 30},
		{"invalid int", "TEST_INT", "invalid", 25, 25},
		{"empty env", "NONEXISTENT", "", 25, 25},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				t.Setenv(tt.key, tt.envValue)
			}
			
			result := envInt(tt.key, tt.defValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}