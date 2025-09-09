package interfaces

import "context"

// VoteRepository define el comportamiento de persistencia para votos.
type VoteRepository interface {
	HasUserVoted(ctx context.Context, videoID, userID uint) (bool, error)
	Create(ctx context.Context, videoID, userID uint) error
}
