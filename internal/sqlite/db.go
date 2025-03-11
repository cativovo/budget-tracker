package sqlite

import (
	"embed"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

type DB struct {
	reader *sqlx.DB
	writer *sqlx.DB
}

func NewDB(dbPath string) (*DB, error) {
	const (
		maxOpenConns int           = 120
		maxIdleConns int           = 15
		maxIdleTime  time.Duration = 3 * time.Minute
	)

	reader, err := connectDB(dbPath)
	if err != nil {
		return nil, err
	}

	reader.SetMaxOpenConns(maxOpenConns)
	reader.SetMaxIdleConns(maxIdleConns)
	reader.SetConnMaxIdleTime(maxIdleTime)

	writer, err := connectDB(dbPath)
	if err != nil {
		return nil, err
	}

	writer.SetMaxOpenConns(1)
	writer.SetMaxIdleConns(1)
	writer.SetConnMaxIdleTime(maxIdleTime)

	return &DB{
		reader: reader,
		writer: writer,
	}, nil
}

func (r *DB) Close() {
	r.reader.Close()
	r.writer.Close()
}

//go:embed all:migrations
var embedMigrations embed.FS

func (r *DB) Migrate(logger *zap.SugaredLogger) error {
	goose.SetLogger(gooseLogger{base: logger})
	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	goose.SetBaseFS(embedMigrations)
	if err := goose.Up(r.writer.DB, "migrations"); err != nil {
		return err
	}

	return nil
}
