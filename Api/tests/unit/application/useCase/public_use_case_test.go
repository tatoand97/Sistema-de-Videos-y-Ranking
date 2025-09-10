package application_test

import (
	usecase "api/internal/application/useCase"
	"api/internal/domain"
	"api/internal/domain/responses"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockPublicRepo struct {
	ListFunc     func(ctx context.Context) ([]responses.PublicVideoResponse, error)
	GetByIDFunc  func(ctx context.Context, id uint) (*responses.PublicVideoResponse, error)
	RankingsFunc func(ctx context.Context, city *string, page, pageSize int) ([]responses.RankingItem, error)
}

func (m *mockPublicRepo) ListPublicVideos(ctx context.Context) ([]responses.PublicVideoResponse, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx)
	}
	return nil, nil
}
func (m *mockPublicRepo) GetPublicByID(ctx context.Context, id uint) (*responses.PublicVideoResponse, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}
func (m *mockPublicRepo) Rankings(ctx context.Context, city *string, page, pageSize int) ([]responses.RankingItem, error) {
	if m.RankingsFunc != nil {
		return m.RankingsFunc(ctx, city, page, pageSize)
	}
	return nil, nil
}

// Satisfy new interface method; not used in these unit tests
func (m *mockPublicRepo) GetUsersBasicByIDs(ctx context.Context, ids []uint) ([]responses.UserBasic, error) {
	return []responses.UserBasic{}, nil
}

type mockVoteRepo struct {
	HasUserVotedFunc func(ctx context.Context, videoID, userID uint) (bool, error)
	CreateFunc       func(ctx context.Context, videoID, userID uint) error
}

func (m *mockVoteRepo) HasUserVoted(ctx context.Context, videoID, userID uint) (bool, error) {
	if m.HasUserVotedFunc != nil {
		return m.HasUserVotedFunc(ctx, videoID, userID)
	}
	return false, nil
}
func (m *mockVoteRepo) Create(ctx context.Context, videoID, userID uint) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, videoID, userID)
	}
	return nil
}

func TestPublicService_VotePublicVideo_NoVoteRepo(t *testing.T) {
	repo := &mockPublicRepo{GetByIDFunc: func(ctx context.Context, id uint) (*responses.PublicVideoResponse, error) {
		return &responses.PublicVideoResponse{VideoID: 1}, nil
	}}
	svc := usecase.NewPublicService(repo, nil)

	err := svc.VotePublicVideo(context.Background(), 1, 2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not configured")
}

func TestPublicService_VotePublicVideo_NotFound(t *testing.T) {
	repo := &mockPublicRepo{GetByIDFunc: func(ctx context.Context, id uint) (*responses.PublicVideoResponse, error) {
		return nil, domain.ErrNotFound
	}}
	votes := &mockVoteRepo{}
	svc := usecase.NewPublicService(repo, votes)

	err := svc.VotePublicVideo(context.Background(), 1, 1)
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestPublicService_VotePublicVideo_AlreadyVoted(t *testing.T) {
	repo := &mockPublicRepo{GetByIDFunc: func(ctx context.Context, id uint) (*responses.PublicVideoResponse, error) {
		return &responses.PublicVideoResponse{VideoID: id}, nil
	}}
	votes := &mockVoteRepo{HasUserVotedFunc: func(ctx context.Context, videoID, userID uint) (bool, error) { return true, nil }}
	svc := usecase.NewPublicService(repo, votes)

	err := svc.VotePublicVideo(context.Background(), 3, 5)
	assert.ErrorIs(t, err, domain.ErrConflict)
}

func TestPublicService_VotePublicVideo_CreateError(t *testing.T) {
	repo := &mockPublicRepo{GetByIDFunc: func(ctx context.Context, id uint) (*responses.PublicVideoResponse, error) {
		return &responses.PublicVideoResponse{VideoID: id}, nil
	}}
	votes := &mockVoteRepo{
		HasUserVotedFunc: func(ctx context.Context, videoID, userID uint) (bool, error) { return false, nil },
		CreateFunc:       func(ctx context.Context, videoID, userID uint) error { return errors.New("db fail") },
	}
	svc := usecase.NewPublicService(repo, votes)

	err := svc.VotePublicVideo(context.Background(), 3, 5)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db fail")
}

func TestPublicService_VotePublicVideo_Success(t *testing.T) {
	repo := &mockPublicRepo{GetByIDFunc: func(ctx context.Context, id uint) (*responses.PublicVideoResponse, error) {
		return &responses.PublicVideoResponse{VideoID: id}, nil
	}}
	called := false
	votes := &mockVoteRepo{
		HasUserVotedFunc: func(ctx context.Context, videoID, userID uint) (bool, error) { return false, nil },
		CreateFunc:       func(ctx context.Context, videoID, userID uint) error { called = true; return nil },
	}
	svc := usecase.NewPublicService(repo, votes)

	err := svc.VotePublicVideo(context.Background(), 10, 20)
	assert.NoError(t, err)
	assert.True(t, called)
}

func TestPublicService_List_Get_Rankings_Delegation(t *testing.T) {
	repo := &mockPublicRepo{
		ListFunc: func(ctx context.Context) ([]responses.PublicVideoResponse, error) {
			return []responses.PublicVideoResponse{{VideoID: 1}}, nil
		},
		GetByIDFunc: func(ctx context.Context, id uint) (*responses.PublicVideoResponse, error) {
			return &responses.PublicVideoResponse{VideoID: id}, nil
		},
		RankingsFunc: func(ctx context.Context, city *string, page, pageSize int) ([]responses.RankingItem, error) {
			return []responses.RankingItem{{Username: "alice", Votes: 3}}, nil
		},
	}
	votes := &mockVoteRepo{}
	svc := usecase.NewPublicService(repo, votes)

	list, err := svc.ListPublicVideos(context.Background())
	assert.NoError(t, err)
	assert.Len(t, list, 1)

	got, err := svc.GetPublicByID(context.Background(), 42)
	assert.NoError(t, err)
	assert.Equal(t, uint(42), got.VideoID)

	r, err := svc.Rankings(context.Background(), nil, 1, 10)
	assert.NoError(t, err)
	assert.Len(t, r, 1)
}
