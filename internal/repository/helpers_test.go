package repository_test

import (
	"os"
	"testing"
	"time"

	"github.com/cativovo/budget-tracker/internal/repository"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type repositoryHelper struct {
	r      *repository.Repository
	t      *testing.T
	dbPath string
}

var logger = zap.NewNop().Sugar()

func newRepositoryHelper(t *testing.T, dbPath string) *repositoryHelper {
	t.Helper()

	r, err := repository.NewRepository(dbPath)
	if err != nil {
		t.Fatal(err)
	}

	if err := r.Migrate(logger); err != nil {
		t.Fatal(err)
	}

	return &repositoryHelper{
		t:      t,
		dbPath: dbPath,
		r:      r,
	}
}

func (r *repositoryHelper) repository() *repository.Repository {
	r.t.Helper()
	return r.r
}

func (r *repositoryHelper) clean() {
	r.t.Helper()
	r.r.Close()
	if err := os.RemoveAll(r.dbPath); err != nil {
		r.t.Fatal(err)
	}
}

func assertEqualPointer[T any](t *testing.T, want, got *T) {
	t.Helper()

	if want != nil && got != nil {
		assert.Equal(t, *want, *got)
		return
	}

	assert.Equal(t, want, got)
}

func assertTimestampWithin(t *testing.T, want time.Time, got string, d time.Duration) {
	t.Helper()

	parsed, err := time.Parse(time.DateTime, got)
	if err != nil {
		t.Fatal(err)
	}
	assert.WithinDuration(t, want, parsed, d)
}
