package user

import (
	"context"
)

// The provided key must be comparable and should not be of type string
// or any other built-in type to avoid collisions between packages using context.
// Users of WithValue should define their own types for keys.
type contextKey string

const contextKeyUser contextKey = "user"

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
