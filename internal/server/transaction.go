package server

import (
	"context"

	"github.com/cativovo/budget-tracker/internal/store"
)

type TransactionStore interface {
	CreateTransaction(ctx context.Context, arg store.CreateTransactionParams) (store.CreateTransactionRow, error)
	ListTransactionsByDate(ctx context.Context, arg store.ListTransactionsByDateParams) ([]byte, error)
}
