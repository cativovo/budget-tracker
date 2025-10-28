package server

import (
	"context"
	"net/http"

	"github.com/cativovo/budget-tracker/internal/category"
	"github.com/cativovo/budget-tracker/internal/log"
	"github.com/danielgtaylor/huma/v2"
)

func mountCategoryRoutes(api *huma.Group, cs *category.Service) {
	cg := huma.NewGroup(api, "/category")
	huma.Register(
		cg,
		huma.Operation{
			OperationID: "get-category-by-id",
			Method:      http.MethodGet,
			Summary:     "Get category by ID",
			Path:        "/{id}",
		},
		getCategory(cs),
	)
	huma.Register(
		cg,
		huma.Operation{
			OperationID:   "create-category",
			Method:        http.MethodPost,
			Summary:       "Create a category",
			DefaultStatus: http.StatusCreated,
		},
		createCategory(cs),
	)
}

type getCategoryInput struct {
	ID string `path:"id" doc:"ID of the category to get" example:"123"`
}
type getCategoryOutput struct {
	Body category.Category `json:"category" doc:"A category"`
}

func getCategory(cs *category.Service) humaHandler[getCategoryInput, getCategoryOutput] {
	return func(ctx context.Context, input *getCategoryInput) (*getCategoryOutput, error) {
		logger := log.FromContext(ctx)
		logger.Info("Getting category by ID", "id", input.ID)

		c, err := cs.GetCategoryByID(ctx, input.ID)
		if err != nil {
			logger.Error("Failed to get the category", log.ErrAttr(err))
			return nil, toHumaStatusError(err)
		}

		res := new(getCategoryOutput)
		res.Body = c
		return res, nil
	}
}

type createCategoryInput struct {
	Body category.CreateCategoryInput
}
type createCategoryOutput struct {
	Body category.Category `json:"category" doc:"A category"`
}

func createCategory(cs *category.Service) humaHandler[createCategoryInput, createCategoryOutput] {
	return func(ctx context.Context, input *createCategoryInput) (*createCategoryOutput, error) {
		logger := log.FromContext(ctx)
		logger.Info("Creating new category", "input", input)

		c, err := cs.CreateCategory(ctx, input.Body)
		if err != nil {
			logger.Error("Failed to create a category")
			return nil, toHumaStatusError(err)
		}

		res := new(createCategoryOutput)
		res.Body = c
		return res, nil
	}
}
