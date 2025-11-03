package service

import (
	"database/sql"
	"errors"

	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/auth/repository"
	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/domain/dto"
	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/domain/entity"
	"github.com/Eksekutif-Mahasiswa-UB/terra-server/pkg/hash"
	"github.com/google/uuid"
)

// AuthService defines the interface for authentication business logic
type AuthService interface {
	Register(request dto.RegisterRequest) (*entity.User, error)
	Login(request dto.LoginRequest) (*entity.User, error)
}

// authService is the concrete implementation of AuthService
type authService struct {
	authRepo repository.AuthRepository
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(authRepo repository.AuthRepository) AuthService {
	return &authService{authRepo: authRepo}
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
		GoogleID:   nil,
	}

	// Save user to database
	if err := s.authRepo.CreateUser(newUser); err != nil {
		return nil, err
	}

	// Clear password before returning
	newUser.Password = ""

	return newUser, nil
}

// Login handles the business logic for user login
func (s *authService) Login(request dto.LoginRequest) (*entity.User, error) {
	// Get user by email
	user, err := s.authRepo.GetUserByEmail(request.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("email or password is incorrect")
		}
		return nil, err
	}

	// Check password
	if !hash.CheckPasswordHash(request.Password, user.Password) {
		return nil, errors.New("email or password is incorrect")
	}

	// Clear password before returning
	user.Password = ""

	return user, nil
}
