package category

import (
	"context"

	"github.com/cativovo/budget-tracker/internal"
	"github.com/cativovo/budget-tracker/internal/validator"
)

type CategoryService struct {
	cr CategoryRepository
	v  *validator.Validator
}

func NewCategoryService(cr CategoryRepository, v *validator.Validator) *CategoryService {
	return &CategoryService{
		cr: cr,
		v:  v,
	}
}

func (cs *CategoryService) ListCategories(ctx context.Context, lo internal.ListOptions) ([]Category, error) {
	return cs.cr.ListCategories(ctx, lo)
}

type CreateCategoryReq struct {
	Name  string `json:"name" validate:"required"`
	Color string `json:"color" validate:"required,hexcolor"`
	Icon  string `json:"icon" validate:"required"`
}

func (cs *CategoryService) CreateCategory(ctx context.Context, c CreateCategoryReq) (Category, error) {
	if err := cs.v.Struct(c); err != nil {
		return Category{}, err
	}
	return cs.cr.CreateCategory(ctx, c)
}

type UpdateCategoryReq struct {
	Name  string `json:"name"`
	Color string `json:"color" validate:"hexcolor"`
	Icon  string `json:"icon"`
}

func (cs *CategoryService) UpdateCategory(ctx context.Context, u UpdateCategoryReq) (Category, error) {
	if err := cs.v.Struct(u); err != nil {
		return Category{}, err
	}
	return cs.cr.UpdateCategory(ctx, u)
}
