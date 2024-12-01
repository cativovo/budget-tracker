package repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/cativovo/budget-tracker/internal/models"
)

type ListEntriesByDateParams struct {
	StartDate string
	EndDate   string
	AccountID string
	EntryType []models.EntryType
	OrderBy   OrderBy
	Limit     int
	Offset    int
}

func (r *Repository) ListEntriesByDate(ctx context.Context, p ListEntriesByDateParams) (models.EntriesWithCount, error) {
	countChan := make(chan result[int])
	entriesChan := make(chan result[[]models.Entry])

	eq := sq.Eq{
		"entry.account_id": p.AccountID,
		"entry.entry_type": p.EntryType,
	}

	between := sq.Expr(
		"entry.date BETWEEN ? AND ?",
		p.StartDate,
		p.EndDate,
	)

	go func() {
		q, args, err := sq.
			Select("COUNT(id)").
			From("entry").
			Where(eq).
			Where(between).
			ToSql()
		if err != nil {
			countChan <- result[int]{
				err: fmt.Errorf("ListEntriesByDate: query builder failed: %w", err),
			}
			return
		}

		r.logger.Infow("Executing query", "query", q, "args", args)

		var count int
		err = r.concurrentDB.GetContext(ctx, &count, q, args...)
		if err != nil {
			countChan <- result[int]{
				err: fmt.Errorf("ListEntriesByDate: query failed: %w", err),
			}
			return
		}

		countChan <- result[int]{
			ok: count,
		}
	}()

	go func() {
		builder := sq.
			Select(
				"entry.amount as amount",
				"entry.created_at as created_at",
				"entry.date as date",
				"entry.description as description",
				"entry.entry_type as entry_type",
				"entry.id as id",
				"entry.name as name",
				"entry.updated_at as updated_at",
				"category.id as category_id",
				"category.name as category_name",
				"category.icon as category_icon",
				"category.color_hex as category_color_hex",
			).
			From("entry").
			LeftJoin("category on entry.category_id = category.id").
			Where(eq).
			Where(between).
			Limit(uint64(p.Limit)).
			Offset(uint64(p.Offset))

		if p.OrderBy == Asc {
			builder = builder.OrderBy("entry.date ASC")
		} else {
			builder = builder.OrderBy("entry.date DESC")
		}

		q, args, err := builder.ToSql()
		if err != nil {
			entriesChan <- result[[]models.Entry]{
				err: fmt.Errorf("ListEntriesByDate: query builder failed: %w", err),
			}
			return
		}

		r.logger.Infow("Executing query", "query", q, "args", args)

		rows, err := r.concurrentDB.QueryxContext(ctx, q, args...)
		if err != nil {
			entriesChan <- result[[]models.Entry]{
				err: fmt.Errorf("ListEntriesByDate: query failed: %w", err),
			}
			return
		}

		entries := make([]models.Entry, 0, p.Limit)

		type dest struct {
			Date             string           `db:"date"`
			CreatedAt        string           `db:"created_at"`
			UpdatedAt        string           `db:"updated_at"`
			Description      *string          `db:"description"`
			CategoryID       *string          `db:"category_id"`
			CategoryName     *string          `db:"category_name"`
			CategoryIcon     *string          `db:"category_icon"`
			CategoryColorHex *string          `db:"category_color_hex"`
			ID               string           `db:"id"`
			Name             string           `db:"name"`
			AccountID        string           `db:"account_id"`
			Amount           int64            `db:"amount"`
			EntryType        models.EntryType `db:"entry_type"`
		}
		for rows.Next() {
			var d dest
			err := rows.StructScan(&d)
			if err != nil {
				entriesChan <- result[[]models.Entry]{
					err: fmt.Errorf("ListEntriesByDate: failed to scan row: %w", err),
				}
				return
			}

			entry := models.Entry{
				Date:        d.Date,
				CreatedAt:   d.CreatedAt,
				UpdatedAt:   d.UpdatedAt,
				Description: d.Description,
				ID:          d.ID,
				Name:        d.Name,
				EntryType:   d.EntryType,
				Amount:      d.Amount,
			}

			if d.CategoryID != nil {
				entry.Category = &models.Category{
					ID:       *d.CategoryID,
					Name:     *d.CategoryName,
					Icon:     *d.CategoryIcon,
					ColorHex: *d.CategoryColorHex,
				}
			}

			entries = append(entries, entry)
		}

		entriesChan <- result[[]models.Entry]{
			ok: entries,
		}
	}()

	countResult := <-countChan
	if countResult.err != nil {
		return models.EntriesWithCount{}, countResult.err
	}

	entriesResult := <-entriesChan
	if entriesResult.err != nil {
		return models.EntriesWithCount{}, entriesResult.err
	}

	return models.EntriesWithCount{
		Entries:    entriesResult.ok,
		TotalCount: countResult.ok,
	}, nil
}

type CreateEntryParams struct {
	Date        string
	Description *string
	CategoryID  *string
	Name        string
	AccountID   string
	Amount      int
	EntryType   models.EntryType
}

func (r *Repository) CreateEntry(ctx context.Context, p CreateEntryParams) (models.Entry, error) {
	q, args, err := sq.
		Insert("entry").
		Columns("entry_type", "name", "amount", "description", "date", "category_id", "account_id").
		Values(p.EntryType, p.Name, p.Amount, p.Description, p.Date, p.CategoryID, p.AccountID).
		Suffix(`RETURNING *`).
		ToSql()
	if err != nil {
		return models.Entry{}, fmt.Errorf("CreateEntry: query builder failed: %w", err)
	}

	var dest struct {
		Date        string           `db:"date"`
		CreatedAt   string           `db:"created_at"`
		UpdatedAt   string           `db:"updated_at"`
		Description *string          `db:"description"`
		CategoryID  *string          `db:"category_id"`
		ID          string           `db:"id"`
		Name        string           `db:"name"`
		AccountID   string           `db:"account_id"`
		Amount      int64            `db:"amount"`
		EntryType   models.EntryType `db:"entry_type"`
	}
	err = r.NonConcurrentDB().GetContext(ctx, &dest, q, args...)
	if err != nil {
		return models.Entry{}, fmt.Errorf("CreateEntry: query failed: %w", err)
	}

	var category *models.Category
	if dest.CategoryID != nil {
		c, err := r.GetCategoryByID(ctx, GetCategoryByIDParams{
			ID:        *dest.CategoryID,
			AccountID: p.AccountID,
		})
		if err != nil {
			return models.Entry{}, fmt.Errorf("CreateEntry: %w", err)
		}
		category = &c
	}

	return models.Entry{
		Date:        dest.Date,
		CreatedAt:   dest.CreatedAt,
		UpdatedAt:   dest.UpdatedAt,
		Description: dest.Description,
		Category:    category,
		ID:          dest.ID,
		Name:        dest.Name,
		EntryType:   dest.EntryType,
		Amount:      dest.Amount,
	}, nil
}