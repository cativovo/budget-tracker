package ctxvalue

import (
	"context"

	"github.com/cativovo/budget-tracker/internal/model"
)

const contextKeyUser ContextKey = "user"

// ContextWithUser returns a context with the given User
func ContextWithUser(ctx context.Context, u model.User) context.Context {
	return context.WithValue(ctx, contextKeyUser, u)
}

// UserFromContext extracts User from context
func UserFromContext(ctx context.Context) model.User {
	u, ok := ctx.Value(contextKeyUser).(model.User)
	if !ok {
		// should panic because auth middleware should put User in the context
		// and prevent from continuing if unauthenticated
		panic("missing or invalid User from context")
	}
	return u
}
