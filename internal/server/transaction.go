package server

import (
	"context"

	"github.com/cativovo/budget-tracker/internal/repository"
)

type TransactionStore interface {
	CreateTransaction(ctx context.Context, arg repository.CreateTransactionParams) (repository.CreateTransactionRow, error)
	ListTransactionsByDate(ctx context.Context, arg repository.ListTransactionsByDateParams) (repository.ListTransactionsByDate, error)
}
