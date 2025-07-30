package testutil

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var delta = time.Second * 2

// WithinDuration checks if the difference between two times is within a delta.
func WithinDuration(want, got time.Time) bool {
	dt := want.Sub(got)
	return dt >= -delta && dt <= delta
}

// AssertWithinDuration asserts that the difference between two times is within a delta.
func AssertWithinDuration(t *testing.T, want, got time.Time) {
	assert.WithinDuration(t, want, got, delta)
}
