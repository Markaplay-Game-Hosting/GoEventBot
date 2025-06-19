package data

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Markaplay-Game-Hosting/GoEventBot/internal/validator"
	"github.com/google/uuid"
	"time"
)

type Tag struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedDate time.Time `json:"created_date"`
	UpdatedDate time.Time `json:"updated_date"`
}

func ValidateTag(v *validator.Validator, tag *Tag) {
	v.Check(tag.Name != "", "name", "must be provided")
	v.Check(len(tag.Name) <= 20, "name", "must not be more than 20 characters long")
	v.Check(tag.Description != "", "description", "must be provided")
	v.Check(len(tag.Description) <= 100, "description", "must not be more than 100 characters long")
}

type TagModel struct {
	DB *sql.DB
}

func (t TagModel) Insert(tag *Tag) error {
	query := `INSERT INTO tags (name, description) VALUES ($1, $2) RETURNING id, created_date, updated_date`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := t.DB.QueryRowContext(ctx, query, tag.Name, tag.Description).Scan(&tag.ID, &tag.CreatedDate, &tag.UpdatedDate)

	if err != nil {
		return err
	}

	return nil
}

func (t TagModel) GetByID(id uuid.UUID) (*Tag, error) {
	query := `SELECT id, name, description, created_date, updated_date FROM tags WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tag := &Tag{}
	err := t.DB.QueryRowContext(ctx, query, id).Scan(&tag.ID, &tag.Name, &tag.Description, &tag.CreatedDate, &tag.UpdatedDate)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // No tag found
		}
		return nil, err // Other error
	}

	return tag, nil
}

func (t TagModel) Update(tag *Tag) error {
	query := `UPDATE tags SET name = $1, description = $2, updated_date = NOW() WHERE id = $3`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := t.DB.ExecContext(ctx, query, tag.Name, tag.Description, tag.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows // No rows were updated
	}

	return nil
}

func (t TagModel) Delete(id uuid.UUID) error {
	query := `DELETE FROM tags WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := t.DB.ExecContext(ctx, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows // No rows were deleted
	}

	return nil
}

func (t TagModel) GetAll() ([]Tag, error) {
	query := `SELECT id, name, description, created_date, updated_date FROM tags`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := t.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
	}(rows)

	var tags []Tag
	for rows.Next() {
		var tag Tag
		err := rows.Scan(&tag.ID, &tag.Name, &tag.Description, &tag.CreatedDate, &tag.UpdatedDate)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}
