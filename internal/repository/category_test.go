package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/cativovo/budget-tracker/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestCreateCategory(t *testing.T) {
	h := newRepositoryHelper(t, "TestCreateCategory.db")
	defer h.clean()
	r := h.repository()

	account, err := r.CreateAccount(context.Background(), logger, repository.CreateAccountParams{
		Name: "Alex Albon",
	})

	assert.Nil(t, err)

	testCases := []struct {
		Name      string
		Icon      string
		ColorHex  string
		AccountID string
	}{
		{
			Name:      "Food",
			Icon:      "apple",
			ColorHex:  "#111111",
			AccountID: account.ID,
		},
		{
			Name:      "Work",
			Icon:      "briefcase",
			ColorHex:  "#222222",
			AccountID: account.ID,
		},
		{
			Name:      "Entertainment",
			Icon:      "clapperboard",
			ColorHex:  "#333333",
			AccountID: account.ID,
		},
	}

	for _, testCase := range testCases {
		inserted, err := r.CreateCategory(context.Background(), logger, repository.CreateCategoryParams{
			Name:      testCase.Name,
			Icon:      testCase.Icon,
			ColorHex:  testCase.ColorHex,
			AccountID: testCase.AccountID,
		})
		assert.Nil(t, err)

		got, err := r.GetCategoryByID(context.Background(), logger, repository.GetCategoryByIDParams{
			ID:        inserted.ID,
			AccountID: testCase.AccountID,
		})
		assert.Nil(t, err)

		assert.Equal(t, testCase.Name, got.Name)
		assert.Equal(t, testCase.Icon, got.Icon)
		assert.Equal(t, testCase.ColorHex, got.ColorHex)
	}
}

func TestGetCategory(t *testing.T) {
	h := newRepositoryHelper(t, "TestGetCategory.db")
	defer h.clean()
	r := h.repository()

	account, err := r.CreateAccount(context.Background(), logger, repository.CreateAccountParams{
		Name: "Logan Sargeant",
	})
	assert.Nil(t, err)

	category, err := r.CreateCategory(context.Background(), logger, repository.CreateCategoryParams{
		Name:      "Food",
		Icon:      "apple",
		ColorHex:  "#111111",
		AccountID: account.ID,
	})
	assert.Nil(t, err)

	testCases := []struct {
		Err       error
		Name      string
		ID        string
		Icon      string
		ColorHex  string
		AccountID string
	}{
		{
			ID:        category.ID,
			Name:      "Food",
			Icon:      "apple",
			ColorHex:  "#111111",
			AccountID: account.ID,
		},
		{
			Err:       sql.ErrNoRows,
			AccountID: "6969",
		},
	}

	for _, testCase := range testCases {
		got, err := r.GetCategoryByID(context.Background(), logger, repository.GetCategoryByIDParams{
			ID:        testCase.ID,
			AccountID: testCase.AccountID,
		})

		assert.Equal(t, testCase.Name, got.Name)
		assert.Equal(t, testCase.Icon, got.Icon)
		assert.Equal(t, testCase.ColorHex, got.ColorHex)
		assert.True(t, errors.Is(err, testCase.Err))
	}
}
