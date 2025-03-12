package sqlite_test

import (
	"os"
	"testing"

	"github.com/cativovo/budget-tracker/internal/sqlite"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type dbHelper struct {
	db     *sqlite.DB
	dbPath string
	t      *testing.T
}

var logger = zap.NewNop().Sugar()

func newDBHelper(t *testing.T, dbPath string) *dbHelper {
	t.Helper()

	db, err := sqlite.NewDB(dbPath)
	assert.Nil(t, err)

	err = db.Migrate(logger)
	assert.Nil(t, err)

	return &dbHelper{
		db:     db,
		t:      t,
		dbPath: dbPath,
	}
}

func (d *dbHelper) clean() {
	d.t.Helper()
	d.db.Close()

	err := os.RemoveAll(d.dbPath)
	assert.Nil(d.t, err)
}
