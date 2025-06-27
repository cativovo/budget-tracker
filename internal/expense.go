package internal

import (
	"context"
	"time"
)

// Amount in cents
type Amount int64

type Expense struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Amount    Amount    `json:"amount"`
	Note      string    `json:"note"`
	Date      time.Time `json:"date"`
	Category  Category  `json:"category"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ExpenseGroup struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Date      time.Time `json:"date"`
	Expenses  []Expense `json:"expenses"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GroupedExpenseSummary struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Amount Amount `json:"amount"`
}

type ExpenseSummary struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Amount  Amount    `json:"amount"`
	Date    time.Time `json:"date"`
	IsGroup bool      `json:"is_group"`
}

type ExpenseService interface {
	GetExpense(ctx context.Context, id string) (Expense, error)
	CreateExpense(ctx context.Context, c ExpenseCreate) (Expense, error)
	UpdateExpense(ctx context.Context, c ExpenseUpdate) (Expense, error)
	DeleteExpense(ctx context.Context, id string) error

	GetExpenseGroup(ctx context.Context, id string) (ExpenseGroup, error)
	CreateExpenseGroup(ctx context.Context, c ExpenseGroupCreate) (ExpenseGroup, error)

	// NOTE: do after doing ExpenseGroup
	// ListExpenseSummaries(ctx context.Context, f ExpenseSummaryFilter) ([]ExpenseSummary, int, error)
}

type ExpenseOrderBy string

const (
	ExpenseOrderByName      ExpenseOrderBy = "name"
	ExpenseOrderByAmount    ExpenseOrderBy = "amount"
	ExpenseOrderByDate      ExpenseOrderBy = "date"
	ExpenseOrderByCategory  ExpenseOrderBy = "category"
	ExpenseOrderByCreatedAt ExpenseOrderBy = "created_at"
	ExpenseOrderByUpdatedAt ExpenseOrderBy = "updated_at"
)

type ExpenseSummaryFilter struct {
	Name      *string
	StartDate *time.Time
	EndDate   *time.Time
	OrderBy   *ExpenseOrderBy
	OrderDesc bool
}

type ExpenseCreate struct {
	Name       string `json:"name" validate:"required"`
	Amount     Amount `json:"amount" validate:"gt=0"`
	Date       string `json:"date" validate:"datetime=2006-01-02"`
	Note       string `json:"note"`
	CategoryID string `json:"category_id" validate:"required"`
}

func (c ExpenseCreate) Validate() error {
	if err := ValidateStruct(c); err != nil {
		return NewError(ErrorCodeInvalid, err.Error())
	}
	return nil
}

type ExpenseUpdate struct {
	ID         string  `json:"id" validate:"required"`
	Name       *string `json:"name" validate:"omitnil,min=1"`
	Amount     *Amount `json:"amount" validate:"omitnil,gt=0"`
	Note       *string `json:"note"`
	Date       *string `json:"date" validate:"omitnil,datetime=2006-01-02"`
	CategoryID *string `json:"category_id" validate:"omitnil,min=1"`
}

func (u ExpenseUpdate) Validate() error {
	c := u.Name == nil &&
		u.Amount == nil &&
		u.Date == nil &&
		u.CategoryID == nil &&
		u.Note == nil
	if c {
		return NewError(ErrorCodeInvalid, "No update fields provided")
	}

	if err := ValidateStruct(u); err != nil {
		return NewError(ErrorCodeInvalid, err.Error())
	}

	return nil
}

type ExpenseGroupExpense struct {
	Name   string `json:"name" validate:"required"`
	Amount Amount `json:"amount" validate:"gt=0"`
	Note   string `json:"note"`
}

type ExpenseGroupCreate struct {
	Name       string                `json:"name" validate:"required"`
	Date       string                `json:"date" validate:"datetime=2006-01-02"`
	Expenses   []ExpenseGroupExpense `json:"expenses" validate:"required"`
	Note       string                `json:"note" validate:"required"`
	CategoryID string                `json:"category_id" validate:"required"`
}

func (c ExpenseGroupCreate) Validate() error {
	if err := ValidateStruct(c); err != nil {
		return NewError(ErrorCodeInvalid, err.Error())
	}
	return nil
}
