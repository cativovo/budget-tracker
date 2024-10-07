package server

import (
	"context"

	"github.com/cativovo/budget-tracker/internal/store"
)

type ExpenseStore interface {
	ListExpenses(ctx context.Context, arg store.ListExpensesParams) ([]store.ListExpensesRow, error)
	CreateExpense(ctx context.Context, arg store.CreateExpenseParams) (store.CreateExpenseRow, error)
}
