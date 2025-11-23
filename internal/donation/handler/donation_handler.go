package handler

import (
	"net/http"
	"strings"

	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/domain/dto"
	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/donation/service"
	"github.com/gin-gonic/gin"
)

// DonationHandler handles donation-related HTTP requests
type DonationHandler struct {
	donationService service.DonationService
}

// NewDonationHandler creates a new instance of DonationHandler
func NewDonationHandler(donationService service.DonationService) *DonationHandler {
	return &DonationHandler{donationService: donationService}
}

// CreateDonation handles the create donation endpoint
// @Summary Create new donation
// @Tags donations
// @Accept json
// @Produce json
// @Param request body dto.CreateDonationRequest true "Create Donation Request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/donations [post]
func (h *DonationHandler) CreateDonation(c *gin.Context) {
	// Extract userID from context (set by authentication middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Authentication required",
			"message": "User must be authenticated to create donation",
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

	var request dto.CreateDonationRequest

	// Bind JSON request body
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	// Additional validation
	if strings.TrimSpace(request.ProgramID) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": "Program ID is required",
		})
		return
	}

	if request.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": "Amount must be greater than 0",
		})
		return
	}

	if strings.TrimSpace(request.PaymentMethod) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": "Payment method is required",
		})
		return
	}

	// Call service to create donation
	donation, err := h.donationService.CreateDonation(userIDStr, request)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "program not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "invalid program ID format" ||
			err.Error() == "amount must be greater than 0" {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, gin.H{
			"error":   "Donation creation failed",
			"message": err.Error(),
		})
		return
	}

	// Convert to DTO response
	donationResponse := dto.ToDonationResponse(donation)

	// Return success response
	c.JSON(http.StatusCreated, gin.H{
		"message": "Donation created successfully",
		"data":    donationResponse,
	})
}

// GetMyDonations handles the get my donations endpoint
// @Summary Get donations of authenticated user
// @Tags donations
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/users/my-donations [get]
func (h *DonationHandler) GetMyDonations(c *gin.Context) {
	// Extract userID from context (set by authentication middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Authentication required",
			"message": "User must be authenticated to view donations",
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

	// Call service to get user's donations
	donations, err := h.donationService.GetMyDonations(userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve your donations",
			"message": err.Error(),
		})
		return
	}

	// Convert to DTO response list
	donationResponses := dto.ToDonationResponseList(donations)

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Your donations retrieved successfully",
		"data":    donationResponses,
	})
}

// GetAllDonations handles the get all donations endpoint
// @Summary Get all donations (Admin only)
// @Tags donations
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/donations [get]
func (h *DonationHandler) GetAllDonations(c *gin.Context) {
	// Call service to get all donations
	donations, err := h.donationService.GetAllDonations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve donations",
			"message": err.Error(),
		})
		return
	}

	// Convert to DTO response list
	donationResponses := dto.ToDonationResponseList(donations)

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Donations retrieved successfully",
		"data":    donationResponses,
	})
}

// GetDonationByID handles the get donation by ID endpoint
// @Summary Get donation detail
// @Tags donations
// @Produce json
// @Param id path string true "Donation ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/donations/{id} [get]
func (h *DonationHandler) GetDonationByID(c *gin.Context) {
	// Get donation ID from URL parameter
	id := c.Param("id")

	if strings.TrimSpace(id) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": "Donation ID is required",
		})
		return
	}

	// Call service to get donation by ID
	donation, err := h.donationService.GetDonationByID(id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "donation not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "invalid donation ID format" {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, gin.H{
			"error":   "Failed to retrieve donation",
			"message": err.Error(),
		})
		return
	}

	// Convert to DTO response
	donationResponse := dto.ToDonationResponse(donation)

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Donation retrieved successfully",
		"data":    donationResponse,
	})
}

// UpdateDonationStatus handles the update donation status endpoint
// @Summary Update donation status (Admin only)
// @Tags donations
// @Accept json
// @Produce json
// @Param id path string true "Donation ID"
// @Param request body dto.UpdateDonationStatusRequest true "Update Status Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/donations/{id}/status [put]
func (h *DonationHandler) UpdateDonationStatus(c *gin.Context) {
	// Get donation ID from URL parameter
	id := c.Param("id")

	if strings.TrimSpace(id) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": "Donation ID is required",
		})
		return
	}

	var request dto.UpdateDonationStatusRequest

	// Bind JSON request body
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	// Call service to update status
	donation, err := h.donationService.UpdateDonationStatus(id, request)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "donation not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "invalid donation ID format" ||
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
	donationResponse := dto.ToDonationResponse(donation)

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Donation status updated successfully",
		"data":    donationResponse,
	})
}
