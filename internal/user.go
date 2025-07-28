package internal

import (
	"context"
	"time"
)

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserService interface {
	GetUser(ctx context.Context, id string) (User, error)
	CreateUser(ctx context.Context, c UserCreate) (User, error)
	UpdateUser(ctx context.Context, u UserUpdate) (User, error)
	DeleteUser(ctx context.Context) error
}

type UserCreate struct {
	ID    string `validate:"required" json:"id"`
	Name  string `validate:"required" json:"name"`
	Email string `validate:"email" json:"email"`
}

func (c UserCreate) Validate() error {
	if err := ValidateStruct(c); err != nil {
		return NewError(ErrorCodeInvalid, err.Error())
	}

	return nil
}

// Doesn't accept ID because it's using the User from ctx
type UserUpdate struct {
	// https://github.com/go-playground/validator/issues/1308
	Name  *string `json:"name" validate:"omitnil,min=1"`
	Email *string `json:"email" validate:"omitnil,email"`
}

func (u UserUpdate) Validate() error {
	if u.Name == nil && u.Email == nil {
		return NewError(ErrorCodeInvalid, "no update fields provided")
	}

	if err := ValidateStruct(u); err != nil {
		return NewError(ErrorCodeInvalid, err.Error())
	}

	return nil
}
