package category

import "context"

type categoryStore interface {
	GetCategoryByID(ctx context.Context, id string) (Category, error)
	CreateCategory(ctx context.Context, input CreateCategoryInput) (Category, error)
	UpdateCategory(ctx context.Context, id string, input UpdateCategoryInput) (Category, error)
}
