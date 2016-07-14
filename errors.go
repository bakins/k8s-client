package client

import "github.com/pkg/errors"

type (
	notFoundError struct {
		message string
	}
)

// NewNotFoundError creates a new error. Generally only used by client implementations
func NewNotFoundError(message string) error {
	return &notFoundError{message: message}
}

func (n *notFoundError) Error() string {
	return "object not found: " + n.message
}

// IsNotFoundError can be used to check if the error was a not found error.
func IsNotFoundError(err error) bool {
	_, ok := errors.Cause(err).(*notFoundError)
	return ok
}
