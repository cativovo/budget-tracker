package repository

import (
	"embed"

	"github.com/pressly/goose/v3"
)

//go:embed all:migrations
var embedMigrations embed.FS

func (r *Repository) Migrate() error {
	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	goose.SetBaseFS(embedMigrations)
	if err := goose.Up(r.NonConcurrentDB().DB, "migrations"); err != nil {
		return err
	}

	return nil
}
