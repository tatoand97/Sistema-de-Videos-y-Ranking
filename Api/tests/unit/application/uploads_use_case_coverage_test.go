package application_test

import (
	"api/internal/application/useCase"
	"api/internal/domain/entities"
	"api/tests/mocks"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUploadsUseCase_Coverage(t *testing.T) {
	mockRepo := &mocks.MockVideoRepository{}
	mockStorage := &mocks.MockVideoStorage{}
	mockPublisher := &mocks.MockMessagePublisher{}
	
	uc := useCase.NewUploadsUseCase(mockRepo, mockStorage, mockPublisher, "test-queue")
	
	// Test constructor
	assert.NotNil(t, uc)
	
	// Test ListUserVideos
	mockRepo.ListByUserFunc = func(ctx context.Context, userID uint) ([]*entities.Video, error) {
		return []*entities.Video{{VideoID: 1, UserID: userID}}, nil
	}
	videos, err := uc.ListUserVideos(context.Background(), 1)
	assert.NoError(t, err)
	assert.Len(t, videos, 1)
	
	// Test GetUserVideoByID
	mockRepo.GetByIDAndUserFunc = func(ctx context.Context, id, userID uint) (*entities.Video, error) {
		return &entities.Video{VideoID: id, UserID: userID}, nil
	}
	video, err := uc.GetUserVideoByID(context.Background(), 1, 1)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), video.VideoID)
	
	// Test DeleteUserVideoIfEligible
	mockRepo.GetByIDAndUserFunc = func(ctx context.Context, id, userID uint) (*entities.Video, error) {
		return &entities.Video{VideoID: id, UserID: userID, Status: string(entities.StatusUploaded)}, nil
	}
	mockRepo.DeleteFunc = func(ctx context.Context, id uint) error {
		return nil
	}
	err = uc.DeleteUserVideoIfEligible(context.Background(), 1, 1)
	assert.NoError(t, err)
}