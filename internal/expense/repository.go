package expense

import (
	"context"

	"github.com/cativovo/budget-tracker/internal"
)

type Repository interface {
	ListExpenseSummaries(ctx context.Context, lo internal.ListOptions) ([]ExpenseSummary, error)
	CreateExpense(ctx context.Context, e CreateExpenseReq) (Expense, error)
	UpdateExpense(ctx context.Context, u UpdateExpenseReq) (Expense, error)
}
