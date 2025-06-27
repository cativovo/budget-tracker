package sqlite_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var wantDelta = time.Second * 2

func ptr[T any](v T) *T {
	return &v
}

func assertUpdatedField[T any](t *testing.T, f *T, before, after T) {
	t.Helper()

	if f != nil {
		assert.NotEqual(t, before, after)
	} else {
		assert.Equal(t, before, after)
	}
}
