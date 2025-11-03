package dto

import "time"

// RefreshTokenResponse represents the refresh token data transfer object for API responses
type RefreshTokenResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}
