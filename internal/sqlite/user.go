package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/cativovo/budget-tracker/internal"
	"github.com/huandu/go-sqlbuilder"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
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
		iErr := internal.NewError(internal.ErrorCodeNotFound, "User not found")
		return internal.User{}, fmt.Errorf("sqlite: %w", iErr)
	}

	logger := internal.LoggerFromContext(ctx)

	sb := sqlbuilder.SQLite.NewSelectBuilder()
	sb.Select("id", "name", "email", "created_at", "updated_at")
	sb.From("user")
	sb.Where(sb.EQ("id", id))
	q, args := sb.Build()

	logger.Infow("Get user", "query", q, "args", args)

	var dst userDst
	if err := us.db.reader.GetContext(ctx, &dst, q, args...); err != nil {
		if err == sql.ErrNoRows {
			iErr := internal.NewError(internal.ErrorCodeNotFound, "User not found")
			return internal.User{}, fmt.Errorf("sqlite: %w", iErr)
		}

		return internal.User{}, fmt.Errorf("sqlite: get user: %w", err)
	}

	return internal.User(dst), nil
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

	var dst userDst
	if err := us.db.writer.GetContext(ctx, &dst, q, args...); err != nil {
		var sqErr *sqlite.Error
		if errors.As(err, &sqErr) && sqErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
			iErr := internal.NewError(internal.ErrorCodeConflict, "email already taken")
			return internal.User{}, fmt.Errorf("sqlite: %w", iErr)
		}
		return internal.User{}, fmt.Errorf("sqlite: insert user: %w", err)
	}

	return internal.User(dst), nil
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

	ub.Where(ub.Equal("id", user.ID))
	// https://github.com/huandu/go-sqlbuilder/issues/142
	ub.SQL("RETURNING id, name, email, created_at, updated_at")

	q, args := ub.Build()

	logger.Infow("Update user", "query", q, "args", args)

	var dst userDst
	if err := us.db.writer.GetContext(ctx, &dst, q, args...); err != nil {
		if err == sql.ErrNoRows {
			iErr := internal.NewError(internal.ErrorCodeNotFound, "User not found")
			return internal.User{}, fmt.Errorf("sqlite: %w", iErr)
		}

		var sqErr *sqlite.Error
		if errors.As(err, &sqErr) && sqErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
			iErr := internal.NewError(internal.ErrorCodeConflict, "email already taken")
			return internal.User{}, fmt.Errorf("sqlite: %w", iErr)
		}

		return internal.User{}, fmt.Errorf("sqlite: update user: %w", err)
	}

	return internal.User(dst), nil
}

func (us *UserService) DeleteUser(ctx context.Context) error {
	user := internal.UserFromContext(ctx)
	logger := internal.LoggerFromContext(ctx)

	db := sqlbuilder.SQLite.NewDeleteBuilder()
	db.DeleteFrom("user")
	db.Where(db.Equal("id", user.ID))

	q, args := db.Build()

	logger.Infow(
		"Delete user",
		"query", q,
		"args", args,
	)

	if _, err := us.db.writer.ExecContext(ctx, q, args...); err != nil {
		return fmt.Errorf("sqlite: delete user: %w", err)
	}

	return nil
}

type userDst struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
