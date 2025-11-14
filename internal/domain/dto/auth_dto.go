package dto

// RegisterRequest represents the data transfer object for user registration
type RegisterRequest struct {
	FullName string `json:"full_name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginRequest represents the data transfer object for user login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents the data transfer object for login response with tokens
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// GoogleLoginRequest represents the data transfer object for Google login
type GoogleLoginRequest struct {
	Credential string `json:"credential" validate:"required"`
}

// ForgotPasswordRequest represents the data transfer object for forgot password
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ResetPasswordRequest represents the data transfer object for reset password
type ResetPasswordRequest struct {
	Token           string `json:"token" validate:"required"`
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

// RefreshTokenRequest represents the data transfer object for refresh token operations
type RefreshTokenRequest struct {
	Token string `json:"refresh_token" validate:"required"`
}

// RefreshAccessTokenResponse represents the response for refresh token endpoint
type RefreshAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}
