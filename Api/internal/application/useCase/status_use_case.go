package useCase

import (
	"api/internal/domain/entities"
	"context"
)

// StatusService provides operations related to static domain statuses.
type StatusService struct{}

func NewStatusService() *StatusService { return &StatusService{} }

// ListVideoStatuses returns all available video processing statuses.
func (s *StatusService) ListVideoStatuses(ctx context.Context) []string {
	statuses := entities.AllVideoStatuses()
	out := make([]string, 0, len(statuses))
	for _, st := range statuses {
		out = append(out, string(st))
	}
	return out
}
