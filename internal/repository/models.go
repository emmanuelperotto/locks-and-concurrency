// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package repository

import ()

type Account struct {
	ID      int32
	Balance string
	Version int32
}

type Transfer struct {
	ID            int32
	Amount        string
	FromAccountID int32
	ToAccountID   int32
}
