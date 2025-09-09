package usecases

import (
	"errors"
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

	useCase := NewOrchestrateVideoUseCase(
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

func TestOrchestrateVideoUseCase_Execute_UpdateStatusError(t *testing.T) {
	videoRepo := &MockVideoRepository{}
	publisher := &MockMessagePublisher{}

	video := &domain.Video{
		ID:           123,
		OriginalFile: "test.mp4",
		Status:       "UPLOADED",
	}

	videoRepo.On("FindByID", uint(123)).Return(video, nil)
	publisher.On("PublishMessage", "trim_video_queue", mock.AnythingOfType("[]uint8")).Return(nil)
	videoRepo.On("UpdateStatus", uint(123), domain.StatusTrimming).Return(errors.New("update failed"))

	useCase := NewOrchestrateVideoUseCase(
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
	assert.Contains(t, err.Error(), "update status")
	videoRepo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestOrchestrateVideoUseCase_HandleTrimCompleted_InvalidVideoID(t *testing.T) {
	videoRepo := &MockVideoRepository{}
	publisher := &MockMessagePublisher{}

	// Mock the PublishMessage call that happens before ID validation
	publisher.On("PublishMessage", "edit-queue", mock.AnythingOfType("[]uint8")).Return(nil)

	useCase := NewOrchestrateVideoUseCase(
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

func TestOrchestrateVideoUseCase_HandleTrimCompleted_PublishError(t *testing.T) {
	videoRepo := &MockVideoRepository{}
	publisher := &MockMessagePublisher{}

	publisher.On("PublishMessage", "edit-queue", mock.AnythingOfType("[]uint8")).Return(errors.New("publish failed"))

	useCase := NewOrchestrateVideoUseCase(
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

func TestOrchestrateVideoUseCase_HandleTrimCompleted_UpdateStatusError(t *testing.T) {
	videoRepo := &MockVideoRepository{}
	publisher := &MockMessagePublisher{}

	publisher.On("PublishMessage", "edit-queue", mock.AnythingOfType("[]uint8")).Return(nil)
	videoRepo.On("UpdateStatus", uint(123), domain.StatusAdjustingRes).Return(errors.New("update failed"))

	useCase := NewOrchestrateVideoUseCase(
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
	assert.Contains(t, err.Error(), "update status")
	videoRepo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestOrchestrateVideoUseCase_HandleEditCompleted_InvalidVideoID(t *testing.T) {
	videoRepo := &MockVideoRepository{}
	publisher := &MockMessagePublisher{}

	// Mock the PublishMessage call that happens before ID validation
	publisher.On("PublishMessage", "audio-queue", mock.AnythingOfType("[]uint8")).Return(nil)

	useCase := NewOrchestrateVideoUseCase(
		videoRepo,
		publisher,
		"edit-queue",
		"audio-queue",
		"watermark-queue",
		3,
		5,
	)

	err := useCase.HandleEditCompleted("invalid-id", "edited.mp4")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid video ID format")
	publisher.AssertExpectations(t)
}

func TestOrchestrateVideoUseCase_HandleEditCompleted_PublishError(t *testing.T) {
	videoRepo := &MockVideoRepository{}
	publisher := &MockMessagePublisher{}

	publisher.On("PublishMessage", "audio-queue", mock.AnythingOfType("[]uint8")).Return(errors.New("publish failed"))

	useCase := NewOrchestrateVideoUseCase(
		videoRepo,
		publisher,
		"edit-queue",
		"audio-queue",
		"watermark-queue",
		3,
		5,
	)

	err := useCase.HandleEditCompleted("123", "edited.mp4")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "publish to audio_removal_queue")
	publisher.AssertExpectations(t)
}

func TestOrchestrateVideoUseCase_HandleEditCompleted_UpdateStatusError(t *testing.T) {
	videoRepo := &MockVideoRepository{}
	publisher := &MockMessagePublisher{}

	publisher.On("PublishMessage", "audio-queue", mock.AnythingOfType("[]uint8")).Return(nil)
	videoRepo.On("UpdateStatus", uint(123), domain.StatusRemovingAudio).Return(errors.New("update failed"))

	useCase := NewOrchestrateVideoUseCase(
		videoRepo,
		publisher,
		"edit-queue",
		"audio-queue",
		"watermark-queue",
		3,
		5,
	)

	err := useCase.HandleEditCompleted("123", "edited.mp4")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "update status")
	videoRepo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestOrchestrateVideoUseCase_HandleAudioRemovalCompleted_InvalidVideoID(t *testing.T) {
	videoRepo := &MockVideoRepository{}
	publisher := &MockMessagePublisher{}

	// Mock the PublishMessage call that happens before ID validation
	publisher.On("PublishMessage", "watermark-queue", mock.AnythingOfType("[]uint8")).Return(nil)

	useCase := NewOrchestrateVideoUseCase(
		videoRepo,
		publisher,
		"edit-queue",
		"audio-queue",
		"watermark-queue",
		3,
		5,
	)

	err := useCase.HandleAudioRemovalCompleted("invalid-id", "no-audio.mp4")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid video ID format")
	publisher.AssertExpectations(t)
}

func TestOrchestrateVideoUseCase_HandleAudioRemovalCompleted_PublishError(t *testing.T) {
	videoRepo := &MockVideoRepository{}
	publisher := &MockMessagePublisher{}

	publisher.On("PublishMessage", "watermark-queue", mock.AnythingOfType("[]uint8")).Return(errors.New("publish failed"))

	useCase := NewOrchestrateVideoUseCase(
		videoRepo,
		publisher,
		"edit-queue",
		"audio-queue",
		"watermark-queue",
		3,
		5,
	)

	err := useCase.HandleAudioRemovalCompleted("123", "no-audio.mp4")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "publish to watermarking_queue")
	publisher.AssertExpectations(t)
}

func TestOrchestrateVideoUseCase_HandleAudioRemovalCompleted_UpdateStatusError(t *testing.T) {
	videoRepo := &MockVideoRepository{}
	publisher := &MockMessagePublisher{}

	publisher.On("PublishMessage", "watermark-queue", mock.AnythingOfType("[]uint8")).Return(nil)
	videoRepo.On("UpdateStatus", uint(123), domain.StatusAddingWatermark).Return(errors.New("update failed"))

	useCase := NewOrchestrateVideoUseCase(
		videoRepo,
		publisher,
		"edit-queue",
		"audio-queue",
		"watermark-queue",
		3,
		5,
	)

	err := useCase.HandleAudioRemovalCompleted("123", "no-audio.mp4")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "update status")
	videoRepo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestOrchestrateVideoUseCase_HandleWatermarkingCompleted_InvalidVideoID(t *testing.T) {
	videoRepo := &MockVideoRepository{}
	publisher := &MockMessagePublisher{}

	// Mock the PublishMessage call that happens before ID validation
	publisher.On("PublishMessage", "gossip_open_close_queue", mock.AnythingOfType("[]uint8")).Return(nil)

	useCase := NewOrchestrateVideoUseCase(
		videoRepo,
		publisher,
		"edit-queue",
		"audio-queue",
		"watermark-queue",
		3,
		5,
	)

	err := useCase.HandleWatermarkingCompleted("invalid-id", "watermarked.mp4")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid video ID format")
	publisher.AssertExpectations(t)
}

func TestOrchestrateVideoUseCase_HandleWatermarkingCompleted_PublishError(t *testing.T) {
	videoRepo := &MockVideoRepository{}
	publisher := &MockMessagePublisher{}

	publisher.On("PublishMessage", "gossip_open_close_queue", mock.AnythingOfType("[]uint8")).Return(errors.New("publish failed"))

	useCase := NewOrchestrateVideoUseCase(
		videoRepo,
		publisher,
		"edit-queue",
		"audio-queue",
		"watermark-queue",
		3,
		5,
	)

	err := useCase.HandleWatermarkingCompleted("123", "watermarked.mp4")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "publish to gossip_open_close_queue")
	publisher.AssertExpectations(t)
}

func TestOrchestrateVideoUseCase_HandleWatermarkingCompleted_UpdateStatusError(t *testing.T) {
	videoRepo := &MockVideoRepository{}
	publisher := &MockMessagePublisher{}

	publisher.On("PublishMessage", "gossip_open_close_queue", mock.AnythingOfType("[]uint8")).Return(nil)
	videoRepo.On("UpdateStatus", uint(123), domain.StatusAddingIntroOutro).Return(errors.New("update failed"))

	useCase := NewOrchestrateVideoUseCase(
		videoRepo,
		publisher,
		"edit-queue",
		"audio-queue",
		"watermark-queue",
		3,
		5,
	)

	err := useCase.HandleWatermarkingCompleted("123", "watermarked.mp4")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "update status")
	videoRepo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestOrchestrateVideoUseCase_HandleGossipOpenCloseCompleted_InvalidVideoID(t *testing.T) {
	videoRepo := &MockVideoRepository{}
	publisher := &MockMessagePublisher{}

	useCase := NewOrchestrateVideoUseCase(
		videoRepo,
		publisher,
		"edit-queue",
		"audio-queue",
		"watermark-queue",
		3,
		5,
	)

	err := useCase.HandleGossipOpenCloseCompleted("invalid-id", "final.mp4")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid video ID format")
}

func TestOrchestrateVideoUseCase_HandleGossipOpenCloseCompleted_UpdateError(t *testing.T) {
	videoRepo := &MockVideoRepository{}
	publisher := &MockMessagePublisher{}

	videoRepo.On("UpdateStatusAndProcessedFile", uint(123), domain.StatusProcessed, "final.mp4").Return(errors.New("update failed"))

	useCase := NewOrchestrateVideoUseCase(
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