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
	dh := newDBHelper(t, "test_create_find_expense.db")
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
			createdExpense, err := er.CreateExpense(ctxWithUser, test.input)
			assert.Nil(t, err)
			assertExpense(t, test.want, createdExpense)

			foundExpense, err := er.ExpenseByID(ctxWithUser, createdExpense.ID)
			assert.Nil(t, err)
			assertExpense(t, test.want, foundExpense)
		})
	}

	t.Run("expense not found", func(t *testing.T) {
		ctxWithUser := user.ContextWithUser(ctxWithLogger, user1)
		foundExpense, err := er.ExpenseByID(ctxWithUser, "123")
		assert.Equal(t, expense.Expense{}, foundExpense)
		assert.Equal(t, internal.NewError(internal.ErrorCodeNotFound, "Expense not found"), err)
	})

	t.Run("can't access expense of other user", func(t *testing.T) {
		ctxWithUser := user.ContextWithUser(ctxWithLogger, user1)
		createdExpense, err := er.CreateExpense(ctxWithUser, tests[0].input)
		assert.Nil(t, err)

		ctxWithUser = user.ContextWithUser(ctxWithLogger, user2)
		foundExpense, err := er.ExpenseByID(ctxWithUser, createdExpense.ID)
		assert.Equal(t, expense.Expense{}, foundExpense)
		assert.Equal(t, internal.NewError(internal.ErrorCodeNotFound, "Expense not found"), err)
	})
}

func TestUpdateExpense(t *testing.T) {
	dh := newDBHelper(t, "test_update_expense.db")
	defer dh.clean()

	cr := sqlite.NewCategoryRepository(dh.db)
	er := sqlite.NewExpenseRepository(dh.db, cr)
	ctxWithLogger := logger.ContextWithLogger(context.Background(), zapLogger)

	users := createUsers(t, dh.db)

	user1 := users[0]
	user1Categories := createCategories(t, dh.db, user1)
	user2 := users[1]

	tests := []struct {
		name    string
		user    user.User
		expense expense.CreateExpenseReq
		input   expense.UpdateExpenseReq
		want    expense.Expense
	}{
		{
			name: fmt.Sprintf("update %s's expense name", user1.Name),
			user: user1,
			expense: expense.CreateExpenseReq{
				Name:       "Expense 1",
				Amount:     6969,
				Date:       "2006-01-02",
				CategoryID: user1Categories[0].ID,
				Note:       "Expense 1 Note",
			},
			input: expense.UpdateExpenseReq{
				Name: toPtr(t, "Expense Uno"),
			},
			want: expense.Expense{
				Name:      "Expense Uno",
				Amount:    6969,
				Date:      time.Date(2006, time.January, 2, 0, 0, 0, 0, time.UTC),
				Note:      "Expense 1 Note",
				Category:  user1Categories[0],
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		{
			name: fmt.Sprintf("update %s's expense amount", user1.Name),
			user: user1,
			expense: expense.CreateExpenseReq{
				Name:       "Expense 1",
				Amount:     6969,
				Date:       "2006-01-02",
				CategoryID: user1Categories[0].ID,
				Note:       "Expense 1 Note",
			},
			input: expense.UpdateExpenseReq{
				Amount: toPtr(t, int64(7000)),
			},
			want: expense.Expense{
				Name:      "Expense 1",
				Amount:    7000,
				Date:      time.Date(2006, time.January, 2, 0, 0, 0, 0, time.UTC),
				Note:      "Expense 1 Note",
				Category:  user1Categories[0],
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		{
			name: fmt.Sprintf("update %s's expense date", user1.Name),
			user: user1,
			expense: expense.CreateExpenseReq{
				Name:       "Expense 1",
				Amount:     6969,
				Date:       "2006-01-02",
				CategoryID: user1Categories[0].ID,
				Note:       "Expense 1 Note",
			},
			input: expense.UpdateExpenseReq{
				Date: toPtr(t, "2006-01-10"),
			},
			want: expense.Expense{
				Name:      "Expense 1",
				Amount:    6969,
				Date:      time.Date(2006, time.January, 10, 0, 0, 0, 0, time.UTC),
				Note:      "Expense 1 Note",
				Category:  user1Categories[0],
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		{
			name: fmt.Sprintf("update %s's expense category id", user1.Name),
			user: user1,
			expense: expense.CreateExpenseReq{
				Name:       "Expense 1",
				Amount:     6969,
				Date:       "2006-01-02",
				CategoryID: user1Categories[0].ID,
				Note:       "Expense 1 Note",
			},
			input: expense.UpdateExpenseReq{
				CategoryID: &user1Categories[1].ID,
			},
			want: expense.Expense{
				Name:      "Expense 1",
				Amount:    6969,
				Date:      time.Date(2006, time.January, 2, 0, 0, 0, 0, time.UTC),
				Note:      "Expense 1 Note",
				Category:  user1Categories[1],
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		{
			name: fmt.Sprintf("update %s's expense note", user1.Name),
			user: user1,
			expense: expense.CreateExpenseReq{
				Name:       "Expense 1",
				Amount:     6969,
				Date:       "2006-01-02",
				CategoryID: user1Categories[0].ID,
				Note:       "Expense 1 Note",
			},
			input: expense.UpdateExpenseReq{
				Note: toPtr(t, "Expense Uno Noto"),
			},
			want: expense.Expense{
				Name:      "Expense 1",
				Amount:    6969,
				Date:      time.Date(2006, time.January, 2, 0, 0, 0, 0, time.UTC),
				Note:      "Expense Uno Noto",
				Category:  user1Categories[0],
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		{
			name: fmt.Sprintf("update %s's expense", user1.Name),
			user: user1,
			expense: expense.CreateExpenseReq{
				Name:       "Expense 1",
				Amount:     6969,
				Date:       "2006-01-02",
				CategoryID: user1Categories[0].ID,
				Note:       "Expense 1 Note",
			},
			input: expense.UpdateExpenseReq{
				Name:       toPtr(t, "Expense Uno"),
				Amount:     toPtr(t, int64(7000)),
				Date:       toPtr(t, "2006-01-10"),
				CategoryID: &user1Categories[1].ID,
				Note:       toPtr(t, "Expense Uno Noto"),
			},
			want: expense.Expense{
				Name:      "Expense Uno",
				Amount:    7000,
				Date:      time.Date(2006, time.January, 10, 0, 0, 0, 0, time.UTC),
				Note:      "Expense Uno Noto",
				Category:  user1Categories[1],
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctxWithUser := user.ContextWithUser(ctxWithLogger, test.user)
			createdExpense, err := er.CreateExpense(ctxWithUser, test.expense)
			assert.Nil(t, err)

			test.input.ID = createdExpense.ID

			updatedExpense, err := er.UpdateExpense(ctxWithUser, test.input)
			assert.Nil(t, err)
			assertExpense(t, test.want, updatedExpense)
		})
	}

	t.Run("expense not found", func(t *testing.T) {
		ctxWithUser := user.ContextWithUser(ctxWithLogger, user1)
		foundExpense, err := er.UpdateExpense(ctxWithUser, expense.UpdateExpenseReq{ID: "123", Name: toPtr(t, "Foo")})
		assert.Equal(t, expense.Expense{}, foundExpense)
		assert.Equal(t, internal.NewError(internal.ErrorCodeNotFound, "Expense not found"), err)
	})

	t.Run("can't update expense of other user", func(t *testing.T) {
		ctxWithUser := user.ContextWithUser(ctxWithLogger, user1)
		createdExpense, err := er.CreateExpense(ctxWithUser, tests[0].expense)
		assert.Nil(t, err)

		ctxWithUser = user.ContextWithUser(ctxWithLogger, user2)
		foundExpense, err := er.UpdateExpense(ctxWithUser, expense.UpdateExpenseReq{ID: createdExpense.ID, Name: toPtr(t, "Foo")})
		assert.Equal(t, expense.Expense{}, foundExpense)
		assert.Equal(t, internal.NewError(internal.ErrorCodeNotFound, "Expense not found"), err)
	})
}

func TestDeleteExpense(t *testing.T) {
	dh := newDBHelper(t, "test_delete_expense.db")
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
		name    string
		user    user.User
		expense expense.CreateExpenseReq
	}{
		{
			name: fmt.Sprintf("%s's expense", user1.Name),
			user: user1,
			expense: expense.CreateExpenseReq{
				Name:       "Expense 1",
				Amount:     6969,
				Date:       "2006-01-02",
				CategoryID: user1Categories[0].ID,
				Note:       "Expense 1 Note",
			},
		},
		{
			name: fmt.Sprintf("%s's expense", user2.Name),
			user: user2,
			expense: expense.CreateExpenseReq{
				Name:       "Expense 1",
				Amount:     6969,
				Date:       "2006-01-02",
				CategoryID: user2Categories[0].ID,
				Note:       "Expense 1 Note",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctxWithUser := user.ContextWithUser(ctxWithLogger, test.user)
			createdExpense, err := er.CreateExpense(ctxWithUser, test.expense)
			assert.Nil(t, err)

			err = er.DeleteExpense(ctxWithUser, createdExpense.ID)
			assert.Nil(t, err)

			_, err = er.ExpenseByID(ctxWithUser, createdExpense.ID)
			assert.Equal(t, err, internal.NewError(internal.ErrorCodeNotFound, "Expense not found"))
		})
	}

	t.Run("can't delete expense of other user", func(t *testing.T) {
		ctxWithUser1 := user.ContextWithUser(ctxWithLogger, user1)
		createdExpense, err := er.CreateExpense(ctxWithUser1, tests[0].expense)
		assert.Nil(t, err)

		ctxWithUser2 := user.ContextWithUser(ctxWithLogger, user2)
		err = er.DeleteExpense(ctxWithUser2, createdExpense.ID)
		assert.Nil(t, err)

		foundExpense, err := er.ExpenseByID(ctxWithUser1, createdExpense.ID)
		assert.Nil(t, err)
		assert.Equal(t, createdExpense, foundExpense)
	})
}

func assertExpense(t *testing.T, want, got expense.Expense) {
	t.Helper()

	assert.True(t, got.ID != "")

	// to make it easier to assert
	got.ID = ""

	assert.WithinDuration(t, want.CreatedAt, got.CreatedAt, time.Second)
	assert.WithinDuration(t, want.UpdatedAt, got.UpdatedAt, time.Second)

	// to make it easier to assert
	got.CreatedAt = time.Time{}
	got.UpdatedAt = time.Time{}
	want.CreatedAt = time.Time{}
	want.UpdatedAt = time.Time{}

	assert.Equal(t, want, got)
}
