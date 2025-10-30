package server_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/cativovo/budget-tracker/internal/category"
	"github.com/cativovo/budget-tracker/internal/testutil"
	"github.com/danielgtaylor/huma/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServer_GetCategoryByID(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		want := newCreateCategory(t, nil)
		got, resp, _ := fetchFromServer[category.Category](t, http.MethodGet, "/api/category/"+want.ID, nil, nil)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, want, got)
	})

	t.Run("not found", func(t *testing.T) {
		id := gofakeit.UUID()
		got, resp, errModel := fetchFromServer[category.Category](t, http.MethodGet, "/api/category/"+id, nil, nil)
		require.Equal(t, http.StatusNotFound, resp.StatusCode)
		assert.Zero(t, got)
		wantErrModel := &huma.ErrorModel{
			Title:  "Not Found",
			Status: http.StatusNotFound,
			Detail: "category not found",
		}
		assert.Equal(t, wantErrModel, errModel)
	})
}

// TODO: create test case for same category name for different user once auth is implemented
func TestServer_CreateCategory(t *testing.T) {
	tests := []struct {
		name                string
		wantStatusCode      int
		wantErrModel        *huma.ErrorModel
		createCategoryInput category.CreateCategoryInput
	}{
		{
			name:           "ok",
			wantStatusCode: http.StatusCreated,
			createCategoryInput: category.CreateCategoryInput{
				Name:  gofakeit.Noun(),
				Color: gofakeit.HexColor(),
				Icon:  gofakeit.Emoji(),
			},
		},
		{
			name:           "empty name",
			wantStatusCode: http.StatusBadRequest,
			wantErrModel: &huma.ErrorModel{
				Title:  "Bad Request",
				Status: http.StatusBadRequest,
				Detail: "name is required",
			},
			createCategoryInput: category.CreateCategoryInput{
				Name:  "",
				Color: gofakeit.HexColor(),
				Icon:  gofakeit.Emoji(),
			},
		},
		{
			name:           "empty color",
			wantStatusCode: http.StatusBadRequest,
			wantErrModel: &huma.ErrorModel{
				Title:  "Bad Request",
				Status: http.StatusBadRequest,
				Detail: "color is required",
			},
			createCategoryInput: category.CreateCategoryInput{
				Name:  gofakeit.Noun(),
				Color: "",
				Icon:  gofakeit.Emoji(),
			},
		},
		{
			name:           "empty icon",
			wantStatusCode: http.StatusBadRequest,
			wantErrModel: &huma.ErrorModel{
				Title:  "Bad Request",
				Status: http.StatusBadRequest,
				Detail: "icon is required",
			},
			createCategoryInput: category.CreateCategoryInput{
				Name:  gofakeit.Noun(),
				Color: gofakeit.HexColor(),
				Icon:  "",
			},
		},
		{
			name:           "invalid hex color",
			wantStatusCode: http.StatusBadRequest,
			wantErrModel: &huma.ErrorModel{
				Title:  "Bad Request",
				Status: http.StatusBadRequest,
				Detail: "color must have a valid hex color value",
			},
			createCategoryInput: category.CreateCategoryInput{
				Name:  gofakeit.Noun(),
				Color: "white",
				Icon:  gofakeit.Emoji(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createdCategory, createResp, createErrModel := fetchFromServer[category.Category](t, http.MethodPost, "/api/category", nil, tt.createCategoryInput)
			require.Equal(t, tt.wantStatusCode, createResp.StatusCode)
			require.Equal(t, tt.wantErrModel, createErrModel)

			id := createdCategory.ID
			if id == "" {
				id = gofakeit.UUID()
			}

			gotCategory, getResp, getModel := fetchFromServer[category.Category](t, http.MethodGet, "/api/category/"+id, nil, nil)

			if createResp.StatusCode != http.StatusCreated {
				wantGetModel := &huma.ErrorModel{
					Title:  "Not Found",
					Status: http.StatusNotFound,
					Detail: "category not found",
				}
				require.Equal(t, wantGetModel, getModel)
				assert.Equal(t, getResp.StatusCode, http.StatusNotFound)
				assert.Zero(t, createdCategory)
				assert.Zero(t, gotCategory)
				return
			}

			require.Equal(t, getResp.StatusCode, http.StatusOK)
			require.Nil(t, getModel)
			assert.NotZero(t, createdCategory.ID)
			assert.NotZero(t, createdCategory.CreatedAt)
			assert.NotZero(t, createdCategory.UpdatedAt)
			assert.Equal(t, tt.createCategoryInput.Name, createdCategory.Name)
			assert.Equal(t, tt.createCategoryInput.Color, createdCategory.Color)
			assert.Equal(t, tt.createCategoryInput.Icon, createdCategory.Icon)
			assert.Equal(t, createdCategory, gotCategory)
		})
	}

	t.Run("category name already exists", func(t *testing.T) {
		createdCategory1 := newCreateCategory(t, nil)
		c, resp, errModel := fetchFromServer[category.Category](t, http.MethodPost, "/api/category", nil, category.CreateCategoryInput{
			Name:  createdCategory1.Name,
			Color: gofakeit.HexColor(),
			Icon:  gofakeit.Emoji(),
		})
		require.Equal(t, http.StatusConflict, resp.StatusCode)
		assert.Zero(t, c)
		assert.Equal(t, errModel, &huma.ErrorModel{
			Title:  "Conflict",
			Status: http.StatusConflict,
			Detail: fmt.Sprintf("%s category already exists", createdCategory1.Name),
		})
	})
}

// TODO: once auth is implemented
// - create test case for updating category of another user
// - create test case for same category name for different user
func TestServer_UpdateCategory(t *testing.T) {
	tests := []struct {
		name                string
		wantStatusCode      int
		wantErrModel        *huma.ErrorModel
		updateCategoryInput category.UpdateCategoryInput
	}{
		{
			name:           "update all",
			wantStatusCode: http.StatusOK,
			updateCategoryInput: category.UpdateCategoryInput{
				Name:  testutil.Ptr(gofakeit.Noun()),
				Color: testutil.Ptr(gofakeit.HexColor()),
				Icon:  testutil.Ptr(gofakeit.Emoji()),
			},
		},
		{
			name:           "update name",
			wantStatusCode: http.StatusOK,
			updateCategoryInput: category.UpdateCategoryInput{
				Name: testutil.Ptr(gofakeit.Noun()),
			},
		},
		{
			name:           "update color",
			wantStatusCode: http.StatusOK,
			updateCategoryInput: category.UpdateCategoryInput{
				Color: testutil.Ptr(gofakeit.HexColor()),
			},
		},
		{
			name:           "update icon",
			wantStatusCode: http.StatusOK,
			updateCategoryInput: category.UpdateCategoryInput{
				Name: testutil.Ptr(gofakeit.Emoji()),
			},
		},
		{
			name:           "invalid  color",
			wantStatusCode: http.StatusBadRequest,
			updateCategoryInput: category.UpdateCategoryInput{
				Color: testutil.Ptr("invalid color"),
			},
			wantErrModel: &huma.ErrorModel{
				Title:  "Bad Request",
				Status: http.StatusBadRequest,
				Detail: "color must have a valid hex color value",
			},
		},
		{
			name:           "empty name",
			wantStatusCode: http.StatusBadRequest,
			updateCategoryInput: category.UpdateCategoryInput{
				Name: testutil.Ptr(""),
			},
			wantErrModel: &huma.ErrorModel{
				Title:  "Bad Request",
				Status: http.StatusBadRequest,
				Detail: "name cannot be empty",
			},
		},
		{
			name:           "empty color",
			wantStatusCode: http.StatusBadRequest,
			updateCategoryInput: category.UpdateCategoryInput{
				Color: testutil.Ptr(""),
			},
			wantErrModel: &huma.ErrorModel{
				Title:  "Bad Request",
				Status: http.StatusBadRequest,
				Detail: "color must have a valid hex color value",
			},
		},
		{
			name:           "empty icon",
			wantStatusCode: http.StatusBadRequest,
			updateCategoryInput: category.UpdateCategoryInput{
				Icon: testutil.Ptr(""),
			},
			wantErrModel: &huma.ErrorModel{
				Title:  "Bad Request",
				Status: http.StatusBadRequest,
				Detail: "icon cannot be empty",
			},
		},
		{
			name:                "no update fields",
			wantStatusCode:      http.StatusBadRequest,
			updateCategoryInput: category.UpdateCategoryInput{},
			wantErrModel: &huma.ErrorModel{
				Title:  "Bad Request",
				Status: http.StatusBadRequest,
				Detail: "no update fields provided",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createdCategory := newCreateCategory(t, nil)
			updatedCategory, updateResp, updateErrModel := fetchFromServer[category.Category](t, http.MethodPatch, "/api/category/"+createdCategory.ID, nil, tt.updateCategoryInput)
			require.Equal(t, tt.wantStatusCode, updateResp.StatusCode)
			require.Equal(t, tt.wantErrModel, updateErrModel)

			gotCategory, _, _ := fetchFromServer[category.Category](t, http.MethodGet, "/api/category/"+createdCategory.ID, nil, nil)

			if tt.wantErrModel != nil {
				assert.Zero(t, updatedCategory)
				assert.Equal(t, createdCategory, gotCategory)
				return
			}

			assert.NotZero(t, updatedCategory.ID)
			assert.NotZero(t, updatedCategory.CreatedAt)
			assert.NotZero(t, updatedCategory.UpdatedAt)
			testutil.AssertFieldUpdate(t, tt.updateCategoryInput.Name, createdCategory.Name, updatedCategory.Name)
			testutil.AssertFieldUpdate(t, tt.updateCategoryInput.Color, createdCategory.Color, updatedCategory.Color)
			testutil.AssertFieldUpdate(t, tt.updateCategoryInput.Icon, createdCategory.Icon, updatedCategory.Icon)
			assert.Equal(t, updatedCategory, gotCategory)
		})
	}

	t.Run("no update", func(t *testing.T) {
		createdCategory := newCreateCategory(t, nil)
		c, resp, _ := fetchFromServer[category.Category](t, http.MethodPatch, "/api/category/"+createdCategory.ID, nil, category.UpdateCategoryInput{
			Name:  &createdCategory.Name,
			Icon:  &createdCategory.Icon,
			Color: &createdCategory.Color,
		})

		require.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, createdCategory.Name, c.Name)
		assert.Equal(t, createdCategory.Color, c.Color)
		assert.Equal(t, createdCategory.Icon, c.Icon)
	})

	t.Run("not found", func(t *testing.T) {
		c, resp, errModel := fetchFromServer[category.Category](t, http.MethodPatch, "/api/category/"+gofakeit.UUID(), nil, category.UpdateCategoryInput{
			Name:  testutil.Ptr(gofakeit.Noun()),
			Color: testutil.Ptr(gofakeit.HexColor()),
			Icon:  testutil.Ptr(gofakeit.Emoji()),
		})
		require.Equal(t, http.StatusNotFound, resp.StatusCode)
		assert.Zero(t, c)
		assert.Equal(t, errModel, &huma.ErrorModel{
			Title:  "Not Found",
			Status: http.StatusNotFound,
			Detail: "category not found",
		})
	})

	t.Run("category name already exists", func(t *testing.T) {
		createdCategory1 := newCreateCategory(t, nil)
		createdCategory2 := newCreateCategory(t, nil)
		c, resp, errModel := fetchFromServer[category.Category](t, http.MethodPatch, "/api/category/"+createdCategory2.ID, nil, category.UpdateCategoryInput{
			Name: &createdCategory1.Name,
		})
		require.Equal(t, http.StatusConflict, resp.StatusCode)
		assert.Zero(t, c)
		assert.Equal(t, errModel, &huma.ErrorModel{
			Title:  "Conflict",
			Status: http.StatusConflict,
			Detail: fmt.Sprintf("%s category already exists", createdCategory1.Name),
		})
	})
}
