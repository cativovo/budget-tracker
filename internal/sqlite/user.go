package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/cativovo/budget-tracker/internal"
	"github.com/huandu/go-sqlbuilder"
	"github.com/mattn/go-sqlite3"
)

type UserService struct {
	db *DB
}

var _ internal.UserService = (*UserService)(nil)

func NewUserService(db *DB) *UserService {
	return &UserService{
		db: db,
	}
}

func (us *UserService) GetUser(ctx context.Context, id string) (internal.User, error) {
	if id == "" {
		err := internal.NewError(internal.ErrorCodeNotFound, "user not found")
		return internal.User{}, fmt.Errorf("sqlite: get user: %w", err)
	}

	logger := internal.LoggerFromContext(ctx)

	sb := sqlbuilder.SQLite.NewSelectBuilder()
	sb.Select("id", "name", "email", "created_at", "updated_at")
	sb.From("user")
	sb.Where(sb.EQ("id", id))
	q, args := sb.Build()

	logger.Infow("Get user", "query", q, "args", args)

	var dest userDest
	if err := us.db.Reader.GetContext(ctx, &dest, q, args...); err != nil {
		if err == sql.ErrNoRows {
			err = internal.NewError(internal.ErrorCodeNotFound, "user not found")
		}
		return internal.User{}, fmt.Errorf("sqlite: get user: %w", err)
	}

	return internal.User(dest), nil
}

func (us *UserService) CreateUser(ctx context.Context, c internal.UserCreate) (internal.User, error) {
	if err := c.Validate(); err != nil {
		return internal.User{}, fmt.Errorf("sqlite: validate create user input: %w", err)
	}

	logger := internal.LoggerFromContext(ctx)

	ib := sqlbuilder.SQLite.NewInsertBuilder()
	ib.InsertInto("user")
	ib.Cols("id", "name", "email")
	ib.Values(c.ID, c.Name, c.Email)
	ib.Returning("id", "name", "email", "created_at", "updated_at")

	q, args := ib.Build()

	logger.Infow("Insert user", "query", q, "args", args)

	var dest userDest
	if err := us.db.Writer.GetContext(ctx, &dest, q, args...); err != nil {
		var sqErr sqlite3.Error
		if errors.As(err, &sqErr) && sqErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			err = internal.NewError(internal.ErrorCodeConflict, "email already taken")
		}
		return internal.User{}, fmt.Errorf("sqlite: insert user: %w", err)
	}

	return internal.User(dest), nil
}

func (us *UserService) UpdateUser(ctx context.Context, u internal.UserUpdate) (internal.User, error) {
	if err := u.Validate(); err != nil {
		return internal.User{}, fmt.Errorf("sqlite: validate update user input: %w", err)
	}

	user := internal.UserFromContext(ctx)
	logger := internal.LoggerFromContext(ctx)

	ub := sqlbuilder.SQLite.NewUpdateBuilder()
	ub.Update("user")

	setMoreIfNotNil(ub, "name", u.Name)
	setMoreIfNotNil(ub, "email", u.Email)
	setUpdatedAt(ub)

	ub.Where(ub.EQ("id", user.ID))
	// https://github.com/huandu/go-sqlbuilder/issues/142
	ub.SQL("RETURNING id, name, email, created_at, updated_at")

	q, args := ub.Build()

	logger.Infow("Update user", "query", q, "args", args)

	var dest userDest
	if err := us.db.Writer.GetContext(ctx, &dest, q, args...); err != nil {
		if err == sql.ErrNoRows {
			err = internal.NewError(internal.ErrorCodeNotFound, "user not found")
		}

		var sqErr sqlite3.Error
		if errors.As(err, &sqErr) && sqErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			err = internal.NewError(internal.ErrorCodeConflict, "email already taken")
		}

		return internal.User{}, fmt.Errorf("sqlite: update user: %w", err)
	}

	return internal.User(dest), nil
}

func (us *UserService) DeleteUser(ctx context.Context) error {
	user := internal.UserFromContext(ctx)
	logger := internal.LoggerFromContext(ctx)

	db := sqlbuilder.SQLite.NewDeleteBuilder()
	db.DeleteFrom("user")
	db.Where(db.EQ("id", user.ID))

	q, args := db.Build()

	logger.Infow(
		"Delete user",
		"query", q,
		"args", args,
	)

	_, err := us.db.Writer.ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("sqlite: delete user: %w", err)
	}

	return nil
}

type userDest struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
