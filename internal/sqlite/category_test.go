package sqlite_test

import (
	"context"
	"fmt"
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
	db, c := newTestDB(t)
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
	db, c := newTestDB(t)
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
			assert.NotZero(t, got.ID)
			assert.NotZero(t, got.CreatedAt)
			assert.NotZero(t, got.UpdatedAt)

			want := mustGetCategory(t, db, u, got.ID)
			assert.Equal(t, want, got)
		})
	}

	t.Run("category name already exists", func(t *testing.T) {
		u := newUser(t, db)
		c := newCategory(t, db, u)
		ctx := user.WithContext(context.Background(), u)
		got, err := cs.CreateCategory(ctx, category.CreateCategoryInput{
			Name:  c.Name,
			Color: gofakeit.HexColor(),
			Icon:  gofakeit.Emoji(),
		})
		testutil.RequireEqualError(t, apperror.New(apperror.ErrorCodeConflict, fmt.Sprintf("%s category already exists", c.Name)), err)
		assert.Zero(t, got)
	})

	t.Run("same category name for 2 users", func(t *testing.T) {
		c := newCategory(t, db, newUser(t, db))
		u := newUser(t, db)
		ctx := user.WithContext(context.Background(), u)
		got, err := cs.CreateCategory(ctx, category.CreateCategoryInput{
			Name:  c.Name,
			Color: gofakeit.HexColor(),
			Icon:  gofakeit.Emoji(),
		})
		require.NoError(t, err)
		assert.NotZero(t, got.ID)
		assert.NotZero(t, got.CreatedAt)
		assert.NotZero(t, got.UpdatedAt)

		want := mustGetCategory(t, db, u, got.ID)
		assert.Equal(t, want, got)
	})
}

func TestCategoryStore_UpdateCategory(t *testing.T) {
	testutil.DisableLogs()
	db, c := newTestDB(t)
	defer c()

	cs := sqlite.NewCategoryStore(db)

	tests := []struct {
		name  string
		input category.UpdateCategoryInput
	}{
		{
			name: "update all",
			input: category.UpdateCategoryInput{
				Name:  testutil.Ptr(gofakeit.Noun()),
				Color: testutil.Ptr(gofakeit.HexColor()),
				Icon:  testutil.Ptr(gofakeit.Emoji()),
			},
		},
		{
			name: "update name",
			input: category.UpdateCategoryInput{
				Name: testutil.Ptr(gofakeit.Noun()),
			},
		},
		{
			name: "update color",
			input: category.UpdateCategoryInput{
				Color: testutil.Ptr(gofakeit.HexColor()),
			},
		},
		{
			name: "update icon",
			input: category.UpdateCategoryInput{
				Icon: testutil.Ptr(gofakeit.Emoji()),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := newUser(t, db)
			ctx := user.WithContext(context.Background(), u)
			c := newCategory(t, db, u)
			got, err := cs.UpdateCategory(ctx, c.ID, tt.input)
			require.NoError(t, err)

			// TODO: find a way to check whether the updated_at field has changed. can't check it properly because sqlite is too fast
			assert.Equal(t, c.CreatedAt, got.CreatedAt)
			testutil.AssertFieldUpdate(t, tt.input.Name, c.Name, got.Name)
			testutil.AssertFieldUpdate(t, tt.input.Color, c.Color, got.Color)
			testutil.AssertFieldUpdate(t, tt.input.Icon, c.Icon, got.Icon)

			want := mustGetCategory(t, db, u, c.ID)
			assert.Equal(t, want, got)
		})
	}

	t.Run("no update", func(t *testing.T) {
		u := newUser(t, db)
		c := newCategory(t, db, u)
		ctx := user.WithContext(context.Background(), u)
		got, err := cs.UpdateCategory(ctx, c.ID, category.UpdateCategoryInput{
			Name:  &c.Name,
			Color: &c.Color,
			Icon:  &c.Icon,
		})
		require.NoError(t, err)
		assert.Equal(t, c.Name, got.Name)
		assert.Equal(t, c.Color, got.Color)
		assert.Equal(t, c.Icon, got.Icon)
	})

	t.Run("invalid access", func(t *testing.T) {
		u := newUser(t, db)
		c := newCategory(t, db, u)
		ctx := user.WithContext(context.Background(), newUser(t, db))
		got, err := cs.UpdateCategory(ctx, c.ID, category.UpdateCategoryInput{
			Icon: testutil.Ptr(gofakeit.Emoji()),
		})
		testutil.AssertEqualError(t, apperror.New(apperror.ErrorCodeNotFound, "category not found"), err)
		assert.Zero(t, got)
	})

	t.Run("category name already exists", func(t *testing.T) {
		u := newUser(t, db)
		c1 := newCategory(t, db, u)
		c2 := newCategory(t, db, u)
		ctx := user.WithContext(context.Background(), u)
		got, err := cs.UpdateCategory(ctx, c2.ID, category.UpdateCategoryInput{
			Name: &c1.Name,
		})
		testutil.RequireEqualError(t, apperror.New(apperror.ErrorCodeConflict, fmt.Sprintf("%s category already exists", c1.Name)), err)
		assert.Zero(t, got)
	})

	t.Run("same category name for 2 users", func(t *testing.T) {
		c1 := newCategory(t, db, newUser(t, db))
		u := newUser(t, db)
		c2 := newCategory(t, db, u)
		ctx := user.WithContext(context.Background(), u)
		got, err := cs.UpdateCategory(ctx, c2.ID, category.UpdateCategoryInput{
			Name: &c1.Name,
		})
		require.NoError(t, err)

		// TODO: find a way to check whether the updated_at field has changed. can't check it properly because sqlite is too fast
		assert.Equal(t, c2.CreatedAt, got.CreatedAt)
		assert.Equal(t, c1.Name, got.Name)

		want := mustGetCategory(t, db, u, c2.ID)
		assert.Equal(t, want, got)
	})
}
