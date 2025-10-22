package user_test

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/cativovo/budget-tracker/internal/user"
	"github.com/stretchr/testify/assert"
)

func TestWithAndFromContext(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		wantUser := user.User{
			ID:        gofakeit.UUID(),
			Name:      gofakeit.Name(),
			Email:     gofakeit.Email(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		ctx := user.WithContext(context.Background(), wantUser)
		gotUser := user.FromContext(ctx)
		assert.Equal(t, wantUser, gotUser)
	})
	t.Run("panic", func(t *testing.T) {
		defer func() {
			msg := recover()
			assert.Equal(t, "missing user in context; possible middleware bypass", msg)
		}()
		user.FromContext(context.Background())
	})
}
