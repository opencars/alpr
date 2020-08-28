package handler

import (
	"net/http"
)

var (
	// ErrRequiredImageURL returned, if image_url paramter does not present in request.
	ErrRequiredImageURL = NewError(http.StatusBadRequest, "request.image_url_required")

	// ErrInvalidImageURL returned, if image_url is not a valid URL.
	ErrInvalidImageURL = NewError(http.StatusBadRequest, "request.image_url_invalid")
)

// Error represents a handler error. It provides methods for a HTTP status
// code and embeds the built-in error interface.
type Error interface {
	error
	Status() int
}

// StatusError represents error with http status code.
type StatusError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error allows StatusError to satisfy the error interface.
func (se StatusError) Error() string {
	return se.Message
}

// Status returns our HTTP status code.
func (se StatusError) Status() int {
	return se.Code
}

// NewError creates new error instance.
func NewError(code int, message string) Error {
	return &StatusError{
		Code:    code,
		Message: message,
	}
}
