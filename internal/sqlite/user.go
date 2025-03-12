package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/cativovo/budget-tracker/internal"
	"github.com/cativovo/budget-tracker/internal/user"
	"github.com/huandu/go-sqlbuilder"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type UserRepository struct {
	db *DB
}

var _ user.UserRepository = (*UserRepository)(nil)

func NewUserRepository(db *DB) UserRepository {
	return UserRepository{
		db: db,
	}
}

func (ur *UserRepository) FindUserByID(ctx context.Context, id string) (user.User, error) {
	logger := internal.LoggerFromCtx(ctx)

	sb := sqlbuilder.SQLite.NewSelectBuilder()
	sb.Select(
		"id",
		"name",
		"email",
	)
	sb.From("user")
	sb.Where(sb.EQ("id", id))
	q, args := sb.Build()

	logger.Infow(
		"Find user by id",
		"query", q,
		"args", args,
	)

	var dst struct {
		ID    string `db:"id"`
		Name  string `db:"name"`
		Email string `db:"email"`
	}
	if err := ur.db.reader.GetContext(ctx, &dst, q, args...); err != nil {
		if err == sql.ErrNoRows {
			return user.User{}, internal.NewError(internal.ErrorCodeNotFound, "User not found")
		}

		return user.User{}, fmt.Errorf("sqlite.UserRepository: %w", err)
	}

	return user.User{
		ID:    dst.ID,
		Name:  dst.Name,
		Email: dst.Email,
	}, nil
}

func (ur *UserRepository) CreateUser(ctx context.Context, u user.CreateUserReq) (user.User, error) {
	logger := internal.LoggerFromCtx(ctx)

	ib := sqlbuilder.SQLite.NewInsertBuilder()
	ib.InsertInto("user")
	ib.Cols(
		"id",
		"name",
		"email",
	)
	ib.Values(
		u.ID,
		u.Name,
		u.Email,
	)

	q, args := ib.Build()

	logger.Infow(
		"Insert new user",
		"query", q,
		"args", args,
	)

	if _, err := ur.db.readerWriter.ExecContext(ctx, q, args...); err != nil {
		var e *sqlite.Error
		if errors.As(err, &e) && e.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
			return user.User{}, internal.NewError(internal.ErrorCodeConflict, "User already exists")
		}

		return user.User{}, fmt.Errorf("sqlite.UserRepository: %w", err)
	}

	return user.User{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}, nil
}

func (ur *UserRepository) DeleteUser(ctx context.Context, id string) error {
	logger := internal.LoggerFromCtx(ctx)

	db := sqlbuilder.SQLite.NewDeleteBuilder()
	db.DeleteFrom("user")
	db.Where(db.EQ("id", id))
	q, args := db.Build()

	logger.Infow(
		"Delete user",
		"query", q,
		"args", args,
	)

	_, err := ur.db.readerWriter.ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("sqlite.UserRepository: %w", err)
	}

	return nil
}
