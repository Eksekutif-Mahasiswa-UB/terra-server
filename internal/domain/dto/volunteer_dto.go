package dto

import (
	"time"

	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/domain/entity"
)

// CreateVolunteerRequest represents the request to submit a volunteer application
type CreateVolunteerRequest struct {
	FullName    string    `json:"full_name" binding:"required"`
	Email       string    `json:"email" binding:"required,email"`
	Phone       string    `json:"phone" binding:"required"`
	DateOfBirth time.Time `json:"date_of_birth" binding:"required"`
	Gender      string    `json:"gender" binding:"required,oneof=Male Female"`
	City        string    `json:"city" binding:"required"`
	Occupation  string    `json:"occupation" binding:"required"`
	Interests   string    `json:"interests" binding:"required"`
	Experience  string    `json:"experience"`
}

// UpdateVolunteerStatusRequest represents the request to update a volunteer application status
type UpdateVolunteerStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending approved rejected"`
}

// VolunteerResponse represents the volunteer data transfer object for API responses
type VolunteerResponse struct {
	ID          string                   `json:"id"`
	UserID      string                   `json:"user_id"`
	FullName    string                   `json:"full_name"`
	Email       string                   `json:"email"`
	Phone       string                   `json:"phone"`
	DateOfBirth time.Time                `json:"date_of_birth"`
	Gender      entity.Gender            `json:"gender"`
	City        string                   `json:"city"`
	Occupation  string                   `json:"occupation"`
	Interests   string                   `json:"interests"`
	Experience  string                   `json:"experience"`
	Status      entity.ApplicationStatus `json:"status"`
	CreatedAt   time.Time                `json:"created_at"`
	UpdatedAt   time.Time                `json:"updated_at"`
}
