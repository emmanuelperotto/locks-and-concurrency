// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: accounts.sql

package repository

import (
	"context"
)

const getAccount = `-- name: GetAccount :one
SELECT id, balance, version FROM account
    WHERE id=$1
`

func (q *Queries) GetAccount(ctx context.Context, id int32) (Account, error) {
	row := q.db.QueryRowContext(ctx, getAccount, id)
	var i Account
	err := row.Scan(&i.ID, &i.Balance, &i.Version)
	return i, err
}

const getAccountForUpdate = `-- name: GetAccountForUpdate :one
SELECT id, balance, version FROM account
WHERE id=$1
FOR UPDATE
`

func (q *Queries) GetAccountForUpdate(ctx context.Context, id int32) (Account, error) {
	row := q.db.QueryRowContext(ctx, getAccountForUpdate, id)
	var i Account
	err := row.Scan(&i.ID, &i.Balance, &i.Version)
	return i, err
}

const optimisticUpdateAccount = `-- name: OptimisticUpdateAccount :one
UPDATE account SET balance=$1, version=version+1
WHERE id=$2 AND version=$3
RETURNING id, balance, version
`

type OptimisticUpdateAccountParams struct {
	Balance string
	ID      int32
	Version int32
}

func (q *Queries) OptimisticUpdateAccount(ctx context.Context, arg OptimisticUpdateAccountParams) (Account, error) {
	row := q.db.QueryRowContext(ctx, optimisticUpdateAccount, arg.Balance, arg.ID, arg.Version)
	var i Account
	err := row.Scan(&i.ID, &i.Balance, &i.Version)
	return i, err
}

const updateAccount = `-- name: UpdateAccount :one
UPDATE account SET balance=$1, version=version+1
    WHERE id=$2
RETURNING id, balance, version
`

type UpdateAccountParams struct {
	Balance string
	ID      int32
}

func (q *Queries) UpdateAccount(ctx context.Context, arg UpdateAccountParams) (Account, error) {
	row := q.db.QueryRowContext(ctx, updateAccount, arg.Balance, arg.ID)
	var i Account
	err := row.Scan(&i.ID, &i.Balance, &i.Version)
	return i, err
}
