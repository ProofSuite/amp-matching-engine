package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIError_Error(t *testing.T) {
	e := APIError{
		Message: "abc",
	}
	assert.Equal(t, "abc", e.Error())
}

func TestAPIError_StatusCode(t *testing.T) {
	e := APIError{
		Status: 400,
	}
	assert.Equal(t, 400, e.StatusCode())
}
