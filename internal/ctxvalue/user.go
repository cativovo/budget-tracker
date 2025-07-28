package ctxvalue

import (
	"context"

	"github.com/cativovo/budget-tracker/internal/domain"
)

const contextKeyUser ContextKey = "user"

// ContextWithUser returns a context with the given User
func ContextWithUser(ctx context.Context, u domain.User) context.Context {
	return context.WithValue(ctx, contextKeyUser, u)
}

// UserFromContext extracts User from context
func UserFromContext(ctx context.Context) domain.User {
	u, ok := ctx.Value(contextKeyUser).(domain.User)
	if !ok {
		// should panic because auth middleware should put User in the context
		// and prevent from continuing if unauthenticated
		panic("missing or invalid User from context")
	}
	return u
}
