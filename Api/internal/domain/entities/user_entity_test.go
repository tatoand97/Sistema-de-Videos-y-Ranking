package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser_TableName(t *testing.T) {
	user := &User{}
	assert.Equal(t, "users", user.TableName())
}

func TestUser_Creation(t *testing.T) {
	user := &User{
		UserID:       1,
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		FirstName:    "Test",
		LastName:     "User",
		CityID:       123,
	}
	
	assert.Equal(t, 1, user.UserID)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "hashedpassword", user.PasswordHash)
	assert.Equal(t, "Test", user.FirstName)
	assert.Equal(t, "User", user.LastName)
	assert.Equal(t, 123, user.CityID)
}