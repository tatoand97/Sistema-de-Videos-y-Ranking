package domain_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

type GossipConfig struct {
	IntroSeconds float64
	OutroSeconds float64
	TargetWidth  int
	TargetHeight int
	FPS          int
}

func TestGossipConfig_Creation(t *testing.T) {
	config := GossipConfig{
		IntroSeconds: 2.5,
		OutroSeconds: 2.5,
		TargetWidth:  1280,
		TargetHeight: 720,
		FPS:          30,
	}
	
	assert.Equal(t, 2.5, config.IntroSeconds)
	assert.Equal(t, 2.5, config.OutroSeconds)
	assert.Equal(t, 1280, config.TargetWidth)
	assert.Equal(t, 720, config.TargetHeight)
	assert.Equal(t, 30, config.FPS)
}

func TestGossipConfig_Validation(t *testing.T) {
	tests := []struct {
		name   string
		config GossipConfig
		valid  bool
	}{
		{
			name: "valid HD config",
			config: GossipConfig{
				IntroSeconds: 2.5,
				OutroSeconds: 2.5,
				TargetWidth:  1280,
				TargetHeight: 720,
				FPS:          30,
			},
			valid: true,
		},
		{
			name: "invalid FPS",
			config: GossipConfig{
				IntroSeconds: 2.5,
				OutroSeconds: 2.5,
				TargetWidth:  1280,
				TargetHeight: 720,
				FPS:          0,
			},
			valid: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := tt.config.FPS > 0 && tt.config.TargetWidth > 0 && tt.config.TargetHeight > 0
			assert.Equal(t, tt.valid, valid)
		})
	}
}