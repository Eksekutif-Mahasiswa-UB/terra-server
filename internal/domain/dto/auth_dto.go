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
