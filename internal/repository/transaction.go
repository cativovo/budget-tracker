package repository

import (
	"context"
	"sort"
	"time"

	"github.com/cativovo/budget-tracker/internal/repository/budget_tracker/public/model"
	"github.com/cativovo/budget-tracker/internal/repository/budget_tracker/public/table"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
)

type TransactionType model.TransactionType

const (
	TransactionTypeExpense TransactionType = TransactionType(model.TransactionType_Expense)
	TransactionTypeIncome  TransactionType = TransactionType(model.TransactionType_Income)
)

type ListTransactionsByDateParams struct {
	StartDate        time.Time
	EndDate          time.Time
	TransactionTypes []TransactionType
	Limit            int
	Offset           int
	AccountID        uuid.UUID
}

type Transaction struct {
	Description     *string         `alias:"transaction.description"`
	Date            *time.Time      `alias:"transaction.date"`
	CreatedAt       *time.Time      `alias:"transaction.created_at"`
	UpdatedAt       *time.Time      `alias:"transaction.updated_at"`
	CategoryID      *uuid.UUID      `alias:"transaction.category_id"`
	Name            string          `alias:"transaction.name"`
	TransactionType TransactionType `alias:"transaction.transaction_type"`
	Amount          int             `alias:"transaction.amount"`
	ID              uuid.UUID       `alias:"transaction.id"`
	AccountID       uuid.UUID       `alias:"transaction.account_id"`
}

type TransactionsByDate struct {
	Date          time.Time
	Transactions  []Transaction
	TotalExpenses int
	TotalIncome   int
}

type ListTransactionsByDate struct {
	TransactionsByDate []TransactionsByDate
	Total              int
}

func (r *Repository) ListTransactionsByDate(ctx context.Context, p ListTransactionsByDateParams) (ListTransactionsByDate, error) {
	accountCondition := table.Transaction.AccountID.EQ(UUID(p.AccountID))
	accountAndDateCondition := accountCondition.
		AND(table.Transaction.Date.BETWEEN(DateT(p.StartDate), DateT(p.EndDate)))

	transactionCondition := Bool(false)
	for _, v := range p.TransactionTypes {
		transactionCondition = transactionCondition.
			OR(table.Transaction.TransactionType.EQ(NewEnumValue(string(v))))
	}

	qDates := SELECT(table.Transaction.Date).
		DISTINCT().
		FROM(table.Transaction).
		WHERE(accountAndDateCondition).
		ORDER_BY(table.Transaction.Date.DESC()).
		LIMIT(int64(p.Limit)).
		OFFSET(int64(p.Offset)).
		AsTable("dates")
	date := table.Transaction.Date.From(qDates)

	qTransactions := SELECT(
		table.Transaction.Description,
		table.Transaction.Date,
		table.Transaction.CreatedAt,
		table.Transaction.UpdatedAt,
		table.Transaction.CategoryID,
		table.Transaction.Name,
		table.Transaction.TransactionType,
		table.Transaction.Amount,
		table.Transaction.ID,
		table.Transaction.AccountID,
	).
		FROM(
			qDates.INNER_JOIN(table.Transaction, table.Transaction.Date.EQ(date)),
		).
		WHERE(accountCondition.AND(transactionCondition))

	qTransactionCount := SELECT(
		COUNT(table.Transaction.ID),
	).
		FROM(table.Transaction).
		WHERE(accountAndDateCondition.AND(transactionCondition))

	type result[T any] struct {
		data T
		err  error
	}

	transactionsChan := make(chan result[[]TransactionsByDate])
	countChan := make(chan result[int])
	db := r.OpenDBFromPool()

	go func() {
		var transactions []Transaction
		err := qTransactions.QueryContext(ctx, db, &transactions)
		if err != nil {
			transactionsChan <- result[[]TransactionsByDate]{
				err: err,
			}
			return
		}

		m := make(map[time.Time]TransactionsByDate)
		for _, v := range transactions {
			date := *v.Date
			transaction, ok := m[date]
			if ok {
				transaction.Transactions = append(transaction.Transactions, v)
			} else {
				transaction.Date = date
				transaction.Transactions = []Transaction{v}
			}

			switch v.TransactionType {
			case TransactionTypeIncome:
				transaction.TotalIncome += v.Amount
			case TransactionTypeExpense:
				transaction.TotalExpenses += v.Amount
			}

			m[date] = transaction
		}

		t := make([]TransactionsByDate, 0, len(m))
		for _, v := range m {
			t = append(t, v)
		}

		sort.Slice(t, func(i, j int) bool {
			return t[i].Date.After(t[j].Date)
		})

		transactionsChan <- result[[]TransactionsByDate]{
			data: t,
		}
	}()

	go func() {
		var dest struct{ Count int }
		err := qTransactionCount.QueryContext(ctx, db, &dest)
		countChan <- result[int]{
			data: dest.Count,
			err:  err,
		}
	}()

	tResult := <-transactionsChan
	if tResult.err != nil {
		return ListTransactionsByDate{}, tResult.err
	}

	cResult := <-countChan
	if cResult.err != nil {
		return ListTransactionsByDate{}, tResult.err
	}

	return ListTransactionsByDate{
		TransactionsByDate: tResult.data,
		Total:              cResult.data,
	}, nil
}

type CreateTransactionParams struct {
	Description     *string
	Date            *time.Time
	CategoryID      *uuid.UUID
	Name            string
	TransactionType TransactionType
	Amount          int
	AccountID       uuid.UUID
}

type CreateTransactionRow struct {
	Description     *string         `alias:"transaction.description"`
	Date            *time.Time      `alias:"transaction.date"`
	CreatedAt       *time.Time      `alias:"transaction.created_at"`
	UpdatedAt       *time.Time      `alias:"transaction.updated_at"`
	CategoryID      *uuid.UUID      `alias:"transaction.category_id"`
	Name            string          `alias:"transaction.name"`
	TransactionType TransactionType `alias:"transaction.transaction_type"`
	Amount          int             `alias:"transaction.amount"`
	ID              uuid.UUID       `alias:"transaction.id"`
	AccountID       uuid.UUID       `alias:"transaction.account_id"`
}

func (r *Repository) CreateTransaction(ctx context.Context, p CreateTransactionParams) (CreateTransactionRow, error) {
	m := model.Transaction{
		TransactionType: model.TransactionType(p.TransactionType),
		Name:            p.Name,
		Amount:          int64(p.Amount),
		Description:     p.Description,
		Date:            p.Date,
		CategoryID:      p.CategoryID,
		AccountID:       p.AccountID,
	}

	q := table.Transaction.
		INSERT(
			table.Transaction.TransactionType,
			table.Transaction.Name,
			table.Transaction.Amount,
			table.Transaction.Description,
			table.Transaction.Date,
			table.Transaction.CategoryID,
			table.Transaction.AccountID,
		).
		MODEL(m).
		RETURNING(
			table.Transaction.Description,
			table.Transaction.Date,
			table.Transaction.CreatedAt,
			table.Transaction.UpdatedAt,
			table.Transaction.CategoryID,
			table.Transaction.Name,
			table.Transaction.Amount,
			table.Transaction.TransactionType,
			table.Transaction.ID,
			table.Transaction.AccountID,
		)

	var row CreateTransactionRow
	err := q.QueryContext(ctx, r.OpenDBFromPool(), &row)
	return row, err
}
