package log

import (
	"context"

	"github.com/cativovo/budget-tracker/internal"
	"go.uber.org/zap"
)

const contextKeyLogger internal.ContextKey = "logger"

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
