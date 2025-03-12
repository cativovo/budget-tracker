package category

import "time"

type Category struct {
	ID        string
	Name      string
	Color     string
	Icon      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
