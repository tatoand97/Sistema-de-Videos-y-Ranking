package handlers_test

import (
	"testing"

	"api/internal/domain"

	"github.com/stretchr/testify/assert"
)

func TestDomainErrors_Exist(t *testing.T) {
	assert.NotNil(t, domain.ErrNotFound)
	assert.NotNil(t, domain.ErrForbidden)
	assert.NotNil(t, domain.ErrConflict)
	assert.NotNil(t, domain.ErrInvalid)
	assert.NotNil(t, domain.ErrIdempotent)
}

func TestDomainErrors_Messages(t *testing.T) {
	assert.Equal(t, "not found", domain.ErrNotFound.Error())
	assert.Equal(t, "forbidden", domain.ErrForbidden.Error())
	assert.Equal(t, "conflict", domain.ErrConflict.Error())
	assert.Equal(t, "invalid input", domain.ErrInvalid.Error())
	assert.Equal(t, "idempotent", domain.ErrIdempotent.Error())
}
