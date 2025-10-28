package category

import (
	"context"

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
	Color string `json:"color" validate:"required" doc:"Color of category in hex" example:"#696969"`
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
