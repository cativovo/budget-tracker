package sqlite_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/cativovo/budget-tracker/internal/category"
	"github.com/cativovo/budget-tracker/internal/sqlite"
	"github.com/cativovo/budget-tracker/internal/user"
	"github.com/stretchr/testify/require"
)

func newTestDB(t *testing.T) (*sqlite.DB, func()) {
	t.Helper()
	db, c := mustConnectDB(t)
	err := db.Migrate(context.Background())
	require.NoError(t, err)
	return db, c
}

func newUser(t *testing.T, db *sqlite.DB) user.User {
	t.Helper()
	return mustCreateUser(t, db, user.CreateUserInput{
		ID:    gofakeit.UUID(),
		Name:  gofakeit.Name(),
		Email: gofakeit.Email(),
	})
}

func newCategory(t *testing.T, db *sqlite.DB, u user.User) category.Category {
	t.Helper()
	return mustCreateCategory(t, db, u, category.CreateCategoryInput{
		Name:  gofakeit.Noun(),
		Color: gofakeit.HexColor(),
		Icon:  gofakeit.Emoji(),
	})
}

func mustCreateCategory(t *testing.T, db *sqlite.DB, u user.User, input category.CreateCategoryInput) category.Category {
	t.Helper()
	cs := sqlite.NewCategoryStore(db)
	c, err := cs.CreateCategory(user.WithContext(context.Background(), u), input)
	require.NoError(t, err)
	return c
}

func mustCreateUser(t *testing.T, db *sqlite.DB, input user.CreateUserInput) user.User {
	t.Helper()
	us := sqlite.NewUserStore(db)
	u, err := us.CreateUser(context.Background(), input)
	require.NoError(t, err)
	return u
}
