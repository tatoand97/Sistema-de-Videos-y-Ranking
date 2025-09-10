package adapters_test

import (
	"editvideo/internal/adapters"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMessageHandler(t *testing.T) {
	handler := adapters.NewMessageHandler(nil)
	assert.NotNil(t, handler)
}

func TestVideoMessage_Structure(t *testing.T) {
	msg := adapters.VideoMessage{
		VideoID:  "123",
		Filename: "test.mp4",
	}
	assert.Equal(t, "123", msg.VideoID)
	assert.Equal(t, "test.mp4", msg.Filename)
}

func TestVideoMessage_JSONMarshaling(t *testing.T) {
	msg := adapters.VideoMessage{
		VideoID:  "video-123",
		Filename: "test.mp4",
	}
	
	data, err := json.Marshal(msg)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "video-123")
}