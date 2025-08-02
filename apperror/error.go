package apperror

import (
	"errors"
	"fmt"
)

// ErrorCode defines a string-based identifier for application-specific error types.
type ErrorCode string

const (
	// ErrorCodeInvalid indicates the input or operation is invalid
	ErrorCodeInvalid ErrorCode = "invalid"
	// ErrorCodeNotFound indicates that the requested resource does not exist
	ErrorCodeNotFound ErrorCode = "not_found"
	// ErrorCodeConflict indicates that a resource conflict occurred (e.g., duplicate data)
	ErrorCodeConflict ErrorCode = "conflict"
	// ErrorCodeUnauthorized indicates that the request lacks valid authentication credentials
	ErrorCodeUnauthorized ErrorCode = "unauthorized"
	// ErrorCodeNotImplemented indicates that the requested functionality is not implemented
	ErrorCodeNotImplemented ErrorCode = "not_implemented"
	// ErrorCodeInternal indicates an internal error occurred
	ErrorCodeInternal ErrorCode = "internal"
)

// Error represents the application-specific error
type Error struct {
	code    ErrorCode
	message string
}

// New creates new Error
func New(e ErrorCode, m string) *Error {
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
	return e.message
}

// Message extracts message from Error
// Non application errors always return "Internal error"
func Message(err error) string {
	var e *Error
	if errors.As(err, &e) {
		return e.message
	}
	return "Internal error"
}

// Code extracts ErrorCode from Error
// Non application errors always return ErrorCodeInternal.
func Code(err error) ErrorCode {
	var e *Error
	if errors.As(err, &e) {
		return e.code
	}
	return ErrorCodeInternal
}
