package store

import (
	"encoding/json"

	"github.com/cativovo/budget-tracker/internal/constants"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

type TransactionRow struct {
	Date            pgtype.Date `json:"date"`
	Name            string      `json:"name"`
	Description     string      `json:"description"`
	ID              string      `json:"id"`
	Amount          float64     `json:"amount"`
	TransactionType int16       `json:"transaction_type"`
}

type TransactionByDateRow struct {
	Date          pgtype.Date
	TotalExpenses decimal.Decimal
	TotalIncome   decimal.Decimal
	Transactions  []TransactionRow
}

type Foo struct {
	Date         pgtype.Date      `json:"date"`
	Transactions []TransactionRow `json:"transactions"`
}

type Bar struct {
	Result []Foo `json:"transactions"`
	Count  int   `json:"count"`
}

type ParseListTransactionsByDateRowsResult struct {
	Transactions []TransactionByDateRow
	Count        int
}

func ParseListTransactionsByDateRows(data []byte) (ParseListTransactionsByDateRowsResult, error) {
	var foo Bar

	if err := json.Unmarshal(data, &foo); err != nil {
		return ParseListTransactionsByDateRowsResult{}, err
	}

	transactions := make([]TransactionByDateRow, 0)

	for _, row := range foo.Result {
		var totalIncome decimal.Decimal
		var totalExpenses decimal.Decimal
		for _, t := range row.Transactions {
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
			Transactions:  row.Transactions,
		}

		transactions = append(transactions, record)
	}

	result := ParseListTransactionsByDateRowsResult{
		Count:        foo.Count,
		Transactions: transactions,
	}

	return result, nil
}
