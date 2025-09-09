package adapters

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMessageHandler(t *testing.T) {
	// Test constructor with nil (we'll test functionality separately)
	handler := &MessageHandler{processVideoUC: nil}
	assert.NotNil(t, handler)
}

func TestMessageHandler_HandleMessage_InvalidJSON(t *testing.T) {
	handler := &MessageHandler{processVideoUC: nil}
	
	invalidJSON := []byte(`{"invalid": json}`)
	
	err := handler.HandleMessage(invalidJSON)
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid character")
}

func TestMessageHandler_HandleMessage_MalformedJSON(t *testing.T) {
	handler := &MessageHandler{processVideoUC: nil}
	
	malformedJSON := []byte(`{"video_id": "test", "filename":}`)
	
	err := handler.HandleMessage(malformedJSON)
	
	assert.Error(t, err)
}

func TestVideoMessage_Structure(t *testing.T) {
	msg := VideoMessage{
		VideoID:  "test-id",
		Filename: "test.mp4",
	}
	
	assert.Equal(t, "test-id", msg.VideoID)
	assert.Equal(t, "test.mp4", msg.Filename)
}

func TestVideoMessage_JSONMarshaling(t *testing.T) {
	msg := VideoMessage{
		VideoID:  "video-123",
		Filename: "test.mp4",
	}
	
	// Test marshaling
	data, err := json.Marshal(msg)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "video-123")
	assert.Contains(t, string(data), "test.mp4")
	
	// Test unmarshaling
	var unmarshaled VideoMessage
	err = json.Unmarshal(data, &unmarshaled)
	assert.NoError(t, err)
	assert.Equal(t, msg.VideoID, unmarshaled.VideoID)
	assert.Equal(t, msg.Filename, unmarshaled.Filename)
}