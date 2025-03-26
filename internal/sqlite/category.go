package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/cativovo/budget-tracker/internal"
	"github.com/cativovo/budget-tracker/internal/category"
	"github.com/cativovo/budget-tracker/internal/logger"
	"github.com/cativovo/budget-tracker/internal/user"
	"github.com/huandu/go-sqlbuilder"
)

type CategoryRepository struct {
	db *DB
}

var _ category.Repository = (*CategoryRepository)(nil)

func NewCategoryRepository(db *DB) CategoryRepository {
	return CategoryRepository{
		db: db,
	}
}

func (cr *CategoryRepository) FindCategoryByID(ctx context.Context, id string) (category.Category, error) {
	u := user.UserFromCtx(ctx)
	logger := logger.LoggerFromCtx(ctx)

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
		sb.And(
			sb.EQ("id", id),
			sb.EQ("user_id", u.ID),
		),
	)

	q, args := sb.Build()

	logger.Infow(
		"Find category by id",
		"query", q,
		"args", args,
	)

	var dst categoryDst
	if err := cr.db.reader.GetContext(ctx, &dst, q, args...); err != nil {
		if err == sql.ErrNoRows {
			return category.Category{}, internal.NewError(internal.ErrorCodeNotFound, "Category not found")
		}

		return category.Category{}, fmt.Errorf("sqlite.CategoryRepository: %w", err)
	}

	return category.Category(dst), nil
}

func (cr *CategoryRepository) ListCategories(ctx context.Context, o internal.ListOptions) ([]category.Category, error) {
	panic("not yet implemented")
}

func (cr *CategoryRepository) findCategoryByName(ctx context.Context, name string) (category.Category, error) {
	u := user.UserFromCtx(ctx)
	logger := logger.LoggerFromCtx(ctx)

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
		sb.And(
			sb.EQ("name", name),
			sb.EQ("user_id", u.ID),
		),
	)

	q, args := sb.Build()

	logger.Infow(
		"Count category by name",
		"query", q,
		"args", args,
	)

	var dst categoryDst
	if err := cr.db.reader.GetContext(ctx, &dst, q, args...); err != nil {
		return category.Category{}, fmt.Errorf("sqlite.CategoryRepository: %w", err)
	}

	return category.Category(dst), nil
}

func (cr *CategoryRepository) CreateCategory(ctx context.Context, c category.CreateCategoryReq) (category.Category, error) {
	_, err := cr.findCategoryByName(ctx, c.Name)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return category.Category{}, err
	}
	if err == nil {
		return category.Category{}, internal.NewErrorf(internal.ErrorCodeConflict, "%s category already exists", c.Name)
	}

	u := user.UserFromCtx(ctx)
	logger := logger.LoggerFromCtx(ctx)

	ib := sqlbuilder.SQLite.NewInsertBuilder()
	ib.InsertInto("category")
	ib.Cols(
		"name",
		"color",
		"icon",
		"user_id",
	)
	ib.Values(
		c.Name,
		c.Color,
		c.Icon,
		u.ID,
	)
	ib.Returning(
		"id",
		"name",
		"color",
		"icon",
		"created_at",
		"updated_at",
	)

	q, args := ib.Build()

	logger.Infow(
		"Insert new category",
		"query", q,
		"args", args,
	)

	var dst categoryDst
	if err := cr.db.readerWriter.GetContext(ctx, &dst, q, args...); err != nil {
		return category.Category{}, fmt.Errorf("sqlite.CategoryRepository: %w", err)
	}

	return category.Category(dst), nil
}

func (cr *CategoryRepository) UpdateCategory(ctx context.Context, c category.UpdateCategoryReq) (category.Category, error) {
	found, err := cr.findCategoryByName(ctx, c.Name)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return category.Category{}, err
	}
	if err == nil && found.ID != c.ID {
		return category.Category{}, internal.NewErrorf(internal.ErrorCodeConflict, "%s category already exists", c.Name)
	}

	u := user.UserFromCtx(ctx)
	logger := logger.LoggerFromCtx(ctx)

	ub := sqlbuilder.SQLite.NewUpdateBuilder()
	ub.Update("category")

	if c.Name != "" {
		ub.SetMore(ub.Assign("name", c.Name))
	}
	if c.Color != "" {
		ub.SetMore(ub.Assign("color", c.Color))
	}
	if c.Icon != "" {
		ub.SetMore(ub.Assign("icon", c.Icon))
	}

	ub.Where(
		ub.And(
			ub.EQ("id", c.ID),
			ub.EQ("user_id", u.ID),
		),
	)
	// https://github.com/huandu/go-sqlbuilder/issues/142
	ub.SQL("RETURNING id, name, color, icon, created_at, updated_at")

	q, args := ub.Build()

	logger.Infow(
		"Insert new category",
		"query", q,
		"args", args,
	)

	var dst categoryDst
	if err := cr.db.readerWriter.GetContext(ctx, &dst, q, args...); err != nil {
		if err == sql.ErrNoRows {
			return category.Category{}, internal.NewError(internal.ErrorCodeNotFound, "Category not found")
		}

		return category.Category{}, fmt.Errorf("sqlite.CategoryRepository: %w", err)
	}

	return category.Category(dst), nil
}

func (cr *CategoryRepository) DeleteCategory(ctx context.Context, id string) error {
	user := user.UserFromCtx(ctx)
	logger := logger.LoggerFromCtx(ctx)

	db := sqlbuilder.SQLite.NewDeleteBuilder()
	db.DeleteFrom("category")
	db.Where(
		db.And(
			db.EQ("id", id),
			db.EQ("user_id", user.ID),
		),
	)

	q, args := db.Build()

	logger.Infow(
		"Delete category",
		"query", q,
		"args", args,
	)

	_, err := cr.db.readerWriter.ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("sqlite.CategoryRepository: %w", err)
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
