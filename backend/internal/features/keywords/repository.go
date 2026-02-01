package keywords

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrDuplicateKeyword = errors.New("duplicate keyword")

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) List(ctx context.Context, query ListQuery) ([]KeywordRow, int, error) {
	countQuery := `SELECT COUNT(*) FROM keywords WHERE is_deleted = FALSE`
	listQuery := `SELECT id, keyword, normalized_value, status, is_deleted, created_at, updated_at FROM keywords WHERE is_deleted = FALSE`

	args := []interface{}{}
	argIndex := 1

	if query.Q != "" {
		filter := fmt.Sprintf(" AND keyword ILIKE $%d", argIndex)
		countQuery += filter
		listQuery += filter
		args = append(args, "%"+query.Q+"%")
		argIndex++
	}

	var total int
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	listQuery += " ORDER BY created_at DESC"
	offset := (query.Page - 1) * query.PageSize
	listQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, query.PageSize, offset)

	rows, err := r.pool.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var keywords []KeywordRow
	for rows.Next() {
		var k KeywordRow
		if err := rows.Scan(&k.ID, &k.Keyword, &k.NormalizedValue, &k.Status, &k.IsDeleted, &k.CreatedAt, &k.UpdatedAt); err != nil {
			return nil, 0, err
		}
		keywords = append(keywords, k)
	}

	return keywords, total, nil
}

func (r *Repository) Create(ctx context.Context, value, normalizedValue string) (*KeywordRow, error) {
	row := &KeywordRow{}
	err := r.pool.QueryRow(ctx, `
		INSERT INTO keywords (keyword, normalized_value, status, is_deleted)
		VALUES ($1, $2, 'active', FALSE)
		RETURNING id, keyword, normalized_value, status, is_deleted, created_at, updated_at
	`, value, normalizedValue).Scan(&row.ID, &row.Keyword, &row.NormalizedValue, &row.Status, &row.IsDeleted, &row.CreatedAt, &row.UpdatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrDuplicateKeyword
		}
		return nil, err
	}

	return row, nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (*KeywordRow, error) {
	row := &KeywordRow{}
	err := r.pool.QueryRow(ctx, `
		SELECT id, keyword, normalized_value, status, is_deleted, created_at, updated_at
		FROM keywords
		WHERE id = $1 AND is_deleted = FALSE
	`, id).Scan(&row.ID, &row.Keyword, &row.NormalizedValue, &row.Status, &row.IsDeleted, &row.CreatedAt, &row.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return row, nil
}

func (r *Repository) SoftDelete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE keywords
		SET status = 'inactive', is_deleted = TRUE, deleted_at = NOW()
		WHERE id = $1
	`, id)
	return err
}

func NormalizeValue(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}
