package sqlite_test

import (
	"context"
	"testing"

	"github.com/cativovo/budget-tracker/internal"
	"github.com/cativovo/budget-tracker/internal/logger"
	"github.com/cativovo/budget-tracker/internal/sqlite"
	"github.com/cativovo/budget-tracker/internal/user"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	dh := newDBHelper(t, "test_create_user.db")
	defer dh.clean()

	ur := sqlite.NewUserRepository(dh.db)
	ctxWithLogger := logger.NewCtxWithLogger(context.Background(), zapLogger)

	tests := map[string]struct {
		input user.CreateUserReq
		want  user.User
		err   error
	}{
		"create user": {
			input: user.CreateUserReq{
				Name:  "Alex Albon",
				ID:    "1",
				Email: "alexalbon@williams.com",
			},
			want: user.User{
				Name:  "Alex Albon",
				ID:    "1",
				Email: "alexalbon@williams.com",
			},
		},
		"duplicate user": {
			input: user.CreateUserReq{
				Name:  "Alex Albon",
				ID:    "1",
				Email: "alexalbon@williams.com",
			},
			err: internal.NewError(internal.ErrorCodeConflict, "User already exists"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := ur.CreateUser(ctxWithLogger, test.input)

			if test.err != nil {
				assert.NotNil(t, err)

				wantCode := internal.GetErrorCode(test.err)
				gotCode := internal.GetErrorCode(err)
				assert.Equal(t, wantCode, gotCode)

				wantMessage := internal.GetErrorMessage(test.err)
				gotMessage := internal.GetErrorMessage(err)
				assert.Equal(t, wantMessage, gotMessage)
				return
			}

			assert.Nil(t, err)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestFindUserByID(t *testing.T) {
	dh := newDBHelper(t, "test_find_user_by_id.db")
	defer dh.clean()

	ur := sqlite.NewUserRepository(dh.db)
	ctxWithLogger := logger.NewCtxWithLogger(context.Background(), zapLogger)

	createUsers(t, ctxWithLogger, dh.db)

	tests := map[string]struct {
		input string
		want  user.User
		err   error
	}{
		"find smooth operator": {
			input: "2",
			want: user.User{
				ID:    "2",
				Name:  "Carlos Sainz Jr.",
				Email: "carlossainzjr@williams.com",
			},
		},
		"find albono": {
			input: "1",
			want: user.User{
				ID:    "1",
				Name:  "Alex Albon",
				Email: "alexalbon@williams.com",
			},
		},
		"user not found": {
			input: "3",
			err:   internal.NewError(internal.ErrorCodeNotFound, "User not found"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := ur.FindUserByID(ctxWithLogger, test.input)

			if test.err != nil {
				assert.NotNil(t, err)

				wantCode := internal.GetErrorCode(test.err)
				gotCode := internal.GetErrorCode(err)
				assert.Equal(t, wantCode, gotCode)

				wantMessage := internal.GetErrorMessage(test.err)
				gotMessage := internal.GetErrorMessage(err)
				assert.Equal(t, wantMessage, gotMessage)
				return
			}

			assert.Nil(t, err)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	dh := newDBHelper(t, "delete_user.db")
	defer dh.clean()

	ur := sqlite.NewUserRepository(dh.db)
	ctxWithLogger := logger.NewCtxWithLogger(context.Background(), zapLogger)

	createUsers(t, ctxWithLogger, dh.db)

	tests := map[string]struct {
		input string
	}{
		"delete albono": {
			input: "1",
		},
		"delete smooth operator": {
			input: "2",
		},
		"user doesn't exist": {
			input: "3",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := ur.DeleteUser(ctxWithLogger, test.input)
			assert.Nil(t, err)

			_, err = ur.FindUserByID(ctxWithLogger, test.input)
			assert.Equal(t, internal.ErrorCodeNotFound, internal.GetErrorCode(err))
			assert.Equal(t, "User not found", internal.GetErrorMessage(err))
		})
	}
}
