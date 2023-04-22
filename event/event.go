package event

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Event struct {
	ID            int64
	Title         string
	Description   string
	Location      string
	StartTime     time.Time
	EndTime       time.Time
	InstagramPage string
}

type Repository interface {
	Create(ctx context.Context, event Event) (*Event, error)
	Migrate() error
}

type SQLRepository struct {
	db *pgxpool.Pool
}

func EventSQLRepository(db *pgxpool.Pool) *SQLRepository {
	return &SQLRepository{db: db}
}

func (r *SQLRepository) Create(ctx context.Context, event Event) (*Event, error) {
	var id int64
	err := r.db.QueryRow(ctx, `
		INSERT INTO events (title, description, location, instagram_page, start_time, end_time) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		event.Title, event.Description, event.Location, event.InstagramPage, event.StartTime, event.EndTime).Scan(&id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	event.ID = id
	return &event, nil
}

func (r *SQLRepository) Migrate() error {
	query := `
		CREATE TABLE events (
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			description TEXT,
			location TEXT,
			start_time TIMESTAMP WITH TIME ZONE NOT NULL,
			end_time TIMESTAMP WITH TIME ZONE NOT NULL,
			instagram_page TEXT
		);
	`
	fmt.Println("Creating events table...")
	_, err := r.db.Exec(context.Background(), query)
	return err
}
