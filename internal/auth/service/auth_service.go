package service

import (
	"database/sql"
	"errors"
	"time"

	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/auth/repository"
	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/config"
	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/domain/dto"
	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/domain/entity"
	"github.com/Eksekutif-Mahasiswa-UB/terra-server/pkg/google"
	"github.com/Eksekutif-Mahasiswa-UB/terra-server/pkg/hash"
	"github.com/Eksekutif-Mahasiswa-UB/terra-server/pkg/jwt"
	"github.com/google/uuid"
)

// AuthService defines the interface for authentication business logic
type AuthService interface {
	Register(request dto.RegisterRequest) (*entity.User, error)
	Login(request dto.LoginRequest) (*dto.LoginResponse, error)
	LoginWithGoogle(request dto.GoogleLoginRequest) (*dto.LoginResponse, error)
}

// authService is the concrete implementation of AuthService
type authService struct {
	authRepo       repository.AuthRepository
	googleClientID string
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(authRepo repository.AuthRepository, googleClientID string) AuthService {
	return &authService{
		authRepo:       authRepo,
		googleClientID: googleClientID,
	}
}

// Register handles the business logic for user registration
func (s *authService) Register(request dto.RegisterRequest) (*entity.User, error) {
	// Check if email already exists
	existingUser, err := s.authRepo.GetUserByEmail(request.Email)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email is registered")
	}

	// Hash password
	hashedPassword, err := hash.HashPassword(request.Password)
	if err != nil {
		return nil, err
	}

	// Create new user
	newUser := &entity.User{
		ID:         uuid.NewString(),
		FullName:   request.FullName,
		Email:      request.Email,
		Password:   hashedPassword,
		Role:       "user",
		AuthMethod: "email",
	}

	// Save user to database
	if err := s.authRepo.CreateUser(newUser); err != nil {
		return nil, err
	}

	// Clear password before returning
	newUser.Password = ""

	return newUser, nil
}

// Login handles the business logic for user login with refresh token generation
func (s *authService) Login(request dto.LoginRequest) (*dto.LoginResponse, error) {
	// Step 1: Get user by email
	user, err := s.authRepo.GetUserByEmail(request.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("email or password is incorrect")
		}
		return nil, err
	}

	// Step 2: Check auth_method - CRUCIAL CHECK
	if user.AuthMethod == "google" {
		return nil, errors.New("this account is registered with Google. Please use Google login")
	}

	// Step 3: Verify auth_method is email
	if user.AuthMethod != "email" {
		return nil, errors.New("invalid authentication method")
	}

	// Step 4: Check password
	if !hash.CheckPasswordHash(request.Password, user.Password) {
		return nil, errors.New("email or password is incorrect")
	}

	// Step 5: Generate tokens (access token 15m, refresh token 7d)
	accessToken, refreshToken, err := jwt.GenerateTokens(user.ID, user.Role, config.AppConfig.JWTSecret)
	if err != nil {
		return nil, errors.New("failed to generate tokens")
	}

	// Step 6: Get refresh token expiry (7 days from now)
	refreshTokenExpiry := time.Now().Add(7 * 24 * time.Hour)

	// Step 7: Save refresh token to database
	refreshTokenEntity := &entity.RefreshToken{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		Token:     refreshToken, // Store the JWT refresh token
		ExpiresAt: refreshTokenExpiry,
	}

	if err := s.authRepo.CreateRefreshToken(refreshTokenEntity); err != nil {
		return nil, errors.New("failed to save refresh token")
	}

	// Step 8: Return login response with both tokens
	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// LoginWithGoogle handles the business logic for Google OAuth login
func (s *authService) LoginWithGoogle(request dto.GoogleLoginRequest) (*dto.LoginResponse, error) {
	// Step 1: Verify Google token
	tokenInfo, err := google.VerifyGoogleToken(request.Credential, s.googleClientID)
	if err != nil {
		return nil, errors.New("invalid Google token")
	}

	// Extract email and name from token
	email := tokenInfo.Email
	fullName := tokenInfo.Email // Use email as fallback if name not available

	// Step 2: Check if user exists
	existingUser, err := s.authRepo.GetUserByEmail(email)

	if err != nil && err != sql.ErrNoRows {
		return nil, errors.New("failed to check user existence")
	}

	var user *entity.User

	// Step 3: User exists - check auth_method
	if existingUser != nil {
		// CRUCIAL CHECK: If user exists with email auth, reject
		if existingUser.AuthMethod == "email" {
			return nil, errors.New("please log in using email and password")
		}

		// If auth_method is google, proceed with login
		if existingUser.AuthMethod == "google" {
			user = existingUser
		} else {
			return nil, errors.New("invalid authentication method")
		}
	} else {
		// Step 4: New user - create account with Google auth
		newUser := &entity.User{
			ID:         uuid.NewString(),
			FullName:   fullName,
			Email:      email,
			Password:   "", // NULL password for Google users
			Role:       "user",
			AuthMethod: "google",
		}

		// Save new user to database
		if err := s.authRepo.CreateUser(newUser); err != nil {
			return nil, errors.New("failed to create user account")
		}

		user = newUser
	}

	// Step 5: Generate tokens (access token 15m, refresh token 7d)
	accessToken, refreshToken, err := jwt.GenerateTokens(user.ID, user.Role, config.AppConfig.JWTSecret)
	if err != nil {
		return nil, errors.New("failed to generate tokens")
	}

	// Step 6: Get refresh token expiry (7 days from now)
	refreshTokenExpiry := time.Now().Add(7 * 24 * time.Hour)

	// Step 7: Save refresh token to database
	refreshTokenEntity := &entity.RefreshToken{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: refreshTokenExpiry,
	}

	if err := s.authRepo.CreateRefreshToken(refreshTokenEntity); err != nil {
		return nil, errors.New("failed to save refresh token")
	}

	// Step 8: Return login response with both tokens
	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
