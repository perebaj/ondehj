package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/perebaj/ondehj/api"
	"golang.org/x/exp/slog"
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
	slog.Info(slog.LevelInfo.String())
	logger := slog.New(slog.NewJSONHandler(os.Stdout))
	slog.SetDefault(logger)

	dbpool, err := pgxpool.New(context.Background(), databaseUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	err = dbpool.Ping(context.Background())
	if err != nil {
		slog.Error(fmt.Sprintf("Unable to ping database: %v\n", err))
		os.Exit(1)
	}
	slog.Info("Connected successfully to database")
	defer dbpool.Close()

	mux := api.HandlerFactory(dbpool)
	slog.Info(fmt.Sprintf("Starting server on port %s", settings.ServicePort))
	err = http.ListenAndServe(fmt.Sprintf(":%s", settings.ServicePort), mux)
	if err != nil {
		slog.Error(fmt.Sprintf("Unable to start server: %v", err))
		os.Exit(1)
	}
}
