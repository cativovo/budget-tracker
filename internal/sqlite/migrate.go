package sqlite

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log/slog"

	"github.com/pressly/goose/v3"
)

type gooseLogger struct {
	logger *slog.Logger
}

func (gl gooseLogger) Printf(format string, v ...any) {
	gl.logger.Info(fmt.Sprintf(format, v...))
}

func (gl gooseLogger) Fatalf(format string, v ...any) {
	gl.logger.Error(fmt.Sprintf(format, v...))
}

//go:embed all:migrations
var embedMigrations embed.FS

func migrate(ctx context.Context, db *sql.DB, l *slog.Logger) error {
	goose.SetLogger(gooseLogger{logger: l})
	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	goose.SetBaseFS(embedMigrations)
	return goose.UpContext(ctx, db, "migrations")
}
