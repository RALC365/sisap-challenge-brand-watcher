package db

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool     *pgxpool.Pool
	poolOnce sync.Once
	poolErr  error
)

func InitPool(ctx context.Context, databaseURL string) error {
	poolOnce.Do(func() {
		config, err := pgxpool.ParseConfig(databaseURL)
		if err != nil {
			poolErr = fmt.Errorf("failed to parse database URL: %w", err)
			return
		}

		config.MaxConns = 10
		config.MinConns = 2
		config.MaxConnLifetime = 30 * time.Minute
		config.MaxConnIdleTime = 5 * time.Minute
		config.HealthCheckPeriod = 1 * time.Minute

		pool, err = pgxpool.NewWithConfig(ctx, config)
		if err != nil {
			poolErr = fmt.Errorf("failed to create connection pool: %w", err)
			return
		}

		if err := pool.Ping(ctx); err != nil {
			poolErr = fmt.Errorf("failed to ping database: %w", err)
			pool.Close()
			pool = nil
			return
		}
	})

	return poolErr
}

func GetPool() *pgxpool.Pool {
	return pool
}

func ClosePool() {
	if pool != nil {
		pool.Close()
	}
}
