package usecases

import (
	"encoding/json"
	"statesmachine/internal/application/usecases"
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
	
	useCase := usecases.NewOrchestrateVideoUseCase(
		videoRepo,
		publisher,
		"edit-queue",
		"audio-queue",
		"watermark-queue",
		3,
		5,
	)
	
	assert.NotNil(t, useCase)
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
	
	assert.NoError(t, err)
	videoRepo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestOrchestrateVideoUseCase_Execute_InvalidVideoID(t *testing.T) {
	videoRepo := &MockVideoRepository{}
	publisher := &MockMessagePublisher{}
	
	useCase := usecases.NewOrchestrateVideoUseCase(
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

func TestOrchestrateVideoUseCase_HandleTrimCompleted_Success(t *testing.T) {
	videoRepo := &MockVideoRepository{}
	publisher := &MockMessagePublisher{}
	
	publisher.On("PublishMessage", "edit-queue", mock.AnythingOfType("[]uint8")).Return(nil)
	videoRepo.On("UpdateStatus", uint(123), domain.StatusAdjustingRes).Return(nil)
	
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
	
	assert.NoError(t, err)
	videoRepo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestOrchestrateVideoUseCase_HandleGossipOpenCloseCompleted_Success(t *testing.T) {
	videoRepo := &MockVideoRepository{}
	publisher := &MockMessagePublisher{}
	
	videoRepo.On("UpdateStatusAndProcessedFile", uint(123), domain.StatusProcessed, "http://localhost:8084/processed-videos/final.mp4").Return(nil)
	
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
	
	assert.NoError(t, err)
	videoRepo.AssertExpectations(t)
}

func TestWorkerMessage_JSONMarshaling(t *testing.T) {
	msg := usecases.WorkerMessage{
		VideoID:  "123",
		Filename: "test.mp4",
	}
	
	data, err := json.Marshal(msg)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "123")
	assert.Contains(t, string(data), "test.mp4")
}