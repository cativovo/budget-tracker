package internal

import (
	"context"

	"go.uber.org/zap"
)

type ctxKey int

const (
	ctxKeyLogger ctxKey = iota
)

func NewCtxWithLogger(ctx context.Context, z *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, ctxKeyLogger, z)
}

func LoggerFromCtx(ctx context.Context) *zap.SugaredLogger {
	l, ok := ctx.Value(ctxKeyLogger).(*zap.SugaredLogger)
	if !ok {
		panic("LoggerFromCtx: missing or invalid logger in context")
	}
	return l
}
