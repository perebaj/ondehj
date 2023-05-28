package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool" // concurrency safe
	"github.com/perebaj/ondehj/api"
	"golang.org/x/exp/slog"
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
	DatabaseName     string
	SSLMode          string
}

// centralize all settings in a single struct
var settings = Settings{
	DatabaseHost:     getEnvWithDefault("POSTGRES_HOST", "localhost"),
	DatabasePort:     getEnvWithDefault("POSTGRES_PORT", "5432"),
	DatabaseUser:     getEnvWithDefault("POSTGRES_USER", "postgres"),
	DatabasePassword: getEnvWithDefault("POSTGRES_PASSWORD", "example_password"),
	ServicePort:      getEnvWithDefault("PORT", "8000"),
	DatabaseName:     getEnvWithDefault("POSTGRES_DB", "example_db"),
	SSLMode:          getEnvWithDefault("POSTGRES_SSLMODE", "disable"),
}

func main() {
	databaseUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		settings.DatabaseUser,
		settings.DatabasePassword,
		settings.DatabaseHost,
		settings.DatabasePort,
		settings.DatabaseName,
		settings.SSLMode,
	)
	slog.Info(slog.LevelInfo.String())
	handler := slog.HandlerOptions{AddSource: true, Level: slog.LevelInfo}.NewJSONHandler(os.Stdout)
	logger := slog.New(handler)
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
	srv := http.Server{
		Addr:         fmt.Sprintf(":%s", settings.ServicePort),
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	err = srv.ListenAndServe()
	if err != nil {
		slog.Error(fmt.Sprintf("Unable to start server: %v", err))
		os.Exit(1)
	}
}
