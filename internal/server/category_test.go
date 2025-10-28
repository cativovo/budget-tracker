package server_test

import (
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/cativovo/budget-tracker/internal/category"
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
	}

	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			createdCategory, createResp, createModel := fetchFromServer[category.Category](t, http.MethodPost, "/api/category", nil, tt.createCategoryInput)
			require.Equal(t, tt.wantStatusCode, createResp.StatusCode)
			require.Equal(t, tt.wantErrModel, createModel)

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
}
