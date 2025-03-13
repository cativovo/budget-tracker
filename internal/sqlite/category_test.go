package sqlite_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/cativovo/budget-tracker/internal"
	"github.com/cativovo/budget-tracker/internal/category"
	"github.com/cativovo/budget-tracker/internal/logger"
	"github.com/cativovo/budget-tracker/internal/sqlite"
	"github.com/cativovo/budget-tracker/internal/user"
	"github.com/stretchr/testify/assert"
)

func TestCreateFindCategory(t *testing.T) {
	dh := newDBHelper(t, "test_create_category.db")
	defer dh.clean()

	cr := sqlite.NewCategoryRepository(dh.db)
	ctxWithLogger := logger.NewCtxWithLogger(context.Background(), zapLogger)

	users := createUsers(t, ctxWithLogger, dh.db)

	tests := []struct {
		name  string
		user  user.User
		input category.CreateCategoryReq
		want  category.Category
		err   error
	}{
		{
			name: fmt.Sprintf("%s create food category", users[0].Name),
			user: users[0],
			input: category.CreateCategoryReq{
				Name:  "food",
				Color: "#696969",
				Icon:  "food-icon",
			},
			want: category.Category{
				Name:      "food",
				Color:     "#696969",
				Icon:      "food-icon",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		{
			name: fmt.Sprintf("%s create duplicate food category", users[0].Name),
			user: users[0],
			input: category.CreateCategoryReq{
				Name:  "food",
				Color: "#696969",
				Icon:  "food-icon",
			},
			err: internal.NewError(internal.ErrorCodeConflict, "food category already exists"),
		},
		{
			name: fmt.Sprintf("%s create food category", users[1].Name),
			user: users[1],
			input: category.CreateCategoryReq{
				Name:  "food",
				Color: "#696969",
				Icon:  "food-icon",
			},
			want: category.Category{
				Name:      "food",
				Color:     "#696969",
				Icon:      "food-icon",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		{
			name: fmt.Sprintf("%s create duplicate food category", users[1].Name),
			user: users[1],
			input: category.CreateCategoryReq{
				Name:  "food",
				Color: "#696969",
				Icon:  "food-icon",
			},
			err: internal.NewError(internal.ErrorCodeConflict, "food category already exists"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctxWithUser := user.NewCtxWithUser(ctxWithLogger, test.user)
			created, createErr := cr.CreateCategory(ctxWithUser, test.input)

			if test.err != nil {
				wantCode := internal.GetErrorCode(test.err)
				gotCode := internal.GetErrorCode(createErr)
				assert.Equal(t, wantCode, gotCode)

				wantMessage := internal.GetErrorMessage(test.err)
				gotMessage := internal.GetErrorMessage(createErr)
				assert.Equal(t, wantMessage, gotMessage)
				return
			}

			assert.Nil(t, createErr)
			assert.True(t, created.ID != "")
			assert.Equal(t, test.want.Name, created.Name)
			assert.Equal(t, test.want.Color, created.Color)
			assert.Equal(t, test.want.Icon, created.Icon)
			assert.WithinDuration(t, test.want.CreatedAt, created.CreatedAt, time.Second*5)
			assert.WithinDuration(t, test.want.UpdatedAt, created.UpdatedAt, time.Second*5)

			found, findErr := cr.FindCategoryByID(ctxWithUser, created.ID)
			assert.Nil(t, findErr)
			assert.Equal(t, created, found)
		})
	}

	t.Run("category not found", func(t *testing.T) {
		ctxWithUser := user.NewCtxWithUser(ctxWithLogger, users[0])
		_, err := cr.FindCategoryByID(ctxWithUser, "123")
		assert.NotNil(t, err)

		wantCode := internal.ErrorCodeNotFound
		wantMessage := "Category not found"

		gotCode := internal.GetErrorCode(err)
		gotMessage := internal.GetErrorMessage(err)

		assert.Equal(t, wantCode, gotCode)
		assert.Equal(t, wantMessage, gotMessage)
	})
}

func TestUpdateCategory(t *testing.T) {
	dh := newDBHelper(t, "test_update_category.db")
	defer dh.clean()

	cr := sqlite.NewCategoryRepository(dh.db)
	ctxWithLogger := logger.NewCtxWithLogger(context.Background(), zapLogger)

	users := createUsers(t, ctxWithLogger, dh.db)

	ccr := []struct {
		user  user.User
		input category.CreateCategoryReq
	}{
		{
			user: users[0],
			input: category.CreateCategoryReq{
				Name:  "food",
				Color: "#696969",
				Icon:  "food-icon",
			},
		},
		{
			user: users[0],
			input: category.CreateCategoryReq{
				Name:  "rent",
				Color: "#ffffff",
				Icon:  "rent-icon",
			},
		},
		{
			user: users[1],
			input: category.CreateCategoryReq{
				Name:  "food",
				Color: "#696969",
				Icon:  "food-icon",
			},
		},
	}

	categories := make([]struct {
		user     user.User
		category category.Category
	}, 0, len(ccr))

	for _, v := range ccr {
		ctxWithUser := user.NewCtxWithUser(ctxWithLogger, v.user)
		c, err := cr.CreateCategory(ctxWithUser, v.input)
		assert.Nil(t, err)

		categories = append(categories, struct {
			user     user.User
			category category.Category
		}{
			user:     v.user,
			category: c,
		})
	}

	tests := []struct {
		name  string
		user  user.User
		input category.UpdateCategoryReq
		want  category.Category
		err   error
	}{
		{
			name: fmt.Sprintf("%s update category to entertainment", categories[0].user.Name),
			user: categories[0].user,
			input: category.UpdateCategoryReq{
				ID:   categories[0].category.ID,
				Name: "entertainment",
			},
			want: category.Category{
				ID:        categories[0].category.ID,
				Name:      "entertainment",
				Color:     "#696969",
				Icon:      "food-icon",
				CreatedAt: categories[0].category.CreatedAt,
				UpdatedAt: time.Now(),
			},
		},
		{
			name: fmt.Sprintf("%s update color to #ffffff", categories[0].user.Name),
			user: categories[0].user,
			input: category.UpdateCategoryReq{
				ID:    categories[0].category.ID,
				Color: "#ffffff",
			},
			want: category.Category{
				ID:        categories[0].category.ID,
				Name:      "entertainment",
				Color:     "#ffffff",
				Icon:      "food-icon",
				CreatedAt: categories[0].category.CreatedAt,
				UpdatedAt: time.Now(),
			},
		},
		{
			name: fmt.Sprintf("%s update icon to entertainment-icon", categories[0].user.Name),
			user: categories[0].user,
			input: category.UpdateCategoryReq{
				ID:   categories[0].category.ID,
				Icon: "entertainment-icon",
			},
			want: category.Category{
				ID:        categories[0].category.ID,
				Name:      "entertainment",
				Color:     "#ffffff",
				Icon:      "entertainment-icon",
				CreatedAt: categories[0].category.CreatedAt,
				UpdatedAt: time.Now(),
			},
		},
		{
			name: fmt.Sprintf("%s update name, color and icon", categories[0].user.Name),
			user: categories[0].user,
			input: category.UpdateCategoryReq{
				ID:    categories[0].category.ID,
				Name:  "drinks",
				Color: "#d4d4d4",
				Icon:  "drinks-icon",
			},
			want: category.Category{
				ID:        categories[0].category.ID,
				Name:      "drinks",
				Color:     "#d4d4d4",
				Icon:      "drinks-icon",
				CreatedAt: categories[0].category.CreatedAt,
				UpdatedAt: time.Now(),
			},
		},
		{
			name: fmt.Sprintf("%s tries to update the category of Smooth Operator", categories[0].user.Name),
			user: categories[0].user,
			input: category.UpdateCategoryReq{
				ID:    categories[2].category.ID,
				Name:  "food",
				Color: "#696969",
				Icon:  "food-icon",
			},
			err: internal.NewError(internal.ErrorCodeNotFound, "Category not found"),
		},
		{
			name: fmt.Sprintf("%s update drinks to gaming", categories[0].user.Name),
			user: categories[0].user,
			input: category.UpdateCategoryReq{
				ID:   categories[0].category.ID,
				Name: "gaming",
			},
			want: category.Category{
				ID:        categories[0].category.ID,
				Name:      "gaming",
				Color:     "#d4d4d4",
				Icon:      "drinks-icon",
				CreatedAt: categories[0].category.CreatedAt,
				UpdatedAt: time.Now(),
			},
		},
		{
			name: fmt.Sprintf("%s update gaming to gaming", categories[0].user.Name),
			user: categories[0].user,
			input: category.UpdateCategoryReq{
				ID:   categories[0].category.ID,
				Name: "gaming",
			},
			want: category.Category{
				ID:        categories[0].category.ID,
				Name:      "gaming",
				Color:     "#d4d4d4",
				Icon:      "drinks-icon",
				CreatedAt: categories[0].category.CreatedAt,
				UpdatedAt: time.Now(),
			},
		},
		{
			name: fmt.Sprintf("%s update gaming to existing name", categories[0].user.Name),
			user: categories[0].user,
			input: category.UpdateCategoryReq{
				ID:   categories[0].category.ID,
				Name: "rent",
			},
			err: internal.NewError(internal.ErrorCodeConflict, "rent category already exists"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctxWithUser := user.NewCtxWithUser(ctxWithLogger, test.user)
			updated, updateErr := cr.UpdateCategory(ctxWithUser, test.input)

			if test.err != nil {
				wantCode := internal.GetErrorCode(test.err)
				gotCode := internal.GetErrorCode(updateErr)
				assert.Equal(t, wantCode, gotCode)

				wantMessage := internal.GetErrorMessage(test.err)
				gotMessage := internal.GetErrorMessage(updateErr)
				assert.Equal(t, wantMessage, gotMessage)
				return
			}

			assert.Nil(t, updateErr)
			assert.Equal(t, test.want.ID, updated.ID)
			assert.Equal(t, test.want.Name, updated.Name)
			assert.Equal(t, test.want.Color, updated.Color)
			assert.Equal(t, test.want.Icon, updated.Icon)
			assert.WithinDuration(t, test.want.CreatedAt, updated.CreatedAt, 0)
			assert.WithinDuration(t, test.want.UpdatedAt, updated.UpdatedAt, time.Second*5)

			found, findErr := cr.FindCategoryByID(ctxWithUser, test.input.ID)
			assert.Nil(t, findErr)
			assert.Equal(t, updated, found)
		})
	}
}
