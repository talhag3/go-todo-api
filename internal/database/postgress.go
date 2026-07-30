package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func New(ctx context.Context, dns string) (*pgxpool.Pool, *sql.DB, error) {
	pgxConf, err := pgxpool.ParseConfig(dns)

	if err != nil {
		return nil, nil, fmt.Errorf("parse dns: %w", err)
	}

	pgxConf.MaxConns = 20
	pgxConf.MinConns = 2
	pgxConf.MaxConnLifetime = time.Hour

	pool, err := pgxpool.NewWithConfig(ctx, pgxConf)

	if err != nil {
		return nil, nil, fmt.Errorf("connect db: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, nil, fmt.Errorf("ping db: %w", err)
	}

	sqlDB := stdlib.OpenDBFromPool(pool)

	return pool, sqlDB, nil
}

func RunMigration(db *sql.DB, dir string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	return goose.Up(db, dir)
}
