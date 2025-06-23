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

func (cs *CategoryService) FindCategories(ctx context.Context, f internal.CategoryFilter) ([]internal.Category, int, error) {
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
		// SQLite doesn't support ILIKE
		fsb.Where(fsb.Like("LOWER(name)", "%"+*f.Name+"%"))
	}

	fsb.OrderBy("updated_at")
	fsb.Limit(f.Limit)
	fsb.Offset(f.Offset)

	q, args := fsb.Build()

	logger.Infow("Find categories", "query", q, "args", args)

	rows, err := cs.db.reader.QueryxContext(ctx, q, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, 0, nil
		}

		return nil, 0, fmt.Errorf("sqlite: find categories: %w", err)
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
			return nil, 0, fmt.Errorf("sqlite: scan category: %w", err)
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
	if err := cs.db.reader.GetContext(ctx, &total, q, args...); err != nil {
		return nil, 0, fmt.Errorf("sqlite: count categories: %w", err)
	}

	return categories, total, nil
}

func (cs *CategoryService) GetCategory(ctx context.Context, id string) (internal.Category, error) {
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
	if err := cs.db.reader.GetContext(ctx, &dst, q, args...); err != nil {
		if err == sql.ErrNoRows {
			iErr := internal.NewError(internal.ErrorCodeNotFound, "Category not found")
			return internal.Category{}, fmt.Errorf("sqlite: %w", iErr)
		}

		return internal.Category{}, fmt.Errorf("sqlite: get category: %w", err)
	}

	return internal.Category(dst), nil
}

func (cs *CategoryService) CreateCategory(ctx context.Context, c internal.CategoryCreate) (internal.Category, error) {
	if err := c.Validate(); err != nil {
		return internal.Category{}, err
	}

	tx, err := cs.db.writer.BeginTxx(ctx, nil)
	if err != nil {
		return internal.Category{}, fmt.Errorf("sqlite: create category begin tx: %w", err)
	}

	logger := internal.LoggerFromContext(ctx)

	defer func() {
		if err := tx.Rollback(); err != nil {
			logger.Error("Failed to rollback", zap.Error(err))
		}
	}()

	if err := validateCategory(ctx, tx, c.Name); err != nil {
		return internal.Category{}, fmt.Errorf("sqlite: validate create category: %w", err)
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

	tx, err := cs.db.writer.BeginTxx(ctx, nil)
	if err != nil {
		return internal.Category{}, fmt.Errorf("sqlite: update category begin tx: %w", err)
	}

	logger := internal.LoggerFromContext(ctx)

	defer func() {
		if err := tx.Rollback(); err != nil {
			logger.Error("Failed to rollback", zap.Error(err))
		}
	}()

	if u.Name != nil {
		if err := validateCategory(ctx, tx, *u.Name); err != nil {
			return internal.Category{}, fmt.Errorf("sqlite: validate update: %w", err)
		}
	}

	user := internal.UserFromContext(ctx)

	ub := sqlbuilder.SQLite.NewUpdateBuilder()
	ub.Update("category")

	setMoreIfNotNil(ub, "name", u.Name)
	setMoreIfNotNil(ub, "color", u.Color)
	setMoreIfNotNil(ub, "icon", u.Icon)

	ub.Where(
		ub.Equal("id", u.ID),
		ub.Equal("user_id", user.ID),
	)

	// https://github.com/huandu/go-sqlbuilder/issues/142
	ub.SQL("RETURNING id, name, color, icon, created_at, updated_at")

	q, args := ub.Build()

	logger.Infow("Update category", "query", q, "args", args)

	var dst categoryDst
	if err := tx.GetContext(ctx, &dst, q, args...); err != nil {
		if err == sql.ErrNoRows {
			iErr := internal.NewError(internal.ErrorCodeNotFound, "Category not found")
			return internal.Category{}, fmt.Errorf("sqlite: %w", iErr)
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

	if _, err := cs.db.writer.ExecContext(ctx, q, args...); err != nil {
		return fmt.Errorf("sqlite: delete category: %w", err)
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

func validateCategory(ctx context.Context, tx *sqlx.Tx, name string) error {
	ok, err := categoryExists(ctx, tx, name)
	if err != nil {
		return err
	}
	if ok {
		return internal.NewError(internal.ErrorCodeConflict, "Category already exists")
	}
	return nil
}

func categoryExists(ctx context.Context, tx *sqlx.Tx, name string) (bool, error) {
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

	var dst bool
	if err := tx.GetContext(ctx, &dst, q, args...); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("category exists: %w", err)
	}

	return dst, nil
}
