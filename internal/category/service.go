package category

import (
	"context"
	"errors"

	"github.com/cativovo/budget-tracker/internal"
	"github.com/cativovo/budget-tracker/internal/apperror"
)

// Service handles categor operations.
type Service struct {
	categoryStore categoryStore
}

// NewService returns a new Service.
func NewService(cs categoryStore) *Service {
	return &Service{categoryStore: cs}
}

// GetCategoryByID returns a category by ID from the categoryStore.
func (s *Service) GetCategoryByID(ctx context.Context, id string) (Category, error) {
	return s.categoryStore.GetCategoryByID(ctx, id)
}

// CreateCategoryInput represents the input data needed to create a category.
type CreateCategoryInput struct {
	Name  string `json:"name" validate:"required" doc:"Name of category" example:"food"`
	Color string `json:"color" validate:"required,hexcolor" doc:"Color of category in hex" example:"#696969"`
	Icon  string `json:"icon" validate:"required" doc:"Icon of category" example:"food-icon"`
}

func (c CreateCategoryInput) validate() error {
	return internal.ValidateStruct(c)
}

// CreateCategory validates the input and creates a new category in the categoryStore.
func (s *Service) CreateCategory(ctx context.Context, input CreateCategoryInput) (Category, error) {
	if err := input.validate(); err != nil {
		return Category{}, apperror.New(apperror.ErrorCodeInvalid, err.Error())
	}
	return s.categoryStore.CreateCategory(ctx, input)
}

type UpdateCategoryInput struct {
	Name  *string `json:"name,omitempty" validate:"omitnil,min=1" doc:"Name of category" example:"food"`
	Color *string `json:"color,omitempty" validate:"omitnil,hexcolor" doc:"Color of category in hex" example:"#696969"`
	Icon  *string `json:"icon,omitempty" validate:"omitnil,min=1" doc:"Icon of category" example:"food-icon"`
}

func (u UpdateCategoryInput) validate() error {
	if err := internal.ValidateStruct(u); err != nil {
		return err
	}
	if u.Name == nil && u.Color == nil && u.Icon == nil {
		return errors.New("no update fields provided")
	}
	return nil
}

// UpdateCategory validates the input and updates the category with the given ID.
func (s *Service) UpdateCategory(ctx context.Context, id string, input UpdateCategoryInput) (Category, error) {
	if id == "" {
		return Category{}, apperror.New(apperror.ErrorCodeInvalid, "no id provided")
	}

	if err := input.validate(); err != nil {
		return Category{}, apperror.New(apperror.ErrorCodeInvalid, err.Error())
	}

	return s.categoryStore.UpdateCategory(ctx, id, input)
}
