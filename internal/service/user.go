package service

import (
	"context"

	"github.com/cativovo/budget-tracker/internal"
	"github.com/cativovo/budget-tracker/internal/domain"
)

type userStore interface {
	GetUser(ctx context.Context, id string) (domain.User, error)
	CreateUser(ctx context.Context, input domain.User) (domain.User, error)
}

type UserService struct {
	userStore userStore
}

func NewUserService(us userStore) *UserService {
	return &UserService{
		userStore: us,
	}
}

func (us *UserService) GetUser(ctx context.Context, id string) (domain.User, error) {
	if id == "" {
		return domain.User{}, internal.NewError(internal.ErrorCodeInvalid, "id is required")
	}
	return us.userStore.GetUser(ctx, id)
}

type UserCreate struct {
	ID    string `json:"id" validate:"required"`
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"email"`
}

func (us *UserService) CreateUser(ctx context.Context, input UserCreate) (domain.User, error) {
	if err := validateStruct(input); err != nil {
		return domain.User{}, internal.NewError(internal.ErrorCodeInvalid, err.Error())
	}
	return us.userStore.CreateUser(ctx, domain.User{
		ID:    input.ID,
		Name:  input.Name,
		Email: input.Email,
	})
}
