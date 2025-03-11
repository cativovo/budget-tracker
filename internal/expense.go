package internal

import (
	"context"
	"time"
)

type Expense struct {
	ID        string
	Name      string
	Amount    int64
	Date      time.Time
	Note      string
	Category  Category
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ExpenseGroup struct {
	ID        string
	Name      string
	Expenses  []Expense
	Note      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ExpenseSummary struct {
	ID      string
	Name    string
	Amount  int64
	Date    time.Time
	IsGroup bool
}

type ExpenseRepository interface {
	ListExpenseSummaries(ctx context.Context, lo ListOptions) ([]ExpenseSummary, error)
	CreateExpense(ctx context.Context, e Expense) (Expense, error)
	UpdateExpense(ctx context.Context, e Expense) (Expense, error)
}
