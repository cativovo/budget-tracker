-- name: ListTransactions :many
SELECT 
	transaction.id as id,
	transaction.transaction_type as transaction_type,
	transaction.name as name,
	transaction.date as date,
	transaction.created_at as created_at,
	transaction.updated_at as updated_at,
	category.name as category,
	category.color_hex as color,
	category.icon as icon
FROM transaction 
LEFT JOIN category ON category.id = transaction.category_id
WHERE transaction.account_id = $1 AND transaction.transaction_type = ANY(@transaction_types::SMALLINT[]) AND transaction.date BETWEEN @start_date AND @end_date
ORDER BY transaction.date
LIMIT $2
OFFSET $3;

-- name: ListTransactionsByDate :many
WITH daily_totals AS (
    SELECT 
        DISTINCT date
    FROM transaction
		WHERE transaction.account_id = $1 AND transaction.date BETWEEN @start_date AND @end_date
		ORDER BY transaction.date
		LIMIT $2
		OFFSET $3
)
SELECT 
    daily_totals.date,
    (
        SELECT 
					COALESCE(
						JSON_AGG(
							JSON_BUILD_OBJECT(
									'id', transaction.id,
									'name', transaction.name,
									'amount', transaction.amount,
									'description', transaction.description,
									'date', transaction.date,
									'transaction_type', transaction.transaction_type
							)
						),
						'[]'
					)::JSON
        FROM transaction
        WHERE transaction.account_id = $1 AND transaction.date = daily_totals.date AND transaction.transaction_type = ANY(@transaction_types::SMALLINT[]) 
    ) AS transactions
FROM daily_totals;

-- name: CountTransactions :one
SELECT 
(
	SELECT
		COUNT(id)
	FROM transaction
	WHERE transaction.account_id = $1 AND transaction.transaction_type = 0 AND transaction.date BETWEEN @start_date AND @end_date
) as expense_count,
(
	SELECT
		COUNT(id)
	FROM transaction
	WHERE transaction.account_id = $1 AND transaction.transaction_type = 1 AND transaction.date BETWEEN @start_date AND @end_date
) as income_count;

-- name: GetTransaction :one
SELECT 
	transaction.id as id,
	transaction.transaction_type as transaction_type,
	transaction.name as name,
	transaction.description as description,
	transaction.date as date,
	transaction.created_at as created_at,
	transaction.updated_at as updated_at,
	category.name as category,
	category.color_hex as color,
	category.icon as icon
FROM transaction 
LEFT JOIN category ON category.id = transaction.category_id
WHERE category.account_id = $1 AND transaction.id = $2;

-- name: CreateTransaction :one
WITH inserted_transaction as (
	INSERT INTO transaction (name, amount, transaction_type, description, date, category_id, account_id)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING *
) SELECT 
	inserted_transaction.id as id,
	inserted_transaction.transaction_type as transaction_type,
	inserted_transaction.name as name,
	inserted_transaction.description as description,
	inserted_transaction.date as date,
	inserted_transaction.created_at as created_at,
	inserted_transaction.updated_at as updated_at,
	category.name as category,
	category.color_hex as color,
	category.icon as icon
FROM inserted_transaction 
LEFT JOIN category ON category.id = inserted_transaction.category_id;

-- name: UpdateTransaction :one
WITH updated_transaction as (
	UPDATE transaction
	SET name = $1, description = $2, date = $3, category_id = $4
	WHERE transaction.account_id = $5 AND transaction.id=$6
	RETURNING *
) 
SELECT 
	updated_transaction.id as id,
	updated_transaction.transaction_type as transaction_type,
	updated_transaction.name as name,
	updated_transaction.description as description,
	updated_transaction.date as date,
	updated_transaction.created_at as created_at,
	updated_transaction.updated_at as updated_at,
	category.name as category,
	category.color_hex as color,
	category.icon as icon
FROM updated_transaction 
LEFT JOIN category ON category.id = updated_transaction.category_id;

-- name: DeleteTransaction :exec
DELETE FROM transaction WHERE account_id = $1 AND id = $2;
