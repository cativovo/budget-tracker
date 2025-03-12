package user

import (
	"context"

	"github.com/cativovo/budget-tracker/internal"
)

const CtxKeyUser internal.CtxKey = "user"

func NewCtxWithUser(ctx context.Context, u User) context.Context {
	return context.WithValue(ctx, CtxKeyUser, u)
}

func UserFromCtx(ctx context.Context) User {
	u, ok := ctx.Value(CtxKeyUser).(User)
	if !ok {
		panic("UserFromContext: missing or invalid user in context")
	}
	return u
}
