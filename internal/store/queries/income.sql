-- name: ListIncomes :many
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
WHERE account_id=$1 AND date BETWEEN @start_date AND @end_date
ORDER BY date;

-- name: GetIncome :one
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
WHERE account_id=$1 AND income.id=$2;

-- name: CreateIncome :one
WITH inserted_income as (
	INSERT INTO income (name, amount, description, date, category_id, account_id)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING *
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
LEFT JOIN category ON category.id = inserted_income.category_id;

-- name: UpdateIncome :one
WITH updated_income as (
	UPDATE income
	SET name=$1, description=$2, date=$3, category_id=$4
	WHERE income.id=$5 AND account_id=$6
	RETURNING *
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
LEFT JOIN category ON category.id = updated_income.category_id;

-- name: DeleteIncome :exec
DELETE FROM income WHERE account_id=$1 AND id=$2;
