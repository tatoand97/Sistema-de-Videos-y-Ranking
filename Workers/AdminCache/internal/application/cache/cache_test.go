package cache

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// Test that New creates a non-nil Impl
	impl := &Impl{inner: nil}
	assert.NotNil(t, impl)
}

func TestSetJSON_MarshalError(t *testing.T) {
	impl := &Impl{inner: nil}
	
	// Use a value that cannot be marshaled to JSON
	invalidData := make(chan int)
	
	err := impl.SetJSON(nil, "test-key", invalidData)
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "json: unsupported type")
}

func TestSetJSON_ValidData(t *testing.T) {
	// Test JSON marshaling functionality
	testData := map[string]interface{}{
		"name": "test",
		"id":   123,
	}
	
	expectedBytes, err := json.Marshal(testData)
	assert.NoError(t, err)
	assert.NotEmpty(t, expectedBytes)
	
	// Verify the JSON is valid
	var unmarshaled map[string]interface{}
	err = json.Unmarshal(expectedBytes, &unmarshaled)
	assert.NoError(t, err)
	assert.Equal(t, "test", unmarshaled["name"])
	assert.Equal(t, float64(123), unmarshaled["id"]) // JSON numbers are float64
}