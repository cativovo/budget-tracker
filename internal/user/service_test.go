package user_test

import (
	"context"
	"testing"
	"time"

	"github.com/cativovo/budget-tracker/internal"
	"github.com/cativovo/budget-tracker/internal/user"
	"github.com/stretchr/testify/assert"
)

func TestUserService_GetUser(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		want    user.User
		wantErr error
	}{
		{
			name: "get user",
			id:   "1",
			want: user.User{
				ID:        "1",
				Name:      "Jose Rizal",
				Email:     "joselrizzal@katipunan.ph",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
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

			us := user.NewService(m)
			got, gotErr := us.GetUser(ctx, tt.id)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, gotErr)
		})
	}
}

func TestUserService_CreateUser(t *testing.T) {
	tests := []struct {
		name    string
		input   user.UserCreate
		want    user.User
		wantErr error
	}{
		{
			name: "create user",
			input: user.UserCreate{
				ID:    "2",
				Name:  "Andres Bonifacio",
				Email: "andres.bonifacio@katipunan.ph",
			},
			want: user.User{
				ID:        "2",
				Name:      "Andres Bonifacio",
				Email:     "andres.bonifacio@katipunan.ph",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		{
			name: "no id",
			input: user.UserCreate{
				Name:  "Andres Bonifacio",
				Email: "andres.bonifacio@katipunan.ph",
			},
			wantErr: internal.NewError(internal.ErrorCodeInvalid, "id is required"),
		},
		{
			name: "no name",
			input: user.UserCreate{
				ID:    "2",
				Email: "andres.bonifacio@katipunan.ph",
			},
			wantErr: internal.NewError(internal.ErrorCodeInvalid, "name is required"),
		},
		{
			name: "no email",
			input: user.UserCreate{
				ID:   "2",
				Name: "Andres Bonifacio",
			},
			wantErr: internal.NewError(internal.ErrorCodeInvalid, "invalid email"),
		},
		{
			name: "invalid email",
			input: user.UserCreate{
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
				m.EXPECT().CreateUser(ctx, tt.input).Return(tt.want, nil).Once()
			}

			us := user.NewService(m)
			got, gotErr := us.CreateUser(ctx, tt.input)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, gotErr)
		})
	}
}

// func TestUserService_UpdateUser(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		user    user.User
// 		input   user.UserUpdate
// 		want    user.User
// 		wantErr error
// 	}{
// 		{
// 			name: "update user",
// 			user: user.User{
// 				ID:        "3",
// 				Name:      "Melchora Aquin",
// 				Email:     "batang.sora@revolution.ph",
// 				CreatedAt: time.Date(2025, 06, 01, 0, 0, 0, 0, time.UTC),
// 				UpdatedAt: time.Date(2025, 06, 01, 0, 0, 0, 0, time.UTC),
// 			},
// 			input: user.UserUpdate{
// 				Name:  testutil.ToPtr("Melchora Aquino"),
// 				Email: testutil.ToPtr("tandang.sora@revolution.ph"),
// 			},
// 			want: user.User{
// 				ID:        "3",
// 				Name:      "Melchora Aquino",
// 				Email:     "tandang.sora@revolution.ph",
// 				CreatedAt: time.Date(2025, 06, 01, 0, 0, 0, 0, time.UTC),
// 				UpdatedAt: time.Date(2025, 06, 01, 0, 0, 0, 50, time.UTC),
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctx := user.WithContext(context.Background(), tt.user)
// 			m := newMockuserStore(t)
//
// 			if tt.wantErr == nil {
// 				u := tt.want
// 				u.CreatedAt = tt.user.CreatedAt
// 				u.UpdatedAt = tt.user.UpdatedAt
// 				m.EXPECT().UpdateUser(ctx, u).Return(tt.want, nil)
// 			}
//
// 			us := user.NewService(m)
// 			got, gotErr := us.UpdateUser(ctx, tt.input)
// 			assert.Equal(t, tt.want, got)
// 			assert.Equal(t, tt.wantErr, gotErr)
// 		})
// 	}
// }
