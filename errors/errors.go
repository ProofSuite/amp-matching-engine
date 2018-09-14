package errors

import (
	"net/http"
	"sort"

	"github.com/go-ozzo/ozzo-validation"
)

type validationError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// InternalServerError creates a new API error representing an internal server error (HTTP 500)
func InternalServerError(err error) *APIError {
	return NewHTTPError(http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", Params{"error": err.Error()})
}

// NotFound creates a new API error representing a resource-not-found error (HTTP 404)
func NotFound(resource string) *APIError {
	return NewHTTPError(http.StatusNotFound, "NOT_FOUND", Params{"resource": resource})
}

// Unauthorized creates a new API error representing an authentication failure (HTTP 401)
func Unauthorized(err string) *APIError {
	return NewHTTPError(http.StatusUnauthorized, "UNAUTHORIZED", Params{"error": err})
}

// InvalidData converts a data validation error into an API error (HTTP 400)
func InvalidData(errs validation.Errors) *APIError {
	result := []validationError{}
	fields := []string{}
	for field := range errs {
		fields = append(fields, field)
	}
	sort.Strings(fields)
	for _, field := range fields {
		err := errs[field]
		result = append(result, validationError{
			Field: field,
			Error: err.Error(),
		})
	}

	err := NewHTTPError(http.StatusBadRequest, "INVALID_DATA", nil)
	err.Details = result

	return err
}
