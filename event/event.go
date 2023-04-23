package event

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrDeleteFailed = errors.New("Delete failed")
)

type Event struct {
	ID            int64     `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Location      string    `json:"location"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	InstagramPage string    `json:"instagram_page"`
}

type Repository interface {
	Create(ctx context.Context, event Event) (*Event, error)
	Migrate() error
	Delete(ctx context.Context, id int64) error
	All(ctx context.Context) ([]Event, error)
	GetByID(ctx context.Context, id int64) (*Event, error)
	Update(ctx context.Context, id int64, newEvent Event) (*Event, error)
}

type SQLRepository struct {
	db *pgxpool.Pool
}

func EventSQLRepository(db *pgxpool.Pool) *SQLRepository {
	return &SQLRepository{db: db}
}

func (r *SQLRepository) Update(ctx context.Context, id int64, newEvent Event) (*Event, error) {

	err := r.db.QueryRow(ctx,
		`UPDATE events SET title = $1, description = $2, location = $3, instagram_page = $4, start_time = $5, end_time = $6 WHERE id = $7 RETURNING id`,
		newEvent.Title, newEvent.Description, newEvent.Location, newEvent.InstagramPage, newEvent.StartTime, newEvent.EndTime, id).Scan(
		&newEvent.ID)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &newEvent, nil

}

func (r *SQLRepository) GetByID(ctx context.Context, id int64) (*Event, error) {
	var event Event
	err := r.db.QueryRow(
		ctx,
		`SELECT id, title, description, location, start_time, end_time, instagram_page FROM events WHERE id = $1`, id).Scan(&event.ID, &event.Title, &event.Description, &event.Location, &event.StartTime, &event.EndTime, &event.InstagramPage)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &event, nil
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

func (r *SQLRepository) Delete(ctx context.Context, id int64) error {
	res, err := r.db.Exec(ctx, `DELETE FROM events WHERE id = $1`, id)
	rowsAffcected := res.RowsAffected()
	if rowsAffcected == 0 {
		return ErrDeleteFailed
	}
	return err
}

func (r *SQLRepository) All(ctx context.Context) ([]Event, error) {
	rows, err := r.db.Query(ctx, `SELECT * FROM events`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var events []Event
	for rows.Next() {
		var event Event
		err = rows.Scan(&event.ID, &event.Title, &event.Description, &event.Location, &event.StartTime, &event.EndTime, &event.InstagramPage)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}
