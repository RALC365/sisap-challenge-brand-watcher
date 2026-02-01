package export

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, filename string, filterParams FilterParams) (string, error) {
	paramsJSON, err := json.Marshal(filterParams)
	if err != nil {
		return "", err
	}

	var id string
	err = r.pool.QueryRow(ctx, `
		INSERT INTO exports (filename, filter_params, status, record_count)
		VALUES ($1, $2, 'pending', 0)
		RETURNING id
	`, filename, paramsJSON).Scan(&id)

	return id, err
}

func (r *Repository) UpdateCompleted(ctx context.Context, id string, recordCount int) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE exports
		SET status = 'completed', record_count = $2
		WHERE id = $1
	`, id, recordCount)
	return err
}

func (r *Repository) UpdateFailed(ctx context.Context, id string, errorMessage string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE exports
		SET status = 'failed', error_message = $2
		WHERE id = $1
	`, id, errorMessage)
	return err
}
