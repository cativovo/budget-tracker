package store

import (
	"context"
	"fmt"
)

type ListTransactionsWithCountRow struct {
	Transactions []ListTransactionsRow
	CountIncome  int
	CountExpense int
	CountTotal   int // total count before limit/offset
}

// workaround for https://github.com/sqlc-dev/sqlc/issues/2761
func (q *Queries) ListTransactionsWithCount(ctx context.Context, params ListTransactionsParams) (ListTransactionsWithCountRow, error) {
	type count struct {
		income  int
		expense int
	}

	transactionsChan := make(chan try[[]ListTransactionsRow])
	transactionsCountChan := make(chan try[count])

	go func() {
		transactions, err := q.ListTransactions(ctx, ListTransactionsParams{
			TransactionTypes: params.TransactionTypes,
			AccountID:        params.AccountID,
			StartDate:        params.StartDate,
			EndDate:          params.EndDate,
			Limit:            params.Limit,
			Offset:           params.Offset,
		})
		transactionsChan <- try[[]ListTransactionsRow]{
			value: transactions,
			err:   err,
		}
	}()

	go func() {
		transactionsCount, err := q.CountTransactions(ctx, CountTransactionsParams{
			AccountID: params.AccountID,
			StartDate: params.StartDate,
			EndDate:   params.EndDate,
		})
		transactionsCountChan <- try[count]{
			value: count{
				income:  int(transactionsCount.IncomeCount),
				expense: int(transactionsCount.ExpenseCount),
			},
			err: err,
		}
	}()

	tryTransactions := <-transactionsChan
	if err := tryTransactions.err; err != nil {
		return ListTransactionsWithCountRow{}, fmt.Errorf("tryTransactions: %w", err)
	}

	tryTransactionsCount := <-transactionsCountChan
	if err := tryTransactionsCount.err; err != nil {
		return ListTransactionsWithCountRow{}, fmt.Errorf("tryTransactionsCount: %w", err)
	}

	return ListTransactionsWithCountRow{
		Transactions: tryTransactions.value,
		CountIncome:  tryTransactionsCount.value.income,
		CountExpense: tryTransactionsCount.value.expense,
		CountTotal:   tryTransactionsCount.value.income + tryTransactionsCount.value.expense,
	}, nil
}
