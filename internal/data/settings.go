package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Setting
// @Description setting model
type Setting struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Value       string    `json:"value"`
} // @name setting

type SettingModel struct {
	DB *sql.DB
}

func (m SettingModel) GetByID(id uuid.UUID) (*Setting, error) {
	query := "SELECT id, name, description, value FROM settings WHERE id=$1"

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	setting := &Setting{}
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&setting.ID, &setting.Name, &setting.Description, &setting.Value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
		}
	}
	return setting, err
}

func (m SettingModel) GetByName(name string) (*Setting, error) {
	query := "SELECT id, name, description, value FROM settings WHERE name=$1"

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	setting := &Setting{}
	err := m.DB.QueryRowContext(ctx, query, name).Scan(&setting.ID, &setting.Name, &setting.Description, &setting.Value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
		}
	}
	return setting, err
}

func (m SettingModel) Insert(setting *Setting) error {
	query := "INSERT INTO settings (name, description, value) VALUES ($1, $2, $3)"

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, setting.Name, setting.Description, setting.Value)
	if err != nil {
		return err
	}
	return nil
}

func (m SettingModel) Update(setting *Setting) error {
	query := "UPDATE settings SET name=$1, description=$2, value=$3 WHERE id=$4"
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, setting.Name, setting.Description, setting.Value, setting.ID)
	if err != nil {
		return err
	}
	return nil
}

func (m SettingModel) Delete(id uuid.UUID) error {
	query := "DELETE FROM settings WHERE id=$1"
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (m SettingModel) GetAll() ([]*Setting, error) {
	query := "SELECT id, name, description, value FROM settings"
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var settings []*Setting
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return settings, err
	}

	defer func(rows *sql.Rows) {
		err = rows.Close()
	}(rows)

	if err != nil {
		return settings, err
	}

	for rows.Next() {
		var setting Setting
		if err := rows.Scan(&setting.ID, &setting.Name, &setting.Description, &setting.Value); err != nil {
			return settings, err
		}
		settings = append(settings, &setting)
	}
	return settings, nil
}
