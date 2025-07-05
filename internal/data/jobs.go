package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Job
// @Description Job model
type Job struct {
	ID            uuid.UUID `json:"id"`
	EventId       uuid.UUID `json:"event_id"`
	ExecutionDate time.Time `json:"execution_date"`
	Status        int       `json:"status"`
} // @name Jobs.Get.Response

type JobStatus int

const (
	Unknown = iota
	Pending
	Running
	Completed
	Failed
	Cancel
)

type JobModel struct {
	DB *sql.DB
}

func (j JobModel) Insert(job *Job) (*Job, error) {
	query := `INSERT INTO jobs (id, event_id, execution_date, status) VALUES ($1, $2, $3, $4) RETURNING id`
	args := []any{job.ID, job.EventId, job.ExecutionDate, job.Status}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := j.DB.QueryRowContext(ctx, query, args...).Scan(&job.ID)
	if err != nil {
		return nil, err
	}
	return job, nil
}

func (j JobModel) Get(ID uuid.UUID) (Job, error) {
	query := `SELECT id, event_id, execution_date, status FROM jobs WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var job Job

	err := j.DB.QueryRowContext(ctx, query, ID).Scan(
		&job.ID,
		&job.EventId,
		&job.ExecutionDate,
		&job.Status,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Job{}, nil // No job found
		}
		return Job{}, err // Other error
	}
	return job, nil
}

func (j JobModel) GetByEventID(eventID uuid.UUID) ([]Job, error) {
	query := `SELECT id, event_id, execution_date, status FROM jobs WHERE event_id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := j.DB.QueryContext(ctx, query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []Job
	for rows.Next() {
		var job Job
		if err := rows.Scan(&job.ID, &job.EventId, &job.ExecutionDate, &job.Status); err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return jobs, nil
}

func (j JobModel) Update(job *Job) error {
	query := `UPDATE jobs SET event_id = $1, execution_date = $2, status = $3 WHERE id = $4`
	args := []any{job.EventId, job.ExecutionDate, job.Status, job.ID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := j.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (j JobModel) Delete(ID uuid.UUID) error {
	query := `DELETE FROM jobs WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := j.DB.ExecContext(ctx, query, ID)
	if err != nil {
		return err
	}
	return nil
}
func (j JobModel) GetAll() ([]Job, error) {
	query := `SELECT id, event_id, execution_date, status FROM jobs`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := j.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err = rows.Close()
	}(rows)

	if err != nil {
		return nil, err
	}

	var jobs []Job

	for rows.Next() {
		var job Job
		if err := rows.Scan(&job.ID, &job.EventId, &job.ExecutionDate, &job.Status); err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return jobs, nil
}

func (j JobModel) SkipJob(ID uuid.UUID) error {
	query := `UPDATE jobs SET status = 5 WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := j.DB.ExecContext(ctx, query, ID)
	if err != nil {
		return err
	}
	return nil
}

func (j JobStatus) String() string {
	return [...]string{"Unknown", "Pending", "Running", "Completed", "Failed"}[j]
}

func (j JobStatus) ToID() int {
	return int(j)
}
