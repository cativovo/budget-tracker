package testutil

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var delta = time.Second * 2

func WithinDuration(want, got time.Time) bool {
	dt := want.Sub(got)
	return dt >= -delta && dt <= delta
}

func AssertWithinDuration(t *testing.T, want, got time.Time) {
	assert.WithinDuration(t, want, got, delta)
}
