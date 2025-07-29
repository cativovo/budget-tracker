package sqlite_test

import (
	"context"
	"testing"
	"time"

	"github.com/cativovo/budget-tracker/internal/log"
	"github.com/cativovo/budget-tracker/internal/sqlite"
	"github.com/cativovo/budget-tracker/internal/testutil"
	"github.com/cativovo/budget-tracker/internal/user"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestUserStore_CreateUser_GetUser(t *testing.T) {
	db, c := mustConnectDB(t)
	defer c()

	us := sqlite.NewUserStore(db)

	ctx := log.WithContext(context.Background(), zap.NewNop().Sugar())
	now := time.Now()

	tests := []struct {
		name    string
		u       user.UserCreate
		want    user.User
		wantErr error
	}{
		{
			name: "create user",
			u: user.UserCreate{
				ID:    "1",
				Name:  "Gregorio del Pilar",
				Email: "g.delpilar@tiradpass.ph",
			},
			want: user.User{
				ID:        "1",
				Name:      "Gregorio del Pilar",
				Email:     "g.delpilar@tiradpass.ph",
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			created, gotErr := us.CreateUser(ctx, tt.u)
			if tt.wantErr != nil {
				return
			}

			assertUser(t, tt.want, created)
			assert.Nil(t, gotErr)

			got, gotErr := us.GetUser(ctx, created.ID)
			assert.Equal(t, created, got)
			assert.Nil(t, gotErr)
			testutil.AssertWithinDuration(t, now, got.CreatedAt)
			testutil.AssertWithinDuration(t, now, got.UpdatedAt)
		})
	}
}

func assertUser(t *testing.T, want, got user.User) {
	testutil.AssertWithinDuration(t, want.CreatedAt, got.CreatedAt)
	testutil.AssertWithinDuration(t, want.UpdatedAt, got.UpdatedAt)

	got.CreatedAt = time.Time{}
	got.UpdatedAt = time.Time{}
	want.CreatedAt = time.Time{}
	want.UpdatedAt = time.Time{}

	assert.Equal(t, want, got)
}
