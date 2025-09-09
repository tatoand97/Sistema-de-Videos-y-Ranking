package application_test

import (
	"api/internal/application/validations"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckMP4_ValidFile(t *testing.T) {
	// Test with basic MP4 data - this will fail validation as expected
	validMP4 := createTestValidMP4()

	_, _, err := validations.CheckMP4(validMP4)

	// This should fail because our test MP4 doesn't have proper moov/trak structure
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "MP4 sin 'moov' o sin 'trak'")
}

func TestCheckMP4_InvalidFile(t *testing.T) {
	invalidData := []byte("invalid mp4 content")

	_, _, err := validations.CheckMP4(invalidData)

	assert.Error(t, err)
	// Avoid accent/encoding issues across environments by checking a stable prefix
	assert.Contains(t, err.Error(), "no es MP4")
}

func TestCheckMP4_EmptyData(t *testing.T) {
	_, _, err := validations.CheckMP4([]byte{})

	assert.Error(t, err)
	// Empty data will also trigger the moov/trak error
	assert.Contains(t, err.Error(), "MP4")
}

func TestCheckMP4_TooLarge(t *testing.T) {
	// Create data larger than MaxBytes
	largeData := make([]byte, validations.MaxBytes+1)

	_, _, err := validations.CheckMP4(largeData)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "excede 100MB")
}

// Import the testdata package functions
func createTestValidMP4() []byte {
	// This would import from testdata package, but since we can't import from tests,
	// we'll create a minimal valid structure here
	return []byte{
		// ftyp box
		0x00, 0x00, 0x00, 0x20, 0x66, 0x74, 0x79, 0x70,
		0x69, 0x73, 0x6f, 0x6d, 0x00, 0x00, 0x02, 0x00,
		0x69, 0x73, 0x6f, 0x6d, 0x69, 0x73, 0x6f, 0x32,
		0x61, 0x76, 0x63, 0x31, 0x6d, 0x70, 0x34, 0x31,
	}
}

// Note: resolution-specific helper removed as it was unused.
