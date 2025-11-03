package dto

import "time"

// EventResponse represents the event data transfer object for API responses
type EventResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	ImageURL    *string   `json:"image_url,omitempty"`
	EventDate   time.Time `json:"event_date"`
	Location    string    `json:"location"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
