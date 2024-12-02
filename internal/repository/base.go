package repository

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	concurrentDB    *sqlx.DB
	nonConcurrentDB *sqlx.DB
}

func NewRepository(dbPath string) (*Repository, error) {
	const (
		maxOpenConns int           = 120
		maxIdleConns int           = 15
		maxIdleTime  time.Duration = 3 * time.Minute
	)

	concurrentDB, err := connectDB(dbPath)
	if err != nil {
		return nil, err
	}

	concurrentDB.SetMaxOpenConns(maxOpenConns)
	concurrentDB.SetMaxIdleConns(maxIdleConns)
	concurrentDB.SetConnMaxIdleTime(maxIdleTime)

	nonConcurrentDB, err := connectDB(dbPath)
	if err != nil {
		return nil, err
	}

	nonConcurrentDB.SetMaxOpenConns(1)
	nonConcurrentDB.SetMaxIdleConns(1)
	nonConcurrentDB.SetConnMaxIdleTime(maxIdleTime)

	return &Repository{
		concurrentDB:    concurrentDB,
		nonConcurrentDB: nonConcurrentDB,
	}, nil
}

func (r *Repository) ConcurrentDB() *sqlx.DB {
	return r.concurrentDB
}

func (r *Repository) NonConcurrentDB() *sqlx.DB {
	return r.nonConcurrentDB
}

func (r *Repository) Close() {
	r.concurrentDB.Close()
	r.nonConcurrentDB.Close()
}
