package application_test

import (
	"api/internal/application/useCase"
	"api/internal/domain"
	"api/internal/domain/entities"
	"api/tests/mocks"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUploadsUseCase_ListUserVideos(t *testing.T) {
	tests := []struct {
		name     string
		userID   uint
		mockRepo *mocks.MockVideoRepository
		want     []*entities.Video
		wantErr  bool
	}{
		{
			name:   "successful list",
			userID: 1,
			mockRepo: &mocks.MockVideoRepository{
				ListByUserFunc: func(ctx context.Context, userID uint) ([]*entities.Video, error) {
					return []*entities.Video{
						{VideoID: 1, UserID: userID, Title: "Video 1"},
						{VideoID: 2, UserID: userID, Title: "Video 2"},
					}, nil
				},
			},
			want: []*entities.Video{
				{VideoID: 1, UserID: 1, Title: "Video 1"},
				{VideoID: 2, UserID: 1, Title: "Video 2"},
			},
			wantErr: false,
		},
		{
			name:   "repository error",
			userID: 1,
			mockRepo: &mocks.MockVideoRepository{
				ListByUserFunc: func(ctx context.Context, userID uint) ([]*entities.Video, error) {
					return nil, errors.New("database error")
				},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usecase := useCase.NewUploadsUseCase(tt.mockRepo, nil, nil, "")

			result, err := usecase.ListUserVideos(context.Background(), tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, result)
			}
		})
	}
}

func TestUploadsUseCase_GetUserVideoByID(t *testing.T) {
	tests := []struct {
		name     string
		userID   uint
		videoID  uint
		mockRepo *mocks.MockVideoRepository
		want     *entities.Video
		wantErr  bool
	}{
		{
			name:    "successful get",
			userID:  1,
			videoID: 123,
			mockRepo: &mocks.MockVideoRepository{
				GetByIDAndUserFunc: func(ctx context.Context, id, userID uint) (*entities.Video, error) {
					return &entities.Video{
						VideoID: id,
						UserID:  userID,
						Title:   "Test Video",
					}, nil
				},
			},
			want: &entities.Video{
				VideoID: 123,
				UserID:  1,
				Title:   "Test Video",
			},
			wantErr: false,
		},
		{
			name:    "video not found",
			userID:  1,
			videoID: 999,
			mockRepo: &mocks.MockVideoRepository{
				GetByIDAndUserFunc: func(ctx context.Context, id, userID uint) (*entities.Video, error) {
					return nil, domain.ErrNotFound
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "forbidden access",
			userID:  1,
			videoID: 123,
			mockRepo: &mocks.MockVideoRepository{
				GetByIDAndUserFunc: func(ctx context.Context, id, userID uint) (*entities.Video, error) {
					return nil, domain.ErrForbidden
				},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usecase := useCase.NewUploadsUseCase(tt.mockRepo, nil, nil, "")

			result, err := usecase.GetUserVideoByID(context.Background(), tt.userID, tt.videoID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, result)
			}
		})
	}
}
