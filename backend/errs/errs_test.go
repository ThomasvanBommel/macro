package errs

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	err := New(400, "Bad Request", nil)
	assert.Equal(t, 400, err.Code())
	assert.Equal(t, "Bad Request", err.Error())
	assert.Nil(t, err.Unwrap())
}

func TestInternal(t *testing.T) {
	internalErr := errors.New("database connection failed")
	err := Internal(internalErr)
	assert.Equal(t, 500, err.Code())
	assert.Equal(t, "An internal error occurred", err.Error())
	assert.ErrorIs(t, err.Unwrap(), internalErr)
}
