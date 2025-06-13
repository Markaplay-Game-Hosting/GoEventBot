package data

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"time"
)

type Event struct {
	ID          uuid.UUID     `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Duration    time.Duration `json:"duration"`
	Interval    time.Duration `json:"repeat_interval,omitempty"`
	StartDate   time.Time     `json:"start_date"`
	EndDate     time.Time     `json:"end_date,omitempty"`
	IsActive    bool          `json:"is_active"`
	WebhookID   uuid.UUID     `json:"webhook_id"`
	CreatedDate time.Time     `json:"created_date"`
	UpdatedDate time.Time     `json:"updated_date"`
}

type EventModel struct {
	DB *sql.DB
}

func (e EventModel) Insert(event *Event) error {
	query := `INSERT INTO events (id, title, description, start_date, end_date, webhook_id, is_active) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_date, updated_date`

	args := []any{event.ID, event.Title, event.Description, event.StartDate, event.EndDate}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := e.DB.QueryRowContext(ctx, query, args...).Scan(&event.ID, &event.CreatedDate, &event.UpdatedDate)

	if err != nil {
		return err
	}
	return nil
}

func (e EventModel) Get(ID string) (Event, error) {
	query := `SELECT id, title, description, start_date, end_date, is_active, created_date, updated_date FROM events WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var event Event

	err := e.DB.QueryRowContext(ctx, query, ID).Scan(
		&event.ID,
		&event.Title,
		&event.Description,
		&event.StartDate,
		&event.EndDate,
		&event.IsActive,
		&event.CreatedDate,
		&event.UpdatedDate,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return Event{}, sql.ErrNoRows
		default:
			return Event{}, err
		}
	}
	return event, nil
}

func (e EventModel) GetAll() ([]Event, error) {
	var events []Event
	query := `SELECT id, title, description, start_date, end_date, is_active FROM events`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := e.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var event Event
		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.StartDate,
			&event.EndDate,
			&event.CreatedDate,
			&event.UpdatedDate,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (e EventModel) Update(event *Event) error {
	query := `UPDATE events SET title = $1, description = $2, start_date = $3, end_date = $4, is_active = $5, updated_date = NOW() WHERE id = $6 RETURNING updated_date`

	args := []any{event.Title, event.Description, event.StartDate, event.EndDate, event.IsActive, event.ID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := e.DB.QueryRowContext(ctx, query, args...).Scan(&event.UpdatedDate)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return sql.ErrNoRows
		default:
			return err
		}
	}
	return nil
}

func (e EventModel) Delete(ID string) error {
	query := `DELETE FROM events WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := e.DB.ExecContext(ctx, query, ID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return sql.ErrNoRows
		default:
			return err
		}
	}
	return nil
}

func (e EventModel) GetActiveEvents() ([]Event, error) {
	query := `SELECT id, title, description, start_date, end_date, duration, interval, webhook_id, FROM events WHERE is_active = true`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := e.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var events []Event
	for rows.Next() {
		var event Event
		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.StartDate,
			&event.EndDate,
			&event.Duration,
			&event.Interval,
			&event.WebhookID,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return events, nil
}
