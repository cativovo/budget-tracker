package testutil

import (
	"testing"

	"github.com/cativovo/budget-tracker/internal/apperror"
	"github.com/stretchr/testify/require"
)

// RequireEqualError fails the test if got is nil or if the two errors differ in code or message.
func RequireEqualError(t *testing.T, want, got error) {
	t.Helper()
	require.Error(t, got)
	require.Equal(t, apperror.Code(want), apperror.Code(got))
	require.Equal(t, apperror.Message(want), apperror.Message(got))
}
