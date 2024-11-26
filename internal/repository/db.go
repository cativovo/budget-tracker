package repository

import (
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

// https://github.com/pocketbase/pocketbase/blob/0ac4a388c000f98bcec830de6819b1776d0ee242/core/db_connect.go#L10
func connectDB(dbPath string) (*sqlx.DB, error) {
	// Note: the busy_timeout pragma must be first because
	// the connection needs to be set to block on busy before WAL mode
	// is set in case it hasn't been already set by another connection.
	pragmas := "?_pragma=busy_timeout(10000)&_pragma=journal_mode(WAL)&_pragma=journal_size_limit(200000000)&_pragma=synchronous(NORMAL)&_pragma=foreign_keys(ON)&_pragma=temp_store(MEMORY)&_pragma=cache_size(-16000)"
	dsn := dbPath + pragmas

	return sqlx.Open("sqlite", dsn)
}
