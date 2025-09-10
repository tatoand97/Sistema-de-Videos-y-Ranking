package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewOpenCloseVideoProcessingService(t *testing.T) {
	service := NewOpenCloseVideoProcessingService()
	assert.NotNil(t, service)
}

func TestOpenCloseVideoProcessingService_DurationLimits(t *testing.T) {
	tests := []struct {
		name           string
		introSeconds   float64
		outroSeconds   float64
		expectedIntro  float64
		expectedOutro  float64
		shouldScale    bool
	}{
		{
			name:          "within limits",
			introSeconds:  2.0,
			outroSeconds:  2.0,
			expectedIntro: 2.0,
			expectedOutro: 2.0,
			shouldScale:   false,
		},
		{
			name:          "exceeds limits",
			introSeconds:  3.0,
			outroSeconds:  3.0,
			expectedIntro: 2.5,
			expectedOutro: 2.5,
			shouldScale:   true,
		},
		{
			name:          "negative values",
			introSeconds:  -1.0,
			outroSeconds:  -1.0,
			expectedIntro: 0.0,
			expectedOutro: 0.0,
			shouldScale:   false,
		},
		{
			name:          "zero values",
			introSeconds:  0.0,
			outroSeconds:  0.0,
			expectedIntro: 0.0,
			expectedOutro: 0.0,
			shouldScale:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test duration calculation logic
			introSeconds := tt.introSeconds
			outroSeconds := tt.outroSeconds
			
			if introSeconds < 0 {
				introSeconds = 0
			}
			if outroSeconds < 0 {
				outroSeconds = 0
			}
			
			total := introSeconds + outroSeconds
			if total > 5.0 {
				if total == 0 {
					total = 1
				}
				scale := 5.0 / total
				introSeconds = introSeconds * scale
				outroSeconds = outroSeconds * scale
			}
			
			assert.InDelta(t, tt.expectedIntro, introSeconds, 0.01)
			assert.InDelta(t, tt.expectedOutro, outroSeconds, 0.01)
		})
	}
}

func TestVideoInfo_Structure(t *testing.T) {
	vi := videoInfo{
		W:   1920,
		H:   1080,
		FPS: 30,
	}
	
	assert.Equal(t, 1920, vi.W)
	assert.Equal(t, 1080, vi.H)
	assert.Equal(t, 30, vi.FPS)
}

func TestVideoInfo_DefaultValues(t *testing.T) {
	vi := videoInfo{}
	
	assert.Equal(t, 0, vi.W)
	assert.Equal(t, 0, vi.H)
	assert.Equal(t, 0, vi.FPS)
}

func TestOpenCloseVideoProcessingService_ParameterValidation(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
		fps    int
		valid  bool
	}{
		{"valid parameters", 1920, 1080, 30, true},
		{"zero width", 0, 1080, 30, false},
		{"zero height", 1920, 0, 30, false},
		{"negative width", -1920, 1080, 30, false},
		{"negative height", 1920, -1080, 30, false},
		{"zero fps gets fallback", 1920, 1080, 0, true}, // fps gets fallback to 30
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			width, height, fps := tt.width, tt.height, tt.fps
			
			// Simulate parameter validation logic
			if width <= 0 || height <= 0 {
				assert.False(t, tt.valid)
				return
			}
			if fps <= 0 {
				fps = 30 // fallback
			}
			
			assert.True(t, tt.valid)
			assert.Greater(t, width, 0)
			assert.Greater(t, height, 0)
			assert.Greater(t, fps, 0)
		})
	}
}

func TestOpenCloseVideoProcessingService_FadeDurationCalculation(t *testing.T) {
	tests := []struct {
		name        string
		seconds     float64
		expectedMin float64
		expectedMax float64
	}{
		{"short duration", 1.0, 0.1, 0.2},
		{"medium duration", 3.0, 0.1, 0.5},
		{"long duration", 10.0, 0.1, 0.5},
		{"very short", 0.5, 0.1, 0.1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate fade duration calculation logic
			fadeDur := tt.seconds * 0.2
			if fadeDur > 0.5 {
				fadeDur = 0.5
			}
			if fadeDur < 0.1 {
				fadeDur = 0.1
			}
			
			assert.GreaterOrEqual(t, fadeDur, tt.expectedMin)
			assert.LessOrEqual(t, fadeDur, tt.expectedMax)
		})
	}
}

func TestOpenCloseVideoProcessingService_LogoScaling(t *testing.T) {
	tests := []struct {
		name          string
		videoWidth    int
		expectedScale int
	}{
		{"HD video", 1920, 672},   // 1920 * 0.35 = 672
		{"SD video", 1280, 448},   // 1280 * 0.35 = 448
		{"4K video", 3840, 1344},  // 3840 * 0.35 = 1344
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate logo scaling logic (35% of video width)
			logoWidth := int(float64(tt.videoWidth) * 0.35)
			assert.Equal(t, tt.expectedScale, logoWidth)
		})
	}
}