package expense

import (
	"context"

	"github.com/cativovo/budget-tracker/internal"
)

type Repository interface {
	ExpenseByID(ctx context.Context, id string) (Expense, error)
	ExpenseGroupByID(ctx context.Context, id string) (ExpenseGroup, error)
	ListExpenseSummaries(ctx context.Context, lo internal.ListOptions) ([]ExpenseSummary, error)
	CreateExpense(ctx context.Context, e CreateExpenseReq) (Expense, error)
	CreateExpenseGroup(ctx context.Context, e CreateExpenseGroupReq) (ExpenseGroup, error)
	UpdateExpense(ctx context.Context, u UpdateExpenseReq) (Expense, error)
	DeleteExpense(ctx context.Context, id string) error
	UpdateExpenseGroup(ctx context.Context, u UpdateExpenseGroupReq) (ExpenseGroup, error)
	DeleteExpenseGroup(ctx context.Context, id string) error
}
