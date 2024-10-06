-- name: GetAccount :one
SELECT id, name 
FROM account
WHERE id=$1;

-- name: CreateAccount :one
INSERT INTO account (name)
VALUES ($1)
RETURNING id, name;
