package internal

import (
	"context"
	"time"
)

type Category struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	Icon      string    `json:"icon"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CategoryService interface {
	ListCategories(ctx context.Context, f CategoryFilter) ([]Category, int, error)
	GetCategory(ctx context.Context, id string) (Category, error)
	CreateCategory(ctx context.Context, c CategoryCreate) (Category, error)
	UpdateCategory(ctx context.Context, u CategoryUpdate) (Category, error)
	DeleteCategory(ctx context.Context, id string) error
}

type CategoryCreate struct {
	Name  string `json:"name" validate:"required"`
	Color string `json:"color" validate:"hexcolor"`
	Icon  string `json:"icon" validate:"required"`
}

func (c CategoryCreate) Validate() error {
	if err := ValidateStruct(c); err != nil {
		return NewError(ErrorCodeInvalid, err.Error())
	}

	return nil
}

type CategoryUpdate struct {
	ID    string  `json:"id" validate:"required"`
	Name  *string `json:"name" validate:"omitnil,min=1"`
	Color *string `json:"color" validate:"omitnil,hexcolor"`
	Icon  *string `json:"icon" validate:"omitnil,min=1"`
}

func (u CategoryUpdate) Validate() error {
	if u.Name == nil && u.Color == nil && u.Icon == nil {
		return NewError(ErrorCodeInvalid, "no update fields provided")
	}

	if err := ValidateStruct(u); err != nil {
		return NewError(ErrorCodeInvalid, err.Error())
	}

	return nil
}

type CategoryFilter struct {
	Name   *string
	Offset int
	Limit  int
}
