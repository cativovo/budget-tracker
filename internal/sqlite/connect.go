package sqlite

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

func buildPragmaQuery(p []string) string {
	const pragma = "_pragma"
	for i := range p {
		p[i] = fmt.Sprintf("%s=%s", pragma, p[i])
	}
	return "?" + strings.Join(p, "&")
}

// https://github.com/pocketbase/pocketbase/blob/391287451729ac2f62a4ca596bb923042b76c213/core/db_connect.go#L10
func connectDB(dbPath string) (*sqlx.DB, error) {
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
	q := buildPragmaQuery(pragmas)
	dsn := dbPath + q
	return sqlx.Open("sqlite", dsn)
}
