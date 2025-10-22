package log_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/cativovo/budget-tracker/internal/log"
	"github.com/stretchr/testify/assert"
)

func TestWithAndFromContext(t *testing.T) {
	t.Run("custom logger", func(t *testing.T) {
		want := slog.New(slog.DiscardHandler)
		ctx := log.WithContext(context.Background(), want)
		got := log.FromContext(ctx)
		assert.Equal(t, want, got)
	})

	t.Run("default logger", func(t *testing.T) {
		want := slog.Default()
		got := log.FromContext(context.Background())
		assert.Equal(t, want, got)
	})
}
