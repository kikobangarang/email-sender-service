package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(path string) (*SQLiteRepository, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	if err := migrate(db); err != nil {
		return nil, err
	}

	return &SQLiteRepository{db: db}, nil
}

func migrate(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS email_jobs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		to_email TEXT NOT NULL,
		subject TEXT NOT NULL,
		body TEXT NOT NULL,
		status TEXT NOT NULL,
		attempts INTEGER NOT NULL DEFAULT 0,
		error_message TEXT,
		created_at DATETIME NOT NULL,
		sent_at DATETIME
	);

	CREATE INDEX IF NOT EXISTS idx_email_jobs_status
	ON email_jobs(status);
	`
	_, err := db.Exec(query)
	return err
}

func (r *SQLiteRepository) Create(job EmailJob) error {
	_, err := r.db.Exec(
		`INSERT INTO email_jobs
		(to_email, subject, body, status, attempts, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		job.To,
		job.Subject,
		job.Body,
		job.Status,
		job.Attempts,
		job.CreatedAt,
	)
	return err
}

func (r *SQLiteRepository) FetchPending(limit int) ([]EmailJob, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Atomically fetch and mark jobs as processing
	rows, err := tx.Query(`
		UPDATE email_jobs
		SET status = ?
		WHERE id IN (
			SELECT id FROM email_jobs
			WHERE status = ?
			ORDER BY created_at
			LIMIT ?
		)
		RETURNING id, to_email, subject, body, attempts, created_at
	`, StatusProcessing, StatusPending, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []EmailJob
	for rows.Next() {
		var job EmailJob
		if err := rows.Scan(
			&job.ID,
			&job.To,
			&job.Subject,
			&job.Body,
			&job.Attempts,
			&job.CreatedAt,
		); err != nil {
			return nil, err
		}

		job.Status = StatusProcessing
		jobs = append(jobs, job)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return jobs, nil
}

func (r *SQLiteRepository) MarkSent(id int64) error {
	res, err := r.db.Exec(
		`UPDATE email_jobs
		 SET status = ?, sent_at = ?
		 WHERE id = ?`,
		StatusSent,
		time.Now(),
		id,
	)
	if err != nil {
		return err
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("job not found")
	}

	return nil
}

func (r *SQLiteRepository) IncrementAttempts(id int64) error {
	res, err := r.db.Exec(
		`UPDATE email_jobs
		 SET attempts = attempts + 1
		 WHERE id = ?`,
		id,
	)
	if err != nil {
		return err
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("job not found")
	}

	return nil
}

func (r *SQLiteRepository) MarkFailed(id int64, sendErr error) error {
	res, err := r.db.Exec(
		`UPDATE email_jobs
		 SET status = ?, error_message = ?
		 WHERE id = ?`,
		StatusFailed,
		sendErr.Error(),
		id,
	)
	if err != nil {
		return err
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("job %d not found", id)
	}

	return nil
}
