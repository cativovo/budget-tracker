package ctxvalue_test

import (
	"context"
	"testing"

	"github.com/cativovo/budget-tracker/internal/ctxvalue"
	"github.com/cativovo/budget-tracker/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestSetAndGetUserFromContext(t *testing.T) {
	u := model.User{
		ID:    "123",
		Name:  "Yuki Tsunoda",
		Email: "yukitsonoda@redbull.com",
	}
	ctxWithUser := ctxvalue.ContextWithUser(context.Background(), u)
	gotUser := ctxvalue.UserFromContext(ctxWithUser)
	assert.Equal(t, u, gotUser)
}
