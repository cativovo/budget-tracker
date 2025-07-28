package sqlite

import (
	"cmp"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
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
		sb.EQ("e.id", id),
		sb.EQ("e.user_id", user.ID),
	)

	q, args := sb.Build()

	logger.Infow("Get expense", "query", q, "args", args)

	var dest struct {
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
	if err := es.db.Reader.GetContext(ctx, &dest, q, args...); err != nil {
		if err == sql.ErrNoRows {
			err = internal.NewError(internal.ErrorCodeNotFound, "expense not found")
		}

		return internal.Expense{}, fmt.Errorf("sqlite: get expense: %w", err)
	}

	return internal.Expense{
		ID:     dest.ID,
		Name:   dest.Name,
		Amount: dest.Amount,
		Note:   dest.Note,
		Date:   dest.Date,
		Category: internal.Category{
			ID:        dest.CategoryID,
			Name:      dest.CategoryName,
			Color:     dest.CategoryColor,
			Icon:      dest.CategoryIcon,
			CreatedAt: dest.CategoryCreatedAt,
			UpdatedAt: dest.CategoryUpdatedAt,
		},
		CreatedAt: dest.CreatedAt,
		UpdatedAt: dest.UpdatedAt,
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
	)
	ib.Values(
		c.Name,
		c.Amount,
		c.Note,
		c.Date,
		c.CategoryID,
		user.ID,
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
	)

	q, args := ib.Build()

	logger.Infow("Create expense", "query", q, "args", args)

	var dest expenseWriteDest
	if err := tx.GetContext(ctx, &dest, q, args...); err != nil {
		var sqErr sqlite3.Error
		if errors.As(err, &sqErr) && sqErr.ExtendedCode == sqlite3.ErrConstraintForeignKey {
			err = internal.NewError(internal.ErrorCodeInvalid, "category not found")
		}
		return internal.Expense{}, fmt.Errorf("sqlite: create expense: %w", err)
	}

	category, err := getCategory(ctx, tx, dest.CategoryID)
	if err != nil {
		return internal.Expense{}, fmt.Errorf("sqlite: create expense: %w", err)
	}

	expense := internal.Expense{
		ID:        dest.ID,
		Name:      dest.Name,
		Amount:    dest.Amount,
		Note:      dest.Note,
		Date:      dest.Date,
		Category:  category,
		CreatedAt: dest.CreatedAt,
		UpdatedAt: dest.UpdatedAt,
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

	tx, err := es.db.Writer.BeginTxx(ctx, nil)
	if err != nil {
		return internal.Expense{}, fmt.Errorf("sqlite: update expense begin: %w", err)
	}

	defer tx.Rollback()

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
		ub.EQ("id", u.ID),
		ub.EQ("user_id", user.ID),
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

	var dest expenseWriteDest
	if err := tx.GetContext(ctx, &dest, q, args...); err != nil {
		if err == sql.ErrNoRows {
			err = internal.NewError(internal.ErrorCodeNotFound, "expense not found")
		}

		var sqErr sqlite3.Error
		if errors.As(err, &sqErr) && sqErr.ExtendedCode == sqlite3.ErrConstraintForeignKey {
			err = internal.NewError(internal.ErrorCodeInvalid, "category not found")
		}

		return internal.Expense{}, fmt.Errorf("sqlite: update expense: %w", err)
	}

	category, err := getCategory(ctx, tx, dest.CategoryID)
	if err != nil {
		return internal.Expense{}, fmt.Errorf("sqlite: update expense get category: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return internal.Expense{}, fmt.Errorf("sqlite: update expense commit: %w", err)
	}

	return internal.Expense{
		ID:        dest.ID,
		Name:      dest.Name,
		Amount:    dest.Amount,
		Note:      dest.Note,
		Date:      dest.Date,
		Category:  category,
		CreatedAt: dest.CreatedAt,
		UpdatedAt: dest.UpdatedAt,
	}, nil
}

func (es *ExpenseService) DeleteExpense(ctx context.Context, id string) error {
	user := internal.UserFromContext(ctx)
	logger := internal.LoggerFromContext(ctx)

	db := sqlbuilder.SQLite.NewDeleteBuilder()
	db.DeleteFrom("expense")
	db.Where(
		db.EQ("id", id),
		db.EQ("user_id", user.ID),
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
	tx, err := es.db.Reader.BeginTxx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return internal.ExpenseGroup{}, fmt.Errorf("sqlite: get expense group begin: %w", err)
	}

	defer tx.Rollback()

	expenseGroup, err := getExpenseGroup(ctx, tx, id)
	if err != nil {
		return internal.ExpenseGroup{}, fmt.Errorf("sqlite: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return internal.ExpenseGroup{}, fmt.Errorf("sqlite: get expense group commit: %w", err)
	}

	return expenseGroup, nil
}

func (es *ExpenseService) CreateExpenseGroup(ctx context.Context, c internal.ExpenseGroupCreate) (internal.ExpenseGroup, error) {
	if err := c.Validate(); err != nil {
		return internal.ExpenseGroup{}, fmt.Errorf("sqlite: %w", err)
	}

	tx, err := es.db.Writer.BeginTxx(ctx, nil)
	if err != nil {
		return internal.ExpenseGroup{}, fmt.Errorf("sqlite: create expense group begin: %w", err)
	}

	defer tx.Rollback()

	user := internal.UserFromContext(ctx)
	logger := internal.LoggerFromContext(ctx)

	egib := builder.NewInsertBuilder()
	egib.InsertInto("expense_group")
	egib.Cols(
		"name",
		"note",
		"date",
		"category_id",
		"user_id",
	)
	egib.Values(
		c.Name,
		c.Note,
		c.Date,
		c.CategoryID,
		user.ID,
	)
	egib.Returning(
		"id",
		"name",
		"note",
		"date",
		"created_at",
		"updated_at",
	)

	q, args := egib.Build()

	logger.Infow("Create expense group", "query", q, "args", args)

	var expenseGroupDest struct {
		ID        string    `db:"id"`
		Name      string    `db:"name"`
		Note      string    `db:"note"`
		Date      time.Time `db:"date"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}
	if err := tx.GetContext(ctx, &expenseGroupDest, q, args...); err != nil {
		return internal.ExpenseGroup{}, fmt.Errorf("sqlite: create expense group: %w", err)
	}

	category, err := getCategory(ctx, tx, c.CategoryID)
	if err != nil {
		return internal.ExpenseGroup{}, fmt.Errorf("create expense group: %w", err)
	}

	expenses := make([]internal.Expense, 0)
	if len(c.Expenses) > 0 {
		eib := sqlbuilder.SQLite.NewInsertBuilder()
		eib.InsertInto("expense")
		eib.Cols(
			"name",
			"amount",
			"note",
			"date",
			"category_id",
			"user_id",
			"expense_group_id",
		)

		for _, v := range c.Expenses {
			eib.Values(
				v.Name,
				v.Amount,
				v.Note,
				c.Date,
				c.CategoryID,
				user.ID,
				expenseGroupDest.ID,
			)
		}

		eib.Returning(
			"id",
			"name",
			"amount",
			"note",
			"date",
			"created_at",
			"updated_at",
			"category_id",
		)

		q, args = eib.Build()
		logger.Infow("Create expense group expenses", "query", q, "args", args)

		rows, err := tx.QueryxContext(ctx, q, args...)
		if err != nil && err != sql.ErrNoRows {
			return internal.ExpenseGroup{}, fmt.Errorf("sqlite: create expense group expenses: %w", err)
		}

		for rows.Next() {
			var expenseDest expenseWriteDest
			if err := rows.StructScan(&expenseDest); err != nil {
				return internal.ExpenseGroup{}, fmt.Errorf("sqlite: create expense group expenses scan : %w", err)
			}

			e := internal.Expense{
				ID:        expenseDest.ID,
				Name:      expenseDest.Name,
				Amount:    expenseDest.Amount,
				Note:      expenseDest.Note,
				Date:      expenseDest.Date,
				CreatedAt: expenseDest.CreatedAt,
				UpdatedAt: expenseDest.UpdatedAt,
				Category:  category,
			}
			expenses = append(expenses, e)
		}

		// SQLite doesn't support using `RETURNING` inside a `WITH` clause
		// Sort the returned rows manually
		slices.SortStableFunc(expenses, func(a, b internal.Expense) int {
			return cmp.Compare(a.Name, b.Name)
		})
	}

	if err := tx.Commit(); err != nil {
		return internal.ExpenseGroup{}, fmt.Errorf("sqlite: create expense group commit: %w", err)
	}

	return internal.ExpenseGroup{
		ID:        expenseGroupDest.ID,
		Name:      expenseGroupDest.Name,
		Date:      expenseGroupDest.Date,
		Category:  category,
		Expenses:  expenses,
		Note:      expenseGroupDest.Note,
		CreatedAt: expenseGroupDest.CreatedAt,
		UpdatedAt: expenseGroupDest.UpdatedAt,
	}, nil
}

func (es *ExpenseService) UpdateExpenseGroup(ctx context.Context, u internal.ExpenseGroupUpdate) (internal.ExpenseGroup, error) {
	if err := u.Validate(); err != nil {
		return internal.ExpenseGroup{}, fmt.Errorf("sqlite: %w", err)
	}

	tx, err := es.db.Writer.BeginTxx(ctx, nil)
	if err != nil {
		return internal.ExpenseGroup{}, fmt.Errorf("sqlite: update expense group begin: %w", err)
	}

	defer tx.Rollback()

	user := internal.UserFromContext(ctx)
	logger := internal.LoggerFromContext(ctx)

	egub := sqlbuilder.SQLite.NewUpdateBuilder()
	egub.Update("expense_group")
	setMoreIfNotNil(egub, "name", u.Name)
	setMoreIfNotNil(egub, "date", u.Date)
	setMoreIfNotNil(egub, "note", u.Note)
	setUpdatedAt(egub)
	egub.Where(
		egub.EQ("id", u.ID),
		egub.EQ("user_id", user.ID),
	)

	q, args := egub.Build()

	logger.Infow("Update expense group", "query", q, "args", args)

	if _, err := tx.ExecContext(ctx, q, args...); err != nil {
		if err == sql.ErrNoRows {
			err = internal.NewError(internal.ErrorCodeNotFound, "expense group not found")
		}
		return internal.ExpenseGroup{}, fmt.Errorf("sqlite: update expense group: %w", err)
	}

	// must be executed after upserting expenses
	if u.CategoryID != nil || u.Date != nil {
		eub := builder.NewUpdateBuilder()
		eub.Update("expense")
		setMoreIfNotNil(eub, "category_id", u.CategoryID)
		setMoreIfNotNil(eub, "date", u.Date)
		eub.Where(
			eub.EQ("expense_group_id", u.ID),
			eub.EQ("user_id", user.ID),
		)

		logger.Infow("Update expense group expenses", "query", q, "args", args)

		if _, err := tx.ExecContext(ctx, q, args...); err != nil {
			if err == sql.ErrNoRows {
				err = internal.NewError(internal.ErrorCodeNotFound, "expense not found")
			}

			var sqErr sqlite3.Error
			if errors.As(err, &sqErr) && sqErr.ExtendedCode == sqlite3.ErrConstraintForeignKey {
				err = internal.NewError(internal.ErrorCodeInvalid, "category not found")
			}

			return internal.ExpenseGroup{}, fmt.Errorf("sqlite: update expense group expenses: %w", err)
		}
	}

	expenseGroup, err := getExpenseGroup(ctx, tx, u.ID)
	if err != nil {
		return internal.ExpenseGroup{}, fmt.Errorf("sqlite: update expense group: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return internal.ExpenseGroup{}, fmt.Errorf("sqlite: update expense group commit: %w", err)
	}

	return expenseGroup, nil
}

func (es *ExpenseService) DeleteExpenseGroup(ctx context.Context, id string) error {
	tx, err := es.db.Writer.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("sqlite: delete expense group begin: %w", err)
	}

	defer tx.Rollback()

	user := internal.UserFromContext(ctx)
	logger := internal.LoggerFromContext(ctx)

	egdb := sqlbuilder.SQLite.NewDeleteBuilder()
	egdb.DeleteFrom("expense_group")
	egdb.Where(
		egdb.EQ("id", id),
		egdb.EQ("user_id", user.ID),
	)

	q, args := egdb.Build()

	logger.Infow("Delete expense group", "query", q, "args", args)

	if _, err := tx.ExecContext(ctx, q, args...); err != nil {
		return fmt.Errorf("sqlite: delete expense group: %w", err)
	}

	edb := sqlbuilder.SQLite.NewDeleteBuilder()
	edb.DeleteFrom("expense")
	edb.Where(
		edb.EQ("expense_group_id", id),
		edb.EQ("user_id", user.ID),
	)

	q, args = edb.Build()

	logger.Infow("Delete expense group expenses", "query", q, "args", args)

	if _, err := tx.ExecContext(ctx, q, args...); err != nil {
		return fmt.Errorf("sqlite: delete expense group expenses: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("sqlite: delete expense group commit: %w", err)
	}

	return nil
}

func getExpenseGroup(ctx context.Context, tx *sqlx.Tx, id string) (internal.ExpenseGroup, error) {
	user := internal.UserFromContext(ctx)
	logger := internal.LoggerFromContext(ctx)

	esgb := sqlbuilder.SQLite.NewSelectBuilder()
	esgb.Select(
		"id",
		"name",
		"note",
		"date",
		"category_id",
		"created_at",
		"updated_at",
	)
	esgb.From("expense_group")
	esgb.Where(
		esgb.EQ("id", id),
		esgb.EQ("user_id", user.ID),
	)

	q, args := esgb.Build()

	logger.Infow("Get expense group", "query", q, "args", args)

	var expenseGroupDest struct {
		ID         string    `db:"id"`
		Name       string    `db:"name"`
		Note       string    `db:"note"`
		Date       time.Time `db:"date"`
		CategoryID string    `db:"category_id"`
		CreatedAt  time.Time `db:"created_at"`
		UpdatedAt  time.Time `db:"updated_at"`
	}
	if err := tx.GetContext(ctx, &expenseGroupDest, q, args...); err != nil {
		if err == sql.ErrNoRows {
			err = internal.NewError(internal.ErrorCodeNotFound, "expense group not found")
		}
		return internal.ExpenseGroup{}, fmt.Errorf("get expense group: %w", err)
	}

	category, err := getCategory(ctx, tx, expenseGroupDest.CategoryID)
	if err != nil {
		return internal.ExpenseGroup{}, fmt.Errorf("get expense group: %w", err)
	}

	esb := sqlbuilder.SQLite.NewSelectBuilder()
	esb.Select(
		"id",
		"name",
		"amount",
		"note",
		"date",
		"created_at",
		"updated_at",
	)
	esb.From("expense")
	esb.Where(
		esb.EQ("expense_group_id", id),
		esb.EQ("user_id", user.ID),
	)
	esb.OrderBy("name")

	q, args = esb.Build()

	logger.Infow("Get expense group expenses", "query", q, "args", args)

	rows, err := tx.QueryxContext(ctx, q, args...)
	if err != nil && err != sql.ErrNoRows {
		return internal.ExpenseGroup{}, fmt.Errorf("get expense group expenses: %w", err)
	}

	expenses := make([]internal.Expense, 0)
	for rows.Next() {
		var expenseDest struct {
			ID        string          `db:"id"`
			Name      string          `db:"name"`
			Amount    internal.Amount `db:"amount"`
			Note      string          `db:"note"`
			Date      time.Time       `db:"date"`
			CreatedAt time.Time       `db:"created_at"`
			UpdatedAt time.Time       `db:"updated_at"`
		}
		if err := rows.StructScan(&expenseDest); err != nil {
			return internal.ExpenseGroup{}, fmt.Errorf("get expense group expenses scan: %w", err)
		}

		e := internal.Expense{
			ID:        expenseDest.ID,
			Name:      expenseDest.Name,
			Amount:    expenseDest.Amount,
			Note:      expenseDest.Note,
			Date:      expenseDest.Date,
			Category:  category,
			CreatedAt: expenseDest.CreatedAt,
			UpdatedAt: expenseDest.UpdatedAt,
		}
		expenses = append(expenses, e)
	}

	return internal.ExpenseGroup{
		ID:        expenseGroupDest.ID,
		Name:      expenseGroupDest.Name,
		Note:      expenseGroupDest.Note,
		Date:      expenseGroupDest.Date,
		Category:  category,
		Expenses:  expenses,
		CreatedAt: expenseGroupDest.CreatedAt,
		UpdatedAt: expenseGroupDest.UpdatedAt,
	}, nil
}

type expenseWriteDest struct {
	ID         string          `db:"id"`
	Name       string          `db:"name"`
	Amount     internal.Amount `db:"amount"`
	Note       string          `db:"note"`
	Date       time.Time       `db:"date"`
	CreatedAt  time.Time       `db:"created_at"`
	UpdatedAt  time.Time       `db:"updated_at"`
	CategoryID string          `db:"category_id"`
}

// TODO: check if should be deleted
type expenseReadDest struct {
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
