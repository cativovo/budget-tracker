package user

import (
	"context"

	"github.com/cativovo/budget-tracker/internal"
)

const ContextKeyUser internal.ContextKey = "user"

func ContextWithUser(ctx context.Context, u User) context.Context {
	return context.WithValue(ctx, ContextKeyUser, u)
}

func FromContext(ctx context.Context) User {
	u, ok := ctx.Value(ContextKeyUser).(User)
	if !ok {
		panic("UserFromContext: missing or invalid user in context")
	}
	return u
}
