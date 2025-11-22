package service

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/domain/dto"
	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/domain/entity"
	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/volunteer/repository"
	"github.com/google/uuid"
)

// VolunteerService defines the interface for volunteer business logic
type VolunteerService interface {
	SubmitApplication(userID string, request dto.CreateVolunteerRequest) (*entity.Volunteer, error)
	GetAllApplications(page, limit int, status string) ([]entity.Volunteer, int, error)
	GetApplicationByID(id string) (*entity.Volunteer, error)
	UpdateApplicationStatus(id string, request dto.UpdateVolunteerStatusRequest) (*entity.Volunteer, error)
}

// volunteerService is the concrete implementation of VolunteerService
type volunteerService struct {
	volunteerRepo repository.VolunteerRepository
}

// NewVolunteerService creates a new instance of VolunteerService
func NewVolunteerService(volunteerRepo repository.VolunteerRepository) VolunteerService {
	return &volunteerService{volunteerRepo: volunteerRepo}
}

// SubmitApplication handles the business logic for submitting a volunteer application
func (s *volunteerService) SubmitApplication(userID string, request dto.CreateVolunteerRequest) (*entity.Volunteer, error) {
	// Validate UUID format
	if _, err := uuid.Parse(userID); err != nil {
		return nil, errors.New("invalid user ID format")
	}

	// Validate required fields
	if strings.TrimSpace(request.FullName) == "" {
		return nil, errors.New("full name is required")
	}

	if strings.TrimSpace(request.Email) == "" {
		return nil, errors.New("email is required")
	}

	if strings.TrimSpace(request.Phone) == "" {
		return nil, errors.New("phone number is required")
	}

	if strings.TrimSpace(request.City) == "" {
		return nil, errors.New("city is required")
	}

	if strings.TrimSpace(request.Occupation) == "" {
		return nil, errors.New("occupation is required")
	}

	if strings.TrimSpace(request.Interests) == "" {
		return nil, errors.New("interests are required")
	}

	// Validate gender
	gender := entity.Gender(request.Gender)
	if gender != entity.GenderMale && gender != entity.GenderFemale {
		return nil, errors.New("gender must be 'Male' or 'Female'")
	}

	// Check if user already has a pending application
	hasPending, err := s.volunteerRepo.HasPendingApplication(userID)
	if err != nil {
		return nil, errors.New("failed to check existing applications")
	}

	if hasPending {
		return nil, errors.New("you already have a pending volunteer application")
	}

	// Create new volunteer application
	volunteer := &entity.Volunteer{
		ID:          uuid.NewString(),
		UserID:      userID,
		FullName:    request.FullName,
		Email:       request.Email,
		Phone:       request.Phone,
		DateOfBirth: request.DateOfBirth,
		Gender:      gender,
		City:        request.City,
		Occupation:  request.Occupation,
		Interests:   request.Interests,
		Experience:  request.Experience,
		Status:      entity.ApplicationStatusPending,
	}

	// Save to database
	if err := s.volunteerRepo.Create(volunteer); err != nil {
		return nil, errors.New("failed to submit volunteer application")
	}

	return volunteer, nil
}

// GetAllApplications retrieves all volunteer applications with pagination and filtering
func (s *volunteerService) GetAllApplications(page, limit int, status string) ([]entity.Volunteer, int, error) {
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	// Validate status filter if provided
	if status != "" {
		appStatus := entity.ApplicationStatus(status)
		if appStatus != entity.ApplicationStatusPending &&
			appStatus != entity.ApplicationStatusApproved &&
			appStatus != entity.ApplicationStatusRejected {
			return nil, 0, errors.New("invalid status filter. Must be 'pending', 'approved', or 'rejected'")
		}
	}

	volunteers, totalCount, err := s.volunteerRepo.GetAll(page, limit, status)
	if err != nil {
		return nil, 0, errors.New("failed to retrieve volunteer applications")
	}

	return volunteers, totalCount, nil
}

// GetApplicationByID retrieves a volunteer application by its ID
func (s *volunteerService) GetApplicationByID(id string) (*entity.Volunteer, error) {
	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return nil, errors.New("invalid application ID format")
	}

	volunteer, err := s.volunteerRepo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("volunteer application not found")
		}
		return nil, errors.New("failed to retrieve volunteer application")
	}

	return volunteer, nil
}

// UpdateApplicationStatus handles the business logic for updating application status
func (s *volunteerService) UpdateApplicationStatus(id string, request dto.UpdateVolunteerStatusRequest) (*entity.Volunteer, error) {
	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return nil, errors.New("invalid application ID format")
	}

	// Validate status
	newStatus := entity.ApplicationStatus(request.Status)
	if newStatus != entity.ApplicationStatusPending &&
		newStatus != entity.ApplicationStatusApproved &&
		newStatus != entity.ApplicationStatusRejected {
		return nil, errors.New("invalid status. Must be 'pending', 'approved', or 'rejected'")
	}

	// Check if application exists
	volunteer, err := s.volunteerRepo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("volunteer application not found")
		}
		return nil, errors.New("failed to retrieve volunteer application")
	}

	// Update status
	if err := s.volunteerRepo.UpdateStatus(id, newStatus); err != nil {
		return nil, errors.New("failed to update application status")
	}

	// Retrieve updated volunteer
	volunteer.Status = newStatus

	return volunteer, nil
}
