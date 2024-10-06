// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: income.sql

package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createIncome = `-- name: CreateIncome :one
WITH inserted_income as (
	INSERT INTO income (name, amount, description, date, category_id, account_id)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, name, amount, description, date, created_at, updated_at, category_id, account_id
) SELECT 
	inserted_income.id as id,
	inserted_income.name as name,
	inserted_income.description as description,
	inserted_income.date as date,
	inserted_income.created_at as created_at,
	inserted_income.updated_at as updated_at,
	category.name as category,
	category.color_hex as color,
	category.icon as icon
FROM inserted_income 
LEFT JOIN category ON category.id = inserted_income.category_id
`

type CreateIncomeParams struct {
	Name        string
	Amount      pgtype.Numeric
	Description pgtype.Text
	Date        pgtype.Date
	CategoryID  pgtype.UUID
	AccountID   pgtype.UUID
}

type CreateIncomeRow struct {
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

func (q *Queries) CreateIncome(ctx context.Context, arg CreateIncomeParams) (CreateIncomeRow, error) {
	row := q.db.QueryRow(ctx, createIncome,
		arg.Name,
		arg.Amount,
		arg.Description,
		arg.Date,
		arg.CategoryID,
		arg.AccountID,
	)
	var i CreateIncomeRow
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

const deleteIncome = `-- name: DeleteIncome :exec
DELETE FROM income WHERE account_id=$1 AND id=$2
`

type DeleteIncomeParams struct {
	AccountID pgtype.UUID
	ID        pgtype.UUID
}

func (q *Queries) DeleteIncome(ctx context.Context, arg DeleteIncomeParams) error {
	_, err := q.db.Exec(ctx, deleteIncome, arg.AccountID, arg.ID)
	return err
}

const getIncome = `-- name: GetIncome :one
SELECT 
	income.id as id,
	income.name as name,
	income.description as description,
	income.date as date,
	income.created_at as created_at,
	income.updated_at as updated_at,
	category.name as category,
	category.color_hex as color,
	category.icon as icon
FROM income 
LEFT JOIN category ON category.id = income.category_id
WHERE income.account_id=$1 AND income.id=$2
`

type GetIncomeParams struct {
	AccountID pgtype.UUID
	ID        pgtype.UUID
}

type GetIncomeRow struct {
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

func (q *Queries) GetIncome(ctx context.Context, arg GetIncomeParams) (GetIncomeRow, error) {
	row := q.db.QueryRow(ctx, getIncome, arg.AccountID, arg.ID)
	var i GetIncomeRow
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

const listIncomes = `-- name: ListIncomes :many
SELECT 
	income.id as id,
	income.name as name,
	income.date as date,
	income.created_at as created_at,
	income.updated_at as updated_at,
	category.name as category,
	category.color_hex as color,
	category.icon as icon
FROM income 
LEFT JOIN category ON category.id = income.category_id
WHERE income.account_id=$1 AND income.date BETWEEN $2 AND $3
ORDER BY income.date
`

type ListIncomesParams struct {
	AccountID pgtype.UUID
	StartDate pgtype.Date
	EndDate   pgtype.Date
}

type ListIncomesRow struct {
	ID        pgtype.UUID
	Name      string
	Date      pgtype.Date
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
	Category  pgtype.Text
	Color     pgtype.Text
	Icon      pgtype.Text
}

func (q *Queries) ListIncomes(ctx context.Context, arg ListIncomesParams) ([]ListIncomesRow, error) {
	rows, err := q.db.Query(ctx, listIncomes, arg.AccountID, arg.StartDate, arg.EndDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListIncomesRow
	for rows.Next() {
		var i ListIncomesRow
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

const updateIncome = `-- name: UpdateIncome :one
WITH updated_income as (
	UPDATE income
	SET name=$1, description=$2, date=$3, category_id=$4
	WHERE income.id=$5 AND income.account_id=$6
	RETURNING id, name, amount, description, date, created_at, updated_at, category_id, account_id
) SELECT 
	updated_income.id as id,
	updated_income.name as name,
	updated_income.description as description,
	updated_income.date as date,
	updated_income.created_at as created_at,
	updated_income.updated_at as updated_at,
	category.name as category,
	category.color_hex as color,
	category.icon as icon
FROM updated_income 
LEFT JOIN category ON category.id = updated_income.category_id
`

type UpdateIncomeParams struct {
	Name        string
	Description pgtype.Text
	Date        pgtype.Date
	CategoryID  pgtype.UUID
	ID          pgtype.UUID
	AccountID   pgtype.UUID
}

type UpdateIncomeRow struct {
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

func (q *Queries) UpdateIncome(ctx context.Context, arg UpdateIncomeParams) (UpdateIncomeRow, error) {
	row := q.db.QueryRow(ctx, updateIncome,
		arg.Name,
		arg.Description,
		arg.Date,
		arg.CategoryID,
		arg.ID,
		arg.AccountID,
	)
	var i UpdateIncomeRow
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
