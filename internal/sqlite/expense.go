package sqlite

import (
	"context"

	"github.com/cativovo/budget-tracker/internal"
)

type ExpenseRepository struct {
	db *DB
}

var _ internal.ExpenseRepository = (*ExpenseRepository)(nil)

func NewExpenseRepository(db *DB) ExpenseRepository {
	return ExpenseRepository{
		db: db,
	}
}

func (er *ExpenseRepository) ListExpenseSummaries(ctx context.Context, lo internal.ListOptions) ([]internal.ExpenseSummary, error) {
	panic("not yet implemented")
}

func (er *ExpenseRepository) CreateExpense(ctx context.Context, e internal.Expense) (internal.Expense, error) {
	panic("not yet implemented")
}

func (er *ExpenseRepository) UpdateExpense(ctx context.Context, e internal.Expense) (internal.Expense, error) {
	panic("not yet implemented")
}
