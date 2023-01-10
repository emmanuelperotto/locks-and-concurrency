-- name: CreateTransfer :one
INSERT INTO transfer (amount, from_account_id, to_account_id)
VALUES ($1, $2, $3)
RETURNING *;