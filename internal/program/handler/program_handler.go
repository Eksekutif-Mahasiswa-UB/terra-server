package handler

import (
	"net/http"
	"strings"

	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/domain/dto"
	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/program/service"
	"github.com/gin-gonic/gin"
)

// ProgramHandler handles program-related HTTP requests
type ProgramHandler struct {
	programService service.ProgramService
}

// NewProgramHandler creates a new instance of ProgramHandler
func NewProgramHandler(programService service.ProgramService) *ProgramHandler {
	return &ProgramHandler{programService: programService}
}

// CreateProgram handles the create program endpoint
// @Summary Create new program
// @Tags programs
// @Accept json
// @Produce json
// @Param request body dto.CreateProgramRequest true "Create Program Request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/programs [post]
func (h *ProgramHandler) CreateProgram(c *gin.Context) {
	var request dto.CreateProgramRequest

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

	if request.TargetAmount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": "Target amount must be greater than 0",
		})
		return
	}

	// Call service to create program
	program, err := h.programService.CreateProgram(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Program creation failed",
			"message": err.Error(),
		})
		return
	}

	// Convert to DTO response
	programResponse := dto.ToProgramResponse(program)

	// Return success response
	c.JSON(http.StatusCreated, gin.H{
		"message": "Program created successfully",
		"data":    programResponse,
	})
}

// GetAllPrograms handles the get all programs endpoint
// @Summary Get all programs
// @Tags programs
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/programs [get]
func (h *ProgramHandler) GetAllPrograms(c *gin.Context) {
	// Call service to get all programs
	programs, err := h.programService.GetAllPrograms()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve programs",
			"message": err.Error(),
		})
		return
	}

	// Convert to DTO response list
	programResponses := dto.ToProgramResponseList(programs)

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Programs retrieved successfully",
		"data":    programResponses,
	})
}

// GetProgramByID handles the get program by ID endpoint
// @Summary Get program by ID
// @Tags programs
// @Produce json
// @Param id path string true "Program ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/programs/{id} [get]
func (h *ProgramHandler) GetProgramByID(c *gin.Context) {
	// Get program ID from URL parameter
	id := c.Param("id")

	if strings.TrimSpace(id) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": "Program ID is required",
		})
		return
	}

	// Call service to get program by ID
	program, err := h.programService.GetProgramByID(id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "program not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "invalid program ID format" {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, gin.H{
			"error":   "Failed to retrieve program",
			"message": err.Error(),
		})
		return
	}

	// Convert to DTO response
	programResponse := dto.ToProgramResponse(program)

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Program retrieved successfully",
		"data":    programResponse,
	})
}

// UpdateProgram handles the update program endpoint
// @Summary Update program
// @Tags programs
// @Accept json
// @Produce json
// @Param id path string true "Program ID"
// @Param request body dto.UpdateProgramRequest true "Update Program Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/programs/{id} [put]
func (h *ProgramHandler) UpdateProgram(c *gin.Context) {
	// Get program ID from URL parameter
	id := c.Param("id")

	if strings.TrimSpace(id) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": "Program ID is required",
		})
		return
	}

	var request dto.UpdateProgramRequest

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

	if request.TargetAmount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": "Target amount must be greater than 0",
		})
		return
	}

	// Call service to update program
	program, err := h.programService.UpdateProgram(id, request)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "program not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "invalid program ID format" {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, gin.H{
			"error":   "Program update failed",
			"message": err.Error(),
		})
		return
	}

	// Convert to DTO response
	programResponse := dto.ToProgramResponse(program)

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Program updated successfully",
		"data":    programResponse,
	})
}

// DeleteProgram handles the delete program endpoint
// @Summary Delete program
// @Tags programs
// @Produce json
// @Param id path string true "Program ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/programs/{id} [delete]
func (h *ProgramHandler) DeleteProgram(c *gin.Context) {
	// Get program ID from URL parameter
	id := c.Param("id")

	if strings.TrimSpace(id) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": "Program ID is required",
		})
		return
	}

	// Call service to delete program
	err := h.programService.DeleteProgram(id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "program not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "invalid program ID format" {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, gin.H{
			"error":   "Program deletion failed",
			"message": err.Error(),
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Program deleted successfully",
	})
}
