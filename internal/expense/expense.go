package expense

import (
	"time"

	"github.com/cativovo/budget-tracker/internal/category"
)

type Expense struct {
	ID        string
	Name      string
	Amount    int64
	Date      time.Time
	Note      string
	Category  category.Category
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ExpenseGroup struct {
	ID        string
	Name      string
	Expenses  []Expense
	Note      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ExpenseSummary struct {
	ID      string
	Name    string
	Amount  int64
	Date    time.Time
	IsGroup bool
}
