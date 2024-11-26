package repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/cativovo/budget-tracker/internal/models"
)

type account struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

func (r *Repository) GetAccountByID(ctx context.Context, id string) (models.Account, error) {
	q, args, err := sq.
		Select("id", "name").
		From("account").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return models.Account{}, fmt.Errorf("GetAccountByID: query builder failed: %w", err)
	}

	r.logger.Infow("Executing query", "query", q, "args", args)

	var dest account
	err = r.ConcurrentDB().GetContext(ctx, &dest, q, args...)
	if err != nil {
		return models.Account{}, fmt.Errorf("GetAccountByID: query failed: %w", err)
	}

	return models.Account{
		ID:   dest.ID,
		Name: dest.Name,
	}, nil
}

type CreateAccountParams struct {
	Name string
}

func (r *Repository) CreateAccount(ctx context.Context, p CreateAccountParams) (models.Account, error) {
	q, args, err := sq.
		Insert("account").
		Columns("name").
		Values(p.Name).
		Suffix(`RETURNING "id", "name"`).
		ToSql()
	if err != nil {
		return models.Account{}, fmt.Errorf("CreateAccount: query builder failed: %w", err)
	}

	r.logger.Infow("Executing query", "query", q, "args", args)

	var dest account
	err = r.NonConcurrentDB().GetContext(ctx, &dest, q, args...)
	if err != nil {
		return models.Account{}, fmt.Errorf("CreateAccount: query failed: %w", err)
	}

	return models.Account{
		ID:   dest.ID,
		Name: dest.Name,
	}, nil
}
