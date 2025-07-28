package sqlite

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

type DB struct {
	Writer *sqlx.DB
	Reader *sqlx.DB
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
		Writer: writer,
		Reader: reader,
	}, nil
}

func (d *DB) Migrate(ctx context.Context, l *zap.SugaredLogger) error {
	if err := migrate(ctx, d.Writer.DB, l); err != nil {
		return fmt.Errorf("sqlite: migrate: %w", err)
	}
	return nil
}

func (d *DB) Close() error {
	if err := d.Reader.Close(); err != nil {
		return fmt.Errorf("sqlite: close reader: %w", err)
	}

	if err := d.Writer.Close(); err != nil {
		return fmt.Errorf("sqlite: close writer: %w", err)
	}

	return nil
}

// https://github.com/pocketbase/pocketbase/blob/391287451729ac2f62a4ca596bb923042b76c213/core/db_connect.go#L10
func setPragmas(ctx context.Context, db *sqlx.DB) error {
	// Note: the busy_timeout pragma must be first because
	// the connection needs to be set to block on busy before WAL mode
	// is set in case it hasn't been already set by another connection.
	q := `
		PRAGMA busy_timeout       = 10000;
		PRAGMA journal_mode       = WAL;
		PRAGMA journal_size_limit = 200000000;
		PRAGMA synchronous        = NORMAL;
		PRAGMA foreign_keys       = ON;
		PRAGMA temp_store         = MEMORY;
		PRAGMA cache_size         = -16000;
	`

	_, err := db.ExecContext(ctx, q)
	if err != nil {
		return fmt.Errorf("set pragma: %w", err)
	}

	return nil
}

func connectDB(ctx context.Context, dbPath string, readonly bool) (*sqlx.DB, error) {
	q := make(url.Values)

	if readonly {
		q.Set("mode", "ro")
	} else {
		q.Set("_txlock", "immediate")
	}

	dsn := fmt.Sprintf("file:%s?%s", dbPath, q.Encode())
	db, err := sqlx.ConnectContext(ctx, "sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("connect to db: %w", err)
	}

	if err := setPragmas(ctx, db); err != nil {
		return nil, err
	}

	return db, nil
}
