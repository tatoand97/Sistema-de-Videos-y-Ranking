package usecases

import (
	"gossipopenclose/internal/application/usecases"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProcessVideoUseCase(t *testing.T) {
	useCase := usecases.NewProcessVideoUseCase(
		nil,
		nil,
		nil,
		nil,
		"input-bucket",
		"output-bucket",
		"logo.png",
	)

	assert.NotNil(t, useCase)
}