package server

import (
	"github.com/cativovo/budget-tracker/internal/apperror"
	"github.com/danielgtaylor/huma/v2"
)

func toHumaStatusError(err error) huma.StatusError {
	code := apperror.Code(err)
	msg := apperror.Message(err)
	switch code {
	case apperror.ErrorCodeInvalid:
		return huma.Error400BadRequest(msg)
	case apperror.ErrorCodeUnauthorized:
		return huma.Error401Unauthorized(msg)
	case apperror.ErrorCodeNotFound:
		return huma.Error404NotFound(msg)
	case apperror.ErrorCodeConflict:
		return huma.Error409Conflict(msg)
	case apperror.ErrorCodeInternal:
		return huma.Error500InternalServerError(msg)
	case apperror.ErrorCodeNotImplemented:
		return huma.Error501NotImplemented(msg)
	default:
		return huma.Error500InternalServerError(msg)
	}
}
