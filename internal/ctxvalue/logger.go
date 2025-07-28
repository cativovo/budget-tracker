package ctxvalue

import (
	"context"

	"go.uber.org/zap"
)

const contextKeyLogger ContextKey = "logger"

// ContextWithLogger returns a context with the given logger
func ContextWithLogger(
	ctx context.Context,
	z *zap.SugaredLogger,
) context.Context {
	return context.WithValue(ctx, contextKeyLogger, z)
}

// LoggerFromContext extracts logger from context
func LoggerFromContext(ctx context.Context) *zap.SugaredLogger {
	l, ok := ctx.Value(contextKeyLogger).(*zap.SugaredLogger)
	if !ok {
		panic("missing or invalid logger from context")
	}
	return l
}
