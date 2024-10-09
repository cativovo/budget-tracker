package store

import (
	"github.com/jackc/pgx/v5/pgtype"
)

func NewUUID(id any) (pgtype.UUID, error) {
	var result pgtype.UUID
	err := result.Scan(id)
	return result, err
}

func NewDate(date any) (pgtype.Date, error) {
	var result pgtype.Date
	err := result.Scan(date)
	return result, err
}

func NewText(text any) (pgtype.Text, error) {
	var result pgtype.Text
	err := result.Scan(text)
	return result, err
}
