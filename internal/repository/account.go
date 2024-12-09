package repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"go.uber.org/zap"
)

type Account struct {
	ID   string
	Name string
}

func (r *Repository) GetAccountByID(ctx context.Context, logger *zap.SugaredLogger, id string) (Account, error) {
	q, args, err := sq.
		Select("id", "name").
		From("account").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return Account{}, fmt.Errorf("GetAccountByID: query builder failed: %w", err)
	}

	logger.Infow(
		"Executing query",
		"query", q,
		"args", []any{id},
	)

	var dest struct {
		ID   string `db:"id"`
		Name string `db:"name"`
	}
	err = r.ConcurrentDB().GetContext(ctx, &dest, q, args...)
	if err != nil {
		return Account{}, fmt.Errorf("GetAccountByID: query failed: %w", err)
	}

	return Account{
		ID:   dest.ID,
		Name: dest.Name,
	}, nil
}

type CreateAccountParams struct {
	Name string
}

func (r *Repository) CreateAccount(ctx context.Context, logger *zap.SugaredLogger, p CreateAccountParams) (Account, error) {
	q, args, err := sq.
		Insert("account").
		Columns("name").
		Values(p.Name).
		Suffix(`RETURNING "id", "name"`).
		ToSql()
	if err != nil {
		return Account{}, fmt.Errorf("CreateAccount: query builder failed: %w", err)
	}

	logger.Infow(
		"Executing query",
		"query", q,
		"args", []any{p.Name},
	)

	var dest struct {
		ID   string `db:"id"`
		Name string `db:"name"`
	}
	err = r.NonConcurrentDB().GetContext(ctx, &dest, q, args...)
	if err != nil {
		return Account{}, fmt.Errorf("CreateAccount: query failed: %w", err)
	}

	return Account{
		ID:   dest.ID,
		Name: dest.Name,
	}, nil
}
