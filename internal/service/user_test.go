package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/cativovo/budget-tracker/internal"
	"github.com/cativovo/budget-tracker/internal/domain"
	"github.com/cativovo/budget-tracker/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestUserService_GetUser(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		want    domain.User
		wantErr error
	}{
		{
			name: "get user",
			id:   "1",
			want: domain.User{
				ID:        "1",
				Name:      "Jose Rizal",
				Email:     "joselrizzal@katipunan.ph",
				CreatedAt: time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name:    "no id",
			wantErr: internal.NewError(internal.ErrorCodeInvalid, "id is required"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			m := newMockuserStore(t)

			if tt.wantErr == nil {
				m.EXPECT().GetUser(ctx, tt.id).Return(tt.want, nil)
			}

			us := service.NewUserService(m)
			got, gotErr := us.GetUser(ctx, tt.id)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, gotErr)
		})
	}
}

func TestUserService_CreateUser(t *testing.T) {
	tests := []struct {
		name    string
		input   service.UserCreate
		want    domain.User
		wantErr error
	}{
		{
			name: "create user",
			input: service.UserCreate{
				ID:    "2",
				Name:  "Andres Bonifacio",
				Email: "andres.bonifacio@katipunan.ph",
			},
			want: domain.User{
				ID:        "2",
				Name:      "Andres Bonifacio",
				Email:     "andres.bonifacio@katipunan.ph",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		{
			name: "no id",
			input: service.UserCreate{
				Name:  "Andres Bonifacio",
				Email: "andres.bonifacio@katipunan.ph",
			},
			wantErr: internal.NewError(internal.ErrorCodeInvalid, "id is required"),
		},
		{
			name: "no name",
			input: service.UserCreate{
				ID:    "2",
				Email: "andres.bonifacio@katipunan.ph",
			},
			wantErr: internal.NewError(internal.ErrorCodeInvalid, "name is required"),
		},
		{
			name: "no email",
			input: service.UserCreate{
				ID:   "2",
				Name: "Andres Bonifacio",
			},
			wantErr: internal.NewError(internal.ErrorCodeInvalid, "invalid email"),
		},
		{
			name: "invalid email",
			input: service.UserCreate{
				ID:    "2",
				Name:  "Andres Bonifacio",
				Email: "andres.bonifaciokatipunan.ph",
			},
			wantErr: internal.NewError(internal.ErrorCodeInvalid, "invalid email"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			m := newMockuserStore(t)

			if tt.wantErr == nil {
				u := domain.User{
					ID:    tt.want.ID,
					Name:  tt.want.Name,
					Email: tt.want.Email,
				}
				m.EXPECT().CreateUser(ctx, u).Return(tt.want, nil)
			}

			us := service.NewUserService(m)
			got, gotErr := us.CreateUser(ctx, tt.input)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, gotErr)
		})
	}
}
