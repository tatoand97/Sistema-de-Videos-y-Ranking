package domain_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

type WatermarkConfig struct {
	LogoPath string
	Position string
	Scale    float64
}

func TestWatermarkConfig_Creation(t *testing.T) {
	config := WatermarkConfig{
		LogoPath: "/assets/logo.png",
		Position: "top-right",
		Scale:    0.1,
	}
	
	assert.Equal(t, "/assets/logo.png", config.LogoPath)
	assert.Equal(t, "top-right", config.Position)
	assert.Equal(t, 0.1, config.Scale)
}

func TestWatermarkConfig_Validation(t *testing.T) {
	tests := []struct {
		name   string
		config WatermarkConfig
		valid  bool
	}{
		{
			name: "valid config",
			config: WatermarkConfig{
				LogoPath: "/assets/logo.png",
				Position: "top-right",
				Scale:    0.1,
			},
			valid: true,
		},
		{
			name: "invalid scale",
			config: WatermarkConfig{
				LogoPath: "/assets/logo.png",
				Position: "top-right",
				Scale:    -0.1,
			},
			valid: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := tt.config.Scale > 0 && tt.config.LogoPath != ""
			assert.Equal(t, tt.valid, valid)
		})
	}
}