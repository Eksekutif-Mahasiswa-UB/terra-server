package entity

import "time"

// Event represents an event entity in the system
type Event struct {
	ID          string    `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Slug        string    `db:"slug" json:"slug"`
	Description string    `db:"description" json:"description"`
	ImageURL    *string   `db:"image_url" json:"image_url,omitempty"`
	EventDate   time.Time `db:"event_date" json:"event_date"`
	Location    string    `db:"location" json:"location"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
