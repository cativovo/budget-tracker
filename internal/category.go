package internal

import (
	"context"
	"time"
)

type Category struct {
	ID        string
	Name      string
	Color     string
	Icon      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CategoryRepository interface {
	ListCategories(ctx context.Context, lo ListOptions) ([]Category, error)
	CreateCategory(ctx context.Context, c Category) (Category, error)
	UpdateCategory(ctx context.Context, c Category) (Category, error)
}
