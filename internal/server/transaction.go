package server

import (
	"context"

	"github.com/cativovo/budget-tracker/internal/store"
)

type TransactionStore interface {
	ListTransactionsWithCount(ctx context.Context, arg store.ListTransactionsParams) (store.ListTransactionsWithCountRow, error)
	CreateTransaction(ctx context.Context, arg store.CreateTransactionParams) (store.CreateTransactionRow, error)
}
