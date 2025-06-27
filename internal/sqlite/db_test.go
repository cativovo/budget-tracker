package sqlite_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/cativovo/budget-tracker/internal/sqlite"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestReader(t *testing.T) {
	db, c := mustConnectDB(t)
	defer c()

	assertPragmas(t, db.Reader)

	_, err := db.Reader.Exec("INSERT INTO user (id, name, email) VALUES ('1', 'test', 'test@test.com')")
	assert.Error(t, err)
}

func TestWriter(t *testing.T) {
	db, c := mustConnectDB(t)
	defer c()

	assertPragmas(t, db.Writer)

	_, err := db.Writer.Exec("INSERT INTO user (id, name, email) VALUES ('1', 'test', 'test@test.com')")
	assert.Nil(t, err)
}

func mustConnectDB(t *testing.T) (*sqlite.DB, func()) {
	t.Helper()

	tmp, err := os.CreateTemp("", "")
	require.Nil(t, err)

	ctx := context.Background()
	db, err := sqlite.NewDB(ctx, tmp.Name())
	require.Nil(t, err)

	err = db.Migrate(ctx, zap.NewNop().Sugar())
	require.Nil(t, err)

	c := func() {
		n := tmp.Name()

		err := os.Remove(n)
		require.Nil(nil, err)

		err = os.Remove(fmt.Sprintf("%s-shm", n))
		require.Nil(nil, err)

		err = os.Remove(fmt.Sprintf("%s-wal", n))
		require.Nil(nil, err)
	}

	return db, c
}

func assertPragmas(t *testing.T, db *sqlx.DB) {
	t.Helper()

	wantMap := map[string]string{
		"busy_timeout":       "10000",     // returns integer as string
		"journal_mode":       "wal",       // returns lowercase string
		"journal_size_limit": "200000000", // returns integer as string
		"synchronous":        "1",         // returns 1 (NORMAL)
		"foreign_keys":       "1",         // returns 1 (ON)
		"temp_store":         "2",         // returns 2 (MEMORY)
		"cache_size":         "-16000",    // returns integer as string
	}

	for k, want := range wantMap {
		var dst string
		if err := db.Get(&dst, fmt.Sprintf("PRAGMA %s", k)); err != nil {
			t.Error(err)
		}

		assert.Equal(t, want, dst)
	}
}
