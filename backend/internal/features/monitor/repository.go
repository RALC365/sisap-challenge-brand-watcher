package monitor

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) GetState(ctx context.Context) (*MonitorStateRow, error) {
	row := &MonitorStateRow{}
	err := r.pool.QueryRow(ctx, `
		SELECT id, state, last_run_at, last_success_at, last_error_code, last_error
		FROM monitor_state
		WHERE id = 1
	`).Scan(&row.ID, &row.State, &row.LastRunAt, &row.LastSuccessAt, &row.LastErrorCode, &row.LastErrorMessage)
	if err != nil {
		return nil, err
	}
	return row, nil
}

func (r *Repository) GetLastCompletedRun(ctx context.Context) (*LastRunRow, error) {
	row := &LastRunRow{}
	err := r.pool.QueryRow(ctx, `
		SELECT certificates_processed, matches_found, parse_error_count, duration_ms, ct_latency_ms, db_latency_ms
		FROM monitor_run
		WHERE state = 'completed'
		ORDER BY ended_at DESC
		LIMIT 1
	`).Scan(&row.CertificatesProcessed, &row.MatchesFound, &row.ParseErrorCount, &row.DurationMs, &row.CtLatencyMs, &row.DbLatencyMs)
	if err != nil {
		return nil, err
	}
	return row, nil
}
