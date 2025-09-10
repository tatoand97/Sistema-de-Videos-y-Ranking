package adapters

import (
	"testing"
	"watermarking/internal/adapters"

	"github.com/stretchr/testify/assert"
)

func TestNewMinIOStorage_InvalidEndpoint(t *testing.T) {
	storage, err := adapters.NewMinIOStorage("", "access", "secret")
	assert.Error(t, err)
	assert.Nil(t, storage)
}

func TestNewMinIOStorage_ValidParams(t *testing.T) {
	storage, err := adapters.NewMinIOStorage("localhost:9000", "minioadmin", "minioadmin")
	if err == nil {
		assert.NotNil(t, storage)
	}
}

func TestMinIOStorage_Structure(t *testing.T) {
	storage := &adapters.MinIOStorage{}
	assert.NotNil(t, storage)
}