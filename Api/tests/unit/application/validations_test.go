package application_test

import (
	"api/internal/application/validations"
	"api/tests/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckMP4(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{
			name: "ftyp only should fail with strict validation",
			data: []byte{
				0x00, 0x00, 0x00, 0x20, 0x66, 0x74, 0x79, 0x70, // ftyp box
				0x69, 0x73, 0x6f, 0x6d, 0x00, 0x00, 0x02, 0x00,
				0x69, 0x73, 0x6f, 0x6d, 0x69, 0x73, 0x6f, 0x32,
				0x61, 0x76, 0x63, 0x31, 0x6d, 0x70, 0x34, 0x31,
			},
			wantErr: true,
		},
		{
			name:    "empty data",
			data:    []byte{},
			wantErr: true,
		},
		{
			name:    "too small data",
			data:    []byte{0x00, 0x00, 0x00},
			wantErr: true,
		},
		{
			name: "invalid MP4 header",
			data: []byte{
				0x00, 0x00, 0x00, 0x20, 0x69, 0x6e, 0x76, 0x61, // invalid header
				0x6c, 0x69, 0x64, 0x00, 0x00, 0x00, 0x02, 0x00,
				0x69, 0x73, 0x6f, 0x6d, 0x69, 0x73, 0x6f, 0x32,
				0x61, 0x76, 0x63, 0x31, 0x6d, 0x70, 0x34, 0x31,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			width, height, err := validations.CheckMP4(tt.data)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.GreaterOrEqual(t, width, 1920)
				assert.GreaterOrEqual(t, height, 1080)
			}
		})
	}
}
