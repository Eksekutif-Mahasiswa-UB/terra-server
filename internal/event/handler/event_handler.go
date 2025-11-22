package handler

import (
	"net/http"
	"strings"

	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/domain/dto"
	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/event/service"
	"github.com/gin-gonic/gin"
)

// EventHandler handles event-related HTTP requests
type EventHandler struct {
	eventService service.EventService
}

// NewEventHandler creates a new instance of EventHandler
func NewEventHandler(eventService service.EventService) *EventHandler {
	return &EventHandler{eventService: eventService}
}

// CreateEvent handles the create event endpoint
// @Summary Create new event
// @Tags events
// @Accept json
// @Produce json
// @Param request body dto.CreateEventRequest true "Create Event Request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/events [post]
func (h *EventHandler) CreateEvent(c *gin.Context) {
	var request dto.CreateEventRequest

	// Bind JSON request body
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	// Additional validation
	if strings.TrimSpace(request.Title) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": "Title cannot be empty",
		})
		return
	}

	if strings.TrimSpace(request.Description) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": "Description cannot be empty",
		})
		return
	}

	if strings.TrimSpace(request.Location) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": "Location cannot be empty",
		})
		return
	}

	if request.Quota <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": "Quota must be greater than 0",
		})
		return
	}

	// Call service to create event
	event, err := h.eventService.CreateEvent(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Event creation failed",
			"message": err.Error(),
		})
		return
	}

	// Convert to DTO response
	eventResponse := dto.ToEventResponse(event, 0)

	// Return success response
	c.JSON(http.StatusCreated, gin.H{
		"message": "Event created successfully",
		"data":    eventResponse,
	})
}

// GetAllEvents handles the get all events endpoint
// @Summary Get all events
// @Tags events
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/events [get]
func (h *EventHandler) GetAllEvents(c *gin.Context) {
	// Call service to get all events
	events, participantCounts, err := h.eventService.GetAllEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve events",
			"message": err.Error(),
		})
		return
	}

	// Convert to DTO response list
	eventResponses := dto.ToEventResponseList(events, participantCounts)

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Events retrieved successfully",
		"data":    eventResponses,
	})
}

// GetEventByID handles the get event by ID endpoint
// @Summary Get event by ID
// @Tags events
// @Produce json
// @Param id path string true "Event ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/events/{id} [get]
func (h *EventHandler) GetEventByID(c *gin.Context) {
	// Get event ID from URL parameter
	id := c.Param("id")

	if strings.TrimSpace(id) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": "Event ID is required",
		})
		return
	}

	// Call service to get event by ID
	event, participantCount, err := h.eventService.GetEventByID(id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "event not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "invalid event ID format" {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, gin.H{
			"error":   "Failed to retrieve event",
			"message": err.Error(),
		})
		return
	}

	// Convert to DTO response
	eventResponse := dto.ToEventResponse(event, participantCount)

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Event retrieved successfully",
		"data":    eventResponse,
	})
}

// UpdateEvent handles the update event endpoint
// @Summary Update event
// @Tags events
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Param request body dto.UpdateEventRequest true "Update Event Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/events/{id} [put]
func (h *EventHandler) UpdateEvent(c *gin.Context) {
	// Get event ID from URL parameter
	id := c.Param("id")

	if strings.TrimSpace(id) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": "Event ID is required",
		})
		return
	}

	var request dto.UpdateEventRequest

	// Bind JSON request body
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	// Additional validation
	if strings.TrimSpace(request.Title) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": "Title cannot be empty",
		})
		return
	}

	if strings.TrimSpace(request.Description) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": "Description cannot be empty",
		})
		return
	}

	if strings.TrimSpace(request.Location) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": "Location cannot be empty",
		})
		return
	}

	if request.Quota <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": "Quota must be greater than 0",
		})
		return
	}

	// Call service to update event
	event, err := h.eventService.UpdateEvent(id, request)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "event not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "invalid event ID format" ||
			err.Error() == "quota cannot be less than current number of participants" {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, gin.H{
			"error":   "Event update failed",
			"message": err.Error(),
		})
		return
	}

	// Get updated participant count
	_, participantCount, err := h.eventService.GetEventByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve updated participant count",
			"message": err.Error(),
		})
		return
	}

	// Convert to DTO response
	eventResponse := dto.ToEventResponse(event, participantCount)

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Event updated successfully",
		"data":    eventResponse,
	})
}

// DeleteEvent handles the delete event endpoint
// @Summary Delete event
// @Tags events
// @Produce json
// @Param id path string true "Event ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/events/{id} [delete]
func (h *EventHandler) DeleteEvent(c *gin.Context) {
	// Get event ID from URL parameter
	id := c.Param("id")

	if strings.TrimSpace(id) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": "Event ID is required",
		})
		return
	}

	// Call service to delete event
	err := h.eventService.DeleteEvent(id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "event not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "invalid event ID format" {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, gin.H{
			"error":   "Event deletion failed",
			"message": err.Error(),
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Event deleted successfully",
	})
}

// JoinEvent handles the join event endpoint
// @Summary Join an event
// @Tags events
// @Produce json
// @Param id path string true "Event ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/events/{id}/join [post]
func (h *EventHandler) JoinEvent(c *gin.Context) {
	// Get event ID from URL parameter
	eventID := c.Param("id")

	if strings.TrimSpace(eventID) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": "Event ID is required",
		})
		return
	}

	// Extract userID from context (set by authentication middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Authentication required",
			"message": "User must be authenticated to join events",
		})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok || strings.TrimSpace(userIDStr) == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Authentication error",
			"message": "Invalid user authentication",
		})
		return
	}

	// Call service to join event
	err := h.eventService.JoinEvent(userIDStr, eventID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorMsg := err.Error()

		switch errorMsg {
		case "event not found":
			statusCode = http.StatusNotFound
		case "invalid event ID format", "invalid user ID format":
			statusCode = http.StatusBadRequest
		case "cannot join event that has already passed",
			"you have already joined this event",
			"event is already full":
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, gin.H{
			"error":   "Failed to join event",
			"message": errorMsg,
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully joined the event",
	})
}

// GetMyEvents handles the get my events endpoint
// @Summary Get events joined by the authenticated user
// @Tags events
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/users/my-events [get]
func (h *EventHandler) GetMyEvents(c *gin.Context) {
	// Extract userID from context (set by authentication middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Authentication required",
			"message": "User must be authenticated to view their events",
		})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok || strings.TrimSpace(userIDStr) == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Authentication error",
			"message": "Invalid user authentication",
		})
		return
	}

	// Call service to get user's events
	events, participantCounts, err := h.eventService.GetMyEvents(userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve your events",
			"message": err.Error(),
		})
		return
	}

	// Convert to DTO response list
	eventResponses := dto.ToEventResponseList(events, participantCounts)

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Your events retrieved successfully",
		"data":    eventResponses,
	})
}
