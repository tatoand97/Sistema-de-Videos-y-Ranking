package usecases

import (
	"encoding/json"
	"errors"
	"statesmachine/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockVideoRepository struct {
	mock.Mock
}

func (m *MockVideoRepository) FindByID(id uint) (*domain.Video, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Video), args.Error(1)
}

func (m *MockVideoRepository) UpdateStatus(id uint, status domain.VideoStatus) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func (m *MockVideoRepository) UpdateStatusAndProcessedFile(id uint, status domain.VideoStatus, processedFile string) error {
	args := m.Called(id, status, processedFile)
	return args.Error(0)
}

type MockMessagePublisher struct {
	mock.Mock
}

func (m *MockMessagePublisher) PublishMessage(queue string, message []byte) error {
	args := m.Called(queue, message)
	return args.Error(0)
}

func TestNewOrchestrateVideoUseCase(t *testing.T) {
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
	
	assert.NotNil(t, useCase)
	assert.Equal(t, "edit-queue", useCase.editVideoQueue)
	assert.Equal(t, "audio-queue", useCase.audioRemovalQueue)
	assert.Equal(t, "watermark-queue", useCase.watermarkingQueue)
	assert.Equal(t, 3, useCase.maxRetries)
	assert.Equal(t, 5, useCase.retryDelayMinutes)
}

func TestOrchestrateVideoUseCase_Execute_Success(t *testing.T) {
	videoRepo := &MockVideoRepository{}
	publisher := &MockMessagePublisher{}
	
	video := &domain.Video{
		ID:           123,
		OriginalFile: "test.mp4",
		Status:       "UPLOADED",
	}
	
	videoRepo.On("FindByID", uint(123)).Return(video, nil)
	publisher.On("PublishMessage", "trim_video_queue", mock.AnythingOfType("[]uint8")).Return(nil)
	videoRepo.On("UpdateStatus", uint(123), domain.StatusTrimming).Return(nil)
	
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
	
	assert.NoError(t, err)
	videoRepo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestOrchestrateVideoUseCase_Execute_InvalidVideoID(t *testing.T) {
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
	
	err := useCase.Execute("invalid-id")
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid video ID format")
}

func TestOrchestrateVideoUseCase_Execute_VideoNotFound(t *testing.T) {
	videoRepo := &MockVideoRepository{}
	publisher := &MockMessagePublisher{}
	
	videoRepo.On("FindByID", uint(999)).Return(nil, errors.New("video not found"))
	
	useCase := NewOrchestrateVideoUseCase(
		videoRepo,
		publisher,
		"edit-queue",
		"audio-queue",
		"watermark-queue",
		3,
		5,
	)
	
	err := useCase.Execute("999")
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "find video")
	videoRepo.AssertExpectations(t)
}

func TestOrchestrateVideoUseCase_HandleTrimCompleted_Success(t *testing.T) {
	videoRepo := &MockVideoRepository{}
	publisher := &MockMessagePublisher{}
	
	publisher.On("PublishMessage", "edit-queue", mock.AnythingOfType("[]uint8")).Return(nil)
	videoRepo.On("UpdateStatus", uint(123), domain.StatusAdjustingRes).Return(nil)
	
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
	
	assert.NoError(t, err)
	videoRepo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestOrchestrateVideoUseCase_HandleEditCompleted_Success(t *testing.T) {
	videoRepo := &MockVideoRepository{}
	publisher := &MockMessagePublisher{}
	
	publisher.On("PublishMessage", "audio-queue", mock.AnythingOfType("[]uint8")).Return(nil)
	videoRepo.On("UpdateStatus", uint(123), domain.StatusRemovingAudio).Return(nil)
	
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
	
	assert.NoError(t, err)
	videoRepo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestOrchestrateVideoUseCase_HandleAudioRemovalCompleted_Success(t *testing.T) {
	videoRepo := &MockVideoRepository{}
	publisher := &MockMessagePublisher{}
	
	publisher.On("PublishMessage", "watermark-queue", mock.AnythingOfType("[]uint8")).Return(nil)
	videoRepo.On("UpdateStatus", uint(123), domain.StatusAddingWatermark).Return(nil)
	
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
	
	assert.NoError(t, err)
	videoRepo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestOrchestrateVideoUseCase_HandleWatermarkingCompleted_Success(t *testing.T) {
	videoRepo := &MockVideoRepository{}
	publisher := &MockMessagePublisher{}
	
	publisher.On("PublishMessage", "gossip_open_close_queue", mock.AnythingOfType("[]uint8")).Return(nil)
	videoRepo.On("UpdateStatus", uint(123), domain.StatusAddingIntroOutro).Return(nil)
	
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
	
	assert.NoError(t, err)
	videoRepo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestOrchestrateVideoUseCase_HandleGossipOpenCloseCompleted_Success(t *testing.T) {
	videoRepo := &MockVideoRepository{}
	publisher := &MockMessagePublisher{}
	
	videoRepo.On("UpdateStatusAndProcessedFile", uint(123), domain.StatusProcessed, "final.mp4").Return(nil)
	
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
	
	assert.NoError(t, err)
	videoRepo.AssertExpectations(t)
}

func TestOrchestrateVideoUseCase_GetRetryDelayMinutes(t *testing.T) {
	videoRepo := &MockVideoRepository{}
	publisher := &MockMessagePublisher{}
	
	useCase := NewOrchestrateVideoUseCase(
		videoRepo,
		publisher,
		"edit-queue",
		"audio-queue",
		"watermark-queue",
		3,
		10,
	)
	
	assert.Equal(t, 10, useCase.GetRetryDelayMinutes())
}

func TestWorkerMessage_Structure(t *testing.T) {
	msg := WorkerMessage{
		VideoID:     "123",
		Filename:    "test.mp4",
		RetryCount:  1,
		MaxRetries:  3,
		LastRetry:   1234567890,
	}
	
	assert.Equal(t, "123", msg.VideoID)
	assert.Equal(t, "test.mp4", msg.Filename)
	assert.Equal(t, 1, msg.RetryCount)
	assert.Equal(t, 3, msg.MaxRetries)
	assert.Equal(t, int64(1234567890), msg.LastRetry)
}

func TestWorkerMessage_JSONMarshaling(t *testing.T) {
	msg := WorkerMessage{
		VideoID:  "123",
		Filename: "test.mp4",
	}
	
	data, err := json.Marshal(msg)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "123")
	assert.Contains(t, string(data), "test.mp4")
	
	var unmarshaled WorkerMessage
	err = json.Unmarshal(data, &unmarshaled)
	assert.NoError(t, err)
	assert.Equal(t, msg.VideoID, unmarshaled.VideoID)
	assert.Equal(t, msg.Filename, unmarshaled.Filename)
}