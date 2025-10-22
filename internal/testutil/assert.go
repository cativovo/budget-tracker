package testutil

import (
	"testing"
	"time"

	"github.com/cativovo/budget-tracker/internal/apperror"
	"github.com/stretchr/testify/assert"
)

var delta = time.Second * 2

// WithinDuration checks if the difference between two times is within a delta.
func WithinDuration(want, got time.Time) bool {
	dt := want.Sub(got)
	return dt >= -delta && dt <= delta
}

// AssertWithinDuration asserts that the difference between two times is within a delta.
func AssertWithinDuration(t *testing.T, want, got time.Time) bool {
	t.Helper()
	return assert.WithinDuration(t, want, got, delta)
}

// AssertEqualError asserts that that an error occurred and that both errors have the same code and message.
func AssertEqualError(t *testing.T, want, got error) bool {
	t.Helper()
	if !assert.Error(t, got) {
		return false
	}
	return assert.Equal(t, apperror.Code(want), apperror.Code(got)) &&
		assert.Equal(t, apperror.Message(want), apperror.Message(got))
}
