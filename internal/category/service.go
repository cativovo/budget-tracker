package category

import (
	"context"

	"github.com/cativovo/budget-tracker/internal"
	"github.com/cativovo/budget-tracker/internal/validator"
)

type Service interface {
	ListCategories(ctx context.Context, lo internal.ListOptions) ([]Category, error)
	CreateCategory(ctx context.Context, c CreateCategoryReq) (Category, error)
	UpdateCategory(ctx context.Context, u UpdateCategoryReq) (Category, error)
	DeleteCategory(ctx context.Context, id string) error
}

type CreateCategoryReq struct {
	Name  string `json:"name" validate:"required"`
	Color string `json:"color" validate:"required,hexcolor"`
	Icon  string `json:"icon" validate:"required"`
}

type UpdateCategoryReq struct {
	ID    string  `json:"id" validate:"required"`
	Name  *string `json:"name"`
	Color *string `json:"color" validate:"hexcolor"`
	Icon  *string `json:"icon"`
}

type service struct {
	r Repository
	v *validator.Validator
}

func NewService(r Repository, v *validator.Validator) Service {
	return &service{
		r: r,
		v: v,
	}
}

func (s *service) ListCategories(ctx context.Context, lo internal.ListOptions) ([]Category, error) {
	return s.r.ListCategories(ctx, lo)
}

func (s *service) CreateCategory(ctx context.Context, c CreateCategoryReq) (Category, error) {
	if err := s.v.Struct(c); err != nil {
		return Category{}, internal.NewError(internal.ErrorCodeInvalid, err.Error())
	}
	return s.r.CreateCategory(ctx, c)
}

func (s *service) UpdateCategory(ctx context.Context, u UpdateCategoryReq) (Category, error) {
	if err := s.v.Struct(u); err != nil {
		return Category{}, internal.NewError(internal.ErrorCodeInvalid, err.Error())
	}

	if u.Name == nil && u.Icon == nil && u.Color == nil {
		return Category{}, internal.NewError(internal.ErrorCodeInvalid, "Must update at least one field")
	}

	return s.r.UpdateCategory(ctx, u)
}

func (s *service) DeleteCategory(ctx context.Context, id string) error {
	return s.r.DeleteCategory(ctx, id)
}
