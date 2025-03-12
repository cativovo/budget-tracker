package category

import (
	"context"

	"github.com/cativovo/budget-tracker/internal"
)

type CategoryRepository interface {
	ListCategories(ctx context.Context, lo internal.ListOptions) ([]Category, error)
	CreateCategory(ctx context.Context, c CreateCategoryReq) (Category, error)
	UpdateCategory(ctx context.Context, u UpdateCategoryReq) (Category, error)
}
