package dto

import "time"

// CreateEventRequest represents the request to create a new event
type CreateEventRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
	ImageURL    string    `json:"image_url"`
	EventDate   time.Time `json:"event_date" binding:"required"`
	Location    string    `json:"location" binding:"required"`
	Quota       int       `json:"quota" binding:"required,gt=0"`
}

// UpdateEventRequest represents the request to update an event
type UpdateEventRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
	ImageURL    string    `json:"image_url"`
	EventDate   time.Time `json:"event_date" binding:"required"`
	Location    string    `json:"location" binding:"required"`
	Quota       int       `json:"quota" binding:"required,gt=0"`
}

// EventResponse represents the event data transfer object for API responses
type EventResponse struct {
	ID                  string    `json:"id"`
	Title               string    `json:"title"`
	Slug                string    `json:"slug"`
	Description         string    `json:"description"`
	ImageURL            string    `json:"image_url"`
	EventDate           time.Time `json:"event_date"`
	Location            string    `json:"location"`
	Quota               int       `json:"quota"`
	CurrentParticipants int       `json:"current_participants"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}
