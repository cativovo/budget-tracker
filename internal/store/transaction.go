package store

import (
	"encoding/json"
	"fmt"

	"github.com/cativovo/budget-tracker/internal/constants"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

type TransactionRow struct {
	Date            pgtype.Date `json:"date"`
	Name            string      `json:"name"`
	Description     string      `json:"description"`
	Amount          float64     `json:"amount"`
	TransactionType int16       `json:"transaction_type"`
	ID              pgtype.UUID `json:"id"`
}

type TransactionByDateRow struct {
	Date          pgtype.Date
	TotalExpenses decimal.Decimal
	TotalIncome   decimal.Decimal
	Transactions  []TransactionRow
}

func ParseListTransactionsByDateRows(rows []ListTransactionsByDateRow) ([]TransactionByDateRow, error) {
	result := make([]TransactionByDateRow, 0, len(rows))

	for _, row := range rows {
		var transactions []TransactionRow
		if err := json.Unmarshal(row.Transactions, &transactions); err != nil {
			return nil, fmt.Errorf("listTransactions: unmarshal: %w - %s", err, string(row.Transactions))
		}

		var totalIncome decimal.Decimal
		var totalExpenses decimal.Decimal
		for _, t := range transactions {
			if t.TransactionType == constants.TransactionTypeIncome {
				totalIncome = totalIncome.Add(decimal.NewFromFloat(t.Amount))
			} else {
				totalExpenses = totalExpenses.Add(decimal.NewFromFloat(t.Amount))
			}
		}

		record := TransactionByDateRow{
			Date:          row.Date,
			TotalIncome:   totalIncome,
			TotalExpenses: totalExpenses,
			Transactions:  transactions,
		}

		result = append(result, record)
	}

	return result, nil
}
