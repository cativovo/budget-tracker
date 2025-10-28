package sqlite_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/cativovo/budget-tracker/internal/apperror"
	"github.com/cativovo/budget-tracker/internal/category"
	"github.com/cativovo/budget-tracker/internal/sqlite"
	"github.com/cativovo/budget-tracker/internal/testutil"
	"github.com/cativovo/budget-tracker/internal/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCategoryStore_GetCategoryByID(t *testing.T) {
	testutil.DisableLogs()
	db, c := testutil.NewTestDB(t)
	defer c()

	u := newUser(t, db)
	cs := sqlite.NewCategoryStore(db)

	t.Run("ok", func(t *testing.T) {
		want := newCategory(t, db, u)
		got, err := cs.GetCategoryByID(user.WithContext(context.Background(), u), want.ID)
		require.NoError(t, err)
		require.Equal(t, want, got)
	})

	t.Run("not found", func(t *testing.T) {
		got, err := cs.GetCategoryByID(user.WithContext(context.Background(), u), "69")
		testutil.RequireEqualError(t, apperror.New(apperror.ErrorCodeNotFound, "category not found"), err)
		assert.Zero(t, got)
	})

	t.Run("invalid access", func(t *testing.T) {
		want := newCategory(t, db, u)
		anotherU := newUser(t, db)
		got, err := cs.GetCategoryByID(user.WithContext(context.Background(), anotherU), want.ID)
		testutil.RequireEqualError(t, apperror.New(apperror.ErrorCodeNotFound, "category not found"), err)
		assert.Zero(t, got)
	})
}

func TestCategoryStore_CreateCategory(t *testing.T) {
	testutil.DisableLogs()
	db, c := testutil.NewTestDB(t)
	defer c()

	u := newUser(t, db)
	cs := sqlite.NewCategoryStore(db)

	tests := []struct {
		name  string
		input category.CreateCategoryInput
	}{
		{
			name: "ok",
			input: category.CreateCategoryInput{
				Name:  gofakeit.Noun(),
				Color: gofakeit.HexColor(),
				Icon:  gofakeit.Emoji(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := user.WithContext(context.Background(), u)
			got, err := cs.CreateCategory(ctx, tt.input)
			require.NoError(t, err)
			require.NotZero(t, got.ID)
			require.NotZero(t, got.CreatedAt)
			require.NotZero(t, got.UpdatedAt)

			want, err := cs.GetCategoryByID(ctx, got.ID)
			require.NoError(t, err)
			assert.Equal(t, want, got)
		})
	}
}
