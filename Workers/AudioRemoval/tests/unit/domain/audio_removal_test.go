package domain_test

import (
	"errors"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type AudioRemovalRequest struct {
	VideoID    string `json:"video_id"`
	InputPath  string `json:"input_path"`
	OutputPath string `json:"output_path"`
}

func (r *AudioRemovalRequest) Validate() error {
	if r.VideoID == "" {
		return errors.New("video_id is required")
	}
	if r.InputPath == "" {
		return errors.New("input_path is required")
	}
	if r.OutputPath == "" {
		return errors.New("output_path is required")
	}
	return nil
}

func TestAudioRemovalRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request AudioRemovalRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request",
			request: AudioRemovalRequest{
				VideoID:    "video-123",
				InputPath:  "/input/video.mp4",
				OutputPath: "/output/video.mp4",
			},
			wantErr: false,
		},
		{
			name: "missing video_id",
			request: AudioRemovalRequest{
				InputPath:  "/input/video.mp4",
				OutputPath: "/output/video.mp4",
			},
			wantErr: true,
			errMsg:  "video_id is required",
		},
		{
			name: "missing input_path",
			request: AudioRemovalRequest{
				VideoID:    "video-123",
				OutputPath: "/output/video.mp4",
			},
			wantErr: true,
			errMsg:  "input_path is required",
		},
		{
			name: "missing output_path",
			request: AudioRemovalRequest{
				VideoID:   "video-123",
				InputPath: "/input/video.mp4",
			},
			wantErr: true,
			errMsg:  "output_path is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}