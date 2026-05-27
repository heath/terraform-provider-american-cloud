package client

import (
	"encoding/json"
	"fmt"
)

// APIError represents a structured error from the American Cloud API.
type APIError struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
	ErrorField string `json:"error"`
}

func (e *APIError) Error() string {
	msg := e.Message
	if msg == "" {
		msg = e.ErrorField
	}
	if msg == "" {
		msg = fmt.Sprintf("HTTP %d", e.StatusCode)
	}
	return fmt.Sprintf("American Cloud API error (status %d): %s", e.StatusCode, msg)
}

// IsNotFound returns true if the error is a 404.
func IsNotFound(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == 404
	}
	return false
}

// parseAPIError attempts to parse a structured error from the response body.
func parseAPIError(statusCode int, body []byte) error {
	apiErr := &APIError{StatusCode: statusCode}

	if len(body) > 0 {
		// Try to parse as JSON error
		_ = json.Unmarshal(body, apiErr)
	}

	if apiErr.Message == "" && apiErr.ErrorField == "" {
		apiErr.Message = string(body)
	}

	return apiErr
}
