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

-- name: DebitAccount :one
UPDATE account SET balance=balance - sqlc.arg(amount)::numeric, version=version+1
WHERE id=sqlc.arg(id)::integer AND balance >= sqlc.arg(amount)::numeric
RETURNING *;

-- name: CreditAccount :one
UPDATE account SET balance=balance+sqlc.arg(amount)::numeric, version=version+1
WHERE id=sqlc.arg(id)::integer
RETURNING *;