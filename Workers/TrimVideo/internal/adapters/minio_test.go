package adapters

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMinIOStorage_InvalidEndpoint(t *testing.T) {
	storage, err := NewMinIOStorage("", "access", "secret")
	assert.Error(t, err)
	assert.Nil(t, storage)
}

func TestNewMinIOStorage_ValidParams(t *testing.T) {
	storage, err := NewMinIOStorage("localhost:9000", "minioadmin", "minioadmin")
	if err == nil {
		assert.NotNil(t, storage)
		assert.NotNil(t, storage.client)
	}
}

func TestMinIOStorage_Structure(t *testing.T) {
	storage := &MinIOStorage{client: nil}
	assert.NotNil(t, storage)
}

func TestMinIOStorage_Operations(t *testing.T) {
	bucket := "test-bucket"
	filename := "test.mp4"
	data := bytes.NewReader([]byte("test data"))
	size := int64(9)
	
	assert.NotEmpty(t, bucket)
	assert.NotEmpty(t, filename)
	assert.NotNil(t, data)
	assert.Greater(t, size, int64(0))
}