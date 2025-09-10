package ffmpeg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessWithTempFiles_PathTraversalPrevention(t *testing.T) {
	inputData := []byte("test data")
	
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "unix path traversal",
			args: []string{"-i", "../../../etc/passwd", "{output}"},
		},
		{
			name: "windows path traversal",
			args: []string{"-i", "..\\..\\windows\\system32\\config\\sam", "{output}"},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ProcessWithTempFiles(inputData, tt.args)
			
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "path traversal")
		})
	}
}

func TestProcessWithTempFiles_BasicValidation(t *testing.T) {
	inputData := []byte("test data")
	args := []string{"-i", "{input}", "-c", "copy", "{output}"}
	
	_, err := ProcessWithTempFiles(inputData, args)
	
	// Expected to fail due to invalid data, but tests the flow
	assert.Error(t, err)
}

func TestValidateFFmpeg(t *testing.T) {
	err := ValidateFFmpeg()
	
	// May pass or fail depending on system
	if err != nil {
		assert.Contains(t, err.Error(), "ffmpeg not found")
	}
}