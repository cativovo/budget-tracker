package internal_test

import (
	"context"
	"testing"

	"github.com/cativovo/budget-tracker/internal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestSetAndGetUserFromContext(t *testing.T) {
	u := internal.User{
		ID:    "123",
		Name:  "Yuki Tsunoda",
		Email: "yukitsonoda@redbull.com",
	}
	ctxWithUser := internal.ContextWithUser(context.Background(), u)
	gotUser := internal.UserFromContext(ctxWithUser)
	assert.Equal(t, u, gotUser)
}

func TestSetAndGetLoggerFromContext(t *testing.T) {
	l := zap.NewNop().Sugar()
	ctxWithLogger := internal.ContextWithLogger(context.Background(), l)
	gotLogger := internal.LoggerFromContext(ctxWithLogger)
	assert.Equal(t, l, gotLogger)
}
