package errors

import (
	errs "errors"
	"net/http"
	"testing"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/stretchr/testify/assert"
)

func TestInternalServerError(t *testing.T) {
	assert.Equal(t, http.StatusInternalServerError, InternalServerError(errs.New("")).Status)
}

func TestUnauthorized(t *testing.T) {
	assert.Equal(t, http.StatusUnauthorized, Unauthorized("t").Status)
}

func TestInvalidData(t *testing.T) {
	err := InvalidData(validation.Errors{
		"abc": errs.New("1"),
		"xyz": errs.New("2"),
	})
	assert.Equal(t, http.StatusBadRequest, err.Status)
	assert.NotNil(t, err.Details)
}

func TestNotFound(t *testing.T) {
	assert.Equal(t, http.StatusNotFound, NotFound("abc").Status)
}
