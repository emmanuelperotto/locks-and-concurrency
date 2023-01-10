package main

import (
    "database/sql"
    "errors"
    "github.com/emmanuelperotto/locks-and-concurrency/internal/handler"
    "github.com/emmanuelperotto/locks-and-concurrency/internal/repository"
    "github.com/gofiber/fiber/v2"
    _ "github.com/lib/pq"
    "log"
)

func main() {
    app := fiber.New(fiber.Config{ErrorHandler: func(ctx *fiber.Ctx, err error) error {
        log.Println("error on request", err)
        code := fiber.StatusInternalServerError

        var e *fiber.Error
        if errors.As(err, &e) {
            code = e.Code
        }
        ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
        return ctx.Status(code).JSON(fiber.Map{
            "error": err.Error(),
        })
    }})
    queries, db := connectDB()
    defer func(d *sql.DB) {
        _ = db.Close()
    }(db)

    h := handler.NewTransfer(queries, db)

    app.Post("/inconsistent-transfer", h.InconsistentTransfer)
    app.Post("/optimistic-transfer", h.OptimisticLockTransfer)
    app.Post("/pessimistic-transfer", h.PessimisticLockTransfer)

    app.Post("/transfer", h.Transfer)

    log.Panic(app.Listen(":3000"))
}

func connectDB() (*repository.Queries, *sql.DB) {

    db, err := sql.Open("postgres", "postgresql://user:example@localhost:5432/ledger?sslmode=disable")
    if err != nil {
        log.Panic("can't connect to DB", err)
    }

    if err := db.Ping(); err != nil {
        log.Panic("ping DB didn't work", err)
    }

    return repository.New(db), db
}
