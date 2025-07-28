package service

import (
	"context"

	"github.com/cativovo/budget-tracker/internal"
	"github.com/cativovo/budget-tracker/internal/ctxvalue"
	"github.com/cativovo/budget-tracker/internal/model"
)

type userStore interface {
	GetUser(ctx context.Context, id string) (model.User, error)
	CreateUser(ctx context.Context, user model.User) (model.User, error)
	UpdateUser(ctx context.Context, user model.User) (model.User, error)
}

type UserService struct {
	userStore userStore
}

func NewUserService(us userStore) *UserService {
	return &UserService{
		userStore: us,
	}
}

func (us *UserService) GetUser(ctx context.Context, id string) (model.User, error) {
	if id == "" {
		return model.User{}, internal.NewError(internal.ErrorCodeInvalid, "id is required")
	}
	return us.userStore.GetUser(ctx, id)
}

type UserCreate struct {
	ID    string `json:"id" validate:"required"`
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"email"`
}

func (us *UserService) CreateUser(ctx context.Context, input UserCreate) (model.User, error) {
	if err := validateStruct(input); err != nil {
		return model.User{}, internal.NewError(internal.ErrorCodeInvalid, err.Error())
	}
	return us.userStore.CreateUser(ctx, model.User{
		ID:    input.ID,
		Name:  input.Name,
		Email: input.Email,
	})
}

type UserUpdate struct {
	// https://github.com/go-playground/validator/issues/1308
	Name  *string `json:"name" validate:"omitnil,min=1"`
	Email *string `json:"email" validate:"omitnil,email"`
}

func (us *UserService) UpdateUser(ctx context.Context, input UserUpdate) (model.User, error) {
	if err := validateStruct(input); err != nil {
		return model.User{}, internal.NewError(internal.ErrorCodeInvalid, err.Error())
	}

	if input.Name == nil && input.Email == nil {
		return model.User{}, internal.NewError(internal.ErrorCodeInvalid, "no update fields provided")
	}

	user := ctxvalue.UserFromContext(ctx)

	if input.Name != nil {
		user.Name = *input.Name
	}

	if input.Email != nil {
		user.Email = *input.Email
	}

	return us.userStore.UpdateUser(ctx, user)
}
