package sqlite

import (
	"context"
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

type gooseLogger struct {
	base *zap.SugaredLogger
}

func (g gooseLogger) Printf(format string, v ...any) {
	g.base.Infof(format, v...)
}

func (g gooseLogger) Fatalf(format string, v ...any) {
	g.base.Fatalf(format, v...)
}

//go:embed all:migrations
var embedMigrations embed.FS

func migrate(ctx context.Context, db *sql.DB, l *zap.SugaredLogger) error {
	goose.SetLogger(gooseLogger{base: l})
	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	goose.SetBaseFS(embedMigrations)
	return goose.UpContext(ctx, db, "migrations")
}
