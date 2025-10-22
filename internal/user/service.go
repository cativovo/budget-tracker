package user

import (
	"context"

	"github.com/cativovo/budget-tracker/internal"
	"github.com/cativovo/budget-tracker/internal/apperror"
)

// Service handles user operations.
type Service struct {
	userStore userStore
}

// NewService returns a new Service.
func NewService(us userStore) *Service {
	return &Service{userStore: us}
}

// GetUserByID returns a user by ID from the userStore.
func (s *Service) GetUserByID(ctx context.Context, id string) (User, error) {
	return s.userStore.GetUserByID(ctx, id)
}

// CreateUserInput represents the input data needed to create a user.
type CreateUserInput struct {
	ID    string `json:"id" validate:"required"`
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"email"`
}

func (c CreateUserInput) validate() error {
	return internal.ValidateStruct(c)
}

// CreateUser validates the input and creates a new user in the userStore.
func (s *Service) CreateUser(ctx context.Context, input CreateUserInput) (User, error) {
	if err := input.validate(); err != nil {
		return User{}, apperror.New(apperror.ErrorCodeInvalid, err.Error())
	}
	return s.userStore.CreateUser(ctx, input)
}
