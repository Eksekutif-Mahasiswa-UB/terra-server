package entity

import "time"

// DonationStatus represents the status of a donation
type DonationStatus string

const (
	DonationStatusPending DonationStatus = "pending"
	DonationStatusSuccess DonationStatus = "success"
	DonationStatusFailed  DonationStatus = "failed"
)

// Donation represents a donation entity in the system
type Donation struct {
	ID        string         `db:"id" json:"id"`
	OrderID   string         `db:"order_id" json:"order_id"`
	UserID    *string        `db:"user_id" json:"user_id,omitempty"`
	Amount    int64          `db:"amount" json:"amount"`
	Status    DonationStatus `db:"status" json:"status"`
	PaidBy    string         `db:"paid_by" json:"paid_by"`
	CreatedAt time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt time.Time      `db:"updated_at" json:"updated_at"`
}
