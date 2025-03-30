package logger

import (
	"context"

	"github.com/cativovo/budget-tracker/internal"
	"go.uber.org/zap"
)

const CtxKeyLogger internal.CtxKey = "logger"

func NewCtxWithLogger(ctx context.Context, z *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, CtxKeyLogger, z)
}

func FromContext(ctx context.Context) *zap.SugaredLogger {
	l, ok := ctx.Value(CtxKeyLogger).(*zap.SugaredLogger)
	if !ok {
		panic("LoggerFromCtx: missing or invalid logger in context")
	}
	return l
}
