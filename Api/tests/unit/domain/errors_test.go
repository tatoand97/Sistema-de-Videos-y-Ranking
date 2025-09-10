package domain

import (
	"api/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDomainErrors_Structure(t *testing.T) {
	// Test ErrInvalid
	assert.NotNil(t, domain.ErrInvalid)
	assert.Contains(t, domain.ErrInvalid.Error(), "invalid")
	
	// Test ErrNotFound
	assert.NotNil(t, domain.ErrNotFound)
	assert.Contains(t, domain.ErrNotFound.Error(), "not found")
	
	// Test ErrConflict
	assert.NotNil(t, domain.ErrConflict)
	assert.Contains(t, domain.ErrConflict.Error(), "conflict")
	
	// Test ErrForbidden
	assert.NotNil(t, domain.ErrForbidden)
	assert.Contains(t, domain.ErrForbidden.Error(), "forbidden")
	
	// Test ErrIdempotent
	assert.NotNil(t, domain.ErrIdempotent)
	assert.Contains(t, domain.ErrIdempotent.Error(), "idempotent")
}

func TestDomainErrors_Types(t *testing.T) {
	// Verify that errors are of error type
	var err error
	
	err = domain.ErrInvalid
	assert.Error(t, err)
	
	err = domain.ErrNotFound
	assert.Error(t, err)
	
	err = domain.ErrConflict
	assert.Error(t, err)
	
	err = domain.ErrForbidden
	assert.Error(t, err)
	
	err = domain.ErrIdempotent
	assert.Error(t, err)
}

func TestDomainErrors_Uniqueness(t *testing.T) {
	// Verify that each error is unique
	errors := []error{
		domain.ErrInvalid,
		domain.ErrNotFound,
		domain.ErrConflict,
		domain.ErrForbidden,
		domain.ErrIdempotent,
	}
	
	for i, err1 := range errors {
		for j, err2 := range errors {
			if i != j {
				assert.NotEqual(t, err1.Error(), err2.Error(), "Errors should have unique messages")
			}
		}
	}
}