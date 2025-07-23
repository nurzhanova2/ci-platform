package repository

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"github.com/nurzhanova2/ci-platform/internal/model"
)

type DBRepository struct {
	db *sql.DB
}

func NewJobDB(dsn string, driver string) (*DBRepository, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("DB connection failed: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("DB ping failed: %w", err)
	}
	return &DBRepository{db: db}, nil
}

func (r *DBRepository) SavePipelineRun(ctx context.Context, run model.PipelineRun) error {
	query := `INSERT INTO pipeline_runs (repo, status, logs, started_at, finished_at)
			  VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, run.Repo, run.Status, run.Logs, run.StartedAt, run.FinishedAt)
	return err
}

func (r *DBRepository) GetLastRuns(ctx context.Context, limit int) ([]model.PipelineRun, error) {
	query := `SELECT id, repo, status, logs, started_at, finished_at
			  FROM pipeline_runs ORDER BY started_at DESC LIMIT $1`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []model.PipelineRun
	for rows.Next() {
		var r model.PipelineRun
		err := rows.Scan(&r.ID, &r.Repo, &r.Status, &r.Logs, &r.StartedAt, &r.FinishedAt)
		if err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, nil
}
