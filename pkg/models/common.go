package models

import (
	"fmt"
	"time"
)

// APIError represents a custom error type for API-related errors
type APIError struct {
	Service string
	Message string
	Code    int
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%s API error (%d): %s", e.Service, e.Code, e.Message)
}

// NewAPIError creates a new API error
func NewAPIError(service, message string, code int) *APIError {
	return &APIError{
		Service: service,
		Message: message,
		Code:    code,
	}
}

// Coordinates represents latitude and longitude
type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// ResponseMetadata contains common response metadata
type ResponseMetadata struct {
	Timestamp time.Time `json:"timestamp"`
	Source    string    `json:"source"`
}
