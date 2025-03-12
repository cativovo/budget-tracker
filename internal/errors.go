package internal

import (
	"errors"
	"fmt"
)

type ErrorCode string

const (
	ErrorCodeInvalid  ErrorCode = "invalid"
	ErrorCodeNotFound ErrorCode = "not_found"
	ErrorCodeConflict ErrorCode = "conflict"
	ErrorCodeInternal ErrorCode = "internal"
)

type Error struct {
	code    ErrorCode
	message string
}

func NewError(e ErrorCode, m string) *Error {
	return &Error{
		code:    e,
		message: m,
	}
}

func NewErrorf(e ErrorCode, format string, args ...any) *Error {
	return &Error{
		code:    ErrorCodeInvalid,
		message: fmt.Sprintf(format, args...),
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("internal error: code=%s, message=%s", e.code, e.message)
}

func GetErrorMessage(err error) string {
	var e *Error
	if errors.As(err, &e) {
		return e.message
	}
	return err.Error()
}

func GetErrorCode(err error) ErrorCode {
	var e *Error
	if errors.As(err, &e) {
		return e.code
	}
	return ErrorCodeInternal
}
