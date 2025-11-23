package service

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/domain/dto"
	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/domain/entity"
	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/donation/repository"
	programRepo "github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/program/repository"
	"github.com/google/uuid"
)

// DonationService defines the interface for donation business logic
type DonationService interface {
	CreateDonation(userID string, request dto.CreateDonationRequest) (*entity.Donation, error)
	GetMyDonations(userID string) ([]entity.Donation, error)
	GetAllDonations() ([]entity.Donation, error)
	GetDonationByID(id string) (*entity.Donation, error)
	UpdateDonationStatus(id string, request dto.UpdateDonationStatusRequest) (*entity.Donation, error)
}

// donationService is the concrete implementation of DonationService
type donationService struct {
	donationRepo repository.DonationRepository
	programRepo  programRepo.ProgramRepository
}

// NewDonationService creates a new instance of DonationService
func NewDonationService(donationRepo repository.DonationRepository, programRepo programRepo.ProgramRepository) DonationService {
	return &donationService{
		donationRepo: donationRepo,
		programRepo:  programRepo,
	}
}

// CreateDonation handles the business logic for creating a new donation
func (s *donationService) CreateDonation(userID string, request dto.CreateDonationRequest) (*entity.Donation, error) {
	// Validate UUID format
	if _, err := uuid.Parse(userID); err != nil {
		return nil, errors.New("invalid user ID format")
	}

	if _, err := uuid.Parse(request.ProgramID); err != nil {
		return nil, errors.New("invalid program ID format")
	}

	// Validate required fields
	if strings.TrimSpace(request.ProgramID) == "" {
		return nil, errors.New("program ID is required")
	}

	if request.Amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	if strings.TrimSpace(request.PaymentMethod) == "" {
		return nil, errors.New("payment method is required")
	}

	// Validate that the program exists
	program, err := s.programRepo.GetProgramByID(request.ProgramID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("program not found")
		}
		return nil, errors.New("failed to validate program")
	}

	if program == nil {
		return nil, errors.New("program not found")
	}

	// Create new donation
	donation := &entity.Donation{
		ID:            uuid.NewString(),
		UserID:        userID,
		ProgramID:     request.ProgramID,
		Amount:        request.Amount,
		PaymentMethod: request.PaymentMethod,
		Status:        entity.DonationStatusPending, // Default status
		ProofImage:    request.ProofImage,
	}

	// Save to database
	if err := s.donationRepo.Create(donation); err != nil {
		return nil, errors.New("failed to create donation")
	}

	return donation, nil
}

// GetMyDonations retrieves all donations made by a specific user
func (s *donationService) GetMyDonations(userID string) ([]entity.Donation, error) {
	// Validate UUID format
	if _, err := uuid.Parse(userID); err != nil {
		return nil, errors.New("invalid user ID format")
	}

	donations, err := s.donationRepo.GetByUserID(userID)
	if err != nil {
		return nil, errors.New("failed to retrieve your donations")
	}

	return donations, nil
}

// GetAllDonations retrieves all donations from the database
func (s *donationService) GetAllDonations() ([]entity.Donation, error) {
	donations, err := s.donationRepo.GetAll()
	if err != nil {
		return nil, errors.New("failed to retrieve donations")
	}

	return donations, nil
}

// GetDonationByID retrieves a donation by its ID
func (s *donationService) GetDonationByID(id string) (*entity.Donation, error) {
	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return nil, errors.New("invalid donation ID format")
	}

	donation, err := s.donationRepo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("donation not found")
		}
		return nil, errors.New("failed to retrieve donation")
	}

	return donation, nil
}

// UpdateDonationStatus handles the business logic for updating donation status
func (s *donationService) UpdateDonationStatus(id string, request dto.UpdateDonationStatusRequest) (*entity.Donation, error) {
	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return nil, errors.New("invalid donation ID format")
	}

	// Validate status
	newStatus := entity.DonationStatus(request.Status)
	if newStatus != entity.DonationStatusPending &&
		newStatus != entity.DonationStatusPaid &&
		newStatus != entity.DonationStatusFailed {
		return nil, errors.New("invalid status. Must be 'pending', 'paid', or 'failed'")
	}

	// Check if donation exists
	donation, err := s.donationRepo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("donation not found")
		}
		return nil, errors.New("failed to retrieve donation")
	}

	// Update status
	if err := s.donationRepo.UpdateStatus(id, newStatus); err != nil {
		return nil, errors.New("failed to update donation status")
	}

	// Retrieve updated donation
	donation.Status = newStatus

	return donation, nil
}
