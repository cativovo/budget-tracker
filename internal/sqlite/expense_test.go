package sqlite_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/cativovo/budget-tracker/internal"
	"github.com/cativovo/budget-tracker/internal/expense"
	"github.com/cativovo/budget-tracker/internal/logger"
	"github.com/cativovo/budget-tracker/internal/sqlite"
	"github.com/cativovo/budget-tracker/internal/user"
	"github.com/stretchr/testify/assert"
)

func TestCreateFindExpense(t *testing.T) {
	dh := newDBHelper(t, "foo.db")
	defer dh.clean()

	cr := sqlite.NewCategoryRepository(dh.db)
	er := sqlite.NewExpenseRepository(dh.db, cr)
	ctxWithLogger := logger.ContextWithLogger(context.Background(), zapLogger)

	users := createUsers(t, dh.db)

	user1 := users[0]
	user1Categories := createCategories(t, dh.db, user1)

	user2 := users[1]
	user2Categories := createCategories(t, dh.db, user2)

	tests := []struct {
		name  string
		user  user.User
		input expense.CreateExpenseReq
		want  expense.Expense
	}{
		{
			name: fmt.Sprintf("%s's expense 1", user1.Name),
			user: user1,
			input: expense.CreateExpenseReq{
				Name:       "Expense 1",
				Amount:     6969,
				Date:       "2006-01-02",
				CategoryID: user1Categories[0].ID,
				Note:       "Expense 1 Note",
			},
			want: expense.Expense{
				Name:      "Expense 1",
				Amount:    6969,
				Date:      time.Date(2006, time.January, 2, 0, 0, 0, 0, time.UTC),
				Note:      "Expense 1 Note",
				Category:  user1Categories[0],
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		{
			name: fmt.Sprintf("%s's expense 2 without note", user1.Name),
			user: user1,
			input: expense.CreateExpenseReq{
				Name:       "Expense 2",
				Amount:     7070,
				Date:       "2006-01-03",
				CategoryID: user1Categories[0].ID,
			},
			want: expense.Expense{
				Name:      "Expense 2",
				Amount:    7070,
				Date:      time.Date(2006, time.January, 3, 0, 0, 0, 0, time.UTC),
				Category:  user1Categories[0],
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		{
			name: fmt.Sprintf("%s's expense 1", user2.Name),
			user: user2,
			input: expense.CreateExpenseReq{
				Name:       "Expense 1",
				Amount:     6969,
				Date:       "2006-01-02",
				CategoryID: user2Categories[0].ID,
				Note:       "Expense 1 Note",
			},
			want: expense.Expense{
				Name:      "Expense 1",
				Amount:    6969,
				Date:      time.Date(2006, time.January, 2, 0, 0, 0, 0, time.UTC),
				Note:      "Expense 1 Note",
				Category:  user2Categories[0],
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctxWithUser := user.ContextWithUser(ctxWithLogger, test.user)
			e, err := er.CreateExpense(ctxWithUser, test.input)

			expenseCopy := e
			wantCopy := test.want

			assert.Nil(t, err)
			assert.True(t, e.ID != "")
			// to make it easier to assert
			e.ID = ""

			assert.WithinDuration(t, test.want.CreatedAt, e.CreatedAt, time.Second)
			assert.WithinDuration(t, test.want.UpdatedAt, e.UpdatedAt, time.Second)

			// to make it easier to assert
			e.CreatedAt = time.Time{}
			e.UpdatedAt = time.Time{}
			test.want.CreatedAt = time.Time{}
			test.want.UpdatedAt = time.Time{}

			assert.Equal(t, test.want, e)

			foundExpense, err := er.ExpenseByID(ctxWithUser, expenseCopy.ID)
			assert.Nil(t, err)
			assert.Nil(t, err)
			assert.True(t, foundExpense.ID != "")
			foundExpense.ID = ""

			assert.WithinDuration(t, wantCopy.CreatedAt, foundExpense.CreatedAt, time.Second)
			assert.WithinDuration(t, wantCopy.UpdatedAt, foundExpense.UpdatedAt, time.Second)

			foundExpense.CreatedAt = time.Time{}
			foundExpense.UpdatedAt = time.Time{}
			wantCopy.CreatedAt = time.Time{}
			wantCopy.UpdatedAt = time.Time{}

			assert.Equal(t, test.want, foundExpense)
		})
	}

	t.Run("expense not found", func(t *testing.T) {
		ctxWithUser := user.ContextWithUser(ctxWithLogger, user1)
		foundExpense, err := er.ExpenseByID(ctxWithUser, "123")
		assert.Equal(t, expense.Expense{}, foundExpense)
		assert.Equal(t, err, internal.NewError(internal.ErrorCodeNotFound, "Expense not found"))
	})

	t.Run("can't access expense of other user", func(t *testing.T) {
		ctxWithUser := user.ContextWithUser(ctxWithLogger, user1)
		e, err := er.CreateExpense(ctxWithUser, tests[0].input)
		assert.Nil(t, err)

		ctxWithUser = user.ContextWithUser(ctxWithLogger, user2)
		foundExpense, err := er.ExpenseByID(ctxWithUser, e.ID)
		assert.Equal(t, expense.Expense{}, foundExpense)
		assert.Equal(t, err, internal.NewError(internal.ErrorCodeNotFound, "Expense not found"))
	})
}
