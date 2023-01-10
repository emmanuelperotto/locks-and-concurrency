-- name: GetAccount :one
SELECT * FROM account
    WHERE id=$1;

-- name: GetAccountForUpdate :one
SELECT * FROM account
WHERE id=$1
FOR UPDATE;

-- name: UpdateAccount :one
UPDATE account SET balance=$1, version=version+1
    WHERE id=$2
RETURNING *;

-- name: OptimisticUpdateAccount :one
UPDATE account SET balance=$1, version=version+1
WHERE id=$2 AND version=$3
RETURNING *;