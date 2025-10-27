package model

import "time"


type Donation struct {
	ID             string    `json:"id" db:"id"`
	OrderID        string    `json:"order_id" db:"order_id"`
	UserID         *string   `json:"user_id" db:"user_id"` 
	Amount         int64     `json:"amount" db:"amount"` 
	Status         string    `json:"status" db:"status"`
	PaymentGateway string    `json:"payment_gateway" db:"payment_gateway"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}