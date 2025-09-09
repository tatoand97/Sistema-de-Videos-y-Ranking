package domain_test

import (
	"errors"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type EditVideoRequest struct {
	VideoID     string  `json:"video_id"`
	InputPath   string  `json:"input_path"`
	OutputPath  string  `json:"output_path"`
	StartTime   float64 `json:"start_time"`
	EndTime     float64 `json:"end_time"`
	Quality     string  `json:"quality"`
	Resolution  string  `json:"resolution"`
}

func (r *EditVideoRequest) Validate() error {
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
	if r.EndTime <= r.StartTime {
		return errors.New("end_time must be greater than start_time")
	}
	return nil
}

func TestEditVideoRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request EditVideoRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request",
			request: EditVideoRequest{
				VideoID:    "video-123",
				InputPath:  "/input/video.mp4",
				OutputPath: "/output/video.mp4",
				StartTime:  0,
				EndTime:    60,
				Quality:    "high",
				Resolution: "1080p",
			},
			wantErr: false,
		},
		{
			name: "missing video_id",
			request: EditVideoRequest{
				InputPath:  "/input/video.mp4",
				OutputPath: "/output/video.mp4",
				StartTime:  0,
				EndTime:    60,
			},
			wantErr: true,
			errMsg:  "video_id is required",
		},
		{
			name: "negative start_time",
			request: EditVideoRequest{
				VideoID:    "video-123",
				InputPath:  "/input/video.mp4",
				OutputPath: "/output/video.mp4",
				StartTime:  -5,
				EndTime:    60,
			},
			wantErr: true,
			errMsg:  "start_time must be non-negative",
		},
		{
			name: "end_time less than start_time",
			request: EditVideoRequest{
				VideoID:    "video-123",
				InputPath:  "/input/video.mp4",
				OutputPath: "/output/video.mp4",
				StartTime:  60,
				EndTime:    30,
			},
			wantErr: true,
			errMsg:  "end_time must be greater than start_time",
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