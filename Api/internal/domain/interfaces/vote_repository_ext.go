package interfaces

import "context"

// VoteRepositoryWithEvent es una extension opcional de VoteRepository
// que permite persistir un eventID unico para idempotencia fuerte.
type VoteRepositoryWithEvent interface {
	VoteRepository
	CreateWithEvent(ctx context.Context, videoID, userID uint, eventID *string) error
}
