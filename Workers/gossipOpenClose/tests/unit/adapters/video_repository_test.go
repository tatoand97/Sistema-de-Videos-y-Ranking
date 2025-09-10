package adapters

import (
	"gossipopenclose/internal/adapters"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewVideoRepository(t *testing.T) {
	repo := adapters.NewVideoRepository()
	assert.NotNil(t, repo)
}

func TestVideoRepository_Structure(t *testing.T) {
	repo := adapters.NewVideoRepository()
	assert.NotNil(t, repo)
}