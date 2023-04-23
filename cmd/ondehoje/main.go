package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/perebaj/ondehj/api"
	// concurrency safe
)

func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

type Settings struct {
	DatabaseHost     string
	DatabasePort     string
	DatabaseUser     string
	DatabasePassword string
	ServicePort      string
}

// centralize all settings in a single struct
var settings = Settings{
	DatabaseHost:     getEnvWithDefault("POSTGRES_HOST", "localhost"),
	DatabasePort:     getEnvWithDefault("POSTGRES_PORT", "5432"),
	DatabaseUser:     getEnvWithDefault("POSTGRES_USER", "postgres"),
	DatabasePassword: getEnvWithDefault("POSTGRES_PASSWORD", "example_password"),
	ServicePort:      getEnvWithDefault("SERVICE_PORT", "8000"),
}

func main() {
	databaseUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/example_db?sslmode=disable",
		settings.DatabaseUser,
		settings.DatabasePassword,
		settings.DatabaseHost,
		settings.DatabasePort,
	)
	fmt.Printf("Connecting to database: %s", databaseUrl)
	dbpool, err := pgxpool.New(context.Background(), databaseUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	mux := api.HandlerFactory(dbpool)
	fmt.Printf("Starting server on port %s \n", settings.ServicePort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", settings.ServicePort), mux))
}
