package domain_test

import (
	"errors"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TrimVideoRequest struct {
	VideoID    string  `json:"video_id"`
	InputPath  string  `json:"input_path"`
	OutputPath string  `json:"output_path"`
	StartTime  float64 `json:"start_time"`
	Duration   float64 `json:"duration"`
}

func (r *TrimVideoRequest) Validate() error {
	if r.VideoID == "" {
		return errors.New("video_id is required")
	}
	if r.InputPath == "" {
		return errors.New("input_path is required")
	}
	if r.OutputPath == "" {
		return errors.New("output_path is required")
	}
	if r.StartTime < 0 {
		return errors.New("start_time must be non-negative")
	}
	if r.Duration <= 0 {
		return errors.New("duration must be positive")
	}
	return nil
}

func TestTrimVideoRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request TrimVideoRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request",
			request: TrimVideoRequest{
				VideoID:    "video-123",
				InputPath:  "/input/video.mp4",
				OutputPath: "/output/video.mp4",
				StartTime:  10.5,
				Duration:   30.0,
			},
			wantErr: false,
		},
		{
			name: "missing video_id",
			request: TrimVideoRequest{
				InputPath:  "/input/video.mp4",
				OutputPath: "/output/video.mp4",
				StartTime:  0,
				Duration:   30,
			},
			wantErr: true,
			errMsg:  "video_id is required",
		},
		{
			name: "negative start_time",
			request: TrimVideoRequest{
				VideoID:    "video-123",
				InputPath:  "/input/video.mp4",
				OutputPath: "/output/video.mp4",
				StartTime:  -5,
				Duration:   30,
			},
			wantErr: true,
			errMsg:  "start_time must be non-negative",
		},
		{
			name: "zero duration",
			request: TrimVideoRequest{
				VideoID:    "video-123",
				InputPath:  "/input/video.mp4",
				OutputPath: "/output/video.mp4",
				StartTime:  10,
				Duration:   0,
			},
			wantErr: true,
			errMsg:  "duration must be positive",
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