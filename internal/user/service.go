package user

import (
	"context"
	"fmt"

	"github.com/cativovo/budget-tracker/internal"
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

func (us *UserService) FindUserByID(ctx context.Context, id string) (User, error) {
	if id == "" {
		return User{}, internal.NewError(internal.ErrorCodeInvalid, "ID is required")
	}
	return us.ur.FindUserByID(ctx, id)
}

type CreateUserReq struct {
	Name  string `json:"name" validate:"required"`
	ID    string `json:"id"`
	Email string `json:"email" validate:"required,email"`
}

func (us *UserService) CreateUser(ctx context.Context, u CreateUserReq) (User, error) {
	if err := us.v.Struct(u); err != nil {
		return User{}, internal.NewError(internal.ErrorCodeInvalid, err.Error())
	}

	result, err := us.ur.CreateUser(ctx, u)
	if err != nil {
		return User{}, fmt.Errorf("user.UserService: %w", err)
	}

	return result, nil
}

func (us *UserService) DeleteUser(ctx context.Context, id string) error {
	if id == "" {
		return internal.NewError(internal.ErrorCodeInvalid, "ID is required")
	}
	return us.ur.DeleteUser(ctx, id)
}
