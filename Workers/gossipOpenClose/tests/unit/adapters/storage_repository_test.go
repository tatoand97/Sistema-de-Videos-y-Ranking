package adapters

import (
	"gossipopenclose/internal/adapters"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStorageRepository(t *testing.T) {
	repo := adapters.NewStorageRepository(nil)
	assert.NotNil(t, repo)
}

func TestStorageRepository_Structure(t *testing.T) {
	repo := adapters.NewStorageRepository(nil)
	assert.NotNil(t, repo)
}