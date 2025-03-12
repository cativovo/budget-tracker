package expense

import (
	"context"

	"github.com/cativovo/budget-tracker/internal"
	"github.com/cativovo/budget-tracker/internal/validator"
)

type ExpenseService struct {
	er ExpenseRepository
	v  *validator.Validator
}

func NewExpenseService(er ExpenseRepository, v *validator.Validator) *ExpenseService {
	return &ExpenseService{
		er: er,
		v:  v,
	}
}

func (es *ExpenseService) ListExpenseSummaries(ctx context.Context, lo internal.ListOptions) ([]ExpenseSummary, error) {
	return es.er.ListExpenseSummaries(ctx, lo)
}

type CreateExpenseReq struct {
	Name   string `json:"name" validate:"required"`
	Amount int64  `json:"amount" validate:"gt=0"`
	Date   string `json:"date" validate:"required,datetime=2006-01-02"`
	Note   string `json:"note"`
}

func (es *ExpenseService) CreateExpense(ctx context.Context, c CreateExpenseReq) (Expense, error) {
	if err := es.v.Struct(c); err != nil {
		return Expense{}, err
	}
	return es.er.CreateExpense(ctx, c)
}

type UpdateExpenseReq struct {
	Name   string `json:"name"`
	Amount int64  `json:"amount" validate:"gt=0"`
	Date   string `json:"date" validate:"datetime=2006-01-02"`
	Note   string `json:"note"`
}

func (es *ExpenseService) UpdateExpense(ctx context.Context, u UpdateExpenseReq) (Expense, error) {
	if err := es.v.Struct(u); err != nil {
		return Expense{}, err
	}
	return es.er.UpdateExpense(ctx, u)
}
