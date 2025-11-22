package dto

import (
	"time"

	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/domain/entity"
)

// CreateDonationRequest represents the request to create a new donation
type CreateDonationRequest struct {
	ProgramID     string  `json:"program_id" binding:"required"`
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	PaymentMethod string  `json:"payment_method" binding:"required"`
	ProofImage    string  `json:"proof_image"`
}

// UpdateDonationStatusRequest represents the request to update a donation status
type UpdateDonationStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending paid failed"`
}

// DonationResponse represents the donation data transfer object for API responses
type DonationResponse struct {
	ID            string                `json:"id"`
	UserID        string                `json:"user_id"`
	ProgramID     string                `json:"program_id"`
	Amount        float64               `json:"amount"`
	PaymentMethod string                `json:"payment_method"`
	Status        entity.DonationStatus `json:"status"`
	ProofImage    string                `json:"proof_image"`
	CreatedAt     time.Time             `json:"created_at"`
	UpdatedAt     time.Time             `json:"updated_at"`
}
