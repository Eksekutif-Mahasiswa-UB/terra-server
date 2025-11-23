package handler

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/auth/service"
	"github.com/Eksekutif-Mahasiswa-UB/terra-server/pkg/google"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

// OAuth2Handler handles OAuth2-related HTTP requests
type OAuth2Handler struct {
	authService  service.AuthService
	oauth2Config *oauth2.Config
}

// NewOAuth2Handler creates a new instance of OAuth2Handler
func NewOAuth2Handler(authService service.AuthService, oauth2Config *oauth2.Config) *OAuth2Handler {
	return &OAuth2Handler{
		authService:  authService,
		oauth2Config: oauth2Config,
	}
}

// generateState generates a random state string for CSRF protection
func generateState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// HandleGoogleLogin initiates the Google OAuth2 authorization flow
// @Summary Initiate Google OAuth2 login
// @Description Redirects user to Google's authorization page for authentication
// @Tags auth
// @Produce html
// @Success 302 {string} string "Redirect to Google authorization page"
// @Router /auth/google/login [get]
func (h *OAuth2Handler) HandleGoogleLogin(c *gin.Context) {
	// Generate random state for CSRF protection
	state := generateState()

	// Store state in session/cookie for validation in callback
	// For now, we'll pass it directly (in production, store in session)
	c.SetCookie(
		"oauth_state",
		state,
		3600, // 1 hour
		"/",
		"",
		false, // Set to true in production with HTTPS
		true,  // HttpOnly
	)

	// Generate authorization URL
	authURL := google.GetAuthURL(h.oauth2Config, state)

	// Redirect user to Google's authorization page
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// HandleGoogleCallback handles the OAuth2 callback from Google
// @Summary Handle Google OAuth2 callback
// @Description Exchanges authorization code for tokens, retrieves user info, and logs in the user
// @Tags auth
// @Produce json
// @Param code query string true "Authorization code from Google"
// @Param state query string true "State parameter for CSRF protection"
// @Success 200 {object} map[string]interface{} "Login successful with tokens"
// @Failure 400 {object} map[string]interface{} "Bad request or missing parameters"
// @Failure 401 {object} map[string]interface{} "Unauthorized or invalid token"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/google/callback [get]
func (h *OAuth2Handler) HandleGoogleCallback(c *gin.Context) {
	// Step 1: Get authorization code from query parameter
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Authorization failed",
			"message": "No authorization code received from Google",
		})
		return
	}

	// Step 2: Validate state parameter (CSRF protection)
	state := c.Query("state")
	storedState, err := c.Cookie("oauth_state")
	if err != nil || state != storedState {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "State validation failed",
			"message": "Invalid state parameter. Possible CSRF attack",
		})
		return
	}

	// Clear the state cookie
	c.SetCookie("oauth_state", "", -1, "/", "", false, true)

	// Step 3: Exchange authorization code for tokens
	token, err := google.ExchangeCode(h.oauth2Config, code)
	if err != nil {
		fmt.Println("Failed to exchange code:", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Token exchange failed",
			"message": "Failed to exchange authorization code for tokens",
		})
		return
	}

	// Step 4: Get user information from Google
	userInfo, err := google.GetUserInfo(token.AccessToken)
	if err != nil {
		fmt.Println("Failed to get user info:", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Failed to get user info",
			"message": "Could not retrieve user information from Google",
		})
		return
	}

	// Step 5: Validate email is verified
	if !userInfo.VerifiedEmail {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Email not verified",
			"message": "Your Google email is not verified. Please verify your email first",
		})
		return
	}

	// Step 6: Find or create user in database
	user, err := h.authService.FindOrCreateGoogleUser(userInfo.Email, userInfo.Name)
	if err != nil {
		// Handle specific error cases
		if err.Error() == "this email is already registered with email/password. Please use email login" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Account exists with different method",
				"message": err.Error(),
			})
			return
		}

		fmt.Println("Failed to find/create user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "User creation failed",
			"message": "Failed to create or retrieve user account",
		})
		return
	}

	// Step 7: Generate application tokens (access + refresh)
	loginResponse, err := h.authService.GenerateAppTokens(user.ID, user.Role)
	if err != nil {
		fmt.Println("Failed to generate tokens:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Token generation failed",
			"message": "Failed to generate authentication tokens",
		})
		return
	}

	// Step 8: Return success response with tokens
	c.JSON(http.StatusOK, gin.H{
		"message": "Google login successful",
		"data":    loginResponse,
		"user": gin.H{
			"id":       user.ID,
			"email":    user.Email,
			"fullName": user.FullName,
			"role":     user.Role,
		},
	})
}
