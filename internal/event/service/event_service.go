package service

import (
	"database/sql"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/domain/dto"
	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/domain/entity"
	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/event/repository"
	"github.com/google/uuid"
)

// EventService defines the interface for event business logic
type EventService interface {
	CreateEvent(request dto.CreateEventRequest) (*entity.Event, error)
	GetAllEvents() ([]entity.Event, map[string]int, error)
	GetEventByID(id string) (*entity.Event, int, error)
	UpdateEvent(id string, request dto.UpdateEventRequest) (*entity.Event, error)
	DeleteEvent(id string) error

	// Participation methods
	JoinEvent(userID, eventID string) error
	GetMyEvents(userID string) ([]entity.Event, map[string]int, error)
}

// eventService is the concrete implementation of EventService
type eventService struct {
	eventRepo repository.EventRepository
}

// NewEventService creates a new instance of EventService
func NewEventService(eventRepo repository.EventRepository) EventService {
	return &eventService{eventRepo: eventRepo}
}

// GenerateSlug converts a title to a URL-friendly slug
func GenerateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	reg := regexp.MustCompile("[^a-z0-9-]+")
	slug = reg.ReplaceAllString(slug, "")
	reg = regexp.MustCompile("-+")
	slug = reg.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	return slug
}

// CreateEvent handles the business logic for creating a new event
func (s *eventService) CreateEvent(request dto.CreateEventRequest) (*entity.Event, error) {
	// Validate required fields
	if strings.TrimSpace(request.Title) == "" {
		return nil, errors.New("title cannot be empty")
	}

	if strings.TrimSpace(request.Description) == "" {
		return nil, errors.New("description cannot be empty")
	}

	if strings.TrimSpace(request.Location) == "" {
		return nil, errors.New("location cannot be empty")
	}

	if request.Quota <= 0 {
		return nil, errors.New("quota must be greater than 0")
	}

	// Validate event date is in the future
	if request.EventDate.Before(time.Now()) {
		return nil, errors.New("event date must be in the future")
	}

	// Generate slug from title
	slug := GenerateSlug(request.Title)
	if slug == "" {
		return nil, errors.New("unable to generate valid slug from title")
	}

	// Check if slug already exists
	existingEvent, err := s.eventRepo.GetEventBySlug(slug)
	if err != nil && err != sql.ErrNoRows {
		return nil, errors.New("failed to check slug uniqueness")
	}

	// If slug exists, append UUID suffix
	if existingEvent != nil {
		slug = slug + "-" + uuid.NewString()[:8]
	}

	// Create new event entity
	newEvent := &entity.Event{
		ID:          uuid.NewString(),
		Title:       request.Title,
		Slug:        slug,
		Description: request.Description,
		ImageURL:    request.ImageURL,
		EventDate:   request.EventDate,
		Location:    request.Location,
		Quota:       request.Quota,
	}

	// Save to database
	if err := s.eventRepo.CreateEvent(newEvent); err != nil {
		return nil, errors.New("failed to create event")
	}

	return newEvent, nil
}

// GetAllEvents retrieves all events with participant counts
func (s *eventService) GetAllEvents() ([]entity.Event, map[string]int, error) {
	events, err := s.eventRepo.GetAllEvents()
	if err != nil {
		return nil, nil, errors.New("failed to retrieve events")
	}

	// Get participant counts for all events
	participantCounts, err := s.eventRepo.GetAllParticipantCounts()
	if err != nil {
		return nil, nil, errors.New("failed to retrieve participant counts")
	}

	return events, participantCounts, nil
}

// GetEventByID retrieves an event by its ID with participant count
func (s *eventService) GetEventByID(id string) (*entity.Event, int, error) {
	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return nil, 0, errors.New("invalid event ID format")
	}

	event, err := s.eventRepo.GetEventByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, 0, errors.New("event not found")
		}
		return nil, 0, errors.New("failed to retrieve event")
	}

	// Get participant count
	count, err := s.eventRepo.GetParticipantCount(id)
	if err != nil {
		return nil, 0, errors.New("failed to retrieve participant count")
	}

	return event, count, nil
}

// UpdateEvent handles the business logic for updating an event
func (s *eventService) UpdateEvent(id string, request dto.UpdateEventRequest) (*entity.Event, error) {
	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return nil, errors.New("invalid event ID format")
	}

	// Validate required fields
	if strings.TrimSpace(request.Title) == "" {
		return nil, errors.New("title cannot be empty")
	}

	if strings.TrimSpace(request.Description) == "" {
		return nil, errors.New("description cannot be empty")
	}

	if strings.TrimSpace(request.Location) == "" {
		return nil, errors.New("location cannot be empty")
	}

	if request.Quota <= 0 {
		return nil, errors.New("quota must be greater than 0")
	}

	// Check if event exists
	existingEvent, err := s.eventRepo.GetEventByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("event not found")
		}
		return nil, errors.New("failed to retrieve event")
	}

	// Check if new quota is less than current participants
	currentParticipants, err := s.eventRepo.GetParticipantCount(id)
	if err != nil {
		return nil, errors.New("failed to check current participants")
	}

	if request.Quota < currentParticipants {
		return nil, errors.New("quota cannot be less than current number of participants")
	}

	// Generate new slug if title changed
	newSlug := existingEvent.Slug
	if request.Title != existingEvent.Title {
		newSlug = GenerateSlug(request.Title)
		if newSlug == "" {
			return nil, errors.New("unable to generate valid slug from title")
		}

		// Check if new slug already exists (and it's not the current event)
		if newSlug != existingEvent.Slug {
			checkEvent, err := s.eventRepo.GetEventBySlug(newSlug)
			if err != nil && err != sql.ErrNoRows {
				return nil, errors.New("failed to check slug uniqueness")
			}
			if checkEvent != nil {
				newSlug = newSlug + "-" + uuid.NewString()[:8]
			}
		}
	}

	// Update event fields
	existingEvent.Title = request.Title
	existingEvent.Slug = newSlug
	existingEvent.Description = request.Description
	existingEvent.ImageURL = request.ImageURL
	existingEvent.EventDate = request.EventDate
	existingEvent.Location = request.Location
	existingEvent.Quota = request.Quota

	// Save updates to database
	if err := s.eventRepo.UpdateEvent(existingEvent); err != nil {
		return nil, errors.New("failed to update event")
	}

	return existingEvent, nil
}

// DeleteEvent handles the business logic for deleting an event
func (s *eventService) DeleteEvent(id string) error {
	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return errors.New("invalid event ID format")
	}

	// Check if event exists
	_, err := s.eventRepo.GetEventByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("event not found")
		}
		return errors.New("failed to retrieve event")
	}

	// Delete event (participants will be deleted via cascade in repository)
	if err := s.eventRepo.DeleteEvent(id); err != nil {
		return errors.New("failed to delete event")
	}

	return nil
}

// JoinEvent handles the business logic for a user joining an event
func (s *eventService) JoinEvent(userID, eventID string) error {
	// Validate UUID formats
	if _, err := uuid.Parse(userID); err != nil {
		return errors.New("invalid user ID format")
	}

	if _, err := uuid.Parse(eventID); err != nil {
		return errors.New("invalid event ID format")
	}

	// Step 1: Check if event exists
	event, err := s.eventRepo.GetEventByID(eventID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("event not found")
		}
		return errors.New("failed to retrieve event")
	}

	// Step 2: Check if event date has passed
	if event.EventDate.Before(time.Now()) {
		return errors.New("cannot join event that has already passed")
	}

	// Step 3: Check if user already joined (prevent double registration)
	alreadyJoined, err := s.eventRepo.IsUserJoined(userID, eventID)
	if err != nil {
		return errors.New("failed to check user participation status")
	}

	if alreadyJoined {
		return errors.New("you have already joined this event")
	}

	// Step 4: Check if quota is full
	currentParticipants, err := s.eventRepo.GetParticipantCount(eventID)
	if err != nil {
		return errors.New("failed to check event capacity")
	}

	if currentParticipants >= event.Quota {
		return errors.New("event is already full")
	}

	// Step 5: Insert to event_participants
	if err := s.eventRepo.JoinEvent(userID, eventID); err != nil {
		return errors.New("failed to join event")
	}

	return nil
}

// GetMyEvents retrieves all events that a user has joined
func (s *eventService) GetMyEvents(userID string) ([]entity.Event, map[string]int, error) {
	// Validate UUID format
	if _, err := uuid.Parse(userID); err != nil {
		return nil, nil, errors.New("invalid user ID format")
	}

	events, err := s.eventRepo.GetEventsByUserID(userID)
	if err != nil {
		return nil, nil, errors.New("failed to retrieve your events")
	}

	// Get participant counts for all events
	participantCounts, err := s.eventRepo.GetAllParticipantCounts()
	if err != nil {
		return nil, nil, errors.New("failed to retrieve participant counts")
	}

	return events, participantCounts, nil
}
