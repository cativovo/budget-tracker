package repository

import (
	"embed"

	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

//go:embed all:migrations
var embedMigrations embed.FS

type gooseLogger struct {
	l *zap.SugaredLogger
}

func (g gooseLogger) Printf(format string, v ...interface{}) {
	g.l.Infof(format, v...)
}

func (g gooseLogger) Fatalf(format string, v ...interface{}) {
	g.l.Fatalf(format, v...)
}

var _ goose.Logger = (*gooseLogger)(nil)

func (r *Repository) Migrate() error {
	logger := gooseLogger{l: r.logger}
	goose.SetLogger(logger)
	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	goose.SetBaseFS(embedMigrations)
	if err := goose.Up(r.NonConcurrentDB().DB, "migrations"); err != nil {
		return err
	}

	return nil
}
