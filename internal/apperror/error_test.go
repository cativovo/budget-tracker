package apperror_test

import (
	"errors"
	"testing"

	"github.com/cativovo/budget-tracker/internal/apperror"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	wantMessage := "name is required"
	wantCode := apperror.ErrorCodeInvalid
	err := apperror.New(apperror.ErrorCodeInvalid, "name is required")
	assert.Equal(t, wantMessage, err.Error())
	assert.Equal(t, wantMessage, apperror.Message(err))
	assert.Equal(t, wantCode, apperror.Code(err))
}

func TestErrorf(t *testing.T) {
	wantMessage := `username "juanusa" already exists`
	wantCode := apperror.ErrorCodeConflict
	err := apperror.Errorf(apperror.ErrorCodeConflict, "username %q already exists", "juanusa")
	assert.Equal(t, wantMessage, err.Error())
	assert.Equal(t, wantMessage, apperror.Message(err))
	assert.Equal(t, wantCode, apperror.Code(err))
}

func TestMessage(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "application error",
			err:  apperror.New(apperror.ErrorCodeInvalid, "name is required"),
			want: "name is required",
		},
		{
			name: "application error with custom format",
			err:  apperror.Errorf(apperror.ErrorCodeConflict, "value %d is invalid", 70),
			want: "value 70 is invalid",
		},
		{
			name: "non application error",
			err:  errors.New("this is a sensitive error"),
			want: "Internal error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := apperror.Message(tt.err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCode(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want apperror.ErrorCode
	}{
		{
			name: "invalid",
			err:  apperror.New(apperror.ErrorCodeInvalid, "invalid"),
			want: apperror.ErrorCodeInvalid,
		},
		{
			name: "not_found",
			err:  apperror.New(apperror.ErrorCodeNotFound, "not_found"),
			want: apperror.ErrorCodeNotFound,
		},
		{
			name: "conflict",
			err:  apperror.New(apperror.ErrorCodeConflict, "conflict"),
			want: apperror.ErrorCodeConflict,
		},
		{
			name: "unauthorized",
			err:  apperror.New(apperror.ErrorCodeUnauthorized, "unauthorized"),
			want: apperror.ErrorCodeUnauthorized,
		},
		{
			name: "not_implemented",
			err:  apperror.New(apperror.ErrorCodeNotImplemented, "not_implemented"),
			want: apperror.ErrorCodeNotImplemented,
		},
		{
			name: "internal",
			err:  apperror.New(apperror.ErrorCodeInternal, "internal"),
			want: apperror.ErrorCodeInternal,
		},
		{
			name: "non-application error",
			err:  errors.New("this is a senstive error"),
			want: apperror.ErrorCodeInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := apperror.Code(tt.err)
			assert.Equal(t, tt.want, got)
		})
	}
}
