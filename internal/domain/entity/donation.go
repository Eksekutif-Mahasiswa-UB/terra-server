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
	ID        string         `db:"id"`
	OrderID   string         `db:"order_id"`
	UserID    *string        `db:"user_id"`
	Amount    int64          `db:"amount"`
	Status    DonationStatus `db:"status"`
	PaidBy    string         `db:"paid_by"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}
