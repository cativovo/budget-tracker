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
	if err := cs.db.Writer.GetContext(ctx, &dest, q, args...); err != nil {
		return category.Category{}, fmt.Errorf("sqlite: create category: %w", err)
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
