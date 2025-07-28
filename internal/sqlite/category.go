package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/cativovo/budget-tracker/internal"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type CategoryService struct {
	db *DB
}

var _ internal.CategoryService = (*CategoryService)(nil)

func NewCategoryService(db *DB) *CategoryService {
	return &CategoryService{
		db: db,
	}
}

func (cs *CategoryService) ListCategories(ctx context.Context, f internal.CategoryFilter) ([]internal.Category, int, error) {
	tx, err := cs.db.Reader.BeginTxx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, 0, fmt.Errorf("sqlite: list categories begin: %w", err)
	}

	defer tx.Rollback()

	user := internal.UserFromContext(ctx)
	logger := internal.LoggerFromContext(ctx)

	fsb := sqlbuilder.SQLite.NewSelectBuilder()
	fsb.Select(
		"id",
		"name",
		"color",
		"icon",
		"created_at",
		"updated_at",
	)
	fsb.From("category")
	fsb.Where(fsb.EQ("user_id", user.ID))

	if f.Name != nil {
		// `LIKE` in SQLite is case insensitive for ASCII
		fsb.Where(fsb.Like("name", "%"+*f.Name+"%"))
	}

	fsb.OrderBy("updated_at")
	fsb.Limit(f.Limit)
	fsb.Offset(f.Offset)

	q, args := fsb.Build()

	logger.Infow("List categories", "query", q, "args", args)

	rows, err := tx.QueryxContext(ctx, q, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, 0, nil
		}

		return nil, 0, fmt.Errorf("sqlite: list categories: %w", err)
	}

	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error("Failed to close the rows", zap.Error(err))
		}
	}()

	categories := make([]internal.Category, 0)
	for rows.Next() {
		var dst categoryDst
		err := rows.StructScan(&dst)
		if err != nil {
			return nil, 0, fmt.Errorf("sqlite: list categories scan: %w", err)
		}

		categories = append(categories, internal.Category(dst))
	}

	csb := sqlbuilder.SQLite.NewSelectBuilder()
	csb.Select(csb.As("COUNT(id)", "total"))
	csb.From("category")
	csb.WhereClause = fsb.WhereClause

	q, args = csb.Build()

	logger.Infow("Count categories", "query", q, "args", args)

	var total int
	if err := tx.GetContext(ctx, &total, q, args...); err != nil {
		return nil, 0, fmt.Errorf("sqlite: list categories count: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, 0, fmt.Errorf("sqlite: list categories commit: %w", err)
	}

	return categories, total, nil
}

func (cs *CategoryService) GetCategory(ctx context.Context, id string) (internal.Category, error) {
	c, err := getCategory(ctx, cs.db.Reader, id)
	if err != nil {
		return internal.Category{}, fmt.Errorf("sqlite: %w", err)
	}

	return c, nil
}

func (cs *CategoryService) CreateCategory(ctx context.Context, c internal.CategoryCreate) (internal.Category, error) {
	if err := c.Validate(); err != nil {
		return internal.Category{}, err
	}

	tx, err := cs.db.Writer.BeginTxx(ctx, nil)
	if err != nil {
		return internal.Category{}, fmt.Errorf("sqlite: create category begin tx: %w", err)
	}

	logger := internal.LoggerFromContext(ctx)

	defer tx.Rollback()

	if err := categoryExists(ctx, tx, c.Name); err != nil {
		return internal.Category{}, fmt.Errorf("sqlite:  %w", err)
	}

	user := internal.UserFromContext(ctx)

	ib := sqlbuilder.SQLite.NewInsertBuilder()
	ib.InsertInto("category")
	ib.Cols("name", "color", "icon", "user_id")
	ib.Values(c.Name, c.Color, c.Icon, user.ID)
	ib.Returning(
		"id",
		"name",
		"color",
		"icon",
		"created_at",
		"updated_at",
	)

	q, args := ib.Build()

	logger.Infow("Insert category", "query", q, "args", args)

	var dst categoryDst
	if err := tx.GetContext(ctx, &dst, q, args...); err != nil {
		return internal.Category{}, fmt.Errorf("sqlite: create category: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return internal.Category{}, fmt.Errorf("sqlite: create category commit: %w", err)
	}

	return internal.Category(dst), nil
}

func (cs *CategoryService) UpdateCategory(ctx context.Context, u internal.CategoryUpdate) (internal.Category, error) {
	if err := u.Validate(); err != nil {
		return internal.Category{}, err
	}

	tx, err := cs.db.Writer.BeginTxx(ctx, nil)
	if err != nil {
		return internal.Category{}, fmt.Errorf("sqlite: update category begin tx: %w", err)
	}

	logger := internal.LoggerFromContext(ctx)

	defer tx.Rollback()

	if u.Name != nil {
		if err := categoryExists(ctx, tx, *u.Name); err != nil {
			return internal.Category{}, fmt.Errorf("sqlite: %w", err)
		}
	}

	user := internal.UserFromContext(ctx)

	ub := sqlbuilder.SQLite.NewUpdateBuilder()
	ub.Update("category")

	setMoreIfNotNil(ub, "name", u.Name)
	setMoreIfNotNil(ub, "color", u.Color)
	setMoreIfNotNil(ub, "icon", u.Icon)
	setUpdatedAt(ub)

	ub.Where(
		ub.EQ("id", u.ID),
		ub.EQ("user_id", user.ID),
	)

	// https://github.com/huandu/go-sqlbuilder/issues/142
	ub.SQL("RETURNING id, name, color, icon, created_at, updated_at")

	q, args := ub.Build()

	logger.Infow("Update category", "query", q, "args", args)

	var dst categoryDst
	if err := tx.GetContext(ctx, &dst, q, args...); err != nil {
		if err == sql.ErrNoRows {
			err = internal.NewError(internal.ErrorCodeNotFound, "category not found")
		}
		return internal.Category{}, fmt.Errorf("sqlite: update category: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return internal.Category{}, fmt.Errorf("sqlite: update category commit: %w", err)
	}

	return internal.Category(dst), nil
}

func (cs *CategoryService) DeleteCategory(ctx context.Context, id string) error {
	user := internal.UserFromContext(ctx)
	logger := internal.LoggerFromContext(ctx)

	db := sqlbuilder.SQLite.NewDeleteBuilder()
	db.DeleteFrom("category")
	db.Where(
		db.EQ("id", id),
		db.EQ("user_id", user.ID),
	)

	q, args := db.Build()

	logger.Infow("Delete category", "query", q, "args", args)

	_, err := cs.db.Writer.ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("sqlite: delete category: %w", err)
	}

	return nil
}

func getCategory(ctx context.Context, queryer sqlx.QueryerContext, id string) (internal.Category, error) {
	user := internal.UserFromContext(ctx)
	logger := internal.LoggerFromContext(ctx)

	sb := sqlbuilder.SQLite.NewSelectBuilder()
	sb.Select(
		"id",
		"name",
		"color",
		"icon",
		"created_at",
		"updated_at",
	)
	sb.From("category")
	sb.Where(
		sb.EQ("id", id),
		sb.EQ("user_id", user.ID),
	)

	q, args := sb.Build()

	logger.Infow("Get category", "query", q, "args", args)

	var dst categoryDst
	// can't use `GetContext` here because it's a part of any interface
	row := queryer.QueryRowxContext(ctx, q, args...)
	if err := row.StructScan(&dst); err != nil {
		if err == sql.ErrNoRows {
			return internal.Category{}, internal.NewError(internal.ErrorCodeNotFound, "category not found")
		}

		return internal.Category{}, fmt.Errorf("get category: %w", err)
	}

	return internal.Category(dst), nil
}

func categoryExists(ctx context.Context, tx *sqlx.Tx, name string) error {
	user := internal.UserFromContext(ctx)
	logger := internal.LoggerFromContext(ctx)

	categorySb := sqlbuilder.SQLite.NewSelectBuilder()
	categorySb.Select("1")
	categorySb.From("category")
	categorySb.Where(
		categorySb.EQ("name", name),
		categorySb.EQ("user_id", user.ID),
	)

	existsSb := sqlbuilder.SQLite.NewSelectBuilder()
	existsSb.Select(existsSb.Exists(categorySb))

	q, args := existsSb.Build()

	logger.Infow("Check category existence", "query", q, "args", args)

	var exists bool
	if err := tx.GetContext(ctx, &exists, q, args...); err != nil {
		return fmt.Errorf("category exists: %w", err)
	}

	if exists {
		return internal.NewError(internal.ErrorCodeConflict, "Category already exists")
	}

	return nil
}

type categoryDst struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Color     string    `db:"color"`
	Icon      string    `db:"icon"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
