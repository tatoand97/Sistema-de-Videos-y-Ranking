package usecases

import (
	"errors"
	"statesmachine/internal/application/usecases"
	"statesmachine/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestOrchestrateVideoUseCase_HandleEditCompleted_Success(t *testing.T) {
	videoRepo := &MockVideoRepository{}
	publisher := &MockMessagePublisher{}

	publisher.On("PublishMessage", "audio-queue", mock.AnythingOfType("[]uint8")).Return(nil)
	videoRepo.On("UpdateStatus", uint(123), domain.StatusRemovingAudio).Return(nil)

	useCase := usecases.NewOrchestrateVideoUseCase(
		videoRepo,
		publisher,
		"edit-queue",
		"audio-queue",
		"watermark-queue",
		3,
		5,
	)

	err := useCase.HandleEditCompleted("123", "edited.mp4")

	assert.NoError(t, err)
	videoRepo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestOrchestrateVideoUseCase_HandleAudioRemovalCompleted_Success(t *testing.T) {
	videoRepo := &MockVideoRepository{}
	publisher := &MockMessagePublisher{}

	publisher.On("PublishMessage", "watermark-queue", mock.AnythingOfType("[]uint8")).Return(nil)
	videoRepo.On("UpdateStatus", uint(123), domain.StatusAddingWatermark).Return(nil)

	useCase := usecases.NewOrchestrateVideoUseCase(
		videoRepo,
		publisher,
		"edit-queue",
		"audio-queue",
		"watermark-queue",
		3,
		5,
	)

	err := useCase.HandleAudioRemovalCompleted("123", "no-audio.mp4")

	assert.NoError(t, err)
	videoRepo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestOrchestrateVideoUseCase_HandleWatermarkingCompleted_Success(t *testing.T) {
	videoRepo := &MockVideoRepository{}
	publisher := &MockMessagePublisher{}

	publisher.On("PublishMessage", "gossip_open_close_queue", mock.AnythingOfType("[]uint8")).Return(nil)
	videoRepo.On("UpdateStatus", uint(123), domain.StatusAddingIntroOutro).Return(nil)

	useCase := usecases.NewOrchestrateVideoUseCase(
		videoRepo,
		publisher,
		"edit-queue",
		"audio-queue",
		"watermark-queue",
		3,
		5,
	)

	err := useCase.HandleWatermarkingCompleted("123", "watermarked.mp4")

	assert.NoError(t, err)
	videoRepo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestOrchestrateVideoUseCase_GetRetryDelayMinutes(t *testing.T) {
	videoRepo := &MockVideoRepository{}
	publisher := &MockMessagePublisher{}

	useCase := usecases.NewOrchestrateVideoUseCase(
		videoRepo,
		publisher,
		"edit-queue",
		"audio-queue",
		"watermark-queue",
		3,
		10,
	)

	delay := useCase.GetRetryDelayMinutes()
	assert.Equal(t, 10, delay)
}

func TestOrchestrateVideoUseCase_Execute_VideoNotFound(t *testing.T) {
	videoRepo := &MockVideoRepository{}
	publisher := &MockMessagePublisher{}

	videoRepo.On("FindByID", uint(123)).Return(nil, errors.New("video not found"))

	useCase := usecases.NewOrchestrateVideoUseCase(
		videoRepo,
		publisher,
		"edit-queue",
		"audio-queue",
		"watermark-queue",
		3,
		5,
	)

	err := useCase.Execute("123")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "find video")
	videoRepo.AssertExpectations(t)
}

func TestOrchestrateVideoUseCase_HandleTrimCompleted_PublishError(t *testing.T) {
	videoRepo := &MockVideoRepository{}
	publisher := &MockMessagePublisher{}

	publisher.On("PublishMessage", "edit-queue", mock.AnythingOfType("[]uint8")).Return(errors.New("publish failed"))

	useCase := usecases.NewOrchestrateVideoUseCase(
		videoRepo,
		publisher,
		"edit-queue",
		"audio-queue",
		"watermark-queue",
		3,
		5,
	)

	err := useCase.HandleTrimCompleted("123", "trimmed.mp4")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "publish to edit_video_queue")
	publisher.AssertExpectations(t)
}