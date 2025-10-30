package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/cativovo/budget-tracker/internal/apperror"
	"github.com/cativovo/budget-tracker/internal/category"
	"github.com/cativovo/budget-tracker/internal/log"
	"github.com/cativovo/budget-tracker/internal/user"
	"github.com/jmoiron/sqlx"
)

// CategoryStore handles category database operations.
type CategoryStore struct {
	db *DB
}

// NewCategoryStore returns a new CategoryStore.
func NewCategoryStore(db *DB) *CategoryStore {
	return &CategoryStore{db: db}
}

// GetCategoryByID returns a category from the database.
func (cs *CategoryStore) GetCategoryByID(ctx context.Context, id string) (category.Category, error) {
	u := user.FromContext(ctx)
	sb := builder.NewSelectBuilder()
	sb.Select(
		"id",
		"name",
		"color",
		"icon",
		"created_at",
		"updated_at",
	)
	sb.From("category")
	sb.Where(sb.Equal("id", id), sb.Equal("user_id", u.ID))

	q, args := sb.Build()

	logger := log.FromContext(ctx)
	logger.Info("Get category by ID", "query", q, "args", args)

	var dest struct {
		ID        string    `db:"id"`
		Name      string    `db:"name"`
		Color     string    `db:"color"`
		Icon      string    `db:"icon"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}
	if err := cs.db.Reader.GetContext(ctx, &dest, q, args...); err != nil {
		if err == sql.ErrNoRows {
			err = apperror.New(apperror.ErrorCodeNotFound, "category not found")
		}
		return category.Category{}, fmt.Errorf("sqlite: get category by id: %w", err)
	}
	return category.Category(dest), nil
}

// CreateCategory creates a new category in database.
func (cs *CategoryStore) CreateCategory(ctx context.Context, input category.CreateCategoryInput) (category.Category, error) {
	tx, err := cs.db.Writer.BeginTxx(ctx, nil)
	if err != nil {
		return category.Category{}, fmt.Errorf("sqlite: create category begin tx: %w", err)
	}

	defer tx.Rollback()

	if err := categoryExists(ctx, tx, input.Name, nil); err != nil {
		return category.Category{}, fmt.Errorf("sqlite: %w", err)
	}

	u := user.FromContext(ctx)
	ib := builder.NewInsertBuilder()
	ib.InsertInto("category")
	ib.Cols(
		"name",
		"color",
		"icon",
		"user_id",
	)
	ib.Values(
		input.Name,
		input.Color,
		input.Icon,
		u.ID,
	)
	ib.Returning("id", "created_at", "updated_at")

	q, args := ib.Build()

	logger := log.FromContext(ctx)
	logger.Info("Create category", "query", q, "args", args)

	var dest struct {
		ID        string    `db:"id"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}
	if err := tx.GetContext(ctx, &dest, q, args...); err != nil {
		return category.Category{}, fmt.Errorf("sqlite: create category: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return category.Category{}, fmt.Errorf("sqlite: create category commit tx: %w", err)
	}

	return category.Category{
		ID:        dest.ID,
		Name:      input.Name,
		Color:     input.Color,
		Icon:      input.Icon,
		CreatedAt: dest.CreatedAt,
		UpdatedAt: dest.UpdatedAt,
	}, nil
}

// UpdateCategory updates the category with the given ID.
func (cs *CategoryStore) UpdateCategory(ctx context.Context, id string, input category.UpdateCategoryInput) (category.Category, error) {
	tx, err := cs.db.Writer.BeginTxx(ctx, nil)
	if err != nil {
		return category.Category{}, fmt.Errorf("sqlite: update category begin tx: %w", err)
	}

	defer tx.Rollback()

	if input.Name != nil {
		if err := categoryExists(ctx, tx, *input.Name, &id); err != nil {
			return category.Category{}, fmt.Errorf("sqlite: update category: %w", err)
		}
	}

	u := user.FromContext(ctx)
	ub := builder.NewUpdateBuilder()
	ub.Update("category")
	setMoreIfNotNil(ub, "name", input.Name)
	setMoreIfNotNil(ub, "color", input.Color)
	setMoreIfNotNil(ub, "icon", input.Icon)
	ub.Where(
		ub.EQ("user_id", u.ID),
		ub.EQ("id", id),
	)
	// https://github.com/huandu/go-sqlbuilder/issues/142
	ub.SQL("RETURNING name, color, icon, created_at, updated_at")

	q, args := ub.Build()

	logger := log.FromContext(ctx)
	logger.Info("Update category", "query", q, log.SafeValuesAttr("args", args))

	var dest struct {
		Name      string    `db:"name"`
		Color     string    `db:"color"`
		Icon      string    `db:"icon"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}
	if err := tx.GetContext(ctx, &dest, q, args...); err != nil {
		if err == sql.ErrNoRows {
			err = apperror.New(apperror.ErrorCodeNotFound, "category not found")
		}
		return category.Category{}, fmt.Errorf("sqlite: update category: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return category.Category{}, fmt.Errorf("sqlite: update category commit tx: %w", err)
	}

	return category.Category{
		ID:        id,
		Name:      dest.Name,
		Color:     dest.Color,
		Icon:      dest.Icon,
		CreatedAt: dest.CreatedAt,
		UpdatedAt: dest.UpdatedAt,
	}, nil
}

func categoryExists(ctx context.Context, tx *sqlx.Tx, name string, id *string) error {
	u := user.FromContext(ctx)
	categorySB := builder.NewSelectBuilder()
	categorySB.Select("1")
	categorySB.From("category")
	categorySB.Where(
		categorySB.EQ("name", name),
		categorySB.EQ("user_id", u.ID),
	)

	if id != nil {
		categorySB.Where(categorySB.NotEqual("id", id))
	}

	existsSB := builder.NewSelectBuilder()
	existsSB.Select(existsSB.Exists(categorySB))

	q, args := existsSB.Build()

	logger := log.FromContext(ctx)
	logger.Info("Check category existence", "query", q, log.SafeValuesAttr("args", args))

	var exists bool
	if err := tx.GetContext(ctx, &exists, q, args...); err != nil {
		return fmt.Errorf("category exists: %w", err)
	}

	if exists {
		return apperror.New(apperror.ErrorCodeConflict, fmt.Sprintf("%s category already exists", name))
	}

	return nil
}
