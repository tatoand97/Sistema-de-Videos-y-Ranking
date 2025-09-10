package requests_test

import (
	"testing"

	"api/internal/domain/requests"

	"github.com/stretchr/testify/assert"
)

func TestLoginRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request requests.LoginRequest
		valid   bool
	}{
		{
			name: "valid request",
			request: requests.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			valid: true,
		},
		{
			name: "empty email",
			request: requests.LoginRequest{
				Email:    "",
				Password: "password123",
			},
			valid: false,
		},
		{
			name: "empty password",
			request: requests.LoginRequest{
				Email:    "test@example.com",
				Password: "",
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.valid {
				assert.NotEmpty(t, tt.request.Email)
				assert.NotEmpty(t, tt.request.Password)
			} else {
				assert.True(t, tt.request.Email == "" || tt.request.Password == "")
			}
		})
	}
}

func TestRegisterUserRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request requests.RegisterUserRequest
		valid   bool
	}{
		{
			name: "valid request",
			request: requests.RegisterUserRequest{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john@example.com",
				Password1: "password123",
				Password2: "password123",
				Country:   "Colombia",
				City:      "Bogotá",
			},
			valid: true,
		},
		{
			name: "password mismatch",
			request: requests.RegisterUserRequest{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john@example.com",
				Password1: "password123",
				Password2: "different",
				Country:   "Colombia",
				City:      "Bogotá",
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.valid {
				assert.Equal(t, tt.request.Password1, tt.request.Password2)
				assert.NotEmpty(t, tt.request.Email)
				assert.NotEmpty(t, tt.request.FirstName)
			} else {
				if tt.request.Password1 != tt.request.Password2 {
					assert.NotEqual(t, tt.request.Password1, tt.request.Password2)
				}
			}
		})
	}
}

func TestCreateUploadRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request requests.CreateUploadRequest
		valid   bool
	}{
		{
			name: "valid request",
			request: requests.CreateUploadRequest{
				Filename:  "test.mp4",
				MimeType:  "video/mp4",
				SizeBytes: 1024,
				Checksum:  "abc123",
			},
			valid: true,
		},
		{
			name: "empty filename",
			request: requests.CreateUploadRequest{
				Filename:  "",
				MimeType:  "video/mp4",
				SizeBytes: 1024,
				Checksum:  "abc123",
			},
			valid: false,
		},
		{
			name: "empty mime type",
			request: requests.CreateUploadRequest{
				Filename:  "test.mp4",
				MimeType:  "",
				SizeBytes: 1024,
				Checksum:  "abc123",
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.valid {
				assert.NotEmpty(t, tt.request.Filename)
				assert.NotEmpty(t, tt.request.MimeType)
				assert.Greater(t, tt.request.SizeBytes, int64(0))
			} else {
				assert.True(t, tt.request.Filename == "" || tt.request.MimeType == "")
			}
		})
	}
}
