package data

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Markaplay-Game-Hosting/GoEventBot/internal/validator"
	"github.com/google/uuid"
	"time"
)

type Event struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Duration    string    `json:"duration"`
	RRule       string    `json:"rrule,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedDate time.Time `json:"created_date"`
	UpdatedDate time.Time `json:"updated_date"`
}

type EventInstance struct {
	EventID     uuid.UUID `json:"event_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Duration    string    `json:"duration"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
}

func ValidateEvent(v *validator.Validator, event *Event) {
	v.Check(event.Title != "", "title", "must be provided")
	v.Check(len(event.Title) <= 100, "title", "must not be more than 100 bytes long")
	v.Check(event.Description != "", "description", "must be provided")
	v.IsValidDurationRule(event.Duration)

	v.IsValidRRule(event.RRule)
}

type EventModel struct {
	DB *sql.DB
}

func (e EventModel) Insert(event *Event) error {
	query := `INSERT INTO events (title, description, duration, rrule, is_active) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_date, updated_date`

	args := []any{event.Title, event.Description, event.Duration, event.RRule, event.IsActive}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := e.DB.QueryRowContext(ctx, query, args...).Scan(&event.ID, &event.CreatedDate, &event.UpdatedDate)

	if err != nil {
		return err
	}
	return nil
}

func (e EventModel) Get(ID uuid.UUID) (Event, error) {
	query := `SELECT id, title, description, duration, rrule, is_active, created_date, updated_date FROM events WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var event Event

	err := e.DB.QueryRowContext(ctx, query, ID).Scan(
		&event.ID,
		&event.Title,
		&event.Description,
		&event.Duration,
		&event.RRule,
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
	query := `SELECT id, title, description, duration, rrule, is_active, created_date, updated_date FROM events`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := e.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err = rows.Close()
	}(rows)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var event Event
		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.Duration,
			&event.RRule,
			&event.IsActive,
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
	query := `UPDATE events SET title = $1, description = $2, is_active = $3, duration = $4, rrule = $5, updated_date = NOW() WHERE id = $6 RETURNING updated_date`

	args := []any{event.Title, event.Description, event.IsActive, event.Duration, event.RRule, event.ID}

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

func (e EventModel) Delete(ID uuid.UUID) error {
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
	query := `SELECT id, title, description, duration, rrule FROM events WHERE is_active = true`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := e.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err = rows.Close()
	}(rows)

	if err != nil {
		return nil, err
	}

	var events []Event
	for rows.Next() {
		var event Event
		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.Duration,
			&event.RRule,
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
