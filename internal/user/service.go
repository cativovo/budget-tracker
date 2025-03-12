package user

import (
	"context"
	"errors"

	"github.com/cativovo/budget-tracker/internal/validator"
)

type UserService struct {
	ur UserRepository
	v  *validator.Validator
}

func NewUserService(ur UserRepository, v *validator.Validator) *UserService {
	return &UserService{
		ur: ur,
		v:  v,
	}
}

type CreateUserReq struct {
	Name  string `json:"name" validate:"required"`
	ID    string `json:"id"`
	Email string `json:"email" validate:"required,email"`
}

func (us *UserService) CreateUser(ctx context.Context, u CreateUserReq) (User, error) {
	if err := us.v.Struct(u); err != nil {
		return User{}, err
	}
	return us.ur.CreateUser(ctx, u)
}

func (us *UserService) DeleteUser(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id is required")
	}
	return us.ur.DeleteUser(ctx, id)
}
