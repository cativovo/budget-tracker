package sqlite_test

import (
	"os"
	"testing"

	"github.com/cativovo/budget-tracker/internal/sqlite"
	"github.com/cativovo/budget-tracker/internal/user"
	"github.com/huandu/go-sqlbuilder"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type dbHelper struct {
	db     *sqlite.DB
	dbPath string
	t      *testing.T
}

var zapLogger = zap.NewNop().Sugar()

func newDBHelper(t *testing.T, dbPath string) *dbHelper {
	t.Helper()

	db, err := sqlite.NewDB(dbPath)
	assert.Nil(t, err)

	err = db.Migrate(zapLogger)
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

func createUsers(t *testing.T, db *sqlite.DB) []user.User {
	t.Helper()

	cuq := []user.CreateUserReq{
		{
			Name:  "Alex Albon",
			ID:    "1",
			Email: "alexalbon@williams.com",
		},
		{
			Name:  "Carlos Sainz Jr.",
			ID:    "2",
			Email: "carlossainzjr@williams.com",
		},
	}

	ib := sqlbuilder.SQLite.NewInsertBuilder()
	ib.InsertInto("user")
	ib.Cols(
		"id",
		"name",
		"email",
	)

	users := make([]user.User, 0, len(cuq))

	for _, v := range cuq {
		ib.Values(
			v.ID,
			v.Name,
			v.Email,
		)
		users = append(users, user.User(v))
	}

	q, args := ib.Build()

	db.ReaderWriter().MustExec(q, args...)

	return users
}

// func createCategorys(t *testing.T, ctx )
