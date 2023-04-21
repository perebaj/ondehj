package event

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgconn"
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

//create a json for the event
 

var ErrEventAlreadyExists = errors.New("event already exists")

type Repository interface {
	Create(ctx context.Context, event Event) (*Event, error)
}

type SQLRepository struct {
	db *sql.DB
}

func NewSQLRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

func (r *SQLRepository) Create(ctx context.Context, event Event) (*Event, error) {
	var id int64
	err := r.db.QueryRowContext(ctx, `
		INSET INTO events (title, description, location, instagram_page, start_time, end_time) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		event.Title, event.Description, event.Location, event.InstagramPage, event.StartTime, event.EndTime).Scan(&id)
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" { // unique_violation
				return nil, ErrEventAlreadyExists
			}
		}

		return nil, err
	}
	event.ID = id
	return &event, nil
}
