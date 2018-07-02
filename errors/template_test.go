package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"net/http"
)

const MESSAGE_FILE = "../config/errors.yaml"

func TestNewAPIError(t *testing.T) {
	defer func() {
		templates = nil
	}()

	assert.Nil(t, LoadMessages(MESSAGE_FILE))

	e := NewAPIError(http.StatusContinue, "xyz", nil)
	assert.Equal(t, http.StatusContinue, e.Status)
	assert.Equal(t, "xyz", e.Message)

	e = NewAPIError(http.StatusNotFound, "NOT_FOUND", nil)
	assert.Equal(t, http.StatusNotFound, e.Status)
	assert.NotEqual(t, "NOT_FOUND", e.Message)
}

func TestLoadMessages(t *testing.T) {
	defer func() {
		templates = nil
	}()

	assert.Nil(t, LoadMessages(MESSAGE_FILE))
	assert.NotNil(t, LoadMessages("xyz"))
}

func Test_replacePlaceholders(t *testing.T) {
	message := replacePlaceholders("abc", nil)
	assert.Equal(t, "abc", message)

	message = replacePlaceholders("abc", Params{"abc": 1})
	assert.Equal(t, "abc", message)

	message = replacePlaceholders("{abc}", Params{"abc": 1})
	assert.Equal(t, "1", message)

	message = replacePlaceholders("123 {abc} xyz {abc} d {xyz}", Params{"abc": 1, "xyz": "t"})
	assert.Equal(t, "123 1 xyz 1 d t", message)
}

func Test_errorTemplate_newAPIError(t *testing.T) {

}
