package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New(ctx context.Context, dns string) (*pgxpool.Pool, error) {
	pgxConf, err := pgxpool.ParseConfig(dns)

	if err != nil {
		return nil, fmt.Errorf("parse dns: %w", err)
	}

	pgxConf.MaxConns = 20
	pgxConf.MinConns = 2
	pgxConf.MaxConnLifetime = time.Hour

	pool, err := pgxpool.NewWithConfig(ctx, pgxConf)

	if err != nil {
		return nil, fmt.Errorf("connect db: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	return pool, nil
}
