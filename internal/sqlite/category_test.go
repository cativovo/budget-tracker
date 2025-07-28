package sqlite_test

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/cativovo/budget-tracker/internal"
	"github.com/cativovo/budget-tracker/internal/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestCreateCategory(t *testing.T) {
	db, c := mustConnectDB(t)
	defer c()

	cs := sqlite.NewCategoryService(db)
	sharedCtx := internal.ContextWithLogger(context.Background(), zap.NewNop().Sugar())

	user := mustCreateUser(sharedCtx, t, db, internal.UserCreate{
		ID:    "1",
		Name:  "Charles Leclerc",
		Email: "charlesleclerc@ferrari.com",
	})

	testCases := []struct {
		name  string
		user  internal.User
		input internal.CategoryCreate
		err   error
	}{
		{
			name: "ok",
			user: user,
			input: internal.CategoryCreate{
				Name:  "Water",
				Color: "#EF1A2D",
				Icon:  "water",
			},
		},
		{
			name: "no name",
			user: user,
			input: internal.CategoryCreate{
				Color: "#EF1A2D",
				Icon:  "test",
			},
			err: internal.NewError(internal.ErrorCodeInvalid, "name is required"),
		},
		{
			name: "no color",
			user: user,
			input: internal.CategoryCreate{
				Name: "test",
				Icon: "test",
			},
			err: internal.NewError(internal.ErrorCodeInvalid, "color must have a valid hex color value"),
		},
		{
			name: "invalid color",
			user: user,
			input: internal.CategoryCreate{
				Name:  "test",
				Color: "test",
				Icon:  "test",
			},
			err: internal.NewError(internal.ErrorCodeInvalid, "color must have a valid hex color value"),
		},
		{
			name: "No icon",
			user: user,
			input: internal.CategoryCreate{
				Name:  "test",
				Color: "#FFF200",
			},
			err: internal.NewError(internal.ErrorCodeInvalid, "icon is required"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := internal.ContextWithUser(sharedCtx, tc.user)

			if tc.err != nil {
				createdCategory, err := cs.CreateCategory(ctx, tc.input)
				assert.Equal(t, internal.GetErrorCode(tc.err), internal.GetErrorCode(err))
				assert.Equal(t, internal.GetErrorMessage(tc.err), internal.GetErrorMessage(err))

				_, err = cs.GetCategory(ctx, createdCategory.ID)
				assert.Equal(t, internal.ErrorCodeNotFound, internal.GetErrorCode(err))
				assert.Equal(t, "category not found", internal.GetErrorMessage(err))
				return
			}

			createdCategory, err := cs.CreateCategory(ctx, tc.input)
			assert.Nil(t, err)
			assert.NotEqual(t, "", createdCategory.ID)
			assert.Equal(t, tc.input.Name, createdCategory.Name)
			assert.Equal(t, tc.input.Color, createdCategory.Color)
			assert.Equal(t, tc.input.Icon, createdCategory.Icon)
			assert.WithinDuration(t, time.Now(), createdCategory.CreatedAt, wantDelta)
			assert.WithinDuration(t, time.Now(), createdCategory.UpdatedAt, wantDelta)

			gotCategory, err := cs.GetCategory(ctx, createdCategory.ID)
			assert.Nil(t, err)
			assert.Equal(t, createdCategory, gotCategory)
		})
	}

	t.Run("validate all fields", func(t *testing.T) {
		ctx := internal.ContextWithUser(sharedCtx, user)
		createdCategory, err := cs.CreateCategory(ctx, internal.CategoryCreate{})
		assert.Equal(t, internal.ErrorCodeInvalid, internal.GetErrorCode(err))

		gotMsgs := strings.Split(internal.GetErrorMessage(err), ", ")
		wantMsgs := []string{
			"name is required",
			"color must have a valid hex color value",
			"icon is required",
		}
		slices.Sort(gotMsgs)
		slices.Sort(wantMsgs)
		assert.Equal(t, wantMsgs, gotMsgs)

		_, err = cs.GetCategory(ctx, createdCategory.ID)
		assert.Equal(t, internal.ErrorCodeNotFound, internal.GetErrorCode(err))
		assert.Equal(t, "category not found", internal.GetErrorMessage(err))
	})

	t.Run("conflict", func(t *testing.T) {
		input := internal.CategoryCreate{
			Name:  "Utilities",
			Color: "#FFFFFF",
			Icon:  "test",
		}
		ctx := internal.ContextWithUser(sharedCtx, user)

		_, err := cs.CreateCategory(ctx, input)
		assert.Nil(t, err)

		createdCategory, err := cs.CreateCategory(ctx, input)
		assert.Equal(t, internal.ErrorCodeConflict, internal.GetErrorCode(err))
		assert.Equal(t, "Category already exists", internal.GetErrorMessage(err))

		_, err = cs.GetCategory(ctx, createdCategory.ID)
		assert.Equal(t, internal.ErrorCodeNotFound, internal.GetErrorCode(err))
		assert.Equal(t, "category not found", internal.GetErrorMessage(err))
	})

	t.Run("invalid access", func(t *testing.T) {
		ctx1 := internal.ContextWithUser(sharedCtx, user)
		createdCategory := mustCreateCategory(ctx1, t, db, internal.CategoryCreate{
			Name:  "Media",
			Color: "#000000",
			Icon:  "media",
		})

		someUser := mustCreateUser(sharedCtx, t, db, internal.UserCreate{
			ID:    "2",
			Name:  "Lewis Hamilton",
			Email: "lewishamilton@ferrari.com",
		})
		_, err := cs.GetCategory(internal.ContextWithUser(sharedCtx, someUser), createdCategory.ID)
		assert.Equal(t, internal.ErrorCodeNotFound, internal.GetErrorCode(err))
		assert.Equal(t, "category not found", internal.GetErrorMessage(err))
	})
}

func TestUpdateCategory(t *testing.T) {
	db, c := mustConnectDB(t)
	defer c()

	cs := sqlite.NewCategoryService(db)
	sharedCtx := internal.ContextWithLogger(context.Background(), zap.NewNop().Sugar())

	user := mustCreateUser(sharedCtx, t, db, internal.UserCreate{
		ID:    "1",
		Name:  "Charles Leclerc",
		Email: "charlesleclerc@ferrari.com",
	})

	testCases := []struct {
		name           string
		user           internal.User
		categoryCreate internal.CategoryCreate
		categoryUpdate internal.CategoryUpdate
		err            error
	}{
		{
			name: "update name",
			user: user,
			categoryCreate: internal.CategoryCreate{
				Name:  "Category1",
				Color: "#FFFFFF",
				Icon:  "category",
			},
			categoryUpdate: internal.CategoryUpdate{
				Name: ptr("Updated Category1"),
			},
		},
		{
			name: "update color",
			user: user,
			categoryCreate: internal.CategoryCreate{
				Name:  "Category2",
				Color: "#EF1A2D",
				Icon:  "category",
			},
			categoryUpdate: internal.CategoryUpdate{
				Color: ptr("#FFF200"),
			},
		},
		{
			name: "update icon",
			user: user,
			categoryCreate: internal.CategoryCreate{
				Name:  "Category3",
				Color: "#EF1A2D",
				Icon:  "category",
			},
			categoryUpdate: internal.CategoryUpdate{
				Icon: ptr("dog"),
			},
		},
		{
			name: "update all",
			user: user,
			categoryCreate: internal.CategoryCreate{
				Name:  "Category4",
				Color: "#EF1A2D",
				Icon:  "category",
			},
			categoryUpdate: internal.CategoryUpdate{
				Name:  ptr("Updated Category4"),
				Color: ptr("#FFF200"),
				Icon:  ptr("dog"),
			},
		},
		{
			name: "empty name",
			user: user,
			categoryCreate: internal.CategoryCreate{
				Name:  "Category5",
				Color: "#FFF200",
				Icon:  "icon",
			},
			categoryUpdate: internal.CategoryUpdate{
				Name: ptr(""),
			},
			err: internal.NewError(internal.ErrorCodeInvalid, "name is required"),
		},
		{
			name: "invalid color",
			user: user,
			categoryCreate: internal.CategoryCreate{
				Name:  "Category6",
				Color: "#FFF200",
				Icon:  "icon",
			},
			categoryUpdate: internal.CategoryUpdate{
				Color: ptr("test"),
			},
			err: internal.NewError(internal.ErrorCodeInvalid, "color must have a valid hex color value"),
		},
		{
			name: "empty icon",
			user: user,
			categoryCreate: internal.CategoryCreate{
				Name:  "Category7",
				Color: "#FFF200",
				Icon:  "icon",
			},
			categoryUpdate: internal.CategoryUpdate{
				Icon: ptr(""),
			},
			err: internal.NewError(internal.ErrorCodeInvalid, "icon is required"),
		},
		{
			name: "no update",
			user: user,
			categoryCreate: internal.CategoryCreate{
				Name:  "Category8",
				Color: "#FFF200",
				Icon:  "icon",
			},
			categoryUpdate: internal.CategoryUpdate{},
			err:            internal.NewError(internal.ErrorCodeInvalid, "no update fields provided"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := internal.ContextWithUser(sharedCtx, tc.user)
			createdCategory := mustCreateCategory(ctx, t, db, tc.categoryCreate)
			tc.categoryUpdate.ID = createdCategory.ID

			if tc.err != nil {
				_, err := cs.UpdateCategory(ctx, tc.categoryUpdate)
				assert.Equal(t, internal.GetErrorCode(tc.err), internal.GetErrorCode(err))
				assert.Equal(t, internal.GetErrorMessage(tc.err), internal.GetErrorMessage(err))

				gotCategory := mustGetCategory(ctx, t, db, createdCategory.ID)
				assert.Equal(t, createdCategory, gotCategory)
				return
			}

			updatedCategory, err := cs.UpdateCategory(ctx, tc.categoryUpdate)
			assert.Nil(t, err)
			assertUpdatedField(t, tc.categoryUpdate.Name, createdCategory.Name, updatedCategory.Name)
			assertUpdatedField(t, tc.categoryUpdate.Color, createdCategory.Color, updatedCategory.Color)
			assertUpdatedField(t, tc.categoryUpdate.Icon, createdCategory.Icon, updatedCategory.Icon)
			assert.Equal(t, createdCategory.CreatedAt, updatedCategory.CreatedAt)

			// TODO: find a way to test UpdatedAt field without slowing down the tests
			// it's too fast that the UpdatedAt is still the same with CreatedAt

			gotCategory := mustGetCategory(ctx, t, db, createdCategory.ID)
			assert.Equal(t, updatedCategory, gotCategory)
		})
	}

	t.Run("update correct category", func(t *testing.T) {
		ctx := internal.ContextWithUser(sharedCtx, user)
		createdCategory1 := mustCreateCategory(ctx, t, db, internal.CategoryCreate{
			Name:  "Update correct category 1",
			Color: "#FFFFFF",
			Icon:  "icon",
		})
		createdCategory2 := mustCreateCategory(ctx, t, db, internal.CategoryCreate{
			Name:  "Update correct category 2",
			Color: "#FFFFFF",
			Icon:  "icon",
		})

		updatedCategory, err := cs.UpdateCategory(ctx, internal.CategoryUpdate{
			ID:   createdCategory1.ID,
			Name: ptr("Update correct category 1 - ok"),
		})
		assert.Nil(t, err)
		assert.NotEqual(t, createdCategory1, updatedCategory)

		gotCategory := mustGetCategory(ctx, t, db, createdCategory2.ID)
		assert.Equal(t, createdCategory2, gotCategory)
	})

	t.Run("no id", func(t *testing.T) {
		_, err := cs.UpdateCategory(sharedCtx, internal.CategoryUpdate{
			Name:  ptr("Category8"),
			Color: ptr("#FFFFFF"),
			Icon:  ptr("icon"),
		})
		assert.Equal(t, internal.ErrorCodeInvalid, internal.GetErrorCode(err))
		assert.Equal(t, "id is required", internal.GetErrorMessage(err))
	})

	t.Run("validate all fields", func(t *testing.T) {
		_, err := cs.UpdateCategory(sharedCtx, internal.CategoryUpdate{
			ID:    "",
			Name:  ptr(""),
			Color: ptr("test"),
			Icon:  ptr(""),
		})

		gotMsgs := strings.Split(internal.GetErrorMessage(err), ", ")
		wantMsgs := []string{
			"id is required",
			"name is required",
			"color must have a valid hex color value",
			"icon is required",
		}
		slices.Sort(gotMsgs)
		slices.Sort(wantMsgs)
		assert.Equal(t, wantMsgs, gotMsgs)
	})

	t.Run("conflict", func(t *testing.T) {
		ctx := internal.ContextWithUser(sharedCtx, user)
		mustCreateCategory(ctx, t, db, internal.CategoryCreate{
			Name:  "Conflict Category 1",
			Color: "#FFFFFF",
			Icon:  "icon",
		})

		createdCategory := mustCreateCategory(ctx, t, db, internal.CategoryCreate{
			Name:  "Conflict Category 2",
			Color: "#FFFFFF",
			Icon:  "icon",
		})

		_, err := cs.UpdateCategory(ctx, internal.CategoryUpdate{
			ID:   createdCategory.ID,
			Name: ptr("Conflict Category 1"),
		})
		assert.Equal(t, internal.ErrorCodeConflict, internal.GetErrorCode(err))
		assert.Equal(t, "Category already exists", internal.GetErrorMessage(err))
	})

	t.Run("invalid access", func(t *testing.T) {
		ctx1 := internal.ContextWithUser(sharedCtx, user)
		createdCategory := mustCreateCategory(ctx1, t, db, internal.CategoryCreate{
			Name:  "Category of user1",
			Color: "#FFFFFF",
			Icon:  "icon",
		})

		someUser := mustCreateUser(sharedCtx, t, db, internal.UserCreate{
			ID:    "2",
			Name:  "Lewis Hamilton",
			Email: "lewishamilton@ferrari.com",
		})
		ctx2 := internal.ContextWithUser(sharedCtx, someUser)
		_, err := cs.UpdateCategory(ctx2, internal.CategoryUpdate{
			ID:   createdCategory.ID,
			Name: ptr("Category of user2"),
		})
		assert.Equal(t, internal.ErrorCodeNotFound, internal.GetErrorCode(err))
		assert.Equal(t, "category not found", internal.GetErrorMessage(err))

		gotCategory := mustGetCategory(ctx1, t, db, createdCategory.ID)
		assert.Equal(t, createdCategory, gotCategory)
	})

	t.Run("not found", func(t *testing.T) {
		ctx := internal.ContextWithUser(sharedCtx, user)
		_, err := cs.UpdateCategory(ctx, internal.CategoryUpdate{
			ID:   "6969",
			Name: ptr("Test"),
		})
		assert.Equal(t, internal.ErrorCodeNotFound, internal.GetErrorCode(err))
		assert.Equal(t, "category not found", internal.GetErrorMessage(err))
	})
}

func TestDeleteCategory(t *testing.T) {
	db, c := mustConnectDB(t)
	defer c()

	cs := sqlite.NewCategoryService(db)
	sharedCtx := internal.ContextWithLogger(context.Background(), zap.NewNop().Sugar())

	user := mustCreateUser(sharedCtx, t, db, internal.UserCreate{
		ID:    "1",
		Name:  "Charles Leclerc",
		Email: "charlesleclerc@ferrari.com",
	})
	sharedCtx = internal.ContextWithUser(sharedCtx, user)

	t.Run("delete correct category", func(t *testing.T) {
		createdCategory1 := mustCreateCategory(sharedCtx, t, db, internal.CategoryCreate{
			Name:  "Food",
			Color: "#FFFFFF",
			Icon:  "icon1",
		})
		createdCategory2 := mustCreateCategory(sharedCtx, t, db, internal.CategoryCreate{
			Name:  "Drinks",
			Color: "#000000",
			Icon:  "icon2",
		})

		err := cs.DeleteCategory(sharedCtx, createdCategory1.ID)
		assert.Nil(t, err)

		_, err = cs.GetCategory(sharedCtx, createdCategory1.ID)
		assert.Equal(t, internal.ErrorCodeNotFound, internal.GetErrorCode(err))
		assert.Equal(t, "category not found", internal.GetErrorMessage(err))

		gotCategory, err := cs.GetCategory(sharedCtx, createdCategory2.ID)
		assert.Nil(t, err)
		assert.Equal(t, createdCategory2, gotCategory)
	})

	t.Run("invalid access", func(t *testing.T) {
		createdCategory := mustCreateCategory(sharedCtx, t, db, internal.CategoryCreate{
			Name:  "Some category",
			Color: "#FFFFFF",
			Icon:  "icon",
		})

		someUser := mustCreateUser(sharedCtx, t, db, internal.UserCreate{
			ID:    "2",
			Name:  "Lewis Hamilton",
			Email: "lewishamilton@ferrari.com",
		})
		err := cs.DeleteCategory(internal.ContextWithUser(sharedCtx, someUser), createdCategory.ID)
		assert.Nil(t, err)

		gotCategory := mustGetCategory(sharedCtx, t, db, createdCategory.ID)
		assert.Equal(t, createdCategory, gotCategory)
	})

	t.Run("not found", func(t *testing.T) {
		err := cs.DeleteCategory(internal.ContextWithUser(sharedCtx, user), "6969")
		assert.Nil(t, err)
	})
}

func TestListCategories(t *testing.T) {
	db, c := mustConnectDB(t)
	defer c()

	cs := sqlite.NewCategoryService(db)
	sharedCtx := internal.ContextWithLogger(context.Background(), zap.NewNop().Sugar())

	user := mustCreateUser(sharedCtx, t, db, internal.UserCreate{
		ID:    "1",
		Name:  "Charles Leclerc",
		Email: "charlesleclerc@ferrari.com",
	})

	for i := range 20 {
		mustCreateCategory(internal.ContextWithUser(sharedCtx, user), t, db, internal.CategoryCreate{
			Name:  fmt.Sprintf("Category%d", i),
			Color: "#FFF200",
			Icon:  "icon1",
		})
	}

	mustCreateCategory(internal.ContextWithUser(sharedCtx, user), t, db, internal.CategoryCreate{
		Name:  "water",
		Color: "#FFF200",
		Icon:  "icon1",
	})

	t.Run("limit", func(t *testing.T) {
		ctx1 := internal.ContextWithUser(sharedCtx, user)
		categories, total, err := cs.ListCategories(ctx1, internal.CategoryFilter{
			Limit: 20,
		})
		assert.Nil(t, err)
		assert.Len(t, categories, 20)
		assert.Equal(t, 21, total)

		for i, v := range categories {
			assert.Equal(t, v.Name, fmt.Sprintf("Category%d", i))
		}
	})

	t.Run("offset", func(t *testing.T) {
		ctx1 := internal.ContextWithUser(sharedCtx, user)
		offset := 10
		categories, total, err := cs.ListCategories(ctx1, internal.CategoryFilter{
			Offset: offset,
			Limit:  10,
		})

		assert.Nil(t, err)
		assert.Len(t, categories, 10)
		assert.Equal(t, 21, total)

		for i, v := range categories {
			assert.Equal(t, v.Name, fmt.Sprintf("Category%d", i+offset))
		}
	})

	t.Run("filter by 'category1'", func(t *testing.T) {
		ctx1 := internal.ContextWithUser(sharedCtx, user)
		categories, total, err := cs.ListCategories(ctx1, internal.CategoryFilter{
			Name:  ptr("category1"),
			Limit: 10,
		})

		assert.Nil(t, err)
		assert.Len(t, categories, 10)
		assert.Equal(t, 11, total)
		assert.Equal(t, categories[0].Name, "Category1")

		for i, v := range categories[1:] {
			assert.Equal(t, v.Name, fmt.Sprintf("Category1%d", i))
		}
	})

	t.Run("filter by 'cateGory15'", func(t *testing.T) {
		ctx1 := internal.ContextWithUser(sharedCtx, user)
		categories, total, err := cs.ListCategories(ctx1, internal.CategoryFilter{
			Name:  ptr("cateGory15"),
			Limit: 10,
		})

		assert.Nil(t, err)
		assert.Len(t, categories, 1)
		assert.Equal(t, 1, total)
		assert.Equal(t, categories[0].Name, "Category15")
	})

	t.Run("filter by 'water'", func(t *testing.T) {
		ctx1 := internal.ContextWithUser(sharedCtx, user)
		categories, total, err := cs.ListCategories(ctx1, internal.CategoryFilter{
			Name:  ptr("water"),
			Limit: 10,
		})
		assert.Nil(t, err)
		assert.Len(t, categories, 1)
		assert.Equal(t, 1, total)
		assert.Equal(t, categories[0].Name, "water")
	})
}

func mustCreateCategory(ctx context.Context, t *testing.T, db *sqlite.DB, c internal.CategoryCreate) internal.Category {
	t.Helper()

	cs := sqlite.NewCategoryService(db)
	o, err := cs.CreateCategory(ctx, c)
	require.Nil(t, err)
	return o
}

func mustGetCategory(ctx context.Context, t *testing.T, db *sqlite.DB, id string) internal.Category {
	t.Helper()

	cs := sqlite.NewCategoryService(db)
	o, err := cs.GetCategory(ctx, id)
	require.Nil(t, err)
	return o
}
