package logger

import (
	"context"

	"github.com/cativovo/budget-tracker/internal"
	"go.uber.org/zap"
)

const ContextKeyLogger internal.ContextKey = "logger"

func ContextWithLogger(ctx context.Context, z *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, ContextKeyLogger, z)
}

func FromContext(ctx context.Context) *zap.SugaredLogger {
	l, ok := ctx.Value(ContextKeyLogger).(*zap.SugaredLogger)
	if !ok {
		panic("LoggerFromCtx: missing or invalid logger in context")
	}
	return l
}
