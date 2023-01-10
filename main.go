package main

import (
    "database/sql"
    "github.com/emmanuelperotto/locks-and-concurrency/internal/handler"
    "github.com/emmanuelperotto/locks-and-concurrency/internal/repository"
    "github.com/gofiber/fiber/v2"
    _ "github.com/lib/pq"
    "log"
)

func main() {
    app := fiber.New()
    queries, db := connectDB()
    defer func(d *sql.DB) {
        _ = db.Close()
    }(db)

    h := handler.NewTransfer(queries, db)

    app.Post("/transfer", h.Transfer)
    app.Post("/optimistic-transfer", h.Transfer)
    app.Post("/pessimistic-transfer", h.Transfer)

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
