package user

import (
	"context"
	"fmt"

	"github.com/cativovo/budget-tracker/internal"
	"github.com/cativovo/budget-tracker/internal/validator"
)

type Service interface {
	UserByID(ctx context.Context, id string) (User, error)
	Create(ctx context.Context, u CreateUserReq) (User, error)
	Delete(ctx context.Context, id string) error
}

type CreateUserReq struct {
	ID    string `json:"id"`
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type service struct {
	r Repository
	v *validator.Validator
}

func NewService(r Repository, v *validator.Validator) Service {
	return &service{
		r: r,
		v: v,
	}
}

func (s *service) UserByID(ctx context.Context, id string) (User, error) {
	if id == "" {
		return User{}, internal.NewError(internal.ErrorCodeInvalid, "ID is required")
	}
	return s.r.UserByID(ctx, id)
}

func (s *service) Create(ctx context.Context, u CreateUserReq) (User, error) {
	if err := s.v.Struct(u); err != nil {
		return User{}, err
	}

	result, err := s.r.CreateUser(ctx, u)
	if err != nil {
		return User{}, fmt.Errorf("user.UserService: %w", err)
	}

	return result, nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	if id == "" {
		return internal.NewError(internal.ErrorCodeInvalid, "ID is required")
	}
	return s.r.DeleteUser(ctx, id)
}
