// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: expense.sql

package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createExpense = `-- name: CreateExpense :one
WITH inserted_expense as (
	INSERT INTO expense (name, amount, description, date, category_id, account_id)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, name, amount, description, date, created_at, updated_at, category_id, account_id
) SELECT 
	inserted_expense.id as id,
	inserted_expense.name as name,
	inserted_expense.description as description,
	inserted_expense.date as date,
	inserted_expense.created_at as created_at,
	inserted_expense.updated_at as updated_at,
	category.name as category,
	category.color_hex as color,
	category.icon as icon
FROM inserted_expense 
LEFT JOIN category ON category.id = inserted_expense.category_id
`

type CreateExpenseParams struct {
	Name        string
	Amount      pgtype.Numeric
	Description pgtype.Text
	Date        pgtype.Date
	CategoryID  pgtype.UUID
	AccountID   pgtype.UUID
}

type CreateExpenseRow struct {
	ID          pgtype.UUID
	Name        string
	Description pgtype.Text
	Date        pgtype.Date
	CreatedAt   pgtype.Timestamp
	UpdatedAt   pgtype.Timestamp
	Category    pgtype.Text
	Color       pgtype.Text
	Icon        pgtype.Text
}

func (q *Queries) CreateExpense(ctx context.Context, arg CreateExpenseParams) (CreateExpenseRow, error) {
	row := q.db.QueryRow(ctx, createExpense,
		arg.Name,
		arg.Amount,
		arg.Description,
		arg.Date,
		arg.CategoryID,
		arg.AccountID,
	)
	var i CreateExpenseRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Date,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Category,
		&i.Color,
		&i.Icon,
	)
	return i, err
}

const deleteExpense = `-- name: DeleteExpense :exec
DELETE FROM expense WHERE account_id=$1 AND id=$2
`

type DeleteExpenseParams struct {
	AccountID pgtype.UUID
	ID        pgtype.UUID
}

func (q *Queries) DeleteExpense(ctx context.Context, arg DeleteExpenseParams) error {
	_, err := q.db.Exec(ctx, deleteExpense, arg.AccountID, arg.ID)
	return err
}

const getExpense = `-- name: GetExpense :one
SELECT 
	expense.id as id,
	expense.name as name,
	expense.description as description,
	expense.date as date,
	expense.created_at as created_at,
	expense.updated_at as updated_at,
	category.name as category,
	category.color_hex as color,
	category.icon as icon
FROM expense 
LEFT JOIN category ON category.id = expense.category_id
WHERE account_id=$1 AND expense.id=$2
`

type GetExpenseParams struct {
	AccountID pgtype.UUID
	ID        pgtype.UUID
}

type GetExpenseRow struct {
	ID          pgtype.UUID
	Name        string
	Description pgtype.Text
	Date        pgtype.Date
	CreatedAt   pgtype.Timestamp
	UpdatedAt   pgtype.Timestamp
	Category    pgtype.Text
	Color       pgtype.Text
	Icon        pgtype.Text
}

func (q *Queries) GetExpense(ctx context.Context, arg GetExpenseParams) (GetExpenseRow, error) {
	row := q.db.QueryRow(ctx, getExpense, arg.AccountID, arg.ID)
	var i GetExpenseRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Date,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Category,
		&i.Color,
		&i.Icon,
	)
	return i, err
}

const listExpenses = `-- name: ListExpenses :many
SELECT 
	expense.id as id,
	expense.name as name,
	expense.date as date,
	expense.created_at as created_at,
	expense.updated_at as updated_at,
	category.name as category,
	category.color_hex as color,
	category.icon as icon
FROM expense 
LEFT JOIN category ON category.id = expense.category_id
WHERE account_id=$1 AND date BETWEEN $2 AND $3
ORDER BY date
`

type ListExpensesParams struct {
	AccountID pgtype.UUID
	StartDate pgtype.Date
	EndDate   pgtype.Date
}

type ListExpensesRow struct {
	ID        pgtype.UUID
	Name      string
	Date      pgtype.Date
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
	Category  pgtype.Text
	Color     pgtype.Text
	Icon      pgtype.Text
}

func (q *Queries) ListExpenses(ctx context.Context, arg ListExpensesParams) ([]ListExpensesRow, error) {
	rows, err := q.db.Query(ctx, listExpenses, arg.AccountID, arg.StartDate, arg.EndDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListExpensesRow
	for rows.Next() {
		var i ListExpensesRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Date,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Category,
			&i.Color,
			&i.Icon,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateExpense = `-- name: UpdateExpense :one
WITH updated_expense as (
	UPDATE expense
	SET name=$1, description=$2, date=$3, category_id=$4
	WHERE expense.id=$5 AND account_id=$6
	RETURNING id, name, amount, description, date, created_at, updated_at, category_id, account_id
) SELECT 
	updated_expense.id as id,
	updated_expense.name as name,
	updated_expense.description as description,
	updated_expense.date as date,
	updated_expense.created_at as created_at,
	updated_expense.updated_at as updated_at,
	category.name as category,
	category.color_hex as color,
	category.icon as icon
FROM updated_expense 
LEFT JOIN category ON category.id = updated_expense.category_id
`

type UpdateExpenseParams struct {
	Name        string
	Description pgtype.Text
	Date        pgtype.Date
	CategoryID  pgtype.UUID
	ID          pgtype.UUID
	AccountID   pgtype.UUID
}

type UpdateExpenseRow struct {
	ID          pgtype.UUID
	Name        string
	Description pgtype.Text
	Date        pgtype.Date
	CreatedAt   pgtype.Timestamp
	UpdatedAt   pgtype.Timestamp
	Category    pgtype.Text
	Color       pgtype.Text
	Icon        pgtype.Text
}

func (q *Queries) UpdateExpense(ctx context.Context, arg UpdateExpenseParams) (UpdateExpenseRow, error) {
	row := q.db.QueryRow(ctx, updateExpense,
		arg.Name,
		arg.Description,
		arg.Date,
		arg.CategoryID,
		arg.ID,
		arg.AccountID,
	)
	var i UpdateExpenseRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Date,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Category,
		&i.Color,
		&i.Icon,
	)
	return i, err
}