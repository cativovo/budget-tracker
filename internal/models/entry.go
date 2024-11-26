package models

type EntryType int8

const (
	EntryTypeExpense EntryType = 0
	EntryTypeIncome  EntryType = 1
)

type Entry struct {
	Date        string
	CreatedAt   string
	UpdatedAt   string
	Description *string
	Category    *Category
	ID          string
	Name        string
	EntryType   EntryType
	Amount      int64
}

type EntriesWithCount struct {
	Entries    []Entry
	TotalCount int
}
