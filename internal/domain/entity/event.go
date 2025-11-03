package entity

import "time"

// Event represents an event entity in the system
type Event struct {
	ID          string    `db:"id"`
	Title       string    `db:"title"`
	Slug        string    `db:"slug"`
	Description string    `db:"description"`
	ImageURL    *string   `db:"image_url"`
	EventDate   time.Time `db:"event_date"`
	Location    string    `db:"location"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
