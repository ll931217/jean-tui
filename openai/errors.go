package openai

import (
	"fmt"
)

// APIError represents an error response from the OpenAI API
type APIError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	return fmt.Sprintf("API error (%s): %s", e.Type, e.Message)
}

// ConfigError represents a configuration error (missing API key, invalid model, etc.)
type ConfigError struct {
	Field   string
	Message string
}

// Error implements the error interface
func (e *ConfigError) Error() string {
	return fmt.Sprintf("configuration error for %s: %s", e.Field, e.Message)
}

// NewConfigError creates a new configuration error
func NewConfigError(field, message string) *ConfigError {
	return &ConfigError{
		Field:   field,
		Message: message,
	}
}

// RequestError represents an error during API request execution
type RequestError struct {
	StatusCode int
	Message    string
}

// Error implements the error interface
func (e *RequestError) Error() string {
	if e.StatusCode > 0 {
		return fmt.Sprintf("request failed with status %d: %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("request failed: %s", e.Message)
}

// NewRequestError creates a new request error
func NewRequestError(statusCode int, message string) *RequestError {
	return &RequestError{
		StatusCode: statusCode,
		Message:    message,
	}
}

// WrapError wraps an error with context about what operation failed
func WrapError(op string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", op, err)
}
