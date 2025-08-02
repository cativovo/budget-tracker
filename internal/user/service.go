package user

import (
	"context"

	"github.com/cativovo/budget-tracker/apperror"
	"github.com/cativovo/budget-tracker/internal"
)

// UserService is a service for managing users.
type UserService struct {
	userStore userStore
}

// NewService creates a new UserService.
func NewService(us userStore) *UserService {
	return &UserService{
		userStore: us,
	}
}

// GetUser gets a user by id.
func (us *UserService) GetUser(ctx context.Context, id string) (User, error) {
	if id == "" {
		return User{}, apperror.New(apperror.ErrorCodeInvalid, "id is required")
	}
	return us.userStore.GetUser(ctx, id)
}

// UserCreate is the input for creating a user.
type UserCreate struct {
	ID    string `json:"id" validate:"required"`
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"email"`
}

func (uc UserCreate) validate() error {
	if err := internal.ValidateStruct(uc); err != nil {
		return apperror.New(apperror.ErrorCodeInvalid, err.Error())
	}
	return nil
}

// CreateUser creates a new user.
func (us *UserService) CreateUser(ctx context.Context, input UserCreate) (User, error) {
	if err := input.validate(); err != nil {
		return User{}, err
	}
	return us.userStore.CreateUser(ctx, input)
}

// UserUpdate is the input for updating a user.
type UserUpdate struct {
	// https://github.com/go-playground/validator/issues/1308
	Name  *string `json:"name" validate:"omitnil,min=1"`
	Email *string `json:"email" validate:"omitnil,email"`
}

func (uu UserUpdate) validate() error {
	if err := internal.ValidateStruct(uu); err != nil {
		return apperror.New(apperror.ErrorCodeInvalid, err.Error())
	}
	if uu.Name == nil && uu.Email == nil {
		return apperror.New(apperror.ErrorCodeInvalid, "no update fields provided")
	}
	return nil
}

// UpdateUser updates a user.
func (us *UserService) UpdateUser(ctx context.Context, input UserUpdate) (User, error) {
	if err := input.validate(); err != nil {
		return User{}, err
	}
	return us.userStore.UpdateUser(ctx, input)
}
