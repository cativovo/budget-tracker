package repository

import (
	"context"

	"github.com/cativovo/budget-tracker/internal/repository/budget_tracker/public/model"
	"github.com/cativovo/budget-tracker/internal/repository/budget_tracker/public/table"
	"github.com/google/uuid"
)

type CreateCategoryParams struct {
	Name      string
	Icon      string
	ColorHex  string
	AccountID uuid.UUID
}

type CreateCategoryRow struct {
	Name     string    `alias:"category.name"`
	Icon     string    `alias:"category.icon"`
	ColorHex string    `alias:"category.color_hex"`
	ID       uuid.UUID `alias:"category.id"`
}

func (r *Repository) CreateCategory(ctx context.Context, p CreateCategoryParams) (CreateCategoryRow, error) {
	m := model.Category{
		Name:      p.Name,
		Icon:      p.Icon,
		ColorHex:  p.ColorHex,
		AccountID: p.AccountID,
	}

	q := table.Category.
		INSERT(
			table.Category.Name,
			table.Category.Icon,
			table.Category.ColorHex,
			table.Category.AccountID,
		).
		MODEL(m).
		RETURNING(
			table.Category.Name,
			table.Category.Icon,
			table.Category.ColorHex,
			table.Category.ID,
		)

	var row CreateCategoryRow
	err := q.QueryContext(ctx, r.OpenDBFromPool(), &row)
	return row, err
}
