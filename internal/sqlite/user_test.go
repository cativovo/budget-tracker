package sqlite_test

import (
	"context"
	"testing"
	"time"

	"github.com/cativovo/budget-tracker/internal/apperror"
	"github.com/cativovo/budget-tracker/internal/sqlite"
	"github.com/cativovo/budget-tracker/internal/testutil"
	"github.com/cativovo/budget-tracker/internal/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserStore_GetUserByID(t *testing.T) {
	testutil.DisableLogs()
	db, c := newTestDB(t)
	defer c()

	us := sqlite.NewUserStore(db)

	tests := []struct {
		name            string
		id              string
		createUserInput user.CreateUserInput
		wantErr         error
	}{
		{
			name: "ok",
			id:   "1",
			createUserInput: user.CreateUserInput{
				ID:    "1",
				Name:  "Juan Usa",
				Email: "juanusa@email.com",
			},
		},
		{
			name:    "not found",
			id:      "2",
			wantErr: apperror.New(apperror.ErrorCodeNotFound, "user not found"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr != nil {
				got, err := us.GetUserByID(context.Background(), tt.id)
				testutil.RequireEqualError(t, tt.wantErr, err)
				assert.Zero(t, got)
				return
			}

			want := mustCreateUser(t, db, tt.createUserInput)
			got, err := us.GetUserByID(context.Background(), tt.id)
			require.NoError(t, err)

			assertUser(t, want, got)
		})
	}
}

func TestUserStore_CreateUser(t *testing.T) {
	testutil.DisableLogs()
	db, c := newTestDB(t)
	defer c()

	us := sqlite.NewUserStore(db)

	tests := []struct {
		name            string
		createUserInput user.CreateUserInput
		wantErr         error
	}{
		{
			name: "ok",
			createUserInput: user.CreateUserInput{
				ID:    "1",
				Name:  "Juan Usa",
				Email: "juanusa@email.com",
			},
		},
		{
			name: "duplicate email",
			createUserInput: user.CreateUserInput{
				ID:    "2",
				Name:  "John Doe",
				Email: "juanusa@email.com",
			},
			wantErr: apperror.New(apperror.ErrorCodeConflict, "email already taken"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := us.CreateUser(context.Background(), tt.createUserInput)

			if tt.wantErr != nil {
				testutil.RequireEqualError(t, tt.wantErr, err)
				assert.Zero(t, got)

				want, err := us.GetUserByID(context.Background(), got.ID)
				assert.Zero(t, want)
				assert.Error(t, err)
				return
			}

			wantUser, err := us.GetUserByID(context.Background(), got.ID)
			require.NoError(t, err)

			assertUser(t, wantUser, got)
		})
	}
}

func assertUser(t *testing.T, want, got user.User) bool {
	testutil.AssertWithinDuration(t, want.CreatedAt, got.CreatedAt)
	testutil.AssertWithinDuration(t, want.UpdatedAt, got.UpdatedAt)

	got.CreatedAt = time.Time{}
	got.UpdatedAt = time.Time{}
	want.CreatedAt = time.Time{}
	want.UpdatedAt = time.Time{}

	return assert.Equal(t, want, got)
}
