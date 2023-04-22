package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/perebaj/ondehj/event"
)

var databaseUrl = "postgres://postgres:example_password@localhost:5432/example_db?sslmode=disable"

func main() {
	dbpool, err := pgxpool.New(context.Background(), databaseUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()
	eventSQLRepo := event.EventSQLRepository(dbpool)
	err = eventSQLRepo.Migrate()
	if err != nil {
		fmt.Println(err)
	}

}
