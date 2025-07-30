package internal

import (
	"errors"
	"fmt"
)

// ErrorCode identifies the error
type ErrorCode string

const (
	// ErrorCodeUnauthenticated indicates an unauthenticated error
	ErrorCodeUnauthenticated ErrorCode = "unauthenticated"
	// ErrorCodeInvalid indicates an invalid request error
	ErrorCodeInvalid ErrorCode = "invalid"
	// ErrorCodeNotFound indicates a resource not found error
	ErrorCodeNotFound ErrorCode = "not_found"
	// ErrorCodeConflict indicates a conflict error
	ErrorCodeConflict ErrorCode = "conflict"
	// ErrorCodeInternal indicates a internal server error
	ErrorCodeInternal ErrorCode = "internal"
)

// Error represents the application-specific error
type Error struct {
	code    ErrorCode
	message string
}

// NewError creates new Error
func NewError(e ErrorCode, m string) *Error {
	return &Error{
		code:    e,
		message: m,
	}
}

// Errorf creates new Error with message
// formatted according to a format specifier
func Errorf(e ErrorCode, format string, args ...any) *Error {
	return &Error{
		code:    e,
		message: fmt.Sprintf(format, args...),
	}
}

// Error implements error interface
func (e *Error) Error() string {
	return fmt.Sprintf("code: %s, message: %s", e.code, e.message)
}

// GetErrorMessage extracts message from Error
func GetErrorMessage(err error) string {
	var e *Error
	if errors.As(err, &e) {
		return e.message
	}
	return err.Error()
}

// GetErrorCode extracts ErrorCode from Error
func GetErrorCode(err error) ErrorCode {
	var e *Error
	if errors.As(err, &e) {
		return e.code
	}
	return ErrorCodeInternal
}
