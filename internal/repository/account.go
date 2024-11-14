package repository

import (
	"context"

	"github.com/cativovo/budget-tracker/internal/repository/budget_tracker/public/model"
	"github.com/cativovo/budget-tracker/internal/repository/budget_tracker/public/table"
	"github.com/google/uuid"
)

type CreateAccountParams struct {
	Name string
}

type CreateAccountRow struct {
	Name string    `alias:"account.name"`
	ID   uuid.UUID `alias:"account.id"`
}

func (r *Repository) CreateAccount(ctx context.Context, p CreateAccountParams) (CreateAccountRow, error) {
	m := model.Account{
		Name: p.Name,
	}
	q := table.Account.
		INSERT(table.Account.Name).
		MODEL(m).
		RETURNING(
			table.Account.ID,
			table.Account.Name,
		)

	var row CreateAccountRow
	err := q.QueryContext(ctx, r.OpenDBFromPool(), &row)
	return row, err
}
