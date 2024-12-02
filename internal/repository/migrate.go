package repository

import (
	"embed"

	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

//go:embed all:migrations
var embedMigrations embed.FS

type gooseLogger struct {
	base *zap.SugaredLogger
}

func (g gooseLogger) Printf(format string, v ...interface{}) {
	g.base.Infof(format, v...)
}

func (g gooseLogger) Fatalf(format string, v ...interface{}) {
	g.base.Fatalf(format, v...)
}

var _ goose.Logger = (*gooseLogger)(nil)

func (r *Repository) Migrate(logger *zap.SugaredLogger) error {
	goose.SetLogger(gooseLogger{base: logger})
	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	goose.SetBaseFS(embedMigrations)
	if err := goose.Up(r.NonConcurrentDB().DB, "migrations"); err != nil {
		return err
	}

	return nil
}
