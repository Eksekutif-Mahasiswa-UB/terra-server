package handler

import (
	"net/http"
	"strings"

	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/auth/service"
	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/domain/dto"
	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler creates a new instance of AuthHandler
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register handles the user registration endpoint
// @Summary Register new user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Register Request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var request dto.RegisterRequest

	// Bind JSON request body
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	// Additional validation
	if strings.TrimSpace(request.Email) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": "Email cannot be empty",
		})
		return
	}

	if strings.TrimSpace(request.Password) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": "The password cannot be empty",
		})
		return
	}

	if strings.TrimSpace(request.FullName) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": "Full name cannot be empty",
		})
		return
	}

	// Call service to register user
	user, err := h.authService.Register(request)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "email sudah terdaftar" {
			statusCode = http.StatusConflict
		}
		c.JSON(statusCode, gin.H{
			"error":   "Registration failed",
			"message": err.Error(),
		})
		return
	}

	// Convert to DTO response
	userResponse := dto.ToUserResponse(user)

	// Return success response
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"data":    userResponse,
	})
}

// Login handles the user login endpoint
// @Summary User login
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var request dto.LoginRequest

	// Bind JSON request body
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	// Call service to login user
	loginResponse, err := h.authService.Login(request)
	if err != nil {
		// Handle different error types
		statusCode := http.StatusInternalServerError

		// Google account error - return 400 Bad Request
		if err.Error() == "this account is registered with Google. Please use Google login" {
			statusCode = http.StatusBadRequest
		} else if err.Error() == "email or password is incorrect" {
			// Invalid credentials - return 401 Unauthorized
			statusCode = http.StatusUnauthorized
		} else if err.Error() == "invalid authentication method" {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, gin.H{
			"error":   "Login failed",
			"message": err.Error(),
		})
		return
	}

	// Return success response with tokens
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"data":    loginResponse,
	})
}
