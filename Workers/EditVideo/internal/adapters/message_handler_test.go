package adapters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMessageHandler(t *testing.T) {
	handler := NewMessageHandler(nil)
	assert.NotNil(t, handler)
}

func TestMessageHandler_HandleMessage_InvalidJSON(t *testing.T) {
	handler := NewMessageHandler(nil)
	err := handler.HandleMessage([]byte("invalid json"))
	assert.Error(t, err)
}

func TestMessageHandler_HandleMessage_EmptyBody(t *testing.T) {
	handler := NewMessageHandler(nil)
	err := handler.HandleMessage([]byte(""))
	assert.Error(t, err)
}

func TestVideoMessage_Structure(t *testing.T) {
	msg := VideoMessage{
		VideoID:  "123",
		Filename: "test.mp4",
	}
	assert.Equal(t, "123", msg.VideoID)
	assert.Equal(t, "test.mp4", msg.Filename)
}