package dto

import (
	"time"

	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/domain/entity"
)

// DonationResponse represents the donation data transfer object for API responses
type DonationResponse struct {
	ID        string                `json:"id"`
	OrderID   string                `json:"order_id"`
	UserID    *string               `json:"user_id,omitempty"`
	Amount    int64                 `json:"amount"`
	Status    entity.DonationStatus `json:"status"`
	PaidBy    string                `json:"paid_by"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
}
