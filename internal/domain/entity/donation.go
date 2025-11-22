package entity

import "time"

// DonationStatus represents the status of a donation
type DonationStatus string

const (
	DonationStatusPending DonationStatus = "pending"
	DonationStatusPaid    DonationStatus = "paid"
	DonationStatusFailed  DonationStatus = "failed"
)

// Donation represents a donation transaction entity in the system
type Donation struct {
	ID            string         `db:"id" json:"id"`
	UserID        string         `db:"user_id" json:"user_id"`
	ProgramID     string         `db:"program_id" json:"program_id"`
	Amount        float64        `db:"amount" json:"amount"`
	PaymentMethod string         `db:"payment_method" json:"payment_method"`
	Status        DonationStatus `db:"status" json:"status"`
	ProofImage    string         `db:"proof_image" json:"proof_image"`
	CreatedAt     time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time      `db:"updated_at" json:"updated_at"`
}
