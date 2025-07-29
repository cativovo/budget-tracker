package user

import (
	"context"

	"github.com/cativovo/budget-tracker/internal"
)

const contextKeyUser internal.ContextKey = "user"

// WithContext returns a context with the given User
func WithContext(ctx context.Context, u User) context.Context {
	return context.WithValue(ctx, contextKeyUser, u)
}

// FromContext extracts User from context
func FromContext(ctx context.Context) User {
	u, ok := ctx.Value(contextKeyUser).(User)
	if !ok {
		panic("missing or invalid User from context")
	}
	return u
}
