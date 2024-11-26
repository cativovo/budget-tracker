package repository_test

import (
	"context"
	"testing"

	"github.com/cativovo/budget-tracker/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestCreateAndGetAccount(t *testing.T) {
	h := newRepositoryHelper(t, "TestCreateAndGetAccount.db")
	r := h.repository()
	defer h.clean()

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
		inserted, err := r.CreateAccount(context.Background(), repository.CreateAccountParams{
			Name: testCase.Name,
		})
		if err != nil {
			t.Fatal(err)
		}

		got, err := r.GetAccountByID(context.Background(), inserted.ID)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, testCase.Name, got.Name)
	}
}
