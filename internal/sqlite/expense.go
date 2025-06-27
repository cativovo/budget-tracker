package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/cativovo/budget-tracker/internal"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
)

type ExpenseService struct {
	db *DB
}

var _ internal.ExpenseService = (*ExpenseService)(nil)

func NewExpenseService(db *DB) *ExpenseService {
	return &ExpenseService{
		db: db,
	}
}

func (es *ExpenseService) GetExpense(ctx context.Context, id string) (internal.Expense, error) {
	user := internal.UserFromContext(ctx)
	logger := internal.LoggerFromContext(ctx)

	sb := sqlbuilder.SQLite.NewSelectBuilder()
	sb.Select(
		"e.id",
		"e.name",
		"e.amount",
		"e.note",
		"e.date",
		"e.created_at",
		"e.updated_at",
		sb.As("c.id", "category_id"),
		sb.As("c.name", "category_name"),
		"c.color",
		"c.icon",
		sb.As("c.created_at", "category_created_at"),
		sb.As("c.updated_at", "category_updated_at"),
	)
	sb.From("expense e")
	sb.Join("category c", "c.id == e.category_id")
	sb.Where(
		sb.Equal("e.id", id),
		sb.Equal("e.user_id", user.ID),
	)

	q, args := sb.Build()

	logger.Infow("Get expense", "query", q, "args", args)

	var dst struct {
		ID                string          `db:"id"`
		Name              string          `db:"name"`
		Amount            internal.Amount `db:"amount"`
		Note              string          `db:"note"`
		Date              time.Time       `db:"date"`
		CreatedAt         time.Time       `db:"created_at"`
		UpdatedAt         time.Time       `db:"updated_at"`
		CategoryID        string          `db:"category_id"`
		CategoryName      string          `db:"category_name"`
		CategoryColor     string          `db:"color"`
		CategoryIcon      string          `db:"icon"`
		CategoryCreatedAt time.Time       `db:"category_created_at"`
		CategoryUpdatedAt time.Time       `db:"category_updated_at"`
	}
	if err := es.db.Reader.GetContext(ctx, &dst, q, args...); err != nil {
		if err == sql.ErrNoRows {
			iErr := internal.NewError(internal.ErrorCodeNotFound, "expense not found")
			return internal.Expense{}, fmt.Errorf("sqlite: %w", iErr)
		}

		return internal.Expense{}, fmt.Errorf("sqlite: get expense: %w", err)
	}

	return internal.Expense{
		ID:     dst.ID,
		Name:   dst.Name,
		Amount: dst.Amount,
		Note:   dst.Note,
		Date:   dst.Date,
		Category: internal.Category{
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

func (es *ExpenseService) CreateExpense(ctx context.Context, c internal.ExpenseCreate) (internal.Expense, error) {
	if err := c.Validate(); err != nil {
		return internal.Expense{}, err
	}

	tx, err := es.db.Writer.BeginTxx(ctx, nil)
	if err != nil {
		return internal.Expense{}, fmt.Errorf("sqlite: create expense begin: %w", err)
	}
	defer tx.Rollback()

	expense, err := createExpense(ctx, tx, c, nil)
	if err != nil {
		return internal.Expense{}, fmt.Errorf("sqlite: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return internal.Expense{}, fmt.Errorf("sqlite: create expense commit: %w", err)
	}

	return expense, nil
}

func (es *ExpenseService) UpdateExpense(ctx context.Context, u internal.ExpenseUpdate) (internal.Expense, error) {
	if err := u.Validate(); err != nil {
		return internal.Expense{}, err
	}

	user := internal.UserFromContext(ctx)
	logger := internal.LoggerFromContext(ctx)

	ub := sqlbuilder.SQLite.NewUpdateBuilder()
	ub.Update("expense")

	setMoreIfNotNil(ub, "name", u.Name)
	setMoreIfNotNil(ub, "amount", u.Amount)
	setMoreIfNotNil(ub, "note", u.Note)
	setMoreIfNotNil(ub, "date", u.Date)
	setMoreIfNotNil(ub, "category_id", u.CategoryID)
	setUpdatedAt(ub)

	ub.Where(
		ub.Equal("id", u.ID),
		ub.Equal("user_id", user.ID),
	)

	ub.SQL(`
		RETURNING
			id,
			name,
			amount,
			note,
			date,
			created_at,
			updated_at,
			category_id
	`)

	q, args := ub.Build()

	logger.Infow("Update expense", "query", q, "args", args)

	tx, err := es.db.Writer.BeginTxx(ctx, nil)
	if err != nil {
		return internal.Expense{}, fmt.Errorf("sqlite: update expense begin: %w", err)
	}

	defer tx.Rollback()

	var dst expenseWriteDst
	if err := tx.GetContext(ctx, &dst, q, args...); err != nil {
		var sqErr sqlite3.Error
		if errors.As(err, &sqErr) && sqErr.ExtendedCode == sqlite3.ErrConstraintForeignKey {
			iErr := internal.NewError(internal.ErrorCodeInvalid, "category not found")
			return internal.Expense{}, fmt.Errorf("sqlite: %w", iErr)
		}
		return internal.Expense{}, fmt.Errorf("sqlite: update expense: %w", err)
	}

	category, err := getCategory(ctx, tx, dst.CategoryID)
	if err != nil {
		return internal.Expense{}, fmt.Errorf("sqlite: update expense get category: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return internal.Expense{}, fmt.Errorf("sqlite: update expense commit: %w", err)
	}

	return internal.Expense{
		ID:        dst.ID,
		Name:      dst.Name,
		Amount:    dst.Amount,
		Note:      dst.Note,
		Date:      dst.Date,
		Category:  category,
		CreatedAt: dst.CreatedAt,
		UpdatedAt: dst.UpdatedAt,
	}, nil
}

func (es *ExpenseService) DeleteExpense(ctx context.Context, id string) error {
	user := internal.UserFromContext(ctx)
	logger := internal.LoggerFromContext(ctx)

	db := sqlbuilder.SQLite.NewDeleteBuilder()
	db.DeleteFrom("expense")
	db.Where(
		db.Equal("id", id),
		db.Equal("user_id", user.ID),
	)

	q, args := db.Build()

	logger.Infow("Delete expense", "query", q, "args", args)

	_, err := es.db.Writer.ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("sqlite: delete expense: %w", err)
	}

	return nil
}

func (es *ExpenseService) GetExpenseGroup(ctx context.Context, id string) (internal.ExpenseGroup, error) {
	user := internal.UserFromContext(ctx)
	logger := internal.LoggerFromContext(ctx)

	esgb := sqlbuilder.SQLite.NewSelectBuilder()
	esgb.Select(
		"id",
		"name",
		"note",
		"date",
		"created_at",
		"updated_at",
	)
	esgb.From("expense_group")
	esgb.Where(
		esgb.Equal("id", id),
		esgb.Equal("user_id", user.ID),
	)

	q, args := esgb.Build()

	logger.Infow("Get expense group", "query", q, "args", args)

	tx, err := es.db.Reader.BeginTxx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return internal.ExpenseGroup{}, fmt.Errorf("sqlite: get expense group begin: %w", err)
	}

	defer tx.Rollback()

	var expenseGroupDst struct {
		ID        string    `db:"id"`
		Name      string    `db:"name"`
		Note      string    `db:"note"`
		Date      time.Time `db:"date"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}
	if err := tx.GetContext(ctx, &expenseGroupDst, q, args...); err != nil {
		if err == sql.ErrNoRows {
			iErr := internal.NewError(internal.ErrorCodeNotFound, "expense group not found")
			return internal.ExpenseGroup{}, fmt.Errorf("sqlite: %w", iErr)
		}
		return internal.ExpenseGroup{}, fmt.Errorf("sqlite: get expense group: %w", err)
	}

	esb := sqlbuilder.SQLite.NewSelectBuilder()
	esb.Select(
		"e.id",
		"e.name",
		"e.amount",
		"e.note",
		"e.date",
		"e.created_at",
		"e.updated_at",
		esb.As("c.id", "category_id"),
		esb.As("c.name", "category_name"),
		"c.color",
		"c.icon",
		esb.As("c.created_at", "category_created_at"),
		esb.As("c.updated_at", "category_updated_at"),
	)
	esb.From("expense e")
	esb.Join("category c", "c.id == e.category_id")
	esb.Where(
		esb.Equal("e.expense_group_id", id),
		esb.Equal("e.user_id", user.ID),
	)

	q, args = esb.Build()

	logger.Infow("Get expense group expenses", "query", q, "args", args)

	rows, err := tx.QueryxContext(ctx, q, args...)
	if err != nil && err != sql.ErrNoRows {
		return internal.ExpenseGroup{}, fmt.Errorf("sqlite: get expense group expenses: %w", err)
	}

	var expenses []internal.Expense
	for rows.Next() {
		var dst expenseReadDst
		if err := rows.StructScan(&dst); err != nil {
			return internal.ExpenseGroup{}, fmt.Errorf("sqlite: get expense group expenses scan: %w", err)
		}
		e := internal.Expense{
			ID:     dst.ID,
			Name:   dst.Name,
			Amount: dst.Amount,
			Note:   dst.Note,
			Date:   dst.Date,
			Category: internal.Category{
				ID:        dst.CategoryID,
				Name:      dst.CategoryName,
				Color:     dst.CategoryColor,
				Icon:      dst.CategoryIcon,
				CreatedAt: dst.CategoryCreatedAt,
				UpdatedAt: dst.CategoryUpdatedAt,
			},
			CreatedAt: dst.CreatedAt,
			UpdatedAt: dst.UpdatedAt,
		}
		expenses = append(expenses, e)
	}

	if err := tx.Commit(); err != nil {
		return internal.ExpenseGroup{}, fmt.Errorf("sqlite: get expense group commit: %w", err)
	}

	return internal.ExpenseGroup{
		ID:        expenseGroupDst.ID,
		Name:      expenseGroupDst.Name,
		Date:      expenseGroupDst.Date,
		Expenses:  expenses,
		Note:      expenseGroupDst.Note,
		CreatedAt: expenseGroupDst.CreatedAt,
		UpdatedAt: expenseGroupDst.UpdatedAt,
	}, nil
}

func (es *ExpenseService) CreateExpenseGroup(ctx context.Context, c internal.ExpenseGroupCreate) (internal.ExpenseGroup, error) {
	if err := c.Validate(); err != nil {
		return internal.ExpenseGroup{}, fmt.Errorf("sqlite: %w", err)
	}

	user := internal.UserFromContext(ctx)
	logger := internal.LoggerFromContext(ctx)

	ib := sqlbuilder.SQLite.NewInsertBuilder()
	ib.InsertInto("expense_group")
	ib.Cols("name", "note", "date", "user_id")
	ib.Values(c.Name, c.Note, c.Date, user.ID)
	ib.Returning(
		"id",
		"name",
		"note",
		"date",
		"created_at",
		"updated_at",
	)

	q, args := ib.Build()

	logger.Infow("Create expense group", "query", q, "args", args)

	tx, err := es.db.Writer.BeginTxx(ctx, nil)
	if err != nil {
		return internal.ExpenseGroup{}, fmt.Errorf("sqlite: create expense group begin: %w", err)
	}

	defer tx.Rollback()

	var dst struct {
		ID        string    `db:"id" json:"id"`
		Name      string    `db:"name" json:"name"`
		Note      string    `db:"note" json:"note"`
		Date      time.Time `db:"date" json:"date"`
		CreatedAt time.Time `db:"created_at" json:"created_at"`
		UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	}
	if err := tx.GetContext(ctx, &dst, q, args...); err != nil {
		return internal.ExpenseGroup{}, fmt.Errorf("sqlite: create expense group: %w", err)
	}

	var expenses []internal.Expense
	for _, v := range c.Expenses {
		e, err := createExpense(ctx, tx, internal.ExpenseCreate{
			Name:       v.Name,
			Amount:     v.Amount,
			Date:       c.Date,
			Note:       v.Note,
			CategoryID: c.CategoryID,
		}, &dst.ID)
		if err != nil {
			return internal.ExpenseGroup{}, fmt.Errorf("sqlite: create expense group: %w", err)
		}

		expenses = append(expenses, e)
	}

	if err := tx.Commit(); err != nil {
		return internal.ExpenseGroup{}, fmt.Errorf("sqlite: create expense group commit: %w", err)
	}

	return internal.ExpenseGroup{
		ID:        dst.ID,
		Name:      dst.Name,
		Date:      dst.Date,
		Expenses:  expenses,
		Note:      dst.Note,
		CreatedAt: dst.CreatedAt,
		UpdatedAt: dst.UpdatedAt,
	}, nil
}

func createExpense(ctx context.Context, tx *sqlx.Tx, c internal.ExpenseCreate, groupID *string) (internal.Expense, error) {
	user := internal.UserFromContext(ctx)
	logger := internal.LoggerFromContext(ctx)

	ib := sqlbuilder.SQLite.NewInsertBuilder()
	ib.InsertInto("expense")
	ib.Cols(
		"name",
		"amount",
		"note",
		"date",
		"category_id",
		"user_id",
		"expense_group_id",
	)
	ib.Values(
		c.Name,
		c.Amount,
		c.Note,
		c.Date,
		c.CategoryID,
		user.ID,
		groupID,
	)
	ib.Returning(
		"id",
		"name",
		"amount",
		"note",
		"date",
		"created_at",
		"updated_at",
		"category_id",
		"expense_group_id",
	)

	q, args := ib.Build()

	logger.Infow("Create expense", "query", q, "args", args)

	var dst expenseWriteDst
	if err := tx.GetContext(ctx, &dst, q, args...); err != nil {
		var sqErr sqlite3.Error
		if errors.As(err, &sqErr) && sqErr.ExtendedCode == sqlite3.ErrConstraintForeignKey {
			return internal.Expense{}, internal.NewError(internal.ErrorCodeInvalid, "category not found")
		}
		return internal.Expense{}, fmt.Errorf("create expense: %w", err)
	}

	category, err := getCategory(ctx, tx, dst.CategoryID)
	if err != nil {
		return internal.Expense{}, fmt.Errorf("create expense: %w", err)
	}

	return internal.Expense{
		ID:        dst.ID,
		Name:      dst.Name,
		Amount:    dst.Amount,
		Note:      dst.Note,
		Date:      dst.Date,
		Category:  category,
		CreatedAt: dst.CreatedAt,
		UpdatedAt: dst.UpdatedAt,
	}, nil
}

type expenseWriteDst struct {
	ID             string          `db:"id"`
	Name           string          `db:"name"`
	Amount         internal.Amount `db:"amount"`
	Note           string          `db:"note"`
	Date           time.Time       `db:"date"`
	CreatedAt      time.Time       `db:"created_at"`
	UpdatedAt      time.Time       `db:"updated_at"`
	CategoryID     string          `db:"category_id"`
	ExpenseGroupID *string         `db:"expense_group_id"`
}

type expenseReadDst struct {
	ID                string          `db:"id"`
	Name              string          `db:"name"`
	Amount            internal.Amount `db:"amount"`
	Note              string          `db:"note"`
	Date              time.Time       `db:"date"`
	CreatedAt         time.Time       `db:"created_at"`
	UpdatedAt         time.Time       `db:"updated_at"`
	CategoryID        string          `db:"category_id"`
	CategoryName      string          `db:"category_name"`
	CategoryColor     string          `db:"color"`
	CategoryIcon      string          `db:"icon"`
	CategoryCreatedAt time.Time       `db:"category_created_at"`
	CategoryUpdatedAt time.Time       `db:"category_updated_at"`
}
