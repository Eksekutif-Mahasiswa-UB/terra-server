package middleware

import (
	"net/http"
	"strings"

	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/config"
	"github.com/Eksekutif-Mahasiswa-UB/terra-server/pkg/jwt"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware checks for a valid JWT and sets the UserID and Role in the context.
// This middleware validates the Bearer token and extracts user information.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Authorization required",
				"message": "Missing Authorization header",
			})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Authentication error",
				"message": "Invalid token format. Expected 'Bearer <token>'",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate the token
		claims, err := jwt.ValidateToken(tokenString, config.AppConfig.JWTSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Invalid or expired access token",
			})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// AdminMiddleware checks if the authenticated user has admin role.
// This middleware must be chained after AuthMiddleware.
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get role from context (set by AuthMiddleware)
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Authentication required",
				"message": "User authentication information not found",
			})
			c.Abort()
			return
		}

		roleStr, ok := role.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Server error",
				"message": "Invalid role format",
			})
			c.Abort()
			return
		}

		// Check if user has admin role
		if roleStr != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Access denied",
				"message": "This endpoint requires admin privileges",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
