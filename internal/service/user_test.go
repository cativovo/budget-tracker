package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/cativovo/budget-tracker/internal"
	"github.com/cativovo/budget-tracker/internal/ctxvalue"
	"github.com/cativovo/budget-tracker/internal/model"
	"github.com/cativovo/budget-tracker/internal/service"
	"github.com/cativovo/budget-tracker/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestUserService_GetUser(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		want    model.User
		wantErr error
	}{
		{
			name: "get user",
			id:   "1",
			want: model.User{
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
		want    model.User
		wantErr error
	}{
		{
			name: "create user",
			input: service.UserCreate{
				ID:    "2",
				Name:  "Andres Bonifacio",
				Email: "andres.bonifacio@katipunan.ph",
			},
			want: model.User{
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
				u := model.User{
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

func TestUserService_UpdateUser(t *testing.T) {
	tests := []struct {
		name    string
		user    model.User
		input   service.UserUpdate
		want    model.User
		wantErr error
	}{
		{
			name: "update user",
			user: model.User{
				ID:        "3",
				Name:      "Melchora Aquin",
				Email:     "batang.sora@revolution.ph",
				CreatedAt: time.Date(2025, 06, 01, 0, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2025, 06, 01, 0, 0, 0, 0, time.UTC),
			},
			input: service.UserUpdate{
				Name:  testutil.ToPtr("Melchora Aquino"),
				Email: testutil.ToPtr("tandang.sora@revolution.ph"),
			},
			want: model.User{
				ID:        "3",
				Name:      "Melchora Aquino",
				Email:     "tandang.sora@revolution.ph",
				CreatedAt: time.Date(2025, 06, 01, 0, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2025, 06, 01, 0, 0, 0, 50, time.UTC),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := ctxvalue.ContextWithUser(context.Background(), tt.user)
			m := newMockuserStore(t)

			if tt.wantErr == nil {
				u := tt.want
				u.CreatedAt = tt.user.CreatedAt
				u.UpdatedAt = tt.user.UpdatedAt
				m.EXPECT().UpdateUser(ctx, u).Return(tt.want, nil)
			}

			us := service.NewUserService(m)
			got, gotErr := us.UpdateUser(ctx, tt.input)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, gotErr)
		})
	}
}
