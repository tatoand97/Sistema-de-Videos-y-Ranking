package application_test

import (
	usecase "api/internal/application/useCase"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatusService_ListVideoStatuses(t *testing.T) {
	svc := usecase.NewStatusService()
	got := svc.ListVideoStatuses(context.Background())

	expected := []string{
		"UPLOADED",
		"TRIMMING",
		"ADJUSTING_RESOLUTION",
		"ADDING_WATERMARK",
		"REMOVING_AUDIO",
		"ADDING_INTRO_OUTRO",
		"PROCESSED",
		"FAILED",
	}
	assert.Equal(t, expected, got)
}
