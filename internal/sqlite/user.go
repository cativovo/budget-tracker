package sqlite

import (
	"context"

	"github.com/cativovo/budget-tracker/internal/user"
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

func (ur *UserRepository) CreateUser(ctx context.Context, u user.CreateUserReq) (user.User, error) {
	panic("not yet implemented")
}

func (ur *UserRepository) DeleteUser(ctx context.Context, id string) error {
	panic("not yet implemented")
}
