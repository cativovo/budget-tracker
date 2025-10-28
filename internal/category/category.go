package category

import "time"

// Category represents the category of expense
type Category struct {
	ID        string    `json:"id" doc:"ID of category" example:"123"`
	Name      string    `json:"name" doc:"Name of category" example:"food"`
	Color     string    `json:"color" doc:"Color of category in hex" example:"#696969"`
	Icon      string    `json:"icon" doc:"Icon of category" example:"food-icon"`
	CreatedAt time.Time `json:"created_at" doc:"Date when the category has been created"`
	UpdatedAt time.Time `json:"updated_at" doc:"Date when the category has been updated"`
}
