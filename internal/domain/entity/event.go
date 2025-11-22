package entity

import "time"

// Event represents an event entity in the system
type Event struct {
	ID          string    `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Slug        string    `db:"slug" json:"slug"`
	Description string    `db:"description" json:"description"`
	ImageURL    string    `db:"image_url" json:"image_url"`
	EventDate   time.Time `db:"event_date" json:"event_date"`
	Location    string    `db:"location" json:"location"`
	Quota       int       `db:"quota" json:"quota"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// EventParticipant represents a user's participation in an event
type EventParticipant struct {
	ID       string    `db:"id" json:"id"`
	UserID   string    `db:"user_id" json:"user_id"`
	EventID  string    `db:"event_id" json:"event_id"`
	JoinedAt time.Time `db:"joined_at" json:"joined_at"`
}
