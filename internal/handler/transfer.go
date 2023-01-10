package handler

import (
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

func strToFloat64(str string) float64 {
    float, _ := strconv.ParseFloat(str, 64)
    return float
}
