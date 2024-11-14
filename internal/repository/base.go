package repository

import (
	"context"
	"database/sql"

	"github.com/cativovo/budget-tracker/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

type Repository struct {
	dbPool *pgxpool.Pool
}

func NewRepository(ctx context.Context, cfg config.DBConfig) (*Repository, error) {
	dbPool, err := newDBPool(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &Repository{
		dbPool: dbPool,
	}, nil
}

func (r *Repository) DBPool() *pgxpool.Pool {
	return r.dbPool
}

func (r *Repository) Close() {
	r.dbPool.Close()
}

func (r *Repository) OpenDBFromPool() *sql.DB {
	return stdlib.OpenDBFromPool(r.dbPool)
}
