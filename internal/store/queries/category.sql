-- name: ListCategories :many
SELECT id, name, icon, color_hex FROM category
WHERE account_id=$1
ORDER BY name;

-- name: CreateCategory :one
INSERT INTO category (name, icon, color_hex, account_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateCategory :one
UPDATE category
SET name=$1, icon=$2, color_hex=$3
WHERE account_id=$4 AND id=$5
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM category WHERE account_id=$1 AND id=$2;
