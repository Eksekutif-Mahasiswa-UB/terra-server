package entity

import "time"

// Program represents a program entity in the system
type Program struct {
	ID           string    `db:"id"`
	Title        string    `db:"title"`
	Description  string    `db:"description"`
	ImageURL     string    `db:"image_url"`
	TargetAmount float64   `db:"target_amount"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
