package usecases_test

import (
	"audioremoval/internal/application/usecases"
	"audioremoval/internal/domain"
	"audioremoval/tests/mocks"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProcessVideoUseCase_Execute_Success(t *testing.T) {
	// Arrange
	videoRepo := mocks.NewVideoRepositoryMock()
	storageRepo := mocks.NewStorageRepositoryMock()
	processingService := mocks.NewVideoProcessingServiceMock()
	notificationService := mocks.NewNotificationServiceMock()

	video := &domain.Video{
		ID:       "video-123",
		Filename: "test.mp4",
		Status:   domain.StatusPending,
		CreatedAt: time.Now(),
	}
	videoRepo.Videos["video-123"] = video

	inputData := []byte("input video data")
	storageRepo.Files["raw-bucket/test.mp4"] = inputData

	useCase := usecases.NewProcessVideoUseCase(
		videoRepo,
		storageRepo,
		processingService,
		notificationService,
		"raw-bucket",
		"processed-bucket",
	)

	// Act
	err := useCase.Execute("video-123", "test.mp4")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, domain.StatusCompleted, videoRepo.StatusUpdates["video-123"])
	assert.Len(t, storageRepo.UploadCalls, 1)
	assert.Equal(t, "processed-bucket", storageRepo.UploadCalls[0].Bucket)
	assert.Equal(t, "test.mp4", storageRepo.UploadCalls[0].Filename)
	assert.Len(t, notificationService.VideoProcessedCalls, 1)
	assert.Equal(t, "video-123", notificationService.VideoProcessedCalls[0].VideoID)
}

func TestProcessVideoUseCase_Execute_VideoNotFound(t *testing.T) {
	// Arrange
	videoRepo := mocks.NewVideoRepositoryMock()
	storageRepo := mocks.NewStorageRepositoryMock()
	processingService := mocks.NewVideoProcessingServiceMock()
	notificationService := mocks.NewNotificationServiceMock()

	useCase := usecases.NewProcessVideoUseCase(
		videoRepo,
		storageRepo,
		processingService,
		notificationService,
		"raw-bucket",
		"processed-bucket",
	)

	// Act
	err := useCase.Execute("video-123", "nonexistent.mp4")

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "video not found")
}

func TestProcessVideoUseCase_Execute_DownloadFails(t *testing.T) {
	// Arrange
	videoRepo := mocks.NewVideoRepositoryMock()
	storageRepo := mocks.NewStorageRepositoryMock()
	processingService := mocks.NewVideoProcessingServiceMock()
	notificationService := mocks.NewNotificationServiceMock()

	video := &domain.Video{
		ID:       "video-123",
		Filename: "test.mp4",
		Status:   domain.StatusPending,
	}
	videoRepo.Videos["video-123"] = video

	storageRepo.DownloadFunc = func(bucket, filename string) ([]byte, error) {
		return nil, errors.New("download failed")
	}

	useCase := usecases.NewProcessVideoUseCase(
		videoRepo,
		storageRepo,
		processingService,
		notificationService,
		"raw-bucket",
		"processed-bucket",
	)

	// Act
	err := useCase.Execute("video-123", "test.mp4")

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to download video")
	assert.Equal(t, domain.StatusFailed, videoRepo.StatusUpdates["video-123"])
}

func TestProcessVideoUseCase_Execute_ProcessingFails(t *testing.T) {
	// Arrange
	videoRepo := mocks.NewVideoRepositoryMock()
	storageRepo := mocks.NewStorageRepositoryMock()
	processingService := mocks.NewVideoProcessingServiceMock()
	notificationService := mocks.NewNotificationServiceMock()

	video := &domain.Video{
		ID:       "video-123",
		Filename: "test.mp4",
		Status:   domain.StatusPending,
	}
	videoRepo.Videos["video-123"] = video

	inputData := []byte("input video data")
	storageRepo.Files["raw-bucket/test.mp4"] = inputData

	processingService.ShouldFail = true

	useCase := usecases.NewProcessVideoUseCase(
		videoRepo,
		storageRepo,
		processingService,
		notificationService,
		"raw-bucket",
		"processed-bucket",
	)

	// Act
	err := useCase.Execute("video-123", "test.mp4")

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to process video")
	assert.Equal(t, domain.StatusFailed, videoRepo.StatusUpdates["video-123"])
}

func TestProcessVideoUseCase_Execute_UploadFails(t *testing.T) {
	// Arrange
	videoRepo := mocks.NewVideoRepositoryMock()
	storageRepo := mocks.NewStorageRepositoryMock()
	processingService := mocks.NewVideoProcessingServiceMock()
	notificationService := mocks.NewNotificationServiceMock()

	video := &domain.Video{
		ID:       "video-123",
		Filename: "test.mp4",
		Status:   domain.StatusPending,
	}
	videoRepo.Videos["video-123"] = video

	inputData := []byte("input video data")
	storageRepo.Files["raw-bucket/test.mp4"] = inputData

	storageRepo.UploadFunc = func(bucket, filename string, data []byte) error {
		return errors.New("upload failed")
	}

	useCase := usecases.NewProcessVideoUseCase(
		videoRepo,
		storageRepo,
		processingService,
		notificationService,
		"raw-bucket",
		"processed-bucket",
	)

	// Act
	err := useCase.Execute("video-123", "test.mp4")

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to upload processed video")
	assert.Equal(t, domain.StatusFailed, videoRepo.StatusUpdates["video-123"])
}