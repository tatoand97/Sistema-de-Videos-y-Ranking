package usecases

import (
	"errors"
	"statesmachine/internal/application/usecases"
	"statesmachine/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestOrchestrateVideoUseCase_Execute_PublishError(t *testing.T) {
	videoRepo := &MockVideoRepository{}
	publisher := &MockMessagePublisher{}

	video := &domain.Video{
		ID:           123,
		OriginalFile: "test.mp4",
		Status:       "UPLOADED",
	}

	videoRepo.On("FindByID", uint(123)).Return(video, nil)
	publisher.On("PublishMessage", "trim_video_queue", mock.AnythingOfType("[]uint8")).Return(errors.New("publish failed"))

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
	assert.Contains(t, err.Error(), "publish to trim_video_queue")
	videoRepo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestOrchestrateVideoUseCase_HandleTrimCompleted_InvalidVideoID(t *testing.T) {
	videoRepo := &MockVideoRepository{}
	publisher := &MockMessagePublisher{}

	publisher.On("PublishMessage", "edit-queue", mock.AnythingOfType("[]uint8")).Return(nil)

	useCase := usecases.NewOrchestrateVideoUseCase(
		videoRepo,
		publisher,
		"edit-queue",
		"audio-queue",
		"watermark-queue",
		3,
		5,
	)

	err := useCase.HandleTrimCompleted("invalid-id", "trimmed.mp4")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid video ID format")
	publisher.AssertExpectations(t)
}

func TestOrchestrateVideoUseCase_HandleGossipOpenCloseCompleted_UpdateError(t *testing.T) {
	videoRepo := &MockVideoRepository{}
	publisher := &MockMessagePublisher{}

	videoRepo.On("UpdateStatusAndProcessedFile", uint(123), domain.StatusProcessed, "http://localhost:8084/processed-videos/final.mp4").Return(errors.New("update failed"))

	useCase := usecases.NewOrchestrateVideoUseCase(
		videoRepo,
		publisher,
		"edit-queue",
		"audio-queue",
		"watermark-queue",
		3,
		5,
	)

	err := useCase.HandleGossipOpenCloseCompleted("123", "final.mp4")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "update final status and processed file")
	videoRepo.AssertExpectations(t)
}