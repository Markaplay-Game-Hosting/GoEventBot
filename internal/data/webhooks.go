package data

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Markaplay-Game-Hosting/GoEventBot/internal/validator"
	"github.com/google/uuid"
	"time"
)

type Webhook struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	URL  string    `json:"url"`
}

func ValidateWebHook(v *validator.Validator, webhook *Webhook) {
	v.Check(webhook.Name != "", "name", "must be provided")
	v.Check(len(webhook.Name) <= 100, "name", "must not be more than 100 characters long")
	v.Check(webhook.URL != "", "url", "must be provided")
	v.Check(len(webhook.URL) <= 150, "url", "must not be more than 255 characters long")
	v.Check(validator.Matches(webhook.URL, validator.UrlWebhookRX), "url", "must be a valid Discord webhook URL")
}

type WebhookModel struct {
	DB *sql.DB
}

func (m WebhookModel) Insert(webhook *Webhook) error {
	query := `INSERT INTO webhooks (name, url) VALUES ($1, $2) RETURNING id`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, webhook.Name, webhook.URL).Scan(&webhook.ID)
	if err != nil {
		return err
	}
	return nil
}

func (m WebhookModel) GetByID(id uuid.UUID) (*Webhook, error) {
	query := `SELECT id, name, url FROM webhooks WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var webhook Webhook
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&webhook.ID, &webhook.Name, &webhook.URL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // No record found
		}
		return nil, err
	}
	return &webhook, nil
}

func (m WebhookModel) GetAll() ([]Webhook, error) {
	query := `SELECT id, name, url FROM webhooks`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var webhooks []Webhook
	for rows.Next() {
		var webhook Webhook
		if err := rows.Scan(&webhook.ID, &webhook.Name, &webhook.URL); err != nil {
			return nil, err
		}
		webhooks = append(webhooks, webhook)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return webhooks, nil
}

func (m WebhookModel) Update(webhook *Webhook) error {
	query := `UPDATE webhooks SET name = $1, url = $2 WHERE id = $3`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, webhook.Name, webhook.URL, webhook.ID)
	if err != nil {
		return err
	}
	return nil
}

func (m WebhookModel) Delete(id uuid.UUID) error {
	query := `DELETE FROM webhooks WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
