package expense

import (
	"context"

	"github.com/cativovo/budget-tracker/internal"
	"github.com/cativovo/budget-tracker/internal/validator"
)

type Service interface {
	ListExpenseSummaries(ctx context.Context, lo internal.ListOptions) ([]ExpenseSummary, error)
	CreateExpense(ctx context.Context, c CreateExpenseReq) (Expense, error)
	UpdateExpense(ctx context.Context, u UpdateExpenseReq) (Expense, error)
}

type service struct {
	r Repository
	v *validator.Validator
}

func NewService(r Repository, v *validator.Validator) Service {
	return &service{
		r: r,
		v: v,
	}
}

func (s *service) ListExpenseSummaries(ctx context.Context, lo internal.ListOptions) ([]ExpenseSummary, error) {
	return s.r.ListExpenseSummaries(ctx, lo)
}

type CreateExpenseReq struct {
	Name       string `json:"name" validate:"required"`
	Amount     int64  `json:"amount" validate:"gt=0"`
	Date       string `json:"date" validate:"required,datetime=2006-01-02"`
	CategoryID string `json:"category_id" validate:"required"`
	Note       string `json:"note"`
}

func (s *service) CreateExpense(ctx context.Context, c CreateExpenseReq) (Expense, error) {
	if err := s.v.Struct(c); err != nil {
		return Expense{}, internal.NewError(internal.ErrorCodeInvalid, err.Error())
	}
	return s.r.CreateExpense(ctx, c)
}

type UpdateExpenseReq struct {
	ID         string  `json:"id" validate:"required"`
	Name       *string `json:"name"`
	Amount     *int64  `json:"amount" validate:"gt=0"`
	Date       *string `json:"date" validate:"datetime=2006-01-02"`
	CategoryID *string `json:"category_id"`
	Note       *string `json:"note"`
}

func (s *service) UpdateExpense(ctx context.Context, u UpdateExpenseReq) (Expense, error) {
	if err := s.v.Struct(u); err != nil {
		return Expense{}, internal.NewError(internal.ErrorCodeInvalid, err.Error())
	}
	return s.r.UpdateExpense(ctx, u)
}

type CreateExpenseGroupReq struct {
	Name     string `json:"name" validate:"required"`
	Expenses []struct {
		Name   string `json:"name" validate:"required"`
		Amount int64  `json:"amount" validate:"required"`
	} `json:"expenses" validate:"required"`
	Date string `json:"date" validate:"required,datetime=2006-01-02"`
	Note string `json:"note"`
}

type UpdateExpenseGroupReq struct {
	ID       string `json:"id"`
	Name     string `json:"name" validate:"required"`
	Expenses []struct {
		ID     string `json:"id"`
		Name   string `json:"name" validate:"required"`
		Amount int64  `json:"amount" validate:"required"`
	} `json:"expenses" validate:"required"`
	Date string `json:"date" validate:"required,datetime=2006-01-02"`
	Note string `json:"note"`
}
