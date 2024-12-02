package repository_test

import (
	"context"
	"testing"

	"github.com/cativovo/budget-tracker/internal/repository"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestCreateAndGetAccount(t *testing.T) {
	h := newRepositoryHelper(t, "TestCreateAndGetAccount.db")
	r := h.repository()
	defer h.clean()
	logger := zap.NewNop().Sugar()

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
		if err != nil {
			t.Fatal(err)
		}

		got, err := r.GetAccountByID(context.Background(), logger, inserted.ID)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, testCase.Name, got.Name)
	}
}
