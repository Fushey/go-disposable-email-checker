package tempmailchecker

import (
	"errors"
	"fmt"
)

// Common errors returned by the client.
var (
	// ErrAPIKeyRequired is returned when an empty API key is provided.
	ErrAPIKeyRequired = errors.New("API key is required")
	// ErrEmailRequired is returned when an empty email is provided to Check.
	ErrEmailRequired = errors.New("email address is required")
	// ErrDomainRequired is returned when an empty domain is provided to CheckDomain.
	ErrDomainRequired = errors.New("domain is required")
	// ErrInvalidEmail is returned when an invalid email format is provided.
	ErrInvalidEmail = errors.New("invalid email address format")
)

// APIError represents an error returned by the API.
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error (HTTP %d): %s", e.StatusCode, e.Message)
}

// RateLimitError is returned when the API rate limit is exceeded.
type RateLimitError struct {
	Message string
}

func (e *RateLimitError) Error() string {
	return fmt.Sprintf("rate limit exceeded: %s", e.Message)
}

// IsRateLimitError checks if an error is a rate limit error.
func IsRateLimitError(err error) bool {
	var rateLimitErr *RateLimitError
	return errors.As(err, &rateLimitErr)
}

// IsAPIError checks if an error is an API error.
func IsAPIError(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr)
}

