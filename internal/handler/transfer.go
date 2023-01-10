package handler

import (
    "context"
    "database/sql"
    "fmt"
    "github.com/emmanuelperotto/locks-and-concurrency/internal/repository"
    "github.com/gofiber/fiber/v2"
    "strconv"
)

type (
    TransferResponse struct {
        From repository.Account `json:"from"`
        To   repository.Account `json:"to"`

        Amount float64 `json:"amount"`
    }

    TransferRequest struct {
        From   int32   `json:"from"`
        To     int32   `json:"to"`
        Amount float64 `json:"amount"`
    }

    Transfer struct {
        queries *repository.Queries
        db      *sql.DB
    }
)

func NewTransfer(queries *repository.Queries, db *sql.DB) Transfer {
    return Transfer{
        queries: queries,
        db:      db,
    }
}

// Transfer is the efficient way to make transfers between accounts. It also avoids deadlocks
//BEGIN
//Create Transfer (INSERT INTO)
//Debit 'from' account balance (UPDATE SET balance= balance - ? WHERE balance >= ?)
//Credit 'to' account balance (UPDATE SET balance= balance + ?)
//COMMIT
func (h Transfer) Transfer(c *fiber.Ctx) error {
    ctx := c.Context()
    req := new(TransferRequest)

    if err := c.BodyParser(req); err != nil {
        return err
    }

    //init tx
    tx, err := h.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := h.queries.WithTx(tx)
    transfer, err := qtx.CreateTransfer(ctx, repository.CreateTransferParams{
        Amount:        fmt.Sprintf("%f", req.Amount),
        FromAccountID: req.From,
        ToAccountID:   req.To,
    })
    if err != nil {
        return err
    }

    fromAcc, toAcc, err := h.updateAccountBalances(ctx, qtx, req)
    if err != nil {
        return err
    }

    if err := tx.Commit(); err != nil {
        return err
    }

    return c.JSON(TransferResponse{
        From:   fromAcc,
        To:     toAcc,
        Amount: strToFloat64(transfer.Amount),
    })
}

func (h Transfer) updateAccountBalances(ctx context.Context, qtx *repository.Queries, req *TransferRequest) (from repository.Account, to repository.Account, err error) {
    if req.From < req.To {
        from, err = qtx.DebitAccount(ctx, repository.DebitAccountParams{
            Amount: fmt.Sprintf("%f", req.Amount),
            ID:     req.From,
        })

        to, err = qtx.CreditAccount(ctx, repository.CreditAccountParams{
            Amount: fmt.Sprintf("%f", req.Amount),
            ID:     req.To,
        })

        return
    }

    to, err = qtx.CreditAccount(ctx, repository.CreditAccountParams{
        Amount: fmt.Sprintf("%f", req.Amount),
        ID:     req.To,
    })

    from, err = qtx.DebitAccount(ctx, repository.DebitAccountParams{
        Amount: fmt.Sprintf("%f", req.Amount),
        ID:     req.From,
    })

    return
}

// InconsistentTransfer is an inefficient way to make a transfer between accounts
// BEGIN
// Get Account 1 ( SELECT WHERE id=?)
// Get Account 2 ( SELECT WHERE id=?)
// Create Transfer (INSERT INTO)
// Update from account balance (UPDATE SET balance=?)
// Update to account balance (UPDATE SET balance=?)
// COMMIT
func (h Transfer) InconsistentTransfer(c *fiber.Ctx) error {
    ctx := c.Context()
    req := new(TransferRequest)

    if err := c.BodyParser(req); err != nil {
        return err
    }

    //init tx
    tx, err := h.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := h.queries.WithTx(tx)
    fromAcc, err := qtx.GetAccount(ctx, req.From)
    if err != nil {
        return err
    }

    toAcc, err := qtx.GetAccount(ctx, req.To)
    if err != nil {
        return err
    }

    transfer, err := qtx.CreateTransfer(ctx, repository.CreateTransferParams{
        Amount:        fmt.Sprintf("%f", req.Amount),
        FromAccountID: fromAcc.ID,
        ToAccountID:   toAcc.ID,
    })
    if err != nil {
        return err
    }

    if fromAcc, err = qtx.UpdateAccount(ctx, repository.UpdateAccountParams{
        Balance: fmt.Sprintf("%f", strToFloat64(fromAcc.Balance)-req.Amount),
        ID:      fromAcc.ID,
    }); err != nil {
        return err
    }

    if toAcc, err = qtx.UpdateAccount(ctx, repository.UpdateAccountParams{
        Balance: fmt.Sprintf("%f", strToFloat64(toAcc.Balance)+req.Amount),
        ID:      toAcc.ID,
    }); err != nil {
        return err
    }

    if err := tx.Commit(); err != nil {
        return err
    }

    return c.JSON(TransferResponse{
        From:   fromAcc,
        To:     toAcc,
        Amount: strToFloat64(transfer.Amount),
    })
}

// PessimisticLockTransfer is still an inefficient way to make a transfer between accounts,
// but it is at least consistent
//
// BEGIN
// Get Account 1 and lock it ( SELECT WHERE id=? FOR UPDATE)
// Get Account 2 and lock it ( SELECT WHERE id=? FOR UPDATE)
// Create Transfer
// Update from account balance (UPDATE SET balance=?)
// Update to account balance (UPDATE SET balance=?)
// COMMIT and release locks
func (h Transfer) PessimisticLockTransfer(c *fiber.Ctx) error {
    ctx := c.Context()
    req := new(TransferRequest)

    if err := c.BodyParser(req); err != nil {
        return err
    }

    //init tx
    tx, err := h.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := h.queries.WithTx(tx)
    fromAcc, err := qtx.GetAccountForUpdate(ctx, req.From)
    if err != nil {
        return err
    }

    toAcc, err := qtx.GetAccountForUpdate(ctx, req.To)
    if err != nil {
        return err
    }

    transfer, err := qtx.CreateTransfer(ctx, repository.CreateTransferParams{
        Amount:        fmt.Sprintf("%f", req.Amount),
        FromAccountID: fromAcc.ID,
        ToAccountID:   toAcc.ID,
    })
    if err != nil {
        return err
    }

    if fromAcc, err = qtx.UpdateAccount(ctx, repository.UpdateAccountParams{
        Balance: fmt.Sprintf("%f", strToFloat64(fromAcc.Balance)-req.Amount),
        ID:      fromAcc.ID,
    }); err != nil {
        return err
    }

    if toAcc, err = qtx.UpdateAccount(ctx, repository.UpdateAccountParams{
        Balance: fmt.Sprintf("%f", strToFloat64(toAcc.Balance)+req.Amount),
        ID:      toAcc.ID,
    }); err != nil {
        return err
    }

    if err := tx.Commit(); err != nil {
        return err
    }

    return c.JSON(TransferResponse{
        From:   fromAcc,
        To:     toAcc,
        Amount: strToFloat64(transfer.Amount),
    })
}

// OptimisticLockTransfer is still an inefficient way to make a transfer between accounts,
// but it is at least consistent
//
// BEGIN
// Get Account 1 ( SELECT WHERE id=?)
// Get Account 2 ( SELECT WHERE id=?)
// Create Transfer (INSERT INTO)
// TRY to Update from account balance (UPDATE SET balance=? WHERE version=?)
// TRY to Update to account balance (UPDATE SET balance=? WHERE version=?)
// COMMIT
func (h Transfer) OptimisticLockTransfer(c *fiber.Ctx) error {
    ctx := c.Context()
    req := new(TransferRequest)

    if err := c.BodyParser(req); err != nil {
        return err
    }

    //init tx
    tx, err := h.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := h.queries.WithTx(tx)
    fromAcc, err := qtx.GetAccount(ctx, req.From)
    if err != nil {
        return err
    }

    toAcc, err := qtx.GetAccount(ctx, req.To)
    if err != nil {
        return err
    }

    transfer, err := qtx.CreateTransfer(ctx, repository.CreateTransferParams{
        Amount:        fmt.Sprintf("%f", req.Amount),
        FromAccountID: fromAcc.ID,
        ToAccountID:   toAcc.ID,
    })
    if err != nil {
        return err
    }

    if fromAcc, err = qtx.OptimisticUpdateAccount(ctx, repository.OptimisticUpdateAccountParams{
        Balance: fmt.Sprintf("%f", strToFloat64(fromAcc.Balance)-req.Amount),
        ID:      fromAcc.ID,
        Version: fromAcc.Version,
    }); err != nil {
        return err
    }

    if toAcc, err = qtx.OptimisticUpdateAccount(ctx, repository.OptimisticUpdateAccountParams{
        Balance: fmt.Sprintf("%f", strToFloat64(toAcc.Balance)+req.Amount),
        ID:      toAcc.ID,
        Version: toAcc.Version,
    }); err != nil {
        return err
    }

    if err := tx.Commit(); err != nil {
        return err
    }

    return c.JSON(TransferResponse{
        From:   fromAcc,
        To:     toAcc,
        Amount: strToFloat64(transfer.Amount),
    })
}

func strToFloat64(str string) float64 {
    float, _ := strconv.ParseFloat(str, 64)
    return float
}
