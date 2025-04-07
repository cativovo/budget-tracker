package sqlite_test

import (
	"os"
	"testing"
	"time"

	"github.com/cativovo/budget-tracker/internal/category"
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

func createCategories(t *testing.T, db *sqlite.DB, u user.User) []category.Category {
	t.Helper()

	ib := sqlbuilder.SQLite.NewInsertBuilder()
	ib.InsertInto("category")
	ib.Cols(
		"name",
		"color",
		"icon",
		"user_id",
	)

	ccr := []category.CreateCategoryReq{
		{
			Name:  "food",
			Color: "#000000",
			Icon:  "food-icon",
		},
		{
			Name:  "rent",
			Color: "#ffffff",
			Icon:  "rent-icon",
		},
		{
			Name:  "gaming",
			Color: "#696969",
			Icon:  "gaming-icon",
		},
	}

	for _, v := range ccr {
		ib.Values(
			v.Name,
			v.Color,
			v.Icon,
			u.ID,
		)
	}

	q, args := ib.Build()

	db.ReaderWriter().MustExec(q, args...)

	var dst []struct {
		ID        string    `db:"id"`
		Name      string    `db:"name"`
		Color     string    `db:"color"`
		Icon      string    `db:"icon"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}

	sb := sqlbuilder.SQLite.NewSelectBuilder()
	sb.Select(
		"id",
		"name",
		"color",
		"icon",
		"created_at",
		"updated_at",
	)
	sb.From("category")
	sb.Where(
		sb.EQ("user_id", u.ID),
	)

	q, args = sb.Build()

	err := db.ReaderWriter().Select(&dst, q, args...)
	assert.Nil(t, err)

	categories := make([]category.Category, len(dst))
	for i := range dst {
		categories[i] = category.Category(dst[i])
	}
	return categories
}

func toPtr[T any](t *testing.T, v T) *T {
	t.Helper()
	return &v
}
