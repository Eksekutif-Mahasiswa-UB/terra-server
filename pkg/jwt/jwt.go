package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims represents the JWT claims structure
type Claims struct {
	UserID  string `json:"user_id"`
	Role    string `json:"role"`
	Email   string `json:"email,omitempty"`
	Purpose string `json:"purpose,omitempty"`
	jwt.RegisteredClaims
}

// GenerateTokens generates both access token and refresh token
// Access token expires in 15 minutes, refresh token expires in 7 days
func GenerateTokens(userID string, role string, jwtSecret string) (accessToken string, refreshToken string, err error) {
	// Generate Access Token (15 minutes)
	accessClaims := Claims{
		UserID:  userID,
		Role:    role,
		Purpose: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.NewString(), // jti claim
		},
	}

	accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = accessTokenObj.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", "", err
	}

	// Generate Refresh Token (7 days)
	refreshClaims := Claims{
		UserID:  userID,
		Role:    role,
		Purpose: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.NewString(), // jti claim (unique token ID)
		},
	}

	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = refreshTokenObj.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// ValidateToken validates a JWT token and returns the claims
func ValidateToken(tokenString string, jwtSecret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// GetTokenExpiry extracts the expiration time from a JWT token
func GetTokenExpiry(tokenString string, jwtSecret string) (time.Time, error) {
	claims, err := ValidateToken(tokenString, jwtSecret)
	if err != nil {
		return time.Time{}, err
	}

	return claims.ExpiresAt.Time, nil
}

// GenerateAccessToken generates only an access token (15 minutes expiry)
func GenerateAccessToken(userID string, role string, jwtSecret string) (string, error) {
	claims := Claims{
		UserID:  userID,
		Role:    role,
		Purpose: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.NewString(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

// GenerateResetToken generates a reset password token with 15 minutes expiry
func GenerateResetToken(email string, jwtSecret string) (string, error) {
	claims := Claims{
		Email:   email,
		Purpose: "reset_password",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.NewString(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}
