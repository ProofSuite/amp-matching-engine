package errors

import (
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

type (
	// Params is used to replace placeholders in an error template with the corresponding values.
	Params map[string]interface{}

	errorTemplate struct {
		Message          string `yaml:"message"`
		DeveloperMessage string `yaml:"developer_message"`
	}
)

var templates map[string]errorTemplate

// LoadMessages reads a YAML file containing error templates.
func LoadMessages(file string) error {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	templates = map[string]errorTemplate{}
	return yaml.Unmarshal(bytes, &templates)
}

// NewHTTPError creates a new APIError with the given HTTP status code, error code, and parameters for replacing placeholders in the error template.
// The param can be nil, indicating there is no need for placeholder replacement.
func NewHTTPError(status int, code string, params Params) *APIError {
	err := &APIError{
		Status:    status,
		ErrorCode: code,
		Message:   code,
	}

	if template, ok := templates[code]; ok {
		err.Message = template.getMessage(params)
		err.DeveloperMessage = template.getDeveloperMessage(params)
	}

	return err
}

// getMessage returns the error message by replacing placeholders in the error template with the actual parameters.
func (e errorTemplate) getMessage(params Params) string {
	return replacePlaceholders(e.Message, params)
}

// getDeveloperMessage returns the developer message by replacing placeholders in the error template with the actual parameters.
func (e errorTemplate) getDeveloperMessage(params Params) string {
	return replacePlaceholders(e.DeveloperMessage, params)
}

func replacePlaceholders(message string, params Params) string {
	if len(message) == 0 {
		return ""
	}
	for key, value := range params {
		message = strings.Replace(message, "{"+key+"}", fmt.Sprint(value), -1)
	}
	return message
}
