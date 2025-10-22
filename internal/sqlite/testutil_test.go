package sqlite_test

import (
	"context"
	"testing"

	"github.com/cativovo/budget-tracker/internal/sqlite"
	"github.com/stretchr/testify/require"
)

func newTestDB(t *testing.T) (*sqlite.DB, func()) {
	t.Helper()
	db, c := mustConnectDB(t)
	err := db.Migrate(context.Background())
	require.NoError(t, err)
	return db, c
}
