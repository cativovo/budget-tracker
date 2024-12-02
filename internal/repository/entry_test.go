package repository_test

import (
	"context"
	"math/rand/v2"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/cativovo/budget-tracker/internal/models"
	"github.com/cativovo/budget-tracker/internal/repository"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestCreateEntry(t *testing.T) {
	h := newRepositoryHelper(t, "TestCreateEntry.db")
	r := h.repository()
	defer h.clean()
	logger := zap.NewNop().Sugar()

	account, err := r.CreateAccount(context.Background(), logger, repository.CreateAccountParams{
		Name: "Zhou Guanyu",
	})
	if err != nil {
		t.Fatal(err)
	}

	category1, err := r.CreateCategory(context.Background(), logger, repository.CreateCategoryParams{
		Name:      "Groceries",
		Icon:      "mdi:cart",
		ColorHex:  "#FF6347",
		AccountID: account.ID,
	})
	if err != nil {
		t.Fatal(err)
	}

	category2, err := r.CreateCategory(context.Background(), logger, repository.CreateCategoryParams{
		Name:      "Utilities",
		Icon:      "mdi:lightbulb",
		ColorHex:  "#32CD32",
		AccountID: account.ID,
	})
	if err != nil {
		t.Fatal(err)
	}

	params := []repository.CreateEntryParams{
		{
			Date:        "2024-11-30",
			CategoryID:  &category1.ID,
			Description: new(string),
			Name:        "Grocery Store Purchase",
			AccountID:   account.ID,
			Amount:      200,
			EntryType:   0,
		},
		{
			Date:        "2024-11-28",
			CategoryID:  &category1.ID,
			Description: new(string),
			Name:        "Utility Bill Payment",
			AccountID:   account.ID,
			Amount:      150,
			EntryType:   1,
		},
		{
			Date:        "2024-11-27",
			CategoryID:  nil,
			Description: new(string),
			Name:        "Electricity Bill Payment",
			AccountID:   account.ID,
			Amount:      120,
			EntryType:   1,
		},
		{
			Date:        "2024-11-26",
			CategoryID:  &category2.ID,
			Description: new(string),
			Name:        "Internet Subscription",
			AccountID:   account.ID,
			Amount:      90,
			EntryType:   0,
		},
		{
			Date:        "2024-11-25",
			CategoryID:  &category2.ID,
			Description: new(string),
			Name:        "Supermarket Shopping",
			AccountID:   account.ID,
			Amount:      100,
			EntryType:   0,
		},
	}

	*params[0].Description = "Purchased groceries for the week, including vegetables and snacks."
	*params[1].Description = "Paid the monthly utility bill covering electricity and water services."
	*params[2].Description = "Paid for electricity consumption for the month."
	*params[3].Description = "Monthly payment for internet subscription, including taxes."
	*params[4].Description = "Bought household essentials, including food and cleaning products."

	expected := []models.Entry{
		{
			Date:        "2024-11-30",
			Category:    &category1,
			Description: new(string),
			Name:        "Grocery Store Purchase",
			EntryType:   0,
			Amount:      200,
		},
		{
			Date:        "2024-11-28",
			Category:    &category1,
			Description: new(string),
			Name:        "Utility Bill Payment",
			EntryType:   1,
			Amount:      150,
		},
		{
			Date:        "2024-11-27",
			Category:    nil,
			Description: new(string),
			Name:        "Electricity Bill Payment",
			EntryType:   1,
			Amount:      120,
		},
		{
			Date:        "2024-11-26",
			Category:    &category2,
			Description: new(string),
			Name:        "Internet Subscription",
			EntryType:   0,
			Amount:      90,
		},
		{
			Date:        "2024-11-25",
			Category:    &category2,
			Description: new(string),
			Name:        "Supermarket Shopping",
			EntryType:   0,
			Amount:      100,
		},
	}

	*expected[0].Description = "Purchased groceries for the week, including vegetables and snacks."
	*expected[1].Description = "Paid the monthly utility bill covering electricity and water services."
	*expected[2].Description = "Paid for electricity consumption for the month."
	*expected[3].Description = "Monthly payment for internet subscription, including taxes."
	*expected[4].Description = "Bought household essentials, including food and cleaning products."

	for _, param := range params {
		_, err := r.CreateEntry(context.Background(), logger, param)
		if err != nil {
			t.Fatal(err)
		}
	}

	entriesWithCount, err := r.ListEntriesByDate(context.Background(), logger, repository.ListEntriesByDateParams{
		StartDate: "2024-11-25",
		EndDate:   "2024-11-30",
		AccountID: account.ID,
		EntryType: []models.EntryType{models.EntryTypeExpense, models.EntryTypeIncome},
		Limit:     10,
		Offset:    0,
		OrderBy:   repository.Desc,
	})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(expected), entriesWithCount.TotalCount)
	assertEntries(t, expected, entriesWithCount.Entries)
}

func TestListEntriesByDate(t *testing.T) {
	h := newRepositoryHelper(t, "TestListEntriesByDate.db")
	r := h.repository()
	defer h.clean()
	logger := zap.NewNop().Sugar()

	account, err := r.CreateAccount(context.Background(), logger, repository.CreateAccountParams{
		Name: "Valterri Bottas",
	})
	if err != nil {
		t.Fatal(err)
	}

	category1, err := r.CreateCategory(context.Background(), logger, repository.CreateCategoryParams{
		Name:      "Groceries",
		Icon:      "mdi:cart",
		ColorHex:  "#FF6347",
		AccountID: account.ID,
	})
	if err != nil {
		t.Fatal(err)
	}

	category2, err := r.CreateCategory(context.Background(), logger, repository.CreateCategoryParams{
		Name:      "Utilities",
		Icon:      "mdi:lightbulb",
		ColorHex:  "#32CD32",
		AccountID: account.ID,
	})
	if err != nil {
		t.Fatal(err)
	}

	params := []repository.CreateEntryParams{
		{
			Date:        "2024-11-30",
			CategoryID:  &category1.ID,
			Description: new(string),
			Name:        "Grocery Store Purchase",
			AccountID:   account.ID,
			Amount:      200,
			EntryType:   0,
		},
		{
			Date:        "2024-11-28",
			CategoryID:  &category1.ID,
			Description: new(string),
			Name:        "Utility Bill Payment",
			AccountID:   account.ID,
			Amount:      150,
			EntryType:   1,
		},
		{
			Date:        "2024-11-27",
			CategoryID:  nil,
			Description: new(string),
			Name:        "Electricity Bill Payment",
			AccountID:   account.ID,
			Amount:      120,
			EntryType:   1,
		},
		{
			Date:        "2024-11-26",
			CategoryID:  &category2.ID,
			Description: new(string),
			Name:        "Internet Subscription",
			AccountID:   account.ID,
			Amount:      90,
			EntryType:   0,
		},
		{
			Date:        "2024-11-25",
			CategoryID:  &category2.ID,
			Description: new(string),
			Name:        "Supermarket Shopping",
			AccountID:   account.ID,
			Amount:      100,
			EntryType:   0,
		},
	}

	*params[0].Description = "Purchased groceries for the week, including vegetables and snacks."
	*params[1].Description = "Paid the monthly utility bill covering electricity and water services."
	*params[2].Description = "Paid for electricity consumption for the month."
	*params[3].Description = "Monthly payment for internet subscription, including taxes."
	*params[4].Description = "Bought household essentials, including food and cleaning products."

	rand.Shuffle(len(params), func(i, j int) {
		params[i], params[j] = params[j], params[i]
	})

	expected := []models.Entry{
		{
			Date:        "2024-11-30",
			Category:    &category1,
			Description: new(string),
			Name:        "Grocery Store Purchase",
			EntryType:   0,
			Amount:      200,
		},
		{
			Date:        "2024-11-28",
			Category:    &category1,
			Description: new(string),
			Name:        "Utility Bill Payment",
			EntryType:   1,
			Amount:      150,
		},
		{
			Date:        "2024-11-27",
			Category:    nil,
			Description: new(string),
			Name:        "Electricity Bill Payment",
			EntryType:   1,
			Amount:      120,
		},
		{
			Date:        "2024-11-26",
			Category:    &category2,
			Description: new(string),
			Name:        "Internet Subscription",
			EntryType:   0,
			Amount:      90,
		},
		{
			Date:        "2024-11-25",
			Category:    &category2,
			Description: new(string),
			Name:        "Supermarket Shopping",
			EntryType:   0,
			Amount:      100,
		},
	}

	*expected[0].Description = "Purchased groceries for the week, including vegetables and snacks."
	*expected[1].Description = "Paid the monthly utility bill covering electricity and water services."
	*expected[2].Description = "Paid for electricity consumption for the month."
	*expected[3].Description = "Monthly payment for internet subscription, including taxes."
	*expected[4].Description = "Bought household essentials, including food and cleaning products."

	for _, param := range params {
		_, err := r.CreateEntry(context.Background(), logger, param)
		if err != nil {
			t.Fatal(err)
		}
	}

	t.Run("DESC", func(t *testing.T) {
		entriesWithCount, err := r.ListEntriesByDate(context.Background(), logger, repository.ListEntriesByDateParams{
			StartDate: "2024-11-25",
			EndDate:   "2024-11-30",
			AccountID: account.ID,
			EntryType: []models.EntryType{models.EntryTypeExpense, models.EntryTypeIncome},
			Limit:     10,
			Offset:    0,
			OrderBy:   repository.Desc,
		})
		if err != nil {
			t.Fatal(err)
		}

		descSortedExpected := make([]models.Entry, len(expected))
		copy(descSortedExpected, expected)

		slices.SortFunc(descSortedExpected, func(a, b models.Entry) int {
			return strings.Compare(b.Date, a.Date)
		})

		assert.Equal(t, len(expected), entriesWithCount.TotalCount)
		assertEntries(t, descSortedExpected, entriesWithCount.Entries)
	})

	t.Run("ASC", func(t *testing.T) {
		entriesWithCount, err := r.ListEntriesByDate(context.Background(), logger, repository.ListEntriesByDateParams{
			StartDate: "2024-11-25",
			EndDate:   "2024-11-30",
			AccountID: account.ID,
			EntryType: []models.EntryType{models.EntryTypeExpense, models.EntryTypeIncome},
			Limit:     10,
			Offset:    0,
			OrderBy:   repository.Asc,
		})
		if err != nil {
			t.Fatal(err)
		}

		ascSortedExpected := make([]models.Entry, len(expected))
		copy(ascSortedExpected, expected)

		slices.SortFunc(ascSortedExpected, func(a, b models.Entry) int {
			return strings.Compare(a.Date, b.Date)
		})

		assert.Equal(t, len(expected), entriesWithCount.TotalCount)
		assertEntries(t, ascSortedExpected, entriesWithCount.Entries)
	})

	t.Run("BETWEEN", func(t *testing.T) {
		entriesWithCount, err := r.ListEntriesByDate(context.Background(), logger, repository.ListEntriesByDateParams{
			StartDate: "2024-11-27",
			EndDate:   "2024-11-30",
			AccountID: account.ID,
			EntryType: []models.EntryType{models.EntryTypeExpense, models.EntryTypeIncome},
			Limit:     10,
			Offset:    0,
			OrderBy:   repository.Desc,
		})
		if err != nil {
			t.Fatal(err)
		}

		expectedCount := 3
		assert.Equal(t, expectedCount, entriesWithCount.TotalCount)
		assertEntries(t, expected[:expectedCount], entriesWithCount.Entries[:expectedCount])
	})

	t.Run("LIMIT", func(t *testing.T) {
		entriesWithCount, err := r.ListEntriesByDate(context.Background(), logger, repository.ListEntriesByDateParams{
			StartDate: "2024-11-25",
			EndDate:   "2024-11-30",
			AccountID: account.ID,
			EntryType: []models.EntryType{models.EntryTypeExpense, models.EntryTypeIncome},
			Limit:     2,
			Offset:    0,
			OrderBy:   repository.Desc,
		})
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, len(expected), entriesWithCount.TotalCount)

		expectedLen := 2
		assert.Equal(t, len(entriesWithCount.Entries), expectedLen)
		assertEntries(t, expected[:expectedLen], entriesWithCount.Entries[:expectedLen])
	})

	t.Run("OFFSET", func(t *testing.T) {
		offset := 2
		entriesWithCount, err := r.ListEntriesByDate(context.Background(), logger, repository.ListEntriesByDateParams{
			StartDate: "2024-11-25",
			EndDate:   "2024-11-30",
			AccountID: account.ID,
			EntryType: []models.EntryType{models.EntryTypeExpense, models.EntryTypeIncome},
			Limit:     2,
			Offset:    offset,
			OrderBy:   repository.Desc,
		})
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, len(expected), entriesWithCount.TotalCount)

		expectedLen := 2
		assert.Equal(t, len(entriesWithCount.Entries), expectedLen)
		assertEntries(t, expected[offset:offset+expectedLen], entriesWithCount.Entries)
	})

	t.Run("EntryTypeExpense", func(t *testing.T) {
		entriesWithCount, err := r.ListEntriesByDate(context.Background(), logger, repository.ListEntriesByDateParams{
			StartDate: "2024-11-25",
			EndDate:   "2024-11-30",
			AccountID: account.ID,
			EntryType: []models.EntryType{models.EntryTypeExpense},
			Limit:     10,
			Offset:    0,
			OrderBy:   repository.Desc,
		})
		if err != nil {
			t.Fatal(err)
		}

		var filteredExpected []models.Entry
		for _, v := range expected {
			if v.EntryType == models.EntryTypeExpense {
				filteredExpected = append(filteredExpected, v)
			}
		}

		assert.Equal(t, len(filteredExpected), entriesWithCount.TotalCount)
		assertEntries(t, filteredExpected, entriesWithCount.Entries)
	})

	t.Run("EntryTypeIncome", func(t *testing.T) {
		entriesWithCount, err := r.ListEntriesByDate(context.Background(), logger, repository.ListEntriesByDateParams{
			StartDate: "2024-11-25",
			EndDate:   "2024-11-30",
			AccountID: account.ID,
			EntryType: []models.EntryType{models.EntryTypeIncome},
			Limit:     10,
			Offset:    0,
			OrderBy:   repository.Desc,
		})
		if err != nil {
			t.Fatal(err)
		}

		var filteredExpected []models.Entry
		for _, v := range expected {
			if v.EntryType == models.EntryTypeIncome {
				filteredExpected = append(filteredExpected, v)
			}
		}

		assert.Equal(t, len(filteredExpected), entriesWithCount.TotalCount)
		assertEntries(t, filteredExpected, entriesWithCount.Entries)
	})

	t.Run("Invalid account ID", func(t *testing.T) {
		entriesWithCount, err := r.ListEntriesByDate(context.Background(), logger, repository.ListEntriesByDateParams{
			StartDate: "2024-11-25",
			EndDate:   "2024-11-30",
			AccountID: "6969",
			EntryType: []models.EntryType{models.EntryTypeIncome},
			Limit:     10,
			Offset:    0,
			OrderBy:   repository.Desc,
		})
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, 0, entriesWithCount.TotalCount)
		assert.Empty(t, entriesWithCount.Entries)
	})
}

func assertEntries(t *testing.T, want, got []models.Entry) {
	t.Helper()

	for i, entry := range got {
		assert.Equal(t, want[i].Date, entry.Date)
		assert.Equal(t, want[i].Name, entry.Name)
		assert.Equal(t, want[i].EntryType, entry.EntryType)
		assert.Equal(t, want[i].Amount, entry.Amount)

		assertEqualPointer(t, want[i].Category, entry.Category)
		assertEqualPointer(t, want[i].Description, entry.Description)

		assert.NotEmpty(t, entry.ID)

		assertTimestampWithin(t, time.Now().UTC(), entry.CreatedAt, time.Minute)
		assertTimestampWithin(t, time.Now().UTC(), entry.UpdatedAt, time.Minute)
	}
}
