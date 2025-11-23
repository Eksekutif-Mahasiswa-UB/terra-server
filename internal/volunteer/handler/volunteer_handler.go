package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/domain/dto"
	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/volunteer/service"
	"github.com/gin-gonic/gin"
)

// VolunteerHandler handles volunteer-related HTTP requests
type VolunteerHandler struct {
	volunteerService service.VolunteerService
}

// NewVolunteerHandler creates a new instance of VolunteerHandler
func NewVolunteerHandler(volunteerService service.VolunteerService) *VolunteerHandler {
	return &VolunteerHandler{volunteerService: volunteerService}
}

// SubmitApplication handles the submit volunteer application endpoint
// @Summary Submit volunteer application
// @Tags volunteers
// @Accept json
// @Produce json
// @Param request body dto.CreateVolunteerRequest true "Volunteer Application Request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/volunteers/apply [post]
func (h *VolunteerHandler) SubmitApplication(c *gin.Context) {
	// Extract userID from context (set by authentication middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Authentication required",
			"message": "User must be authenticated to submit volunteer application",
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

	var request dto.CreateVolunteerRequest

	// Bind JSON request body
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	// Additional validation
	if strings.TrimSpace(request.FullName) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": "Full name is required",
		})
		return
	}

	if strings.TrimSpace(request.Email) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": "Email is required",
		})
		return
	}

	if strings.TrimSpace(request.Phone) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": "Phone number is required",
		})
		return
	}

	if strings.TrimSpace(request.City) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": "City is required",
		})
		return
	}

	if strings.TrimSpace(request.Occupation) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": "Occupation is required",
		})
		return
	}

	if strings.TrimSpace(request.Interests) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": "Interests are required",
		})
		return
	}

	// Call service to submit application
	volunteer, err := h.volunteerService.SubmitApplication(userIDStr, request)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "you already have a pending volunteer application" {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, gin.H{
			"error":   "Application submission failed",
			"message": err.Error(),
		})
		return
	}

	// Convert to DTO response
	volunteerResponse := dto.ToVolunteerResponse(volunteer)

	// Return success response
	c.JSON(http.StatusCreated, gin.H{
		"message": "Volunteer application submitted successfully",
		"data":    volunteerResponse,
	})
}

// GetAllApplications handles the get all volunteer applications endpoint
// @Summary Get all volunteer applications (Admin only)
// @Tags volunteers
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param status query string false "Filter by status" Enums(pending, approved, rejected)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/volunteers [get]
func (h *VolunteerHandler) GetAllApplications(c *gin.Context) {
	// Get pagination parameters
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	status := c.Query("status")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	// Call service to get all applications
	volunteers, totalCount, err := h.volunteerService.GetAllApplications(page, limit, status)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "invalid status filter") {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, gin.H{
			"error":   "Failed to retrieve volunteer applications",
			"message": err.Error(),
		})
		return
	}

	// Convert to DTO response list
	volunteerResponses := dto.ToVolunteerResponseList(volunteers)

	// Calculate pagination metadata
	totalPages := (totalCount + limit - 1) / limit

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Volunteer applications retrieved successfully",
		"data":    volunteerResponses,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total_items": totalCount,
			"total_pages": totalPages,
		},
	})
}

// GetApplicationByID handles the get volunteer application by ID endpoint
// @Summary Get volunteer application detail (Admin only)
// @Tags volunteers
// @Produce json
// @Param id path string true "Application ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/volunteers/{id} [get]
func (h *VolunteerHandler) GetApplicationByID(c *gin.Context) {
	// Get application ID from URL parameter
	id := c.Param("id")

	if strings.TrimSpace(id) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": "Application ID is required",
		})
		return
	}

	// Call service to get application by ID
	volunteer, err := h.volunteerService.GetApplicationByID(id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "volunteer application not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "invalid application ID format" {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, gin.H{
			"error":   "Failed to retrieve volunteer application",
			"message": err.Error(),
		})
		return
	}

	// Convert to DTO response
	volunteerResponse := dto.ToVolunteerResponse(volunteer)

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Volunteer application retrieved successfully",
		"data":    volunteerResponse,
	})
}

// UpdateApplicationStatus handles the update volunteer application status endpoint
// @Summary Update volunteer application status (Admin only)
// @Tags volunteers
// @Accept json
// @Produce json
// @Param id path string true "Application ID"
// @Param request body dto.UpdateVolunteerStatusRequest true "Update Status Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/volunteers/{id}/status [put]
func (h *VolunteerHandler) UpdateApplicationStatus(c *gin.Context) {
	// Get application ID from URL parameter
	id := c.Param("id")

	if strings.TrimSpace(id) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": "Application ID is required",
		})
		return
	}

	var request dto.UpdateVolunteerStatusRequest

	// Bind JSON request body
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	// Call service to update status
	volunteer, err := h.volunteerService.UpdateApplicationStatus(id, request)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "volunteer application not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "invalid application ID format" ||
			strings.Contains(err.Error(), "invalid status") {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, gin.H{
			"error":   "Status update failed",
			"message": err.Error(),
		})
		return
	}

	// Convert to DTO response
	volunteerResponse := dto.ToVolunteerResponse(volunteer)

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Volunteer application status updated successfully",
		"data":    volunteerResponse,
	})
}
