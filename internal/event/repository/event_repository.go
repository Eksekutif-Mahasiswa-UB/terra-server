package repository

import (
	"database/sql"
	"time"

	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/domain/entity"
	"github.com/jmoiron/sqlx"
)

// EventRepository defines the interface for event data operations
type EventRepository interface {
	CreateEvent(event *entity.Event) error
	GetEventByID(id string) (*entity.Event, error)
	GetEventBySlug(slug string) (*entity.Event, error)
	GetAllEvents() ([]entity.Event, error)
	UpdateEvent(event *entity.Event) error
	DeleteEvent(id string) error

	// Participation methods
	JoinEvent(userID, eventID string) error
	IsUserJoined(userID, eventID string) (bool, error)
	GetParticipantCount(eventID string) (int, error)
	GetAllParticipantCounts() (map[string]int, error)
	GetEventsByUserID(userID string) ([]entity.Event, error)
}

// eventRepository is the concrete implementation of EventRepository
type eventRepository struct {
	db *sqlx.DB
}

// NewEventRepository creates a new instance of EventRepository
func NewEventRepository(db *sqlx.DB) EventRepository {
	return &eventRepository{db: db}
}

// CreateEvent inserts a new event into the database
func (r *eventRepository) CreateEvent(event *entity.Event) error {
	event.CreatedAt = time.Now()
	event.UpdatedAt = time.Now()

	query := `INSERT INTO events (id, title, slug, description, image_url, event_date, location, quota, created_at, updated_at) 
			  VALUES (:id, :title, :slug, :description, :image_url, :event_date, :location, :quota, :created_at, :updated_at)`

	_, err := r.db.NamedExec(query, event)
	return err
}

// GetEventByID retrieves an event by its ID
func (r *eventRepository) GetEventByID(id string) (*entity.Event, error) {
	var event entity.Event
	query := `SELECT * FROM events WHERE id = ?`

	err := r.db.Get(&event, query, id)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

// GetEventBySlug retrieves an event by its slug
func (r *eventRepository) GetEventBySlug(slug string) (*entity.Event, error) {
	var event entity.Event
	query := `SELECT * FROM events WHERE slug = ?`

	err := r.db.Get(&event, query, slug)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

// GetAllEvents retrieves all events from the database
func (r *eventRepository) GetAllEvents() ([]entity.Event, error) {
	var events []entity.Event
	query := `SELECT * FROM events ORDER BY event_date ASC`

	err := r.db.Select(&events, query)
	if err != nil {
		return nil, err
	}

	return events, nil
}

// UpdateEvent updates an existing event in the database
func (r *eventRepository) UpdateEvent(event *entity.Event) error {
	event.UpdatedAt = time.Now()

	query := `UPDATE events 
			  SET title = :title, slug = :slug, description = :description, image_url = :image_url, 
			      event_date = :event_date, location = :location, quota = :quota, updated_at = :updated_at 
			  WHERE id = :id`

	_, err := r.db.NamedExec(query, event)
	return err
}

// DeleteEvent deletes an event from the database
func (r *eventRepository) DeleteEvent(id string) error {
	// Delete participants first (cascade)
	_, err := r.db.Exec(`DELETE FROM event_participants WHERE event_id = ?`, id)
	if err != nil {
		return err
	}

	// Delete the event
	query := `DELETE FROM events WHERE id = ?`
	_, err = r.db.Exec(query, id)
	return err
}

// JoinEvent adds a user to an event's participants
func (r *eventRepository) JoinEvent(userID, eventID string) error {
	participant := map[string]interface{}{
		"id":        generateParticipantID(userID, eventID),
		"user_id":   userID,
		"event_id":  eventID,
		"joined_at": time.Now(),
	}

	query := `INSERT INTO event_participants (id, user_id, event_id, joined_at) 
			  VALUES (:id, :user_id, :event_id, :joined_at)`

	_, err := r.db.NamedExec(query, participant)
	return err
}

// IsUserJoined checks if a user has already joined an event
func (r *eventRepository) IsUserJoined(userID, eventID string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM event_participants WHERE user_id = ? AND event_id = ?`

	err := r.db.Get(&count, query, userID, eventID)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetParticipantCount returns the number of participants for a specific event
func (r *eventRepository) GetParticipantCount(eventID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM event_participants WHERE event_id = ?`

	err := r.db.Get(&count, query, eventID)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// GetAllParticipantCounts returns a map of event IDs to participant counts
func (r *eventRepository) GetAllParticipantCounts() (map[string]int, error) {
	type CountResult struct {
		EventID string `db:"event_id"`
		Count   int    `db:"count"`
	}

	var results []CountResult
	query := `SELECT event_id, COUNT(*) as count FROM event_participants GROUP BY event_id`

	err := r.db.Select(&results, query)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	counts := make(map[string]int)
	for _, result := range results {
		counts[result.EventID] = result.Count
	}

	return counts, nil
}

// GetEventsByUserID retrieves all events that a user has joined
func (r *eventRepository) GetEventsByUserID(userID string) ([]entity.Event, error) {
	var events []entity.Event
	query := `SELECT e.* FROM events e
			  INNER JOIN event_participants ep ON e.id = ep.event_id
			  WHERE ep.user_id = ?
			  ORDER BY e.event_date ASC`

	err := r.db.Select(&events, query, userID)
	if err != nil {
		return nil, err
	}

	return events, nil
}

// Helper function to generate a deterministic ID for participant records
func generateParticipantID(userID, eventID string) string {
	return userID + "-" + eventID
}
