package sqlite

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	_ "modernc.org/sqlite"
)

type DB struct {
	writer *sqlx.DB
	reader *sqlx.DB
}

func NewDB(ctx context.Context, dbPath string) (*DB, error) {
	const maxIdleTime time.Duration = 3 * time.Minute

	writer, err := connectDB(ctx, dbPath, false)
	if err != nil {
		return nil, fmt.Errorf("sqlite: create writer conn: %w", err)
	}

	writer.SetMaxOpenConns(1)
	writer.SetMaxIdleConns(1)
	writer.SetConnMaxIdleTime(maxIdleTime)

	reader, err := connectDB(ctx, dbPath, true)
	if err != nil {
		return nil, fmt.Errorf("sqlite: create reader conn: %w", err)
	}

	const (
		maxOpenConns int = 120
		maxIdleConns int = 15
	)
	reader.SetMaxOpenConns(maxOpenConns)
	reader.SetMaxIdleConns(maxIdleConns)
	reader.SetConnMaxIdleTime(maxIdleTime)

	return &DB{
		writer: writer,
		reader: reader,
	}, nil
}

func (d *DB) Migrate(ctx context.Context, l *zap.SugaredLogger) error {
	if err := migrate(ctx, d.writer.DB, l); err != nil {
		return fmt.Errorf("sqlite: migrate: %w", err)
	}
	return nil
}

func (d *DB) Close() error {
	if err := d.reader.Close(); err != nil {
		return fmt.Errorf("sqlite: close reader: %w", err)
	}

	if err := d.writer.Close(); err != nil {
		return fmt.Errorf("sqlite: close writer: %w", err)
	}

	return nil
}

func BuildPragmaQuery(p []string) string {
	const pragma = "_pragma"
	for i := range p {
		p[i] = fmt.Sprintf("%s=%s", pragma, p[i])
	}
	return "?" + strings.Join(p, "&")
}

// https://github.com/pocketbase/pocketbase/blob/391287451729ac2f62a4ca596bb923042b76c213/core/db_connect.go#L10
func connectDB(ctx context.Context, dbPath string, readonly bool) (*sqlx.DB, error) {
	// Note: the busy_timeout pragma must be first because
	// the connection needs to be set to block on busy before WAL mode
	// is set in case it hasn't been already set by another connection.
	pragmas := []string{
		"busy_timeout(10000)",
		"journal_mode(WAL)",
		"journal_size_limit(200000000)",
		"synchronous(NORMAL)",
		"foreign_keys(ON)",
		"temp_store(MEMORY)",
		"cache_size(-16000)",
	}
	q := BuildPragmaQuery(pragmas)
	dsn := "file:" + dbPath + q

	if readonly {
		dsn += "&mode=ro"
	}

	return sqlx.ConnectContext(ctx, "sqlite", dsn)
}
