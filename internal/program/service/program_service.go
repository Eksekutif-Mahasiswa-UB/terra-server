package service

import (
	"database/sql"
	"errors"

	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/domain/dto"
	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/domain/entity"
	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/program/repository"
	"github.com/google/uuid"
)

// ProgramService defines the interface for program business logic
type ProgramService interface {
	CreateProgram(request dto.CreateProgramRequest) (*entity.Program, error)
	GetAllPrograms() ([]entity.Program, error)
	GetProgramByID(id string) (*entity.Program, error)
	UpdateProgram(id string, request dto.UpdateProgramRequest) (*entity.Program, error)
	DeleteProgram(id string) error
}

// programService is the concrete implementation of ProgramService
type programService struct {
	programRepo repository.ProgramRepository
}

// NewProgramService creates a new instance of ProgramService
func NewProgramService(programRepo repository.ProgramRepository) ProgramService {
	return &programService{programRepo: programRepo}
}

// CreateProgram handles the business logic for creating a new program
func (s *programService) CreateProgram(request dto.CreateProgramRequest) (*entity.Program, error) {
	// Validate target amount is positive
	if request.TargetAmount <= 0 {
		return nil, errors.New("target amount must be greater than 0")
	}

	// Create new program entity with generated UUID
	newProgram := &entity.Program{
		ID:           uuid.NewString(),
		Title:        request.Title,
		Description:  request.Description,
		ImageURL:     request.ImageURL,
		TargetAmount: request.TargetAmount,
	}

	// Save to database
	if err := s.programRepo.CreateProgram(newProgram); err != nil {
		return nil, errors.New("failed to create program")
	}

	return newProgram, nil
}

// GetAllPrograms retrieves all programs
func (s *programService) GetAllPrograms() ([]entity.Program, error) {
	programs, err := s.programRepo.GetAllPrograms()
	if err != nil {
		return nil, errors.New("failed to retrieve programs")
	}

	return programs, nil
}

// GetProgramByID retrieves a program by its ID
func (s *programService) GetProgramByID(id string) (*entity.Program, error) {
	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return nil, errors.New("invalid program ID format")
	}

	program, err := s.programRepo.GetProgramByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("program not found")
		}
		return nil, errors.New("failed to retrieve program")
	}

	return program, nil
}

// UpdateProgram handles the business logic for updating a program
func (s *programService) UpdateProgram(id string, request dto.UpdateProgramRequest) (*entity.Program, error) {
	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return nil, errors.New("invalid program ID format")
	}

	// Validate target amount is positive
	if request.TargetAmount <= 0 {
		return nil, errors.New("target amount must be greater than 0")
	}

	// Check if program exists
	existingProgram, err := s.programRepo.GetProgramByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("program not found")
		}
		return nil, errors.New("failed to retrieve program")
	}

	// Update program fields
	existingProgram.Title = request.Title
	existingProgram.Description = request.Description
	existingProgram.ImageURL = request.ImageURL
	existingProgram.TargetAmount = request.TargetAmount

	// Save updates to database
	if err := s.programRepo.UpdateProgram(existingProgram); err != nil {
		return nil, errors.New("failed to update program")
	}

	return existingProgram, nil
}

// DeleteProgram handles the business logic for deleting a program
func (s *programService) DeleteProgram(id string) error {
	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return errors.New("invalid program ID format")
	}

	// Check if program exists
	_, err := s.programRepo.GetProgramByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("program not found")
		}
		return errors.New("failed to retrieve program")
	}

	// Delete program
	if err := s.programRepo.DeleteProgram(id); err != nil {
		return errors.New("failed to delete program")
	}

	return nil
}
