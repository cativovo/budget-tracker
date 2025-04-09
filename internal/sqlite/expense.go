package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/cativovo/budget-tracker/internal"
	"github.com/cativovo/budget-tracker/internal/category"
	"github.com/cativovo/budget-tracker/internal/expense"
	"github.com/cativovo/budget-tracker/internal/logger"
	"github.com/cativovo/budget-tracker/internal/user"
	"github.com/huandu/go-sqlbuilder"
)

type ExpenseRepository struct {
	db *DB
	cr CategoryRepository
}

var _ expense.Repository = (*ExpenseRepository)(nil)

func NewExpenseRepository(db *DB, cr CategoryRepository) ExpenseRepository {
	return ExpenseRepository{
		db: db,
		cr: cr,
	}
}

func (er *ExpenseRepository) ExpenseByID(ctx context.Context, id string) (expense.Expense, error) {
	u := user.FromContext(ctx)
	logger := logger.FromContext(ctx)

	sb := sqlbuilder.SQLite.NewSelectBuilder()
	sb.Select(
		"e.id",
		"e.name",
		"e.amount",
		"e.date",
		"e.note",
		"e.created_at",
		"e.updated_at",
		sb.As("c.id", "category_id"),
		sb.As("c.name", "category_name"),
		sb.As("c.color", "category_color"),
		sb.As("c.icon", "category_icon"),
		sb.As("c.created_at", "category_created_at"),
		sb.As("c.updated_at", "category_updated_at"),
	)
	sb.From("expense e")
	sb.Join(
		"category c",
		"c.id = e.category_id",
	)
	sb.Where(
		sb.And(
			sb.EQ("e.id", id),
			sb.EQ("e.user_id", u.ID),
		),
	)

	q, args := sb.Build()

	logger.Infow(
		"Find expense by id",
		"query", q,
		"args", args,
	)

	var dst struct {
		ID                string    `db:"id"`
		Name              string    `db:"name"`
		Amount            int64     `db:"amount"`
		Date              time.Time `db:"date"`
		Note              string    `db:"note"`
		Category          string    `db:"category"`
		CreatedAt         time.Time `db:"created_at"`
		UpdatedAt         time.Time `db:"updated_at"`
		CategoryID        string    `db:"category_id"`
		CategoryName      string    `db:"category_name"`
		CategoryColor     string    `db:"category_color"`
		CategoryIcon      string    `db:"category_icon"`
		CategoryCreatedAt time.Time `db:"category_created_at"`
		CategoryUpdatedAt time.Time `db:"category_updated_at"`
	}
	if err := er.db.readerWriter.GetContext(ctx, &dst, q, args...); err != nil {
		if err == sql.ErrNoRows {
			return expense.Expense{}, internal.NewError(internal.ErrorCodeNotFound, "Expense not found")
		}

		return expense.Expense{}, fmt.Errorf("sqlite.ExpenseRepository.ExpenseByID: %w", err)
	}

	return expense.Expense{
		ID:     dst.ID,
		Name:   dst.Name,
		Amount: dst.Amount,
		Date:   dst.Date,
		Note:   dst.Note,
		Category: category.Category{
			ID:        dst.CategoryID,
			Name:      dst.CategoryName,
			Color:     dst.CategoryColor,
			Icon:      dst.CategoryIcon,
			CreatedAt: dst.CategoryCreatedAt,
			UpdatedAt: dst.CategoryUpdatedAt,
		},
		CreatedAt: dst.CreatedAt,
		UpdatedAt: dst.UpdatedAt,
	}, nil
}

func (er *ExpenseRepository) ListExpenseSummaries(ctx context.Context, lo internal.ListOptions) ([]expense.ExpenseSummary, error) {
	panic("not yet implemented")
}

func (er *ExpenseRepository) CreateExpense(ctx context.Context, e expense.CreateExpenseReq) (expense.Expense, error) {
	category, err := er.cr.CategoryByID(ctx, e.CategoryID)
	if err != nil {
		return expense.Expense{}, fmt.Errorf("sqlite.ExpenseRepository.CreateExpense: %w", err)
	}

	u := user.FromContext(ctx)
	logger := logger.FromContext(ctx)

	ib := sqlbuilder.SQLite.NewInsertBuilder()
	ib.InsertInto("expense")
	ib.Cols(
		"name",
		"amount",
		"date",
		"category_id",
		"note",
		"user_id",
	)
	ib.Values(
		e.Name,
		e.Amount,
		e.Date,
		e.CategoryID,
		e.Note,
		u.ID,
	)
	ib.Returning(
		"id",
		"date",
		"created_at",
		"updated_at",
	)

	q, args := ib.Build()

	logger.Infow(
		"Insert expense",
		"query", q,
		"args", args,
	)

	var dst struct {
		ID        string    `db:"id"`
		Date      time.Time `db:"date"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}
	if err := er.db.readerWriter.GetContext(ctx, &dst, q, args...); err != nil {
		return expense.Expense{}, fmt.Errorf("sqlite.ExpenseRepository.CreateExpense: %w", err)
	}

	return expense.Expense{
		ID:        dst.ID,
		Name:      e.Name,
		Amount:    e.Amount,
		Date:      dst.Date,
		Note:      e.Note,
		Category:  category,
		CreatedAt: dst.CreatedAt,
		UpdatedAt: dst.UpdatedAt,
	}, nil
}

func (er *ExpenseRepository) UpdateExpense(ctx context.Context, e expense.UpdateExpenseReq) (expense.Expense, error) {
	u := user.FromContext(ctx)
	logger := logger.FromContext(ctx)

	ub := sqlbuilder.SQLite.NewUpdateBuilder()
	ub.Update("expense")

	if e.Name != nil {
		ub.SetMore(ub.Assign("name", e.Name))
	}
	if e.Amount != nil {
		ub.SetMore(ub.Assign("amount", e.Amount))
	}
	if e.Date != nil {
		ub.SetMore(ub.Assign("date", e.Date))
	}
	if e.CategoryID != nil {
		ub.SetMore(ub.Assign("category_id", e.CategoryID))
	}
	if e.Note != nil {
		ub.SetMore(ub.Assign("note", e.Note))
	}

	ub.Where(
		ub.And(
			ub.EQ("id", e.ID),
			ub.EQ("user_id", u.ID),
		),
	)

	// https://github.com/huandu/go-sqlbuilder/issues/142
	ub.SQL("RETURNING name, amount, date, category_id, note, created_at, updated_at")

	q, args := ub.Build()

	logger.Infow(
		"Update expense",
		"query", q,
		"args", args,
	)

	var dst struct {
		Name       string    `db:"name"`
		Amount     int64     `db:"amount"`
		Date       time.Time `db:"date"`
		CategoryID string    `db:"category_id"`
		Note       string    `db:"note"`
		CreatedAt  time.Time `db:"created_at"`
		UpdatedAt  time.Time `db:"updated_at"`
	}
	if err := er.db.readerWriter.GetContext(ctx, &dst, q, args...); err != nil {
		if err == sql.ErrNoRows {
			return expense.Expense{}, internal.NewError(internal.ErrorCodeNotFound, "Expense not found")
		}

		return expense.Expense{}, fmt.Errorf("sqlite.ExpenseRepository.UpdateExpense: %w", err)
	}

	c, err := er.cr.CategoryByID(ctx, dst.CategoryID)
	if err != nil {
		return expense.Expense{}, fmt.Errorf("sqlite.ExpenseRepository.UpdateExpense: %w", err)
	}

	return expense.Expense{
		ID:        e.ID,
		Name:      dst.Name,
		Amount:    dst.Amount,
		Date:      dst.Date,
		Note:      dst.Note,
		Category:  c,
		CreatedAt: dst.CreatedAt,
		UpdatedAt: dst.UpdatedAt,
	}, nil
}

func (er *ExpenseRepository) DeleteExpense(ctx context.Context, id string) error {
	u := user.FromContext(ctx)
	logger := logger.FromContext(ctx)

	db := sqlbuilder.SQLite.NewDeleteBuilder()
	db.DeleteFrom("expense")
	db.Where(
		db.And(
			db.EQ("id", id),
			db.EQ("user_id", u.ID),
		),
	)

	q, args := db.Build()

	logger.Infow(
		"Delete expense",
		"query", q,
		"args", args,
	)

	_, err := er.db.readerWriter.ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("sqlite.ExpenseRepository.DeleteExpense: ExecContext: %w", err)
	}

	return nil
}

func (er *ExpenseRepository) ExpenseGroupByID(ctx context.Context, id string) (expense.ExpenseGroup, error) {
	panic("not yet implemented")
}

func (er *ExpenseRepository) CreateExpenseGroup(ctx context.Context, e expense.CreateExpenseGroupReq) (expense.ExpenseGroup, error) {
	panic("not yet implemented")
}

func (er *ExpenseRepository) UpdateExpenseGroup(ctx context.Context, e expense.UpdateExpenseGroupReq) (expense.ExpenseGroup, error) {
	panic("not yet implemented")
}

func (er *ExpenseRepository) DeleteExpenseGroup(ctx context.Context, id string) error {
	panic("not yet implemented")
}
