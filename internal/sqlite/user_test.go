package sqlite_test

import (
	"context"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/cativovo/budget-tracker/internal"
	"github.com/cativovo/budget-tracker/internal/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestCreateUser(t *testing.T) {
	db, c := mustConnectDB(t)
	defer c()

	us := sqlite.NewUserService(db)
	sharedCtx := internal.ContextWithLogger(context.Background(), zap.NewNop().Sugar())

	testCases := []struct {
		name  string
		input internal.UserCreate
		err   error
	}{
		{
			name: "ok 1",
			input: internal.UserCreate{
				ID:    "1",
				Name:  "Alex Albon",
				Email: "alexalbon@williams.com",
			},
		},
		{
			name: "ok 2",
			input: internal.UserCreate{
				ID:    "2",
				Name:  "Carlos Sainz",
				Email: "carlossainz@williams.com",
			},
		},
		{
			name: "no id",
			input: internal.UserCreate{
				Name:  "Isack Hadjar",
				Email: "isackhadjar@racingbulls.com",
			},
			err: internal.NewError(internal.ErrorCodeInvalid, "id is required"),
		},
		{
			name: "no name",
			input: internal.UserCreate{
				ID:    "1",
				Name:  "",
				Email: "isackhadjar@racingbulls.com",
			},
			err: internal.NewError(internal.ErrorCodeInvalid, "name is required"),
		},
		{
			name: "no email",
			input: internal.UserCreate{
				ID:   "1",
				Name: "Isack Hadjar",
			},
			err: internal.NewError(internal.ErrorCodeInvalid, "invalid email"),
		},
		{
			name: "invalid email",
			input: internal.UserCreate{
				ID:    "1",
				Name:  "Isack Hadjar",
				Email: "isackhadjar",
			},
			err: internal.NewError(internal.ErrorCodeInvalid, "invalid email"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.err != nil {
				createdUser, err := us.CreateUser(sharedCtx, tc.input)
				assert.Equal(t, internal.GetErrorCode(tc.err), internal.GetErrorCode(err))
				assert.Equal(t, internal.GetErrorMessage(tc.err), internal.GetErrorMessage(err))

				_, err = us.GetUser(sharedCtx, createdUser.ID)
				assert.Equal(t, internal.ErrorCodeNotFound, internal.GetErrorCode(err))
				assert.Equal(t, "user not found", internal.GetErrorMessage(err))
				return
			}

			createdUser, err := us.CreateUser(sharedCtx, tc.input)
			assert.Nil(t, err)
			assert.Equal(t, tc.input.ID, createdUser.ID)
			assert.Equal(t, tc.input.Name, createdUser.Name)
			assert.Equal(t, tc.input.Email, createdUser.Email)
			assert.WithinDuration(t, time.Now(), createdUser.CreatedAt, wantDelta)
			assert.WithinDuration(t, time.Now(), createdUser.UpdatedAt, wantDelta)

			foundUser, err := us.GetUser(sharedCtx, createdUser.ID)
			assert.Nil(t, err)
			assert.Equal(t, createdUser, foundUser)
		})
	}

	t.Run("validate all fields", func(t *testing.T) {
		createdUser, err := us.CreateUser(sharedCtx, internal.UserCreate{})
		assert.Equal(t, internal.ErrorCodeInvalid, internal.GetErrorCode(err))

		gotMsgs := strings.Split(internal.GetErrorMessage(err), ", ")
		wantMsgs := []string{
			"id is required",
			"name is required",
			"invalid email",
		}
		slices.Sort(gotMsgs)
		slices.Sort(wantMsgs)
		assert.Equal(t, wantMsgs, gotMsgs)

		_, err = us.GetUser(sharedCtx, createdUser.ID)
		assert.Equal(t, internal.ErrorCodeNotFound, internal.GetErrorCode(err))
		assert.Equal(t, "user not found", internal.GetErrorMessage(err))
	})

	t.Run("conflict", func(t *testing.T) {
		mustCreateUser(sharedCtx, t, db, internal.UserCreate{
			ID:    "9",
			Name:  "Liam Lawson",
			Email: "liamlawson@racingbulls.com",
		})

		createdUser, err := us.CreateUser(sharedCtx, internal.UserCreate{
			ID:    "10",
			Name:  "Isack Hadjar",
			Email: "liamlawson@racingbulls.com",
		})
		assert.Equal(t, internal.ErrorCodeConflict, internal.GetErrorCode(err))
		assert.Equal(t, "email already taken", internal.GetErrorMessage(err))

		_, err = us.GetUser(sharedCtx, createdUser.ID)
		assert.Equal(t, internal.ErrorCodeNotFound, internal.GetErrorCode(err))
		assert.Equal(t, "user not found", internal.GetErrorMessage(err))
	})
}

func TestUpdateUser(t *testing.T) {
	db, c := mustConnectDB(t)
	defer c()

	us := sqlite.NewUserService(db)
	sharedCtx := internal.ContextWithLogger(context.Background(), zap.NewNop().Sugar())

	testCases := []struct {
		name            string
		createUserInput internal.UserCreate
		updateUserInput internal.UserUpdate
		err             error
	}{
		{
			name: "update name",
			createUserInput: internal.UserCreate{
				ID:    "1",
				Name:  "Alex Albon",
				Email: "alexalbon@williams.com",
			},
			updateUserInput: internal.UserUpdate{
				Name: ptr("Alex Albono"),
			},
		},
		{
			name: "update email",
			createUserInput: internal.UserCreate{
				ID:    "2",
				Name:  "Carlos Sainz",
				Email: "carlossainz@ferrari.com",
			},
			updateUserInput: internal.UserUpdate{
				Email: ptr("carlossainz@williams.com"),
			},
		},
		{
			name: "update all",
			createUserInput: internal.UserCreate{
				ID:    "3",
				Name:  "Logans Sargeant",
				Email: "loganssargeant@ferrari.com",
			},
			updateUserInput: internal.UserUpdate{
				Name:  ptr("Logan Sargeant"),
				Email: ptr("logansargeant@ferrari.com"),
			},
		},
		{
			name: "empty name",
			createUserInput: internal.UserCreate{
				ID:    "4",
				Name:  "Liam Lawson",
				Email: "liamlawson@racingbulls.com",
			},
			updateUserInput: internal.UserUpdate{
				Name: ptr(""),
			},
			err: internal.NewError(internal.ErrorCodeInvalid, "name is required"),
		},
		{
			name: "empty email",
			createUserInput: internal.UserCreate{
				ID:    "5",
				Name:  "Lance Stroll",
				Email: "lancestroll@astonmartin.com",
			},
			updateUserInput: internal.UserUpdate{
				Email: ptr(""),
			},
			err: internal.NewError(internal.ErrorCodeInvalid, "invalid email"),
		},
		{
			name: "invalid email",
			createUserInput: internal.UserCreate{
				ID:    "6",
				Name:  "Fernando Alonso",
				Email: "fernandoalonso@astonmartin.com",
			},
			updateUserInput: internal.UserUpdate{
				Email: ptr("nando@"),
			},
			err: internal.NewError(internal.ErrorCodeInvalid, "invalid email"),
		},
		{
			name: "no update",
			createUserInput: internal.UserCreate{
				ID:    "7",
				Name:  "Nico Hulkenberg",
				Email: "nicohulkenberg@audi.com",
			},
			updateUserInput: internal.UserUpdate{},
			err:             internal.NewError(internal.ErrorCodeInvalid, "no update fields provided"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			createdUser := mustCreateUser(sharedCtx, t, db, tc.createUserInput)
			ctx := internal.ContextWithUser(sharedCtx, createdUser)

			if tc.err != nil {
				_, err := us.UpdateUser(ctx, tc.updateUserInput)
				assert.Equal(t, internal.GetErrorCode(tc.err), internal.GetErrorCode(err))
				assert.Equal(t, internal.GetErrorMessage(tc.err), internal.GetErrorMessage(err))

				foundUser := mustGetUser(ctx, t, db, createdUser.ID)
				assert.Equal(t, createdUser, foundUser)
				return
			}

			updatedUser, err := us.UpdateUser(ctx, tc.updateUserInput)
			assert.Nil(t, err)
			assert.NotEqual(t, createdUser, updatedUser)

			assertUpdatedField(t, tc.updateUserInput.Name, createdUser.Name, updatedUser.Name)
			assertUpdatedField(t, tc.updateUserInput.Email, createdUser.Email, updatedUser.Email)
			assert.Equal(t, createdUser.CreatedAt, updatedUser.UpdatedAt)

			// TODO: find a way to test UpdatedAt field without slowing down the tests
			// it's too fast that the UpdatedAt is still the same with CreatedAt

			foundUser := mustGetUser(ctx, t, db, updatedUser.ID)
			assert.Equal(t, updatedUser, foundUser)
		})
	}

	t.Run("update correct user", func(t *testing.T) {
		createdUser1 := mustCreateUser(sharedCtx, t, db, internal.UserCreate{
			ID:    "9",
			Name:  "Oliver Bearman",
			Email: "oliverbearman@haas.com",
		})
		createdUser2 := mustCreateUser(sharedCtx, t, db, internal.UserCreate{
			ID:    "10",
			Name:  "Esteban Ocon",
			Email: "estebanocon@haas.com",
		})

		ctx := internal.ContextWithUser(sharedCtx, createdUser1)
		updatedUser, err := us.UpdateUser(ctx, internal.UserUpdate{
			Email: ptr("manbearpig@haas.com"),
		})
		assert.Nil(t, err)
		assert.NotEqual(t, createdUser1, updatedUser)

		foundUser := mustGetUser(ctx, t, db, createdUser2.ID)
		assert.Equal(t, createdUser2, foundUser)
	})

	t.Run("validate all fields", func(t *testing.T) {
		createdUser := mustCreateUser(sharedCtx, t, db, internal.UserCreate{
			ID:    "8",
			Name:  "Gabriel Bortoleto",
			Email: "gabrielbortoleto@audi.com",
		})
		ctx := internal.ContextWithUser(sharedCtx, createdUser)
		_, err := us.UpdateUser(ctx, internal.UserUpdate{
			Name:  ptr(""),
			Email: ptr("test"),
		})
		assert.Equal(t, internal.ErrorCodeInvalid, internal.GetErrorCode(err))

		gotMsgs := strings.Split(internal.GetErrorMessage(err), ", ")
		wantMsgs := []string{
			"name is required",
			"invalid email",
		}
		slices.Sort(gotMsgs)
		slices.Sort(wantMsgs)
		assert.Equal(t, wantMsgs, gotMsgs)
	})

	t.Run("conflict", func(t *testing.T) {
		mustCreateUser(sharedCtx, t, db, internal.UserCreate{
			ID:    "M-1",
			Name:  "George Russell",
			Email: "georgerussel@mercedes.com",
		})
		createdUser := mustCreateUser(sharedCtx, t, db, internal.UserCreate{
			ID:    "M-2",
			Name:  "Andrea Kimi Antonelli",
			Email: "andreakimiantonelli@mercedes.com",
		})
		ctx := internal.ContextWithUser(sharedCtx, createdUser)
		_, err := us.UpdateUser(ctx, internal.UserUpdate{
			Email: ptr("georgerussel@mercedes.com"),
		})
		assert.Equal(t, internal.ErrorCodeConflict, internal.GetErrorCode(err))
		assert.Equal(t, "email already taken", internal.GetErrorMessage(err))

		foundUser := mustGetUser(ctx, t, db, createdUser.ID)
		assert.Equal(t, createdUser, foundUser)
	})

	t.Run("not found", func(t *testing.T) {
		ctx := internal.ContextWithUser(sharedCtx, internal.User{ID: "69"})
		_, err := us.UpdateUser(ctx, internal.UserUpdate{
			Email: ptr("test@test.com"),
		})
		assert.Equal(t, internal.ErrorCodeNotFound, internal.GetErrorCode(err))
		assert.Equal(t, "user not found", internal.GetErrorMessage(err))
	})
}

func TestDeleteUser(t *testing.T) {
	db, c := mustConnectDB(t)
	defer c()

	us := sqlite.NewUserService(db)
	sharedCtx := internal.ContextWithLogger(context.Background(), zap.NewNop().Sugar())

	t.Run("delete correct user", func(t *testing.T) {
		createdUser1 := mustCreateUser(sharedCtx, t, db, internal.UserCreate{
			ID:    "1",
			Name:  "Alex Albon",
			Email: "alexalbon@williams.com",
		})
		createdUser2 := mustCreateUser(sharedCtx, t, db, internal.UserCreate{
			ID:    "2",
			Name:  "Carlos Sainz",
			Email: "carlossainz@williams.com",
		})

		ctx := internal.ContextWithUser(sharedCtx, createdUser1)
		err := us.DeleteUser(ctx)
		assert.Nil(t, err)

		_, err = us.GetUser(ctx, createdUser1.ID)
		assert.Equal(t, internal.ErrorCodeNotFound, internal.GetErrorCode(err))
		assert.Equal(t, "user not found", internal.GetErrorMessage(err))

		gotUser, err := us.GetUser(ctx, createdUser2.ID)
		assert.Nil(t, err)
		assert.Equal(t, createdUser2, gotUser)
	})

	t.Run("not found", func(t *testing.T) {
		ctx := internal.ContextWithUser(sharedCtx, internal.User{ID: "69"})
		err := us.DeleteUser(ctx)
		assert.Nil(t, err)
	})
}

func mustCreateUser(ctx context.Context, t *testing.T, db *sqlite.DB, c internal.UserCreate) internal.User {
	t.Helper()

	us := sqlite.NewUserService(db)
	u, err := us.CreateUser(ctx, c)
	require.Nil(t, err)
	return u
}

func mustGetUser(ctx context.Context, t *testing.T, db *sqlite.DB, id string) internal.User {
	t.Helper()

	us := sqlite.NewUserService(db)
	u, err := us.GetUser(ctx, id)
	require.Nil(t, err)
	return u
}
