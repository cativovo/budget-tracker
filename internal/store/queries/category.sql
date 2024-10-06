-- name: ListCategories :many
SELECT * FROM category ORDER BY name;

-- name: CreateCategory :one
INSERT INTO category (name, icon, color_hex)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateCategory :one
UPDATE category
SET name=$1, icon=$2, color_hex=$3
WHERE id=$4
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM category WHERE id=$1;
