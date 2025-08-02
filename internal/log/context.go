package log

import (
	"context"

	"go.uber.org/zap"
)

// The provided key must be comparable and should not be of type string
// or any other built-in type to avoid collisions between packages using context.
// Users of WithValue should define their own types for keys.
type contextKey string

const contextKeyLogger contextKey = "logger"

// WithContext returns a context with the given logger
func WithContext(
	ctx context.Context,
	logger *zap.SugaredLogger,
) context.Context {
	return context.WithValue(ctx, contextKeyLogger, logger)
}

// FromContext extracts logger from context
func FromContext(ctx context.Context) *zap.SugaredLogger {
	l, ok := ctx.Value(contextKeyLogger).(*zap.SugaredLogger)
	if !ok {
		panic("missing or invalid logger from context")
	}
	return l
}
