package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/cativovo/budget-tracker/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestCreateAccount(t *testing.T) {
	h := newRepositoryHelper(t, "TestCreateAccount.db")
	defer h.clean()
	r := h.repository()

	testCases := []struct {
		Name string
	}{
		{
			Name: "Bulbasaur",
		},
		{
			Name: "Squirtle",
		},
		{
			Name: "Charmander",
		},
	}

	for _, testCase := range testCases {
		inserted, err := r.CreateAccount(context.Background(), logger, repository.CreateAccountParams{
			Name: testCase.Name,
		})
		assert.Nil(t, err)

		got, err := r.GetAccountByID(context.Background(), logger, inserted.ID)
		assert.Nil(t, err)
		assert.Equal(t, testCase.Name, got.Name)
	}
}

func TestGetAccount(t *testing.T) {
	h := newRepositoryHelper(t, "TestGetAccount.db")
	defer h.clean()
	r := h.repository()

	inserted, err := r.CreateAccount(context.Background(), logger, repository.CreateAccountParams{
		Name: "Franco Colapinto",
	})
	assert.Nil(t, err)

	testCases := []struct {
		Err  error
		ID   string
		Name string
	}{
		{
			ID:   inserted.ID,
			Name: "Franco Colapinto",
		},
		{
			Err: sql.ErrNoRows,
		},
	}

	for _, testCase := range testCases {
		assert.Nil(t, err)

		got, err := r.GetAccountByID(context.Background(), logger, testCase.ID)
		assert.True(t, errors.Is(err, testCase.Err))
		assert.Equal(t, testCase.Name, got.Name)
	}
}
