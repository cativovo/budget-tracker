package user

import (
	"context"
)

type ctxKey string

const ctxKeyUser ctxKey = "user"

// WithContext returns a context with the given user.
func WithContext(ctx context.Context, u User) context.Context {
	return context.WithValue(ctx, ctxKeyUser, u)
}

// FromContext extracts the user from context.
// Panics if there's no user.
func FromContext(ctx context.Context) User {
	v, ok := ctx.Value(ctxKeyUser).(User)
	if !ok {
		panic("missing user in context; possible middleware bypass")
	}
	return v
}
