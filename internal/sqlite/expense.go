package sqlite

import (
	"context"

	"github.com/cativovo/budget-tracker/internal"
	"github.com/cativovo/budget-tracker/internal/expense"
)

type ExpenseRepository struct {
	db *DB
}

var _ expense.Repository = (*ExpenseRepository)(nil)

func NewExpenseRepository(db *DB) ExpenseRepository {
	return ExpenseRepository{
		db: db,
	}
}

func (er *ExpenseRepository) ExpenseByID(ctx context.Context, id string) (expense.Expense, error) {
	panic("not yet implemented")
}

func (er *ExpenseRepository) ExpenseGroupByID(ctx context.Context, id string) (expense.ExpenseGroup, error) {
	panic("not yet implemented")
}

func (er *ExpenseRepository) ListExpenseSummaries(ctx context.Context, lo internal.ListOptions) ([]expense.ExpenseSummary, error) {
	panic("not yet implemented")
}

func (er *ExpenseRepository) CreateExpense(ctx context.Context, e expense.CreateExpenseReq) (expense.Expense, error) {
	panic("not yet implemented")
}

func (er *ExpenseRepository) CreateExpenseGroup(ctx context.Context, e expense.CreateExpenseGroupReq) (expense.ExpenseGroup, error) {
	panic("not yet implemented")
}

func (er *ExpenseRepository) UpdateExpense(ctx context.Context, e expense.UpdateExpenseReq) (expense.Expense, error) {
	panic("not yet implemented")
}

func (er *ExpenseRepository) UpdateExpenseGroup(ctx context.Context, e expense.UpdateExpenseGroupReq) (expense.ExpenseGroup, error) {
	panic("not yet implemented")
}
