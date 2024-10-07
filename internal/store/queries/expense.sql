-- name: ListExpenses :many
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
WHERE expense.account_id=$1 AND expense.date BETWEEN @start_date AND @end_date
ORDER BY expense.date
LIMIT $2
OFFSET $3;

-- name: GetExpense :one
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
WHERE category.account_id=$1 AND expense.id=$2;

-- name: CreateExpense :one
WITH inserted_expense as (
	INSERT INTO expense (name, amount, description, date, category_id, account_id)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING *
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
LEFT JOIN category ON category.id = inserted_expense.category_id;

-- name: UpdateExpense :one
WITH updated_expense as (
	UPDATE expense
	SET name=$1, description=$2, date=$3, category_id=$4
	WHERE expense.id=$5 AND expense.account_id=$6
	RETURNING *
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
LEFT JOIN category ON category.id = updated_expense.category_id;

-- name: DeleteExpense :exec
DELETE FROM expense WHERE account_id=$1 AND id=$2;
