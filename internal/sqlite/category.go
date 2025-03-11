package sqlite

import (
	"context"

	"github.com/cativovo/budget-tracker/internal"
)

type CategoryRepository struct {
	db *DB
}

var _ internal.CategoryRepository = (*CategoryRepository)(nil)

func NewCategoryRepository(db *DB) CategoryRepository {
	return CategoryRepository{
		db: db,
	}
}

func (cr *CategoryRepository) CreateCategory(ctx context.Context, c internal.Category) (internal.Category, error) {
	panic("not yet implemented")
}

func (cr *CategoryRepository) UpdateCategory(ctx context.Context, c internal.Category) (internal.Category, error) {
	panic("not yet implemented")
}

func (cr *CategoryRepository) ListCategories(ctx context.Context, o internal.ListOptions) ([]internal.Category, error) {
	panic("not yet implemented")
}
