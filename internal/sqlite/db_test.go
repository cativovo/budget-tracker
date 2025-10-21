package sqlite_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/cativovo/budget-tracker/internal/sqlite"
	"github.com/cativovo/budget-tracker/internal/testutil"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDB(t *testing.T) {
	db, c := mustConnectDB(t)
	defer c()

	assertPragmas(t, db.Reader)
	assertPragmas(t, db.Writer)

	_, err := db.Writer.Exec(`
		CREATE TABLE test(
			key TEXT NOT NULL,
			value TEXT NOT NULL
		);
		`)
	require.NoError(t, err)

	_, err = db.Reader.Exec("INSERT INTO test (key, value) VALUES ('key', 'value')")
	assert.Error(t, err)

	_, err = db.Writer.Exec("INSERT INTO test (key, value) VALUES ('key', 'value')")
	require.NoError(t, err)

	row := db.Reader.QueryRow("SELECT key, value FROM test WHERE key = ?", "key")
	var key string
	var value string
	err = row.Scan(&key, &value)
	require.NoError(t, err)

	assert.Equal(t, "key", key)
	assert.Equal(t, "value", value)
}

func TestMigrate(t *testing.T) {
	testutil.DisableLogs()

	db, c := mustConnectDB(t)
	defer c()

	err := db.Migrate(context.Background())
	assert.NoError(t, err)

	rows, err := db.Reader.Query("SELECT name FROM sqlite_master WHERE type = 'table'")
	require.NoError(t, err)

	var got []string
	for rows.Next() {
		var v string
		err := rows.Scan(&v)
		require.NoError(t, err)
		got = append(got, v)
	}

	want := []string{
		"user",
		"category",
		"expense_group",
		"expense",
	}

	for _, v := range want {
		// check contains instead of equal to ignore unneeded tables
		assert.Contains(t, got, v)
	}
}

func mustConnectDB(t *testing.T) (*sqlite.DB, func()) {
	t.Helper()

	tmp, err := os.CreateTemp("", "")
	assert.NoError(t, err)
	defer func() {
		err := tmp.Close()
		assert.NoError(t, err)
	}()

	ctx := context.Background()
	db, err := sqlite.NewDB(ctx, tmp.Name())
	assert.NoError(t, err)

	c := func() {
		n := tmp.Name()

		err := os.Remove(n)
		assert.NoError(t, err)

		err = os.Remove(fmt.Sprintf("%s-shm", n))
		assert.NoError(t, err)

		err = os.Remove(fmt.Sprintf("%s-wal", n))
		assert.NoError(t, err)
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
		err := db.Get(&dst, fmt.Sprintf("PRAGMA %s", k))
		if assert.NoError(t, err) {
			assert.Equal(t, want, dst)
		}
	}
}
