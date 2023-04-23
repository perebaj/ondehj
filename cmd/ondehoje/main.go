package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool" // concurrency safe
	"github.com/perebaj/ondehj/api"
)

var databaseUrl = "postgres://postgres:example_password@localhost:5432/example_db?sslmode=disable"

func main() {
	dbpool, err := pgxpool.New(context.Background(), databaseUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	mux := api.HandlerFactory(dbpool)
	fmt.Println("Starting server on port 8000")
	log.Fatal(http.ListenAndServe(":8000", mux))
}
