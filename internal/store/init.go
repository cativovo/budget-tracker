package store

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	"github.com/cativovo/budget-tracker/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed all:migrations
var embedMigrations embed.FS

func migrate(db *sql.DB) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	goose.SetBaseFS(embedMigrations)
	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}

	return nil
}

func InitDB(ctx context.Context, cfg config.DBConfig) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB, cfg.SSL,
	)
	dbpool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}

	db := stdlib.OpenDBFromPool(dbpool)
	if err := migrate(db); err != nil {
		return nil, err
	}

	return dbpool, nil
}
