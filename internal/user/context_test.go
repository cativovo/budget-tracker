package user_test

import (
	"context"
	"testing"

	"github.com/cativovo/budget-tracker/internal/user"
	"github.com/stretchr/testify/assert"
)

func TestSetAndGetUserFromContext(t *testing.T) {
	u := user.User{
		ID:    "123",
		Name:  "Yuki Tsunoda",
		Email: "yukitsonoda@redbull.com",
	}

	ctxWithUser := user.ContextWithUser(context.Background(), u)
	gotUser := user.FromContext(ctxWithUser)

	assert.Equal(t, u, gotUser)
}
