package sqlite_test

import (
	"context"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/cativovo/budget-tracker/internal"
	"github.com/cativovo/budget-tracker/internal/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestCreateExpense(t *testing.T) {
	db, c := mustConnectDB(t)
	defer c()

	es := sqlite.NewExpenseService(db)
	sharedCtx := internal.ContextWithLogger(context.Background(), zap.NewNop().Sugar())

	user := mustCreateUser(sharedCtx, t, db, internal.UserCreate{
		ID:    "RB-1",
		Name:  "Max Verstappen",
		Email: "maxverstappen@redbull.com",
	})

	sharedCtx = internal.ContextWithUser(sharedCtx, user)

	category1 := mustCreateCategory(internal.ContextWithUser(sharedCtx, user), t, db, internal.CategoryCreate{
		Name:  "user1-category1",
		Color: "#FFFFFF",
		Icon:  "icon1",
	})
	category2 := mustCreateCategory(internal.ContextWithUser(sharedCtx, user), t, db, internal.CategoryCreate{
		Name:  "user1-category2",
		Color: "#000000",
		Icon:  "icon2",
	})

	testCases := []struct {
		name          string
		user          internal.User
		category      internal.Category
		expenseCreate internal.ExpenseCreate
		err           error
	}{
		{
			name:     "user category 1 ok",
			user:     user,
			category: category1,
			expenseCreate: internal.ExpenseCreate{
				Name:       "Cat food",
				Amount:     2000,
				Date:       "2025-06-20",
				Note:       "Test note",
				CategoryID: category1.ID,
			},
		},
		{
			name:     "user category 2 ok",
			user:     user,
			category: category2,
			expenseCreate: internal.ExpenseCreate{
				Name:       "IRacing subscription",
				Amount:     681893,
				Date:       "2025-01-01",
				CategoryID: category2.ID,
			},
		},
		{
			name: "invalid category",
			user: user,
			expenseCreate: internal.ExpenseCreate{
				Name:       "Invalid",
				Amount:     3000,
				Date:       "2025-03-15",
				CategoryID: "6969",
			},
			err: internal.NewError(internal.ErrorCodeInvalid, "category not found"),
		},
		{
			name: "empty name",
			user: user,
			expenseCreate: internal.ExpenseCreate{
				Amount:     3000,
				Date:       "2025-03-15",
				CategoryID: category1.ID,
			},
			err: internal.NewError(internal.ErrorCodeInvalid, "name is required"),
		},
		{
			name: "amount <= 0",
			user: user,
			expenseCreate: internal.ExpenseCreate{
				Name:       "Test",
				Amount:     0,
				Date:       "2025-03-15",
				CategoryID: category1.ID,
			},
			err: internal.NewError(internal.ErrorCodeInvalid, "amount must be greater than 0"),
		},
		{
			name: "empty date",
			user: user,
			expenseCreate: internal.ExpenseCreate{
				Name:       "Test",
				Amount:     6000,
				CategoryID: category1.ID,
			},
			err: internal.NewError(internal.ErrorCodeInvalid, "invalid date"),
		},
		{
			name: "Invalid date",
			user: user,
			expenseCreate: internal.ExpenseCreate{
				Name:       "Test",
				Date:       "1-3891",
				Amount:     6000,
				CategoryID: category1.ID,
			},
			err: internal.NewError(internal.ErrorCodeInvalid, "invalid date"),
		},
		{
			name: "empty category ID",
			user: user,
			expenseCreate: internal.ExpenseCreate{
				Name:   "Test",
				Date:   "2025-03-15",
				Amount: 6000,
			},
			err: internal.NewError(internal.ErrorCodeInvalid, "category_id is required"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.err != nil {
				createdExpense, err := es.CreateExpense(sharedCtx, tc.expenseCreate)
				assert.Equal(t, internal.GetErrorCode(tc.err), internal.GetErrorCode(err))
				assert.Equal(t, internal.GetErrorMessage(tc.err), internal.GetErrorMessage(err))

				_, err = es.GetExpense(sharedCtx, createdExpense.ID)
				assert.Equal(t, internal.ErrorCodeNotFound, internal.GetErrorCode(err))
				assert.Equal(t, "expense not found", internal.GetErrorMessage(err))
				return
			}

			createdExpense, err := es.CreateExpense(sharedCtx, tc.expenseCreate)
			assert.Nil(t, err)
			assert.NotEqual(t, "", createdExpense.ID)
			assert.Equal(t, tc.expenseCreate.Name, createdExpense.Name)
			assert.Equal(t, tc.expenseCreate.Amount, createdExpense.Amount)
			assert.Equal(t, tc.expenseCreate.Note, createdExpense.Note)
			assert.Equal(t, tc.expenseCreate.Date, createdExpense.Date.Format(time.DateOnly))
			assert.WithinDuration(t, createdExpense.CreatedAt, time.Now(), wantDelta)
			assert.WithinDuration(t, createdExpense.UpdatedAt, time.Now(), wantDelta)
			assert.Equal(t, tc.category, createdExpense.Category)

			gotExpense, err := es.GetExpense(sharedCtx, createdExpense.ID)
			assert.Nil(t, err)
			assert.Equal(t, createdExpense, gotExpense)
		})
	}

	t.Run("validate all fields", func(t *testing.T) {
		createdExpense, err := es.CreateExpense(sharedCtx, internal.ExpenseCreate{})
		assert.Equal(t, internal.ErrorCodeInvalid, internal.GetErrorCode(err))

		gotMsgs := strings.Split(internal.GetErrorMessage(err), ", ")
		wantMsgs := []string{
			"name is required",
			"amount must be greater than 0",
			"invalid date",
			"category_id is required",
		}
		slices.Sort(gotMsgs)
		slices.Sort(wantMsgs)
		assert.Equal(t, wantMsgs, gotMsgs)

		_, err = es.GetExpense(sharedCtx, createdExpense.ID)
		assert.Equal(t, internal.ErrorCodeNotFound, internal.GetErrorCode(err))
		assert.Equal(t, "expense not found", internal.GetErrorMessage(err))
	})

	t.Run("invalid access", func(t *testing.T) {
		createdExpense := mustCreateExpense(sharedCtx, t, db, internal.ExpenseCreate{
			Name:       "Drinks",
			Amount:     5000,
			Date:       "2025-10-11",
			Note:       "Some note",
			CategoryID: category1.ID,
		})

		someUser := mustCreateUser(sharedCtx, t, db, internal.UserCreate{
			ID:    "2",
			Name:  "Yuki Tsunoda",
			Email: "yukitsunoda@redbull.com",
		})
		_, err := es.GetExpense(internal.ContextWithUser(sharedCtx, someUser), createdExpense.ID)
		assert.Equal(t, internal.ErrorCodeNotFound, internal.GetErrorCode(err))
		assert.Equal(t, "expense not found", internal.GetErrorMessage(err))
	})
}

func TestUpdateExpense(t *testing.T) {
	db, c := mustConnectDB(t)
	defer c()

	es := sqlite.NewExpenseService(db)
	sharedCtx := internal.ContextWithLogger(context.Background(), zap.NewNop().Sugar())

	user := mustCreateUser(sharedCtx, t, db, internal.UserCreate{
		ID:    "RB-1",
		Name:  "Max Verstappen",
		Email: "maxverstappen@redbull.com",
	})

	sharedCtx = internal.ContextWithUser(sharedCtx, user)

	category1 := mustCreateCategory(internal.ContextWithUser(sharedCtx, user), t, db, internal.CategoryCreate{
		Name:  "user1-category1",
		Color: "#FFFFFF",
		Icon:  "icon1",
	})
	category2 := mustCreateCategory(internal.ContextWithUser(sharedCtx, user), t, db, internal.CategoryCreate{
		Name:  "user1-category2",
		Color: "#000000",
		Icon:  "icon2",
	})

	testCases := []struct {
		name          string
		user          internal.User
		category      internal.Category
		expenseCreate internal.ExpenseCreate
		expenseUpdate internal.ExpenseUpdate
		err           error
	}{
		{
			name:     "update name",
			user:     user,
			category: category1,
			expenseCreate: internal.ExpenseCreate{
				Name:       "Food",
				Amount:     4000,
				Date:       "2025-01-20",
				Note:       "Some food",
				CategoryID: category1.ID,
			},
			expenseUpdate: internal.ExpenseUpdate{
				Name: ptr("Updated name"),
			},
		},
		{
			name:     "update amount",
			user:     user,
			category: category1,
			expenseCreate: internal.ExpenseCreate{
				Name:       "Food",
				Amount:     4000,
				Date:       "2025-01-20",
				Note:       "Some food",
				CategoryID: category1.ID,
			},
			expenseUpdate: internal.ExpenseUpdate{
				Amount: ptr(internal.Amount(5000)),
			},
		},
		{
			name:     "update note",
			user:     user,
			category: category1,
			expenseCreate: internal.ExpenseCreate{
				Name:       "Food",
				Amount:     4000,
				Date:       "2025-01-20",
				Note:       "Some food",
				CategoryID: category1.ID,
			},
			expenseUpdate: internal.ExpenseUpdate{
				Note: ptr("Updated note"),
			},
		},
		{
			name:     "update date",
			user:     user,
			category: category1,
			expenseCreate: internal.ExpenseCreate{
				Name:       "Food",
				Amount:     4000,
				Date:       "2025-01-20",
				Note:       "Some food",
				CategoryID: category1.ID,
			},
			expenseUpdate: internal.ExpenseUpdate{
				Date: ptr("2025-02-01"),
			},
		},
		{
			name:     "update category",
			user:     user,
			category: category1,
			expenseCreate: internal.ExpenseCreate{
				Name:       "Food",
				Amount:     4000,
				Date:       "2025-01-20",
				Note:       "Some food",
				CategoryID: category1.ID,
			},
			expenseUpdate: internal.ExpenseUpdate{
				CategoryID: ptr(category2.ID),
			},
		},
		{
			name:     "empty name",
			user:     user,
			category: category1,
			expenseCreate: internal.ExpenseCreate{
				Name:       "Food",
				Amount:     4000,
				Date:       "2025-01-20",
				Note:       "Some food",
				CategoryID: category1.ID,
			},
			expenseUpdate: internal.ExpenseUpdate{
				Name: ptr(""),
			},
			err: internal.NewError(internal.ErrorCodeInvalid, "name is required"),
		},
		{
			name:     "empty note",
			user:     user,
			category: category1,
			expenseCreate: internal.ExpenseCreate{
				Name:       "Food",
				Amount:     4000,
				Date:       "2025-01-20",
				Note:       "Some food",
				CategoryID: category1.ID,
			},
			expenseUpdate: internal.ExpenseUpdate{
				Note: ptr(""),
			},
		},
		{
			name:     "amount <= 0",
			user:     user,
			category: category1,
			expenseCreate: internal.ExpenseCreate{
				Name:       "Food",
				Amount:     4000,
				Date:       "2025-01-20",
				Note:       "Some food",
				CategoryID: category1.ID,
			},
			expenseUpdate: internal.ExpenseUpdate{
				Amount: ptr(internal.Amount(-1)),
			},
			err: internal.NewError(internal.ErrorCodeInvalid, "amount must be greater than 0"),
		},
		{
			name:     "empty date",
			user:     user,
			category: category1,
			expenseCreate: internal.ExpenseCreate{
				Name:       "Food",
				Amount:     4000,
				Date:       "2025-01-20",
				Note:       "Some food",
				CategoryID: category1.ID,
			},
			expenseUpdate: internal.ExpenseUpdate{
				Date: ptr(""),
			},
			err: internal.NewError(internal.ErrorCodeInvalid, "invalid date"),
		},
		{
			name:     "invalid date",
			user:     user,
			category: category1,
			expenseCreate: internal.ExpenseCreate{
				Name:       "Food",
				Amount:     4000,
				Date:       "2025-01-20",
				Note:       "Some food",
				CategoryID: category1.ID,
			},
			expenseUpdate: internal.ExpenseUpdate{
				Date: ptr("-871-234"),
			},
			err: internal.NewError(internal.ErrorCodeInvalid, "invalid date"),
		},
		{
			name:     "empty category ID",
			user:     user,
			category: category1,
			expenseCreate: internal.ExpenseCreate{
				Name:       "Food",
				Amount:     4000,
				Date:       "2025-01-20",
				Note:       "Some food",
				CategoryID: category1.ID,
			},
			expenseUpdate: internal.ExpenseUpdate{
				CategoryID: ptr(""),
			},
			err: internal.NewError(internal.ErrorCodeInvalid, "category_id is required"),
		},
		{
			name:     "invaild category ID",
			user:     user,
			category: category1,
			expenseCreate: internal.ExpenseCreate{
				Name:       "Food",
				Amount:     4000,
				Date:       "2025-01-20",
				Note:       "Some food",
				CategoryID: category1.ID,
			},
			expenseUpdate: internal.ExpenseUpdate{
				CategoryID: ptr("123456789"),
			},
			err: internal.NewError(internal.ErrorCodeInvalid, "category not found"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			createdExpense := mustCreateExpense(sharedCtx, t, db, tc.expenseCreate)
			tc.expenseUpdate.ID = createdExpense.ID

			if tc.err != nil {
				_, err := es.UpdateExpense(sharedCtx, tc.expenseUpdate)
				assert.Equal(t, internal.GetErrorCode(tc.err), internal.GetErrorCode(err))
				assert.Equal(t, internal.GetErrorMessage(tc.err), internal.GetErrorMessage(err))

				gotExpense := mustGetExpense(sharedCtx, t, db, createdExpense.ID)
				assert.Equal(t, createdExpense, gotExpense)
				return
			}

			updatedExpense, err := es.UpdateExpense(sharedCtx, tc.expenseUpdate)
			assert.Nil(t, err)

			assertUpdatedField(t, tc.expenseUpdate.Name, createdExpense.Name, updatedExpense.Name)
			assertUpdatedField(t, tc.expenseUpdate.Amount, createdExpense.Amount, updatedExpense.Amount)
			assertUpdatedField(t, tc.expenseUpdate.Note, createdExpense.Note, updatedExpense.Note)
			assertUpdatedField(t, tc.expenseUpdate.Date, createdExpense.Date.Format(time.DateOnly), updatedExpense.Date.Format(time.DateOnly))
			assertUpdatedField(t, tc.expenseUpdate.CategoryID, createdExpense.Category.ID, updatedExpense.Category.ID)
			assert.Equal(t, createdExpense.CreatedAt, updatedExpense.CreatedAt)

			// TODO: find a way to test UpdatedAt field without slowing down the tests
			// it's too fast that the UpdatedAt is still the same with CreatedAt

			gotExpense := mustGetExpense(sharedCtx, t, db, createdExpense.ID)
			assert.Equal(t, updatedExpense, gotExpense)
		})
	}

	t.Run("update correct expense", func(t *testing.T) {
		expense1 := internal.ExpenseCreate{
			Name:       "Drinks",
			Amount:     4000,
			Date:       "2024-05-04",
			Note:       "some drinks",
			CategoryID: category1.ID,
		}
		expense2 := internal.ExpenseCreate{
			Name:       "Entertainment",
			Amount:     80000,
			Date:       "2024-08-08",
			Note:       "",
			CategoryID: category1.ID,
		}
		createdExpense1 := mustCreateExpense(sharedCtx, t, db, expense1)
		createdExpense2 := mustCreateExpense(sharedCtx, t, db, expense2)

		updatedExpense, err := es.UpdateExpense(sharedCtx, internal.ExpenseUpdate{
			ID:   createdExpense1.ID,
			Name: ptr(createdExpense1.ID),
		})
		assert.Nil(t, err)
		assert.NotEqual(t, createdExpense1, updatedExpense)

		gotExpense := mustGetExpense(sharedCtx, t, db, createdExpense2.ID)
		assert.Equal(t, createdExpense2, gotExpense)
	})

	t.Run("no id", func(t *testing.T) {
		_, err := es.UpdateExpense(sharedCtx, internal.ExpenseUpdate{
			Name: ptr("sample name"),
		})
		assert.Equal(t, internal.ErrorCodeInvalid, internal.GetErrorCode(err))
		assert.Equal(t, "id is required", internal.GetErrorMessage(err))
	})

	t.Run("validate all fields", func(t *testing.T) {
		_, err := es.UpdateExpense(sharedCtx, internal.ExpenseUpdate{
			ID:         "",
			Name:       ptr(""),
			Amount:     ptr(internal.Amount(0)),
			Note:       ptr(""), // no error
			Date:       ptr("-878"),
			CategoryID: ptr(""),
		})

		gotMsgs := strings.Split(internal.GetErrorMessage(err), ", ")
		wantMsgs := []string{
			"id is required",
			"name is required",
			"amount must be greater than 0",
			"invalid date",
			"category_id is required",
		}
		slices.Sort(gotMsgs)
		slices.Sort(wantMsgs)
		assert.Equal(t, wantMsgs, gotMsgs)
	})

	t.Run("invalid access", func(t *testing.T) {
		createdExpense := mustCreateExpense(sharedCtx, t, db, internal.ExpenseCreate{
			Name:       "Drinks",
			Amount:     5000,
			Date:       "2025-10-11",
			Note:       "Some note",
			CategoryID: category1.ID,
		})

		someUser := mustCreateUser(sharedCtx, t, db, internal.UserCreate{
			ID:    "2",
			Name:  "Yuki Tsunoda",
			Email: "yukitsunoda@redbull.com",
		})
		_, err := es.UpdateExpense(internal.ContextWithUser(sharedCtx, someUser), internal.ExpenseUpdate{
			ID:   createdExpense.ID,
			Name: ptr("Updated name"),
		})
		assert.Equal(t, internal.ErrorCodeNotFound, internal.GetErrorCode(err))
		assert.Equal(t, "expense not found", internal.GetErrorMessage(err))

		gotExpense := mustGetExpense(sharedCtx, t, db, createdExpense.ID)
		assert.Equal(t, createdExpense, gotExpense)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := es.UpdateExpense(sharedCtx, internal.ExpenseUpdate{
			ID:   "6969",
			Name: ptr("Update"),
		})
		assert.Equal(t, internal.ErrorCodeNotFound, internal.GetErrorCode(err))
		assert.Equal(t, "expense not found", internal.GetErrorMessage(err))
	})
}

func TestDeleteExpense(t *testing.T) {
	db, c := mustConnectDB(t)
	defer c()

	es := sqlite.NewExpenseService(db)
	sharedCtx := internal.ContextWithLogger(context.Background(), zap.NewNop().Sugar())

	user := mustCreateUser(sharedCtx, t, db, internal.UserCreate{
		ID:    "RB-1",
		Name:  "Max Verstappen",
		Email: "maxverstappen@redbull.com",
	})

	sharedCtx = internal.ContextWithUser(sharedCtx, user)

	category := mustCreateCategory(internal.ContextWithUser(sharedCtx, user), t, db, internal.CategoryCreate{
		Name:  "user1-category1",
		Color: "#FFFFFF",
		Icon:  "icon1",
	})

	t.Run("delete correct expense", func(t *testing.T) {
		createdExpense1 := mustCreateExpense(sharedCtx, t, db, internal.ExpenseCreate{
			Name:       "Expense 1",
			Amount:     7000,
			Date:       "2024-10-25",
			Note:       "some note",
			CategoryID: category.ID,
		})
		createdExpense2 := mustCreateExpense(sharedCtx, t, db, internal.ExpenseCreate{
			Name:       "Expense 2",
			Amount:     1000,
			Date:       "2024-06-24",
			CategoryID: category.ID,
		})

		err := es.DeleteExpense(sharedCtx, createdExpense1.ID)
		assert.Nil(t, err)

		_, err = es.GetExpense(sharedCtx, createdExpense1.ID)
		assert.Equal(t, internal.ErrorCodeNotFound, internal.GetErrorCode(err))
		assert.Equal(t, "expense not found", internal.GetErrorMessage(err))

		gotExpense, err := es.GetExpense(sharedCtx, createdExpense2.ID)
		assert.Nil(t, err)
		assert.Equal(t, createdExpense2, gotExpense)
	})

	t.Run("invalid access", func(t *testing.T) {
		createdExpense := mustCreateExpense(sharedCtx, t, db, internal.ExpenseCreate{
			Name:       "Expense 1",
			Amount:     7000,
			Date:       "2024-10-25",
			Note:       "some note",
			CategoryID: category.ID,
		})

		someUser := mustCreateUser(sharedCtx, t, db, internal.UserCreate{
			ID:    "RB-2",
			Name:  "Yuki Tsunoda",
			Email: "yukitsunoda@redbull.com",
		})
		err := es.DeleteExpense(internal.ContextWithUser(sharedCtx, someUser), createdExpense.ID)
		assert.Nil(t, err)

		gotExpense, err := es.GetExpense(sharedCtx, createdExpense.ID)
		assert.Nil(t, err)
		assert.Equal(t, createdExpense, gotExpense)
	})

	t.Run("not found", func(t *testing.T) {
		err := es.DeleteExpense(sharedCtx, "69696")
		assert.Nil(t, err)
	})
}

func TestCreateExpenseGroup(t *testing.T) {
	db, c := mustConnectDB(t)
	defer c()

	es := sqlite.NewExpenseService(db)
	sharedCtx := internal.ContextWithLogger(context.Background(), zap.NewNop().Sugar())

	user := mustCreateUser(sharedCtx, t, db, internal.UserCreate{
		ID:    "RB-1",
		Name:  "Max Verstappen",
		Email: "maxverstappen@redbull.com",
	})

	sharedCtx = internal.ContextWithUser(sharedCtx, user)

	category := mustCreateCategory(internal.ContextWithUser(sharedCtx, user), t, db, internal.CategoryCreate{
		Name:  "user1-category1",
		Color: "#FFFFFF",
		Icon:  "icon1",
	})

	testCases := []struct {
		name               string
		expenseGroupCreate internal.ExpenseGroupCreate
		category           internal.Category
		err                error
	}{
		{
			name: "ok",
			expenseGroupCreate: internal.ExpenseGroupCreate{
				Name: "group1",
				Date: "2025-01-24",
				Expenses: []internal.ExpenseGroupExpenseCreate{
					{
						Name:   "expense 1",
						Amount: 6969,
						Note:   "expense note 1",
					},
					{
						Name:   "expense 2",
						Amount: 420,
						Note:   "expense note 2",
					},
				},
				Note:       "group note",
				CategoryID: category.ID,
			},
			category: category,
		},
		{
			name: "empty expenses",
			expenseGroupCreate: internal.ExpenseGroupCreate{
				Name:       "group1",
				Date:       "2025-01-24",
				Expenses:   []internal.ExpenseGroupExpenseCreate{},
				Note:       "group note",
				CategoryID: category.ID,
			},
			category: category,
		},
		{
			name: "empty name",
			expenseGroupCreate: internal.ExpenseGroupCreate{
				Date: "2025-01-24",
				Expenses: []internal.ExpenseGroupExpenseCreate{
					{
						Name:   "expense 1",
						Amount: 6969,
						Note:   "expense note 1",
					},
					{
						Name:   "expense 2",
						Amount: 420,
						Note:   "expense note 2",
					},
				},
				Note:       "group note",
				CategoryID: category.ID,
			},
			err: internal.NewError(internal.ErrorCodeInvalid, "name is required"),
		},
		{
			name: "invalid date",
			expenseGroupCreate: internal.ExpenseGroupCreate{
				Name: "group1",
				Date: "202501-24",
				Expenses: []internal.ExpenseGroupExpenseCreate{
					{
						Name:   "expense 1",
						Amount: 6969,
						Note:   "expense note 1",
					},
					{
						Name:   "expense 2",
						Amount: 420,
						Note:   "expense note 2",
					},
				},
				Note:       "group note",
				CategoryID: category.ID,
			},
			err: internal.NewError(internal.ErrorCodeInvalid, "invalid date"),
		},
		{
			name: "empty category id",
			expenseGroupCreate: internal.ExpenseGroupCreate{
				Name: "group1",
				Date: "2025-01-24",
				Expenses: []internal.ExpenseGroupExpenseCreate{
					{
						Name:   "expense 1",
						Amount: 6969,
						Note:   "expense note 1",
					},
					{
						Name:   "expense 2",
						Amount: 420,
						Note:   "expense note 2",
					},
				},
				Note: "group note",
			},
			err: internal.NewError(internal.ErrorCodeInvalid, "category_id is required"),
		},
		{
			name: "nil expenses",
			expenseGroupCreate: internal.ExpenseGroupCreate{
				Name:       "group1",
				Date:       "2025-01-24",
				Expenses:   nil,
				Note:       "group note",
				CategoryID: category.ID,
			},
			err: internal.NewError(internal.ErrorCodeInvalid, "expenses is required"),
		},
		{
			name: "empty expense name",
			expenseGroupCreate: internal.ExpenseGroupCreate{
				Name: "group1",
				Date: "2025-01-24",
				Expenses: []internal.ExpenseGroupExpenseCreate{
					{
						Amount: 6969,
						Note:   "expense note 1",
					},
				},
				Note:       "group note",
				CategoryID: category.ID,
			},
			err: internal.NewError(internal.ErrorCodeInvalid, "name is required"),
		},
		{
			name: "amount <= 0",
			expenseGroupCreate: internal.ExpenseGroupCreate{
				Name: "group1",
				Date: "2025-01-24",
				Expenses: []internal.ExpenseGroupExpenseCreate{
					{
						Name: "expense 1",
						Note: "expense note 1",
					},
				},
				Note:       "group note",
				CategoryID: category.ID,
			},
			err: internal.NewError(internal.ErrorCodeInvalid, "amount must be greater than 0"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.err != nil {
				createdExpenseGroup, err := es.CreateExpenseGroup(sharedCtx, tc.expenseGroupCreate)
				assert.Equal(t, internal.GetErrorCode(tc.err), internal.GetErrorCode(err))
				assert.Equal(t, internal.GetErrorMessage(tc.err), internal.GetErrorMessage(err))

				_, err = es.GetExpenseGroup(sharedCtx, createdExpenseGroup.ID)
				assert.Equal(t, internal.ErrorCodeNotFound, internal.GetErrorCode(err))
				assert.Equal(t, "expense group not found", internal.GetErrorMessage(err))
				return
			}

			createdExpenseGroup, err := es.CreateExpenseGroup(sharedCtx, tc.expenseGroupCreate)
			assert.Nil(t, err)
			assert.Equal(t, tc.expenseGroupCreate.Name, createdExpenseGroup.Name)
			assert.Equal(t, tc.expenseGroupCreate.Date, createdExpenseGroup.Date.Format(time.DateOnly))
			assert.Equal(t, tc.expenseGroupCreate.Note, createdExpenseGroup.Note)
			assert.Equal(t, tc.category, createdExpenseGroup.Category)
			assert.WithinDuration(t, createdExpenseGroup.CreatedAt, time.Now(), wantDelta)
			assert.WithinDuration(t, createdExpenseGroup.UpdatedAt, time.Now(), wantDelta)

			assert.Len(t, createdExpenseGroup.Expenses, len(tc.expenseGroupCreate.Expenses))

			for i, v := range tc.expenseGroupCreate.Expenses {
				createdExpense := createdExpenseGroup.Expenses[i]
				assert.Equal(t, v.Name, createdExpense.Name)
				assert.Equal(t, v.Amount, createdExpense.Amount)
				assert.Equal(t, v.Note, createdExpense.Note)
				assert.Equal(t, tc.expenseGroupCreate.Date, createdExpense.Date.Format(time.DateOnly))
				assert.Equal(t, tc.expenseGroupCreate.CategoryID, createdExpense.Category.ID)
				assert.WithinDuration(t, createdExpense.CreatedAt, time.Now(), wantDelta)
				assert.WithinDuration(t, createdExpense.UpdatedAt, time.Now(), wantDelta)

				gotExpense, err := es.GetExpense(sharedCtx, createdExpense.ID)
				assert.Nil(t, err)
				assert.Equal(t, createdExpense, gotExpense)
			}

			gotExpenseGroup, err := es.GetExpenseGroup(sharedCtx, createdExpenseGroup.ID)
			assert.Nil(t, err)
			assert.Equal(t, createdExpenseGroup, gotExpenseGroup)
		})
	}

	t.Run("validate all top level fields", func(t *testing.T) {
		createdExpenseGroup, err := es.CreateExpenseGroup(sharedCtx, internal.ExpenseGroupCreate{})
		assert.Equal(t, internal.ErrorCodeInvalid, internal.GetErrorCode(err))

		gotMsgs := strings.Split(internal.GetErrorMessage(err), ", ")
		wantMsgs := []string{
			"name is required",
			"invalid date",
			"category_id is required",
			"expenses is required",
		}

		slices.Sort(gotMsgs)
		slices.Sort(wantMsgs)
		assert.Equal(t, wantMsgs, gotMsgs)

		_, err = es.GetExpenseGroup(sharedCtx, createdExpenseGroup.ID)
		assert.Equal(t, internal.ErrorCodeNotFound, internal.GetErrorCode(err))
	})

	t.Run("validate all expenses fields", func(t *testing.T) {
		createdExpenseGroup, err := es.CreateExpenseGroup(sharedCtx, internal.ExpenseGroupCreate{
			Name: "group1",
			Date: "2025-01-24",
			Expenses: []internal.ExpenseGroupExpenseCreate{
				{},
			},
			Note:       "group note",
			CategoryID: category.ID,
		})

		assert.Equal(t, internal.ErrorCodeInvalid, internal.GetErrorCode(err))

		gotMsgs := strings.Split(internal.GetErrorMessage(err), ", ")
		wantMsgs := []string{
			"name is required",
			"amount must be greater than 0",
		}

		slices.Sort(gotMsgs)
		slices.Sort(wantMsgs)
		assert.Equal(t, wantMsgs, gotMsgs)

		_, err = es.GetExpenseGroup(sharedCtx, createdExpenseGroup.ID)
		assert.Equal(t, internal.ErrorCodeNotFound, internal.GetErrorCode(err))
	})

	t.Run("invalid access", func(t *testing.T) {
		createdExpenseGroup, err := es.CreateExpenseGroup(sharedCtx, internal.ExpenseGroupCreate{
			Name: "group1",
			Date: "2025-01-24",
			Expenses: []internal.ExpenseGroupExpenseCreate{
				{
					Name:   "expense 1",
					Amount: 6969,
					Note:   "expense note 1",
				},
				{
					Name:   "expense 2",
					Amount: 420,
					Note:   "expense note 2",
				},
			},
			Note:       "group note",
			CategoryID: category.ID,
		})
		assert.Nil(t, err)

		someUser := mustCreateUser(sharedCtx, t, db, internal.UserCreate{
			ID:    "RB-2",
			Name:  "Yuki Tsunoda",
			Email: "yukitsunoda@redbull.com",
		})
		_, err = es.GetExpenseGroup(internal.ContextWithUser(sharedCtx, someUser), createdExpenseGroup.ID)
		assert.Equal(t, internal.ErrorCodeNotFound, internal.GetErrorCode(err))
	})
}

//	func TestUpdateExpenseGroup(t *testing.T) {
//		db, c := mustConnectDB(t)
//		defer c()
//
//		es := sqlite.NewExpenseService(db)
//		sharedCtx := internal.ContextWithLogger(context.Background(), zap.NewNop().Sugar())
//
//		user := mustCreateUser(sharedCtx, t, db, internal.UserCreate{
//			ID:    "RB-1",
//			Name:  "Max Verstappen",
//			Email: "maxverstappen@redbull.com",
//		})
//
//		sharedCtx = internal.ContextWithUser(sharedCtx, user)
//
//		category := mustCreateCategory(internal.ContextWithUser(sharedCtx, user), t, db, internal.CategoryCreate{
//			Name:  "user1-category1",
//			Color: "#FFFFFF",
//			Icon:  "icon1",
//		})
//
//		testCases := []struct {
//			name               string
//			expenseGroupCreate internal.ExpenseGroupCreate
//			expenseGroupUpdate internal.ExpenseGroupUpdate
//			err                error
//		}{
//			{
//				name: "update name",
//				expenseGroupCreate: internal.ExpenseGroupCreate{
//					Name: "group1",
//					Date: "2025-01-24",
//					Expenses: []internal.ExpenseGroupExpenseCreate{
//						{
//							Name:   "expense 1",
//							Amount: 6969,
//							Note:   "expense note 1",
//						},
//						{
//							Name:   "expense 2",
//							Amount: 420,
//							Note:   "expense note 2",
//						},
//					},
//					Note:       "group note",
//					CategoryID: category.ID,
//				},
//				expenseGroupUpdate: internal.ExpenseGroupUpdate{
//					Name: ptr("updated name"),
//				},
//			},
//			{
//				name: "update date",
//				expenseGroupCreate: internal.ExpenseGroupCreate{
//					Name: "group1",
//					Date: "2025-01-24",
//					Expenses: []internal.ExpenseGroupExpenseCreate{
//						{
//							Name:   "expense 1",
//							Amount: 6969,
//							Note:   "expense note 1",
//						},
//						{
//							Name:   "expense 2",
//							Amount: 420,
//							Note:   "expense note 2",
//						},
//					},
//					Note:       "group note",
//					CategoryID: category.ID,
//				},
//				expenseGroupUpdate: internal.ExpenseGroupUpdate{
//					Date: ptr("2025-03-21"),
//				},
//			},
//			{
//				name: "update note",
//				expenseGroupCreate: internal.ExpenseGroupCreate{
//					Name: "group1",
//					Date: "2025-01-24",
//					Expenses: []internal.ExpenseGroupExpenseCreate{
//						{
//							Name:   "expense 1",
//							Amount: 6969,
//							Note:   "expense note 1",
//						},
//						{
//							Name:   "expense 2",
//							Amount: 420,
//							Note:   "expense note 2",
//						},
//					},
//					Note:       "group note",
//					CategoryID: category.ID,
//				},
//				expenseGroupUpdate: internal.ExpenseGroupUpdate{
//					Note: ptr("updated note"),
//				},
//			},
//		}
//
//		for _, tc := range testCases {
//			t.Run(tc.name, func(t *testing.T) {
//				createdExpenseGroup := mustCreateExpenseGroup(sharedCtx, t, db, tc.expenseGroupCreate)
//				if tc.err != nil {
//					return
//				}
//
//				tc.expenseGroupUpdate.ID = createdExpenseGroup.ID
//
//				updatedExpenseGroup, err := es.UpdateExpenseGroup(sharedCtx, tc.expenseGroupUpdate)
//				assert.Nil(t, err)
//
//				assertUpdatedField(t, tc.expenseGroupUpdate.Name, createdExpenseGroup.Name, updatedExpenseGroup.Name)
//				assertUpdatedField(t, tc.expenseGroupUpdate.Date, createdExpenseGroup.Date.Format(time.DateOnly), updatedExpenseGroup.Date.Format(time.DateOnly))
//				assertUpdatedField(t, tc.expenseGroupUpdate.Note, createdExpenseGroup.Note, updatedExpenseGroup.Note)
//				assert.Equal(t, createdExpenseGroup.CreatedAt, updatedExpenseGroup.CreatedAt)
//
//				// TODO: find a way to test UpdatedAt field without slowing down the tests
//				// it's too fast that the UpdatedAt is still the same with CreatedAt
//
//				for i, createdExpense := range createdExpenseGroup.Expenses {
//					updatedExpense := updatedExpenseGroup.Expenses[i]
//					assertUpdatedField(t, tc.expenseGroupUpdate.Date, createdExpense.Date.Format(time.DateOnly), updatedExpense.Date.Format(time.DateOnly))
//					assertUpdatedField(t, tc.expenseGroupUpdate.CategoryID, createdExpense.Category.ID, updatedExpense.Category.ID)
//				}
//
//				// for _, v := range tc.expenseGroupUpdate.Expenses {
//				// 	var createdExpense *internal.Expense
//				// 	for _, expense := range createdExpenseGroup.Expenses {
//				// 		if v.ID == expense.ID {
//				// 			createdExpense = &expense
//				// 		}
//				// 	}
//				//
//				// 	if createdExpense == nil {
//				// 		t.Fatalf("created expense not found, id: %s", v.ID)
//				// 	}
//				//
//				// 	gotExpense := mustGetExpense(sharedCtx, t, db, v.ID)
//				// 	assertUpdatedField(t, v.Name, createdExpense.Name, gotExpense.Name)
//				// 	assertUpdatedField(t, v.Amount, createdExpense.Amount, gotExpense.Amount)
//				// 	assertUpdatedField(t, v.Note, createdExpense.Note, gotExpense.Note)
//				// }
//
//				gotExpenseGroup := mustGetExpenseGroup(sharedCtx, t, db, updatedExpenseGroup.ID)
//				assert.Equal(t, updatedExpenseGroup, gotExpenseGroup)
//			})
//		}
//	}
func mustCreateExpense(ctx context.Context, t *testing.T, db *sqlite.DB, c internal.ExpenseCreate) internal.Expense {
	t.Helper()

	es := sqlite.NewExpenseService(db)
	e, err := es.CreateExpense(ctx, c)
	require.Nil(t, err)
	return e
}

func mustGetExpense(ctx context.Context, t *testing.T, db *sqlite.DB, id string) internal.Expense {
	t.Helper()

	es := sqlite.NewExpenseService(db)
	e, err := es.GetExpense(ctx, id)
	require.Nil(t, err)
	return e
}

func mustCreateExpenseGroup(ctx context.Context, t *testing.T, db *sqlite.DB, c internal.ExpenseGroupCreate) internal.ExpenseGroup {
	t.Helper()

	es := sqlite.NewExpenseService(db)
	eg, err := es.CreateExpenseGroup(ctx, c)
	require.Nil(t, err)
	return eg
}

func mustGetExpenseGroup(ctx context.Context, t *testing.T, db *sqlite.DB, id string) internal.ExpenseGroup {
	t.Helper()

	es := sqlite.NewExpenseService(db)
	eg, err := es.GetExpenseGroup(ctx, id)
	require.Nil(t, err)
	return eg
}
