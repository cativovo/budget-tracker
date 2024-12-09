package repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"go.uber.org/zap"
)

type Category struct {
	ID       string
	Name     string
	Icon     string
	ColorHex string
}

type GetCategoryByIDParams struct {
	ID        string
	AccountID string
}

func (r *Repository) GetCategoryByID(ctx context.Context, logger *zap.SugaredLogger, p GetCategoryByIDParams) (Category, error) {
	q, args, err := sq.
		Select("id", "name", "icon", "color_hex").
		From("category").
		Where(sq.Eq{"id": p.ID, "account_id": p.AccountID}).
		ToSql()
	if err != nil {
		return Category{}, fmt.Errorf("GetCategoryByID: query builder failed: %w", err)
	}

	logger.Infow(
		"Executing query",
		"query", q,
		"args", []any{p.ID, p.AccountID},
	)

	var dest struct {
		CreatedAt string `db:"created_at"`
		UpdatedAt string `db:"updated_at"`
		ID        string `db:"id"`
		Name      string `db:"name"`
		Icon      string `db:"icon"`
		ColorHex  string `db:"color_hex"`
	}
	err = r.ConcurrentDB().GetContext(ctx, &dest, q, args...)
	if err != nil {
		return Category{}, fmt.Errorf("GetCategoryByID: query failed: %w", err)
	}

	return Category{
		ID:       dest.ID,
		Name:     dest.Name,
		Icon:     dest.Icon,
		ColorHex: dest.ColorHex,
	}, nil
}

type CreateCategoryParams struct {
	Name      string
	Icon      string
	ColorHex  string
	AccountID string
}

func (r *Repository) CreateCategory(ctx context.Context, logger *zap.SugaredLogger, p CreateCategoryParams) (Category, error) {
	q, args, err := sq.
		Insert("category").
		Columns("name", "icon", "color_hex", "account_id").
		Values(p.Name, p.Icon, p.ColorHex, p.AccountID).
		Suffix(`RETURNING "id", "name", "icon", "color_hex"`).
		ToSql()
	if err != nil {
		return Category{}, fmt.Errorf("CreateCategory: query builder failed: %w", err)
	}

	logger.Infow(
		"Executing query",
		"query", q,
		"args", []any{p.Name, p.Icon, p.ColorHex, p.AccountID},
	)

	var dest struct {
		CreatedAt string `db:"created_at"`
		UpdatedAt string `db:"updated_at"`
		ID        string `db:"id"`
		Name      string `db:"name"`
		Icon      string `db:"icon"`
		ColorHex  string `db:"color_hex"`
	}
	err = r.NonConcurrentDB().GetContext(ctx, &dest, q, args...)
	if err != nil {
		return Category{}, fmt.Errorf("CreateCategory: query failed: %w", err)
	}

	return Category{
		ID:       dest.ID,
		Name:     dest.Name,
		Icon:     dest.Icon,
		ColorHex: dest.ColorHex,
	}, nil
}
