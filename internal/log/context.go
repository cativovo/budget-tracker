package log

import (
	"context"
	"log/slog"
)

// The provided key must be comparable and should not be of type string
// or any other built-in type to avoid collisions between packages using context.
// Users of WithValue should define their own types for keys.
type contextKey string

const contextKeyLogger contextKey = "logger"

// WithContext returns a context with the given logger
func WithContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, contextKeyLogger, logger)
}

// FromContext extracts logger from context
func FromContext(ctx context.Context) *slog.Logger {
	l, ok := ctx.Value(contextKeyLogger).(*slog.Logger)
	if !ok {
		return slog.Default()
	}
	return l
}
