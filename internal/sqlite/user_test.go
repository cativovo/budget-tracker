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
				gotUser, err := us.GetUserByID(context.Background(), tt.id)
				testutil.RequireEqualError(t, tt.wantErr, err)
				assert.Zero(t, gotUser)
				return
			}

			wantUser := mustCreateUser(t, db, tt.createUserInput)
			gotUser, err := us.GetUserByID(context.Background(), tt.id)
			require.NoError(t, err)

			assertUser(t, wantUser, gotUser)
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
			gotUser, err := us.CreateUser(context.Background(), tt.createUserInput)

			if tt.wantErr != nil {
				testutil.RequireEqualError(t, tt.wantErr, err)
				assert.Zero(t, gotUser)

				wantUser, err := us.GetUserByID(context.Background(), gotUser.ID)
				assert.Zero(t, wantUser)
				assert.Error(t, err)
				return
			}

			wantUser, err := us.GetUserByID(context.Background(), gotUser.ID)
			require.NoError(t, err)

			assertUser(t, wantUser, gotUser)
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

func mustCreateUser(t *testing.T, db *sqlite.DB, input user.CreateUserInput) user.User {
	t.Helper()
	us := sqlite.NewUserStore(db)
	u, err := us.CreateUser(context.Background(), input)
	require.NoError(t, err)
	return u
}
