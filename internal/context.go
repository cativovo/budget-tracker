package internal

import (
	"context"

	"go.uber.org/zap"
)

// ContextKey defines a context key type
type ContextKey string

const contextKeyUser ContextKey = "user"

// ContextWithUser returns a context with the given User
func ContextWithUser(ctx context.Context, u User) context.Context {
	return context.WithValue(ctx, contextKeyUser, u)
}

// UserFromContext extracts User from context
func UserFromContext(ctx context.Context) User {
	u, ok := ctx.Value(contextKeyUser).(User)
	if !ok {
		// should panic because auth middleware should put User in the context
		// and prevent from continuing if unauthenticated
		panic("missing or invalid User from context")
	}
	return u
}

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
