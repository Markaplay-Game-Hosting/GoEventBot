package data

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Markaplay-Game-Hosting/GoEventBot/internal/validator"
	"github.com/google/uuid"
	"time"
)

// Event
// @Description Event Response object
type Event struct {
	// ID of the event
	ID uuid.UUID `json:"id"`
	// Title of the event
	Title string `json:"title"`
	// Description/Content of the event
	Description string `json:"description"`
	// Duration of the event following the ISO 8601 standard
	Duration string `json:"duration" example:"PT30M"`
	// Recurrence rule following RFC 5545 https://icalendar.org/rrule-tool.html
	RRule string `json:"rrule,omitempty" example:"FREQ=WEEKLY;INTERVAL=1;BYDAY=MO;UNTIL=20250731T000000Z"`
	// Discord ID of the channel the bot will post the event
	ChannelID int `json:"channel_id,omitempty"`
	// Discord Guild ID/server the bot will publish on
	GuildID int `json:"guild_id,omitempty"`
	// Tell if the event is currently active
	IsActive bool `json:"is_active"`
	// Creation date of the event
	CreatedDate time.Time `json:"created_date"`
	// Last modification of the event
	UpdatedDate time.Time `json:"updated_date"`
} // @name Event.Get

// EventInstance
// @Description Event Instance Response object
type EventInstance struct {
	// ID of the event
	EventID uuid.UUID `json:"event_id"`
	// Title of the event
	Title string `json:"title"`
	// Description of the event
	Description string `json:"description"`
	// Duration of the event following the ISO 8601 standard
	Duration string `json:"duration" example:"PT30M"`
	// Start date of the event
	StartDate time.Time `json:"start_date"`
	// End date of the event
	EndDate time.Time `json:"end_date"`
} // @name Event.Instance.Get

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
