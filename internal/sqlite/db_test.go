package sqlite_test

import (
	"context"
	"os"
	"testing"

	"github.com/cativovo/budget-tracker/internal/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestBuildPragmaQuery(t *testing.T) {
	pragmas := []string{
		"busy_timeout(10000)",
		"journal_mode(WAL)",
		"journal_size_limit(200000000)",
		"synchronous(NORMAL)",
		"foreign_keys(ON)",
		"temp_store(MEMORY)",
		"cache_size(-16000)",
	}
	// https://github.com/pocketbase/pocketbase/blob/391287451729ac2f62a4ca596bb923042b76c213/core/db_connect.go#L14
	want := "?_pragma=busy_timeout(10000)&_pragma=journal_mode(WAL)&_pragma=journal_size_limit(200000000)&_pragma=synchronous(NORMAL)&_pragma=foreign_keys(ON)&_pragma=temp_store(MEMORY)&_pragma=cache_size(-16000)"
	got := sqlite.BuildPragmaQuery(pragmas)

	assert.Equal(t, want, got)
}

func TestDB(t *testing.T) {
	_, c := MustConnectDB(t, "db.db")
	defer c()
}

func MustConnectDB(t *testing.T, dbPath string) (*sqlite.DB, func()) {
	t.Helper()

	tmp, err := os.CreateTemp("", dbPath)
	require.Nil(t, err)

	ctx := context.Background()
	db, err := sqlite.NewDB(ctx, tmp.Name())
	require.Nil(t, err)

	err = db.Migrate(ctx, zap.NewNop().Sugar())
	require.Nil(t, err)

	c := func() {
		err := os.Remove(tmp.Name())
		require.Nil(nil, err)
	}

	return db, c
}
