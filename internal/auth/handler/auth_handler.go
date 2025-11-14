package handler

import (
	"fmt"
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

// LoginWithGoogle handles the Google OAuth login endpoint
// @Summary Login with Google
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.GoogleLoginRequest true "Google Login Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/auth/login/google [post]
func (h *AuthHandler) LoginWithGoogle(c *gin.Context) {
	var request dto.GoogleLoginRequest

	// Bind JSON request body
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	// Validate credential is not empty
	if strings.TrimSpace(request.Credential) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": "Google credential cannot be empty",
		})
		return
	}

	// Call service to login with Google
	loginResponse, err := h.authService.LoginWithGoogle(request)
	if err != nil {
		// Handle different error types
		statusCode := http.StatusInternalServerError

		// Email account error - return 400 Bad Request
		if err.Error() == "please log in using email and password" {
			statusCode = http.StatusBadRequest
		} else if err.Error() == "invalid Google token" {
			// Invalid token - return 401 Unauthorized
			statusCode = http.StatusUnauthorized
		} else if err.Error() == "invalid authentication method" {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, gin.H{
			"error":   "Google login failed",
			"message": err.Error(),
		})
		return
	}

	// Return success response with tokens
	c.JSON(http.StatusOK, gin.H{
		"message": "Google login successful",
		"data":    loginResponse,
	})
}

// ForgotPassword handles the forgot password endpoint
// @Summary Request password reset
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.ForgotPasswordRequest true "Forgot Password Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var request dto.ForgotPasswordRequest

	// Bind JSON request body
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	// Validate email is not empty
	if strings.TrimSpace(request.Email) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": "Email cannot be empty",
		})
		return
	}

	// Call service to handle forgot password
	// Note: This always succeeds to prevent email enumeration
	err := h.authService.ForgotPassword(request)
	if err != nil {
		// Log the error internally but do not reveal to the client
		fmt.Println("Error handling forgot password:", err)
	}

	// Always return success message regardless of whether email exists
	c.JSON(http.StatusOK, gin.H{
		"message": "If your email is registered, you will receive a password reset link.",
	})
}

// RefreshToken handles the refresh token endpoint to get a new access token
// @Summary Refresh access token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh Token Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var request dto.RefreshTokenRequest

	// Bind JSON request body
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	// Validate token is not empty
	if strings.TrimSpace(request.Token) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": "Refresh token cannot be empty",
		})
		return
	}

	// Call service to refresh token
	response, err := h.authService.RefreshToken(request)
	if err != nil {
		// Handle different error types
		statusCode := http.StatusUnauthorized

		// Token validation errors - return 401 Unauthorized
		if err.Error() == "invalid or expired token" ||
			err.Error() == "invalid token: not a refresh token" ||
			err.Error() == "token is invalid or has been revoked" ||
			err.Error() == "refresh token has expired" {
			statusCode = http.StatusUnauthorized
		} else {
			statusCode = http.StatusInternalServerError
		}

		c.JSON(statusCode, gin.H{
			"error":   "Token refresh failed",
			"message": err.Error(),
		})
		return
	}

	// Return success response with new access token
	c.JSON(http.StatusOK, gin.H{
		"message": "Token refreshed successfully",
		"data":    response,
	})
}

// Logout handles the logout endpoint by invalidating the refresh token
// @Summary Logout user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh Token Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	var request dto.RefreshTokenRequest

	// Bind JSON request body
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	// Validate token is not empty
	if strings.TrimSpace(request.Token) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": "Refresh token cannot be empty",
		})
		return
	}

	// Call service to logout
	err := h.authService.Logout(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Logout failed",
			"message": err.Error(),
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}

// ResetPassword handles the reset password endpoint
// @Summary Reset password using token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.ResetPasswordRequest true "Reset Password Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var request dto.ResetPasswordRequest

	// Bind JSON request body
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	// Validate token is not empty
	if strings.TrimSpace(request.Token) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": "Token cannot be empty",
		})
		return
	}

	// Validate password is not empty
	if strings.TrimSpace(request.Password) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": "Password cannot be empty",
		})
		return
	}

	// Validate password confirmation
	if request.Password != request.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": "Passwords do not match",
		})
		return
	}

	// Validate password length
	if len(request.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": "Password must be at least 8 characters",
		})
		return
	}

	// Call service to reset password
	err := h.authService.ResetPassword(request)
	if err != nil {
		// Handle different error types
		statusCode := http.StatusInternalServerError

		// Token validation errors - return 401 Unauthorized
		if err.Error() == "invalid or expired token" || err.Error() == "invalid token purpose" {
			statusCode = http.StatusUnauthorized
		} else if err.Error() == "passwords do not match" || err.Error() == "invalid token: missing email" {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, gin.H{
			"error":   "Password reset failed",
			"message": err.Error(),
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Password has been reset successfully.",
	})
}
