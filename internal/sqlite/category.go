package sqlite

import (
	"context"

	"github.com/cativovo/budget-tracker/internal"
	"github.com/cativovo/budget-tracker/internal/category"
)

type CategoryRepository struct {
	db *DB
}

var _ category.CategoryRepository = (*CategoryRepository)(nil)

func NewCategoryRepository(db *DB) CategoryRepository {
	return CategoryRepository{
		db: db,
	}
}

func (cr *CategoryRepository) CreateCategory(ctx context.Context, c category.CreateCategoryReq) (category.Category, error) {
	panic("not yet implemented")
}

func (cr *CategoryRepository) UpdateCategory(ctx context.Context, c category.UpdateCategoryReq) (category.Category, error) {
	panic("not yet implemented")
}

func (cr *CategoryRepository) ListCategories(ctx context.Context, o internal.ListOptions) ([]category.Category, error) {
	panic("not yet implemented")
}
